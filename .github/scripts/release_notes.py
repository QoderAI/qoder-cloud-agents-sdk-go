#!/usr/bin/env python3
"""Validate changelogs and publish their approved version entry to GitHub."""

import argparse
import datetime
import json
import os
import re
import subprocess
import sys
import tempfile
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

VERSION = r"[0-9][0-9A-Za-z.!+-]*"
HEADING = re.compile(rf"## \[(Unreleased|{VERSION})\](?: - (\d{{4}}-\d{{2}}-\d{{2}}))?")


def outside_fences(text):
    fence = None
    for line in text.splitlines():
        marker = re.match(r"^\s{0,3}(`{3,}|~{3,})", line)
        if fence:
            if re.fullmatch(rf"\s{{0,3}}{re.escape(fence[0])}{{{len(fence)},}}\s*", line):
                fence = None
            yield ""
        elif marker:
            fence = marker[1]
            yield ""
        else:
            yield line
    if fence:
        raise ValueError("Unclosed Markdown code fence")


def validate_entry(version, body):
    visible = "\n".join(outside_fences(body))
    if re.search(r"\b(?:TODO|TBD)\b|<[^>]*(?:describe|placeholder)[^>]*>", visible, re.I):
        raise ValueError(f"{version}: replace placeholder text with release notes")
    items = re.findall(r"^\s*[-*] (.+)$", visible, re.M)
    meaningful = [re.sub(r"[^\w]+", "", item).lower() for item in items]
    if (
        not meaningful
        or not all(meaningful)
        or any(item in {"placeholder", "pending", "comingsoon", "none", "na"} for item in meaningful)
    ):
        raise ValueError(f"{version}: add at least one substantive changelog bullet")


def parse_changelog(text):
    # Comments may guide authors but must not appear in the published notes.
    text = re.sub(r"<!--.*?-->", "", text, flags=re.S)
    if "<!--" in text or "-->" in text:
        raise ValueError("Unclosed or unmatched Markdown comment")
    sections = {}
    current = None
    for line, visible in zip(text.splitlines(), list(outside_fences(text))):
        if re.match(r"^##(?:\s|$)", visible):
            match = HEADING.fullmatch(line)
            if not match:
                raise ValueError(f"Expected '## [version]' or '## [version] - YYYY-MM-DD': {line}")
            current, date = match.groups()
            if current in sections:
                raise ValueError(f"Duplicate changelog section: {current}")
            if date:
                datetime.date.fromisoformat(date)
            if current == "Unreleased" and date:
                raise ValueError("Unreleased must not have a release date")
            sections[current] = []
        elif current is not None:
            sections[current].append(line)
    if not sections or next(iter(sections)) != "Unreleased":
        raise ValueError("The first version section must be '## [Unreleased]'")
    entries = {version: "\n".join(lines).strip() + "\n" for version, lines in sections.items()}
    for version, body in entries.items():
        if version != "Unreleased" or body.strip():
            validate_entry(version, body)
    return entries


def extract_notes(text, version):
    if not re.fullmatch(VERSION, version):
        raise ValueError("Use a version without the v prefix, not Unreleased")
    entries = parse_changelog(text)
    if version not in entries:
        raise ValueError(f"CHANGELOG.md has no entry for {version}; include it in the release PR")
    return entries[version]


def existing_release(repository, tag):
    url = f"https://api.github.com/repos/{repository}/releases/tags/{urllib.parse.quote(tag, safe='')}"
    request = urllib.request.Request(
        url,
        headers={
            "Authorization": f"Bearer {os.environ['GH_TOKEN']}",
            "Accept": "application/vnd.github+json",
            "User-Agent": "sdk-release-notes",
        },
    )
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            return json.load(response)
    except urllib.error.HTTPError as exc:
        if exc.code == 404:
            return None
        raise


def publish_notes(version, commit_sha, prerelease, notes):
    repository = os.environ["GITHUB_REPOSITORY"]
    if not re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", repository):
        raise ValueError("Invalid GITHUB_REPOSITORY")
    if not re.fullmatch(r"[0-9a-f]{40}", commit_sha):
        raise ValueError("commit_sha must be a full lowercase commit SHA")
    tag = f"v{version}"

    def git(*args):
        return subprocess.check_output(["git", *args], text=True).strip()

    if git("rev-parse", "HEAD") != commit_sha:
        raise ValueError("Release notes must be read from the approved commit")
    subprocess.run(["git", "fetch", "--force", "origin", f"refs/tags/{tag}:refs/tags/{tag}"], check=True)
    if git("cat-file", "-t", f"refs/tags/{tag}") != "tag":
        raise ValueError(f"{tag} must be an annotated release tag")
    if git("rev-parse", f"refs/tags/{tag}^{{commit}}") != commit_sha:
        raise ValueError(f"{tag} does not point to the approved commit")

    existing = existing_release(repository, tag)
    if existing is not None:
        if (
            existing.get("tag_name") != tag
            or existing.get("name") != tag
            or existing.get("draft") is not False
            or existing.get("prerelease") is not prerelease
            or (existing.get("body") or "").strip() != notes.strip()
        ):
            raise ValueError(f"Existing GitHub Release {tag} differs from the approved notes or metadata")
        url = existing["html_url"]
        print(f"Reusing GitHub Release: {url}")
    else:
        with tempfile.TemporaryDirectory() as directory:
            notes_file = Path(directory) / "release-notes.md"
            notes_file.write_text(notes, encoding="utf-8")
            command = [
                "gh",
                "release",
                "create",
                tag,
                "--repo",
                repository,
                "--verify-tag",
                "--target",
                commit_sha,
                "--title",
                tag,
                "--notes-file",
                str(notes_file),
            ]
            if prerelease:
                command += ["--prerelease", "--latest=false"]
            url = subprocess.check_output(command, text=True).strip()
        print(f"Created GitHub Release: {url}")
    if os.environ.get("GITHUB_STEP_SUMMARY"):
        with open(os.environ["GITHUB_STEP_SUMMARY"], "a", encoding="utf-8") as summary:
            summary.write(f"\nGitHub Release: {url}\n")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=["check", "extract", "publish"])
    parser.add_argument("--changelog", type=Path, default=Path("CHANGELOG.md"))
    parser.add_argument("--version")
    parser.add_argument("--output", type=Path)
    parser.add_argument("--commit-sha")
    parser.add_argument("--prerelease", choices=["true", "false"])
    args = parser.parse_args()
    text = args.changelog.read_text(encoding="utf-8")
    if args.command == "check":
        parse_changelog(text)
        print("Changelog is valid")
        return
    if not args.version:
        parser.error("--version is required")
    notes = extract_notes(text, args.version)
    if args.command == "extract":
        if args.output:
            args.output.write_text(notes, encoding="utf-8")
        else:
            print(notes, end="")
    else:
        if not args.commit_sha or args.prerelease is None:
            parser.error("publish requires --commit-sha and --prerelease")
        publish_notes(args.version, args.commit_sha, args.prerelease == "true", notes)


if __name__ == "__main__":
    try:
        main()
    except (ValueError, KeyError, OSError, urllib.error.URLError, subprocess.CalledProcessError) as exc:
        sys.exit(f"Release notes: {exc}")
