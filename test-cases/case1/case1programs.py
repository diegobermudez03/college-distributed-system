import subprocess
import os

##################################################################################
#### 380 SALONES 60 LABS
#### SEMESTRE 2023-10


facAddress = os.getenv('FAC_ADDRESS', '127.0.0.1')

program_commands = [
    f"./program --name=psicologia --semester=2023-10 --classrooms=10 --faculty-server={facAddress}:5001",
    f"./program --name=sociologia --semester=2023-10 --labs=4 --faculty-server={facAddress}:5001",
]


processes = [subprocess.Popen(cmd, shell=True) for cmd in program_commands]

for p in processes:
    p.wait()
