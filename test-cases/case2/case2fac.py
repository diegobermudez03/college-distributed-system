import subprocess
import os

##################################################################################
#### 380 SALONES 60 LABS
#### SEMESTRE 2023-30


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

faculties_commands = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2023-30 --min-programs=2", 
]

processes = [subprocess.Popen(cmd, shell=True) for cmd in faculties_commands]

for p in processes:
    p.wait()
