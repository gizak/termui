#!/usr/bin/env python3

import signal
import sys
from pathlib import Path
from subprocess import call

if __name__ == '__main__':
    signal.signal(signal.SIGINT, lambda sig, frame: sys.exit(0))

    p = Path('.') / '_example'
    files = p.glob('*.go')
    for file in files:
        command = f'go run {file}'
        print(command)
        call(command.split())
