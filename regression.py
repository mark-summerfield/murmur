#!/usr/bin/env python3
# Copyright © 2023 Mark Summerfield. All rights reserved.
# License: GPLv3

import contextlib
import filecmp
import os
import random
import re
import shutil
import subprocess
import sys

ROOT = os.path.dirname(os.path.realpath(__file__))
EXE = 'murmur'

ALLARGS_FOR_NAME = {
    'add.urm': (('19', '27', '46'), ('91', '17', '108')),
    'sub.urm': (('43', '19', '24'), ('66', '59', '7')),
    'mul.urm': (('7', '8', '56'), ('9', '9', '81')),
    'lt.urm': (('99', '5', '5', '0'), ('99', '13', '17', '1'),
               ('99', '14', '12', '0')),
    'lte.urm': (('99', '5', '5', '1'), ('99', '13', '17', '1'),
                ('99', '14', '12', '0')),
    'max.urm': (('99', '52', '52', '52'), ('99', '23', '81', '81'),
                ('99', '75', '74', '75'))
    }


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
        if name.endswith('.urm'):
            if name[0].isdecimal():
                count += 1
                ok += check(name, verbose)
            elif name == 'index.urm':
                count += 1
                ok += check_index(name, verbose)
            else:
                allargs = ALLARGS_FOR_NAME.get(name)
                if allargs is None:
                    continue
                for args in allargs:
                    count += 1
                    ok += check_func(name, args, verbose)
    if ok == count:
        print(f'All {count} OK')
        shutil.rmtree(actual_dir)
        subprocess.run(['bash', 'deploy.sh']) # redeploy if good
    else:
        print(f'{ok}/{count} FAIL')


def check(name, verbose):
    tdata = name.replace('.urm', '.txt')
    filename = os.path.join(ROOT, 'eg', name)
    actual = os.path.join(ROOT, 'tdata/actual', tdata)
    expected = os.path.join(ROOT, 'tdata/expected', tdata)
    delete_file(actual)
    cmd = [EXE, '-m', '10000', '-ds', filename]
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


def check_func(name, args, verbose):
    rx = re.compile(r'\bA:(?P<a>\d+)')
    filename = os.path.join(ROOT, 'eg', name)
    a = args[-1]
    args = args[:-1]
    cmd = ['./murmur', '-wA', filename, *args]
    output = subprocess.check_output(cmd)
    ans = ''
    for line in output.decode('utf-8').splitlines():
        if line.startswith('#0'):
            continue
        match = rx.search(line)
        if match is not None:
            ans = match.group('a')
    ok = ans == a
    if verbose:
        if ok:
            print(f'{name} OK')
        else:
            print(f'{name} FAIL actual != expected', cmd)
    return ok


def check_index(name, verbose):
    array = list(range(1, 11))
    random.shuffle(array)
    expected = [x + 2 for x in array]
    filename = os.path.join(ROOT, 'eg', name)
    # 2 is the initial value for I, i.e., the first reg index of the array
    cmd = ['./murmur', '-w2-11', filename, '2', *(str(x) for x in array)]
    output = subprocess.check_output(cmd)
    rx = re.compile(r'#\d+ DATA:(\d+) 3:(\d+) 4:(\d+) 5:(\d+) 6:(\d+) '
                    r'7:(\d+) 8:(\d+) 9:(\d+) 10:(\d+) 11:(\d+)')
    ans = [-1 for _ in range(10)] 
    for line in output.decode('utf-8').splitlines():
        if line.startswith('#0'):
            continue
        match = rx.search(line)
        if match is not None:
            for i in range(10):
                ans[i] = int(match.group(i + 1))
    ok = ans == expected
    if verbose:
        if ok:
            print(f'{name} OK')
        else:
            print(f'{name} FAIL actual {ans} != expected {expected}', cmd)
    return ok


def delete_file(filename):
    with contextlib.suppress(FileNotFoundError):
        os.remove(filename)


if __name__ == '__main__':
    main()
