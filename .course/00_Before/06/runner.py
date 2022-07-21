import json
import re
import subprocess
import sys

process = subprocess.Popen(['go', 'test', '-failfast', '-timeout', '1s'],
                     stdout=subprocess.PIPE,
                     stderr=subprocess.PIPE)
stdout, stderr = process.communicate()

if stderr:
    r = re.compile("(.*\.go):(\d+):\d+:(.*)", re.MULTILINE)
    filename, line, message = r.search(stderr.decode("utf-8")).groups()

    print(f"Hmm, looks like there is an error in `{filename}` on line `{line}`: {message.strip()} <fail>")

if stdout:
    failed = re.compile("(?:\w+\.go:\d+|Messages):\s*~\d+\|(.*)~", re.MULTILINE)
    failed_match = failed.search(stdout.decode("utf-8"))

    if failed_match:
        message = failed_match.group(1)
        print(f"{message} <fail>")

    if re.search("PASS", stdout.decode("utf-8")):
        print(f"Well Done! <pass>")

