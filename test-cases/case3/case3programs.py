import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS
#### SEMESTRE 2024-10

facAddrress = os.getenv('FAC_ADDRESS', '127.0.0.1')

##we will handle the program commands in a separate list, is so that we can re order it random, to check the system with random order
program_commands = [
    f"./program --name=psicologia --semester=2024-10 --classrooms=5 --labs=2 --faculty-server={facAddrress}:5001",

    f"./program --name=ingenieria-civil --semester=2024-10 --classrooms=5 --labs=2 --faculty-server={facAddrress}:5003",
    f"./program --name=ingenieria-electronica --semester=2024-10 --classrooms=5 --labs=2 --faculty-server={facAddrress}:5003",
    f"./program --name=ingenieria-de-sistemas --semester=2024-10 --classrooms=5 --faculty-server={facAddrress}:5003",

    f"./program --name=sociologia --semester=2024-10 --classrooms=5 --labs=2 --faculty-server={facAddrress}:5001",
    f"./program --name=biologia --semester=2024-10 --classrooms=5 --labs=2 --faculty-server={facAddrress}:5002",
]


processes = [subprocess.Popen(cmd, shell=True) for cmd in program_commands]

for p in processes:
    p.wait()
