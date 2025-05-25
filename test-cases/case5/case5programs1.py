import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS 5 segundos fail
#### SEMESTRE 2025-10


facAddress = os.getenv('FAC_ADDRESS', '127.0.0.1') 


program_commands = [
    f"./program --name=psicologia --semester=2025-10 --classrooms=5 --labs=2 --faculty-server={facAddress}:5001",
    f"./program --name=sociologia --semester=2025-10 --classrooms=5 --labs=2 --faculty-server={facAddress}:5001",

    f"./program --name=biologia --semester=2025-10 --classrooms=5 --labs=2 --faculty-server={facAddress}:5002",
]




processes = [subprocess.Popen(cmd, shell=True) for cmd in program_commands]

for p in processes:
    p.wait()
