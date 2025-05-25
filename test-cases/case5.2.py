import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS
#### SEMESTRE 2025-10


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

faculties_commands = [
    f"./fac --name=ingenieria --dti-server={dtiServerAddress} --listen-port=5003 --semester=2025-10 --min-programs=3", 
]


program_commands = [
    "./program --name=ingenieria-civil --semester=2025-10 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-electronica --semester=2025-10 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-de-sistemas --semester=2025-10 --classrooms=5 --labs=2 --faculty-server=127.0.0.1:5003",
]



commands = faculties_commands + program_commands

processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

for p in processes:
    p.wait()
