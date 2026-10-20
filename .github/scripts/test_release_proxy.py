"""Exercise release shell steps without network access, Git writes or real delays."""

import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import textwrap
import unittest


WORKFLOW = Path(__file__).resolve().parents[1] / "workflows" / "release.yml"
SHA = "a" * 40
MODULE = "github.com/QoderAI/qoder-cloud-agents-sdk-go"
PROXY_MODULE = "github.com/!qoder!a!i/qoder-cloud-agents-sdk-go"
CHECK_STEPS = (
    "Validate main, module version, tag, and proxy state",
    "Revalidate main and publish tag",
)


def step_script(name):
    lines = WORKFLOW.read_text().splitlines()
    start = lines.index(f"      - name: {name}")
    start = lines.index("        run: |", start) + 1
    end = start
    while end < len(lines) and (not lines[end].strip() or lines[end].startswith("          ")):
        end += 1
    return textwrap.dedent("\n".join(lines[start:end]))


FAKE_COMMAND = r"""
import json
import os
from pathlib import Path
import sys

command = Path(sys.argv[0]).name
args = sys.argv[1:]
root = Path(os.environ["RUNNER_TEMP"])
with (root / "commands.jsonl").open("a") as log:
    log.write(json.dumps([command, *args]) + "\n")
tag_exists = os.environ["TEST_TAG_EXISTS"] == "true"
if command == "git":
    if args[0] == "rev-parse":
        if "-q" in args and not tag_exists:
            sys.exit(1)
        print(os.environ["COMMIT_SHA"])
    elif args[0] == "show-ref":
        sys.exit(0 if tag_exists else 1)
    elif args[0] == "cat-file":
        print("tag")
    elif args[0] == "rev-list":
        print(os.environ["COMMIT_SHA"])
    elif args[0] == "for-each-ref":
        print(f'Release v{os.environ["RELEASE_VERSION"]} (batch_id: {os.environ["BATCH_ID"]})')
elif command == "go":
    print(os.environ["MODULE_PATH"])
elif command == "curl":
    count_file = root / "requests"
    count = int(count_file.read_text()) + 1 if count_file.exists() else 1
    count_file.write_text(str(count))
    ready = int(os.environ["TEST_READY_AFTER"])
    status = "200" if ready and count > ready else os.environ["TEST_HTTP_STATUS"]
    Path(args[args.index("--output") + 1]).write_text(
        '{"Version":"v0.7.0"}' if status == "200" else "not found: unknown revision v0.7.0"
    )
    if "--dump-header" in args:
        Path(args[args.index("--dump-header") + 1]).write_text(
            f"HTTP/2 {status}\ncache-control: public, max-age=1800\n"
        )
    print(status, end="")
"""


class ReleaseProxyTests(unittest.TestCase):
    def run_step(self, name, *, tag_exists=False, status="404", ready_after=0, version="0.7.0"):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            for command in ("git", "go", "make", "curl", "sleep"):
                path = root / command
                path.write_text(f"#!{sys.executable}\n{FAKE_COMMAND}")
                path.chmod(0o755)
            env = os.environ | {
                "PATH": f"{root}{os.pathsep}{os.environ['PATH']}",
                "RUNNER_TEMP": directory,
                "COMMIT_SHA": SHA,
                "RELEASE_VERSION": version,
                "BATCH_ID": "proxy-test",
                "MODULE_PATH": MODULE,
                "PROXY_MODULE_PATH": PROXY_MODULE,
                "TEST_TAG_EXISTS": str(tag_exists).lower(),
                "TEST_HTTP_STATUS": status,
                "TEST_READY_AFTER": str(ready_after),
            }
            result = subprocess.run(
                ["bash", "-c", step_script(name)], env=env, cwd=root,
                capture_output=True, text=True, timeout=30,
            )
            calls = [json.loads(line) for line in (root / "commands.jsonl").read_text().splitlines()]
            return result, calls

    def test_unpublished_version_uses_cached_only_and_can_create_tag(self):
        for name in CHECK_STEPS:
            with self.subTest(step=name):
                result, calls = self.run_step(name, version="0.8.0-RC.1")
                self.assertEqual(result.returncode, 0, result.stderr)
                urls = [call[-1] for call in calls if call[0] == "curl"]
                self.assertEqual(urls, [
                    f"https://proxy.golang.org/cached-only/{PROXY_MODULE}/@v/v0.8.0-!r!c.1.info"
                ])
                if name == CHECK_STEPS[1]:
                    self.assertTrue(any(call[:2] == ["git", "push"] for call in calls))

    def test_cached_publication_without_tag_is_rejected(self):
        for name in CHECK_STEPS:
            with self.subTest(step=name):
                result, calls = self.run_step(name, status="200")
                self.assertNotEqual(result.returncode, 0)
                self.assertIn("remote tag is absent", result.stderr)
                self.assertFalse(any(call[:2] == ["git", "push"] for call in calls))

    def test_matching_tag_is_reused_without_push(self):
        for status in ("200", "404", "410"):
            with self.subTest(status=status):
                result, calls = self.run_step(CHECK_STEPS[1], tag_exists=True, status=status)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertIn("reusing it", result.stdout)
                self.assertFalse(any(call[:2] == ["git", "push"] for call in calls))

    def test_unexpected_proxy_error_stops_before_publish(self):
        for name in CHECK_STEPS:
            with self.subTest(step=name):
                result, calls = self.run_step(name, status="503")
                self.assertNotEqual(result.returncode, 0)
                self.assertIn("HTTP 503", result.stderr)
                self.assertFalse(any(call[:2] == ["git", "push"] for call in calls))

    def test_verification_survives_thirty_minute_negative_cache(self):
        result, calls = self.run_step("Wait for the public Go proxy", ready_after=61)
        self.assertEqual(result.returncode, 0, result.stderr)
        urls = [call[-1] for call in calls if call[0] == "curl"]
        self.assertEqual(len(urls), 62)
        self.assertEqual(set(urls), {
            f"https://proxy.golang.org/{PROXY_MODULE}/@v/v0.7.0.info"
        })
        self.assertGreater(sum(int(call[1]) for call in calls if call[0] == "sleep"), 1800)

    def test_unavailable_version_fails_with_proxy_diagnostics(self):
        result, calls = self.run_step("Wait for the public Go proxy")
        self.assertNotEqual(result.returncode, 0)
        delays = [int(call[1]) for call in calls if call[0] == "sleep"]
        self.assertEqual(sum(delays), 2100)
        self.assertIn("cache-control: public, max-age=1800", result.stderr)
        self.assertIn("unknown revision v0.7.0", result.stderr)


if __name__ == "__main__":
    unittest.main()
