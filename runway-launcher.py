import os
import os.path as p
import subprocess
import sys
from subprocess import Popen
from typing import Optional

executable = "runway"
makefile = "Makefile"


def launch_query(args: list[str], program_name: str) -> Optional[Popen[bytes]]:
    if not check_executable(program_name):
        if not check_executable(makefile):
            sys.stderr.write("Missing make file")
            return
        val = subprocess.call(["make", "build"])
        if val != 0 or not check_executable(program_name):
            sys.stderr.write("error in build")
            print("build failed")
            return

    return Popen([os.getcwd()+"/"+program_name] + args)


def check_executable(name: str) -> bool:
    return p.isfile(os.getcwd()+"/"+name)


if __name__ == "__main__":
    proc = launch_query(args=[], program_name=executable)
    print(proc.returncode)
    proc.kill()
