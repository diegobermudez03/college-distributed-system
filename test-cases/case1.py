import subprocess
import os

##################################################################################
#### 380 SALONES 60 LABS
#### SEMESTRE 2023-10


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

faculties_commands = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2023-10 --min-programs=2", 
]

program_commands = [
    "./program --name=psicologia --semester=2023-10 --classrooms=10 --faculty-server=127.0.0.1:5001",
    "./program --name=sociologia --semester=2023-10 --labs=4 --faculty-server=127.0.0.1:5001",
]


commands = faculties_commands + program_commands

processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

for p in processes:
    p.wait()
