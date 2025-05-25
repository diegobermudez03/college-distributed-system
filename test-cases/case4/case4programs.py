import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS
#### SEMESTRE 2024-30


facAddress = os.getenv('FAC_ADDRESS', '127.0.0.1')


program_commands = [
    f"./program --name=psicologia --semester=2024-30 --classrooms=5 --labs=2 --faculty-server={facAddress}:5001",
   
    f"./program --name=biologia --semester=2024-30 --classrooms=5 --labs=2 --faculty-server={facAddress}:5002",

    f"./program --name=ingenieria-civil --semester=2024-30 --classrooms=5 --labs=2 --faculty-server={facAddress}:5003",
    f"./program --name=ingenieria-electronica --semester=2024-30 --classrooms=5 --labs=2 --faculty-server={facAddress}:5003",
    f"./program --name=ingenieria-de-sistemas --semester=2024-30 --labs=2 --faculty-server={facAddress}:5003",

    f"./program --name=sociologia --semester=2024-30 --classrooms=5 --labs=2 --faculty-server={facAddress}:5001",
    f"./program --name=quimica --semester=2024-30 --classrooms=5 --labs=2 --faculty-server={facAddress}:5002",
]



processes = [subprocess.Popen(cmd, shell=True) for cmd in program_commands]

for p in processes:
    p.wait()
