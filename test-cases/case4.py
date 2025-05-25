import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS
#### SEMESTRE 2024-30


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

faculties_commands = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2024-30 --min-programs=2" , 
    f"./fac --name=ciencias-naturales --dti-server={dtiServerAddress} --listen-port=5002 --semester=2024-30 --min-programs=2", 
    f"./fac --name=ingenieria --dti-server={dtiServerAddress} --listen-port=5003 --semester=2024-30 --min-programs=3", 
]


program_commands = [
    "./program --name=psicologia --semester=2024-30 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5001",
   
    "./program --name=biologia --semester=2024-30 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5002",

    "./program --name=ingenieria-civil --semester=2024-30 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-electronica --semester=2024-30 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-de-sistemas --semester=2024-30 --labs=2 --faculty-server=127.0.0.1:5003",

    "./program --name=sociologia --semester=2024-30 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5001",
    "./program --name=quimica --semester=2024-30 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5002",
]



commands = faculties_commands + program_commands

processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

for p in processes:
    p.wait()
