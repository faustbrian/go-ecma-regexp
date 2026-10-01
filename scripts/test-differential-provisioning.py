#!/usr/bin/env python3
"""Exercise the actual provisioning script's optional Deno version boundary.

Downloads and Go execution are stubbed: this proves guard reachability and
cleanup, not checksum authenticity or differential matching semantics.
"""
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest

SCRIPT = Path(__file__).with_name("run-differential.sh").resolve()
STUB = r'''import os
from pathlib import Path
import sys
name = Path(sys.argv[0]).name
args = sys.argv[1:]
if name == "curl":
    Path(args[args.index("--output") + 1]).touch()
elif name == "shasum":
    digest = ("55a772a600b7bdafb4b35945b3935090e27aff9934b4c11b281220fcd99139d7"
              if args[-1].endswith("node.tar.gz") else
              "d8b96221828ad6f97ac7ac0ab7e95872341af763001e8803e8267652c2652620")
    print(digest + "  " + args[-1])
elif name in ("tar", "unzip"):
    destination = Path(args[args.index("-C" if name == "tar" else "-d") + 1])
    peer = destination / ("node" if name == "tar" else "bun")
    version = "v24.4.1" if name == "tar" else "1.3.14"
    peer.write_text("#!/bin/sh\nprintf '%s\\n' '" + version + "'\n")
    peer.chmod(0o700)
elif name == "uname":
    print("Darwin" if args == ["-s"] else "arm64")
elif name == "deno":
    print(os.environ["DENO_HEADER"], end="")
    sys.exit(int(os.environ.get("DENO_EXIT", "0")))
elif name == "go":
    Path(os.environ["GO_MARKER"]).write_text("reached")
else:
    raise SystemExit("unexpected fixture command")
'''


class DenoProvisioningTest(unittest.TestCase):
    def invoke(self, header, status=0):
        with tempfile.TemporaryDirectory(prefix="ecma-peer-guard-test-") as root:
            root = Path(root)
            commands = root / "bin"
            commands.mkdir()
            temporary = root / "tmp"
            temporary.mkdir()
            marker = root / "go-reached"
            for name in ("mktemp", "chmod", "find", "mkdir", "head"):
                target = shutil.which(name)
                self.assertIsNotNone(target)
                (commands / name).symlink_to(target)
            names = ["curl", "shasum", "tar", "unzip", "uname", "go"]
            if header is not None:
                names.append("deno")
            for name in names:
                stub = commands / name
                stub.write_text("#!" + sys.executable + "\n" + STUB)
                stub.chmod(0o700)
            environment = dict(os.environ, PATH=str(commands), TMPDIR=str(temporary),
                               DENO_HEADER=header or "", DENO_EXIT=str(status),
                               GO_MARKER=str(marker))
            result = subprocess.run([shutil.which("bash"), str(SCRIPT)],
                                    env=environment, capture_output=True, text=True)
            reached = marker.exists()
            self.assertEqual(list(temporary.iterdir()), [], "task directory retained")
            return result, reached

    def test_pinned_official_and_bare_headers_reach_go(self):
        for header in ("deno 2.9.3 (stable, release, aarch64-apple-darwin)\nv8 14.9.207.2-rusty\ntypescript 6.0.3\n",
                       "deno 2.9.3\n"):
            with self.subTest(header=header):
                result, reached = self.invoke(header)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertTrue(reached)

    def test_absent_deno_reaches_go(self):
        result, reached = self.invoke(None)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue(reached)

    def test_wrong_malformed_and_failed_headers_do_not_reach_go(self):
        for header, status in (("deno 2.9.30\n", 0), ("deno 2.9.2\n", 0),
                               ("other 2.9.3\n", 0), ("deno\n", 0), ("", 0),
                               ("deno 2.9.3junk\n", 0),
                               ("deno 2.9.3 (stable, release, broken\n", 0),
                               ("deno 2.9.3\n", 7)):
            with self.subTest(header=header, status=status):
                result, reached = self.invoke(header, status)
                self.assertNotEqual(result.returncode, 0)
                self.assertFalse(reached)


if __name__ == "__main__":
    unittest.main()
