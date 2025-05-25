import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS 5 segundos fail
#### SEMESTRE 2025-10


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

faculties_commands = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2025-10 --min-programs=2", 
    f"./fac --name=ciencias-naturales --dti-server={dtiServerAddress} --listen-port=5002 --semester=2025-10 --min-programs=1", 
]


program_commands = [
    "./program --name=psicologia --semester=2025-10 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5001",
    "./program --name=sociologia --semester=2025-10 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5001",

    "./program --name=biologia --semester=2025-10 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5002",
]



commands = faculties_commands + program_commands

processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

for p in processes:
    p.wait()
