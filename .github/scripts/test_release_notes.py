import io
import os
import subprocess
import tempfile
import unittest
import urllib.error
from pathlib import Path
from unittest.mock import patch

import release_notes as notes

CHANGELOG = """# Changelog

## [Unreleased]

### Added
- Forward: support the next API.

## [0.2.0] - 2026-09-23

### Fixed
- Managed: preserve streamed events after reconnecting. (#42)

### Migration
- Replace the old option with `new_option`.

```python
## [0.1.0]
client = Client(new_option=True)
```

## [0.1.0]

### Added
- Initial SDK clients.
"""
SHA = "a" * 40
BODY = "### Fixed\n\n- Preserve streamed events.\n"
URL = "https://github.com/example/sdk/releases/tag/v0.2.0"


class ChangelogTests(unittest.TestCase):
    def test_extracts_only_target_and_keeps_migration_code(self):
        body = notes.extract_notes(CHANGELOG, "0.2.0")
        self.assertIn("Managed: preserve streamed events", body)
        self.assertIn("## [0.1.0]\nclient = Client", body)
        self.assertNotIn("Initial SDK clients", body)
        self.assertNotIn("support the next API", body)
        self.assertNotIn("2026-09-23", body)

    def test_empty_unreleased_and_optional_date(self):
        entries = notes.parse_changelog("## [Unreleased]\n\n## [1.0.0]\n- Add clients.\n")
        self.assertEqual(entries["1.0.0"], "- Add clients.\n")

    def test_supports_each_ecosystems_canonical_version(self):
        for version in ("0.2.0-rc.1", "0.2.0rc1", "0.2.0.dev1", "0.2.0.post1", "1!0.2.0"):
            with self.subTest(version=version):
                text = f"## [Unreleased]\n\n## [{version}]\n- Fix a regression.\n"
                self.assertEqual(notes.extract_notes(text, version), "- Fix a regression.\n")

    def test_rejects_missing_target_or_unreleased_request(self):
        for version in ("0.3.0", "Unreleased", "v0.2.0", "0.2.0\n--help"):
            with self.subTest(version=version), self.assertRaises(ValueError):
                notes.extract_notes(CHANGELOG, version)

    def test_rejects_duplicate_versions_and_unreleased(self):
        for heading in ("0.2.0", "Unreleased"):
            with self.subTest(heading=heading), self.assertRaisesRegex(ValueError, "Duplicate"):
                notes.parse_changelog(CHANGELOG + f"\n## [{heading}]\n- Another note.\n")

    def test_rejects_empty_and_placeholder_entries(self):
        for body in (
            "",
            "### Fixed",
            "<!-- - Fix a bug. -->",
            "- TODO",
            "- TBD: explain",
            "- ...",
            "- Pending",
            "- <describe the change>",
            "```\n- Only an example.\n```",
            "- Real change.\n- TBD",
        ):
            with self.subTest(body=body), self.assertRaises(ValueError):
                notes.parse_changelog(f"## [Unreleased]\n\n## [0.2.0]\n{body}\n")

    def test_omits_author_comments(self):
        body = notes.extract_notes(
            CHANGELOG.replace("### Fixed", "<!-- TODO: private drafting hint -->\n### Fixed"), "0.2.0"
        )
        self.assertNotIn("TODO", body)

    def test_rejects_bad_structure_and_dates(self):
        for text in (
            "## [0.2.0]\n- Fix a regression.",
            "## [Unreleased]\n## 0.2.0\n- Fix a regression.",
            "## [Unreleased] - 2026-09-23",
            "## [Unreleased]\n## [0.2.0] - 2026-02-30\n- Fix a regression.",
            "## [Unreleased]\n<!-- unfinished comment",
            "## [Unreleased]\n```\nunclosed example",
        ):
            with self.subTest(text=text), self.assertRaises(ValueError):
                notes.parse_changelog(text)

    def test_cli_extract_and_missing_version_exit_status(self):
        script = str(Path(notes.__file__).resolve())
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory)
            (path / "CHANGELOG.md").write_text(CHANGELOG, encoding="utf-8")
            command = [os.sys.executable, script, "extract", "--version"]
            result = subprocess.run(
                command + ["0.2.0", "--output", "notes.md"], cwd=path, capture_output=True, text=True
            )
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertEqual((path / "notes.md").read_text(), notes.extract_notes(CHANGELOG, "0.2.0"))
            result = subprocess.run(command + ["0.3.0"], cwd=path, capture_output=True, text=True)
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("has no entry", result.stderr)


