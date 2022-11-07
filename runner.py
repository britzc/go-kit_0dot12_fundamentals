import json
import re
import subprocess
import sys

process = subprocess.Popen(['go', 'test', './...', '-timeout', '1s'],
                     stdout=subprocess.PIPE,
                     stderr=subprocess.PIPE)
stdout, stderr = process.communicate()

if stderr:
    r = re.compile("(.*\.go):(\d+):\d+:(.*)", re.MULTILINE)
    filename, line, message = r.search(stderr.decode("utf-8")).groups()

    print(f"Hmm, it looks like there is an error in `{filename}` on line `{line}`: {message.strip()} <fail>")

if stdout:
    r = re.compile("(?:\w+\.go:\d+|Messages):\s*~\d+\|(.*)~", re.MULTILINE)
    result = r.search(stdout.decode("utf-8"))

    if result:
        message = result.group(1)
        print(f"{message} <fail>")

    if result == None:
        print(f"Well Done! <pass>")

