#!/usr/bin/python

import getopt
import os
import shutil 
import subprocess
import sys

def main(argv):
    module = '00'

    try:
        opts, args = getopt.getopt(argv,"hm:",["module="])
    except getopt.GetoptError:
        print('test.py -m <module>')
        sys.exit(2)
    for opt, arg in opts:
        if opt == '-h':
            print('test.py -m <module>')
            sys.exit()
        elif opt in ("-m", "--module"):
            module = arg.rjust(2, '0')

    if module == '00':
        print(f'Invalid module number specified')
        print('test.py -m <module>')
        sys.exit()

    print(f'Setting up module {module} in "current" directory')

    src = f'.course/00_Before/{module}' 
    dest = 'current'

    if os.path.exists(dest) and os.path.isdir(dest):
        shutil.rmtree(dest)   
        
    shutil.copytree(src, dest) 
    shutil.copy('runner.py', 'current/.')

    # p = subprocess.Popen(['go', 'mod', 'init'], cwd='current')
    # p.wait()

    p = subprocess.Popen(['go', 'mod', 'tidy'], cwd='current')
    p.wait()

    print(f'Completed setup of module {module} in "current" directory')


if __name__ == "__main__":
    main(sys.argv[1:])
