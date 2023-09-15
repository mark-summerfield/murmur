#!/usr/bin/env python3
# Copyright © 2023 Mark Summerfield. All rights reserved.
# License: GPLv3

import contextlib
import filecmp
import os
import re
import shutil
import subprocess
import sys

ROOT = os.path.dirname(os.path.realpath(__file__))
EXE = 'murmur'


def main():
    verbose = True
    if len(sys.argv) > 1 and sys.argv[1] in {'-q', '--quiet'}:
        verbose = False

    actual_dir = 'tdata/actual'
    os.makedirs(actual_dir, exist_ok=True)
    subprocess.run(['bash', 'build.sh']) # always build fresh
    count = 0
    ok = 0
    for name in sorted(os.listdir(os.path.join(ROOT, 'eg'))):
        if name.endswith('.urm') and name[0].isdecimal():
            count += 1
            ok += check(name, verbose)
        elif name == "lt.urm":
            count += 1
            ok += check_lt(name, verbose)
    if ok == count:
        print(f'All {count} OK')
        shutil.rmtree(actual_dir)
    else:
        print(f'{ok}/{count} FAIL')


def check(name, verbose):
    tdata = name.replace('.urm', '.txt')
    filename = os.path.join(ROOT, 'eg', name)
    actual = os.path.join(ROOT, 'tdata/actual', tdata)
    expected = os.path.join(ROOT, 'tdata/expected', tdata)
    delete_file(actual)
    cmd = ['./murmur', '-m', '10000', '-ds', filename]
    output = subprocess.check_output(cmd)
    with open(actual, 'wb') as file:
        file.write(output)
    ok = filecmp.cmp(actual, expected, shallow=False)
    if verbose:
        if ok:
            print(f'{name} OK')
        else:
            print(f'{name} FAIL actual != expected')
    return int(ok)


def check_lt(name, verbose):
    rx = re.compile(r'\bA:(?P<a>\d+)')
    filename = os.path.join(ROOT, 'eg', name)
    for args in (['5', '5', '0'], ['13', '17', '1'], ['14', '12', '0']):
        a = args[-1]
        args = ['99'] + args[:-1]
    cmd = ['./murmur', '-wA', filename, *args]
    output = subprocess.check_output(cmd)
    ans = ''
    for line in output.decode('utf-8').splitlines():
        match = rx.search(line)
        if match is not None:
            ans = match.group('a')
    ok = ans == a
    if verbose:
        if ok:
            print(f'{name} OK')
        else:
            print(f'{name} FAIL actual != expected')
    return ok


def delete_file(filename):
    with contextlib.suppress(FileNotFoundError):
        os.remove(filename)


if __name__ == '__main__':
    main()
