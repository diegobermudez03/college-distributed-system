import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS
#### SEMESTRE 2024-10

dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

faculties_commands = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2024-10 --min-programs=2", 
    f"./fac --name=ciencias-naturales --dti-server={dtiServerAddress} --listen-port=5002 --semester=2024-10 --min-programs=1", 
    f"./fac --name=ingenieria --dti-server={dtiServerAddress} --listen-port=5003 --semester=2024-10 --min-programs=3",
]


processes = [subprocess.Popen(cmd, shell=True) for cmd in faculties_commands]

for p in processes:
    p.wait()
