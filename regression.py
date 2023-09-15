#!/usr/bin/env python3
# Copyright © 2023 Mark Summerfield. All rights reserved.
# License: GPLv3

import contextlib
import filecmp
import os
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
    if not os.path.isfile(EXE):
        subprocess.run(['bash', 'build.sh'])
    count = 0
    ok = 0
    for name in sorted(os.listdir(os.path.join(ROOT, 'eg'))):
        if name.endswith('.urm'):
            count += 1
            ok += check(name, verbose)
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


def delete_file(filename):
    with contextlib.suppress(FileNotFoundError):
        os.remove(filename)


if __name__ == '__main__':
    main()