class PublishTests(unittest.TestCase):
    def setUp(self):
        self.env = patch.dict(os.environ, {"GITHUB_REPOSITORY": "example/sdk", "GH_TOKEN": "test-token"})
        self.env.start()
        self.addCleanup(self.env.stop)
        self.lookup = patch.object(notes, "existing_release", return_value=None).start()
        self.run = patch.object(notes.subprocess, "run").start()
        self.output = patch.object(notes.subprocess, "check_output", side_effect=self.command_output).start()
        patch.dict(os.environ, {"GITHUB_STEP_SUMMARY": ""}).start()
        self.addCleanup(patch.stopall)
        self.created_notes = None

    def command_output(self, command, **kwargs):
        if command[0] == "gh":
            self.created_notes = Path(command[command.index("--notes-file") + 1]).read_text()
            return URL + "\n"
        if command[1] == "cat-file":
            return "tag\n"
        return SHA + "\n"

    def release(self, **changes):
        return dict(tag_name="v0.2.0", name="v0.2.0", body=BODY, draft=False, prerelease=False, html_url=URL) | changes

    def gh_calls(self):
        return [call.args[0] for call in self.output.call_args_list if call.args[0][0] == "gh"]

    def test_creates_release_using_exact_notes_and_existing_tag(self):
        notes.publish_notes("0.2.0", SHA, False, BODY)
        (command,) = self.gh_calls()
        self.assertIn("--verify-tag", command)
        self.assertEqual(command[command.index("--target") + 1], SHA)
        self.assertNotIn("--prerelease", command)
        self.assertEqual(self.created_notes, BODY)
        self.run.assert_called_once_with(
            ["git", "fetch", "--force", "origin", "refs/tags/v0.2.0:refs/tags/v0.2.0"],
            check=True,
        )

    def test_marks_prerelease_without_promoting_latest(self):
        notes.publish_notes("0.2.0-rc.1", SHA, True, BODY)
        (command,) = self.gh_calls()
        self.assertIn("--prerelease", command)
        self.assertIn("--latest=false", command)

    def test_matching_rerun_has_no_release_write(self):
        self.lookup.return_value = self.release()
        notes.publish_notes("0.2.0", SHA, False, BODY)
        self.assertFalse(self.gh_calls())

    def test_conflicts_fail_without_overwriting_release(self):
        for change in (
            {"body": "Different notes"},
            {"draft": True},
            {"prerelease": True},
            {"tag_name": "v0.3.0"},
            {"name": "Other"},
        ):
            with self.subTest(change=change):
                self.lookup.return_value = self.release(**change)
                with self.assertRaisesRegex(ValueError, "differs"):
                    notes.publish_notes("0.2.0", SHA, False, BODY)
                self.assertFalse(self.gh_calls())

    def test_wrong_checkout_stops_before_tag_fetch_or_api(self):
        self.output.return_value = "b" * 40
        self.output.side_effect = None
        with self.assertRaisesRegex(ValueError, "approved commit"):
            notes.publish_notes("0.2.0", SHA, False, BODY)
        self.run.assert_not_called()
        self.lookup.assert_not_called()

    def test_wrong_tag_commit_stops_before_api(self):
        self.output.side_effect = [SHA, "tag", "b" * 40]
        with self.assertRaisesRegex(ValueError, "does not point"):
            notes.publish_notes("0.2.0", SHA, False, BODY)
        self.lookup.assert_not_called()

    def test_missing_remote_tag_stops_before_api(self):
        self.run.side_effect = subprocess.CalledProcessError(1, ["git", "fetch"])
        with self.assertRaises(subprocess.CalledProcessError):
            notes.publish_notes("0.2.0", SHA, False, BODY)
        self.lookup.assert_not_called()

    def test_lightweight_tag_stops_before_api(self):
        self.output.side_effect = [SHA, "commit"]
        with self.assertRaisesRegex(ValueError, "annotated"):
            notes.publish_notes("0.2.0", SHA, False, BODY)
        self.lookup.assert_not_called()

    def test_lookup_failure_does_not_create_release(self):
        self.lookup.side_effect = urllib.error.URLError("connection failed")
        with self.assertRaises(urllib.error.URLError):
            notes.publish_notes("0.2.0", SHA, False, BODY)
        self.assertFalse(self.gh_calls())


class LookupTests(unittest.TestCase):
    @patch.dict(os.environ, {"GH_TOKEN": "test-token"})
    @patch.object(notes.urllib.request, "urlopen")
    def test_only_404_means_release_missing(self, urlopen):
        for status in (404, 401, 403, 500):
            with self.subTest(status=status):
                urlopen.side_effect = urllib.error.HTTPError(
                    "https://api.github.com", status, "error", {}, io.BytesIO()
                )
                if status == 404:
                    self.assertIsNone(notes.existing_release("example/sdk", "v0.2.0"))
                else:
                    with self.assertRaises(urllib.error.HTTPError):
                        notes.existing_release("example/sdk", "v0.2.0")


if __name__ == "__main__":
    unittest.main()
