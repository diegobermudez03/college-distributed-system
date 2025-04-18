import os
import subprocess

facultiesAddress = os.getenv('FAC_ADDRESS')

commands = [
    ##start prograsms, 1 for each faculty
    f"./program --name=Ingenieria-de-Sistemas --semester=2025-10 --classrooms=6 --labs=3 --faculty-server={facultiesAddress}:5001",
    f"./program --name=Derecho-Penal --semester=2025-10 --classrooms=3 --labs=3 --faculty-server={facultiesAddress}:5002",
    f"./program --name=Teatro --semester=2025-10 --classrooms=3 --labs=2 --faculty-server={facultiesAddress}:5003",
]

#we execute all the commands
processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

#wait to all processes to end
for p in processes:
    p.wait()
