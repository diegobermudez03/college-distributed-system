import subprocess
import os

##################################################################################
#### 380 SALONES 60 LABS
#### SEMESTRE 2023-10


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1') + ":6666"

dtiServerAddress = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2023-10 --min-programs=2", 
]

processes = [subprocess.Popen(cmd, shell=True) for cmd in dtiServerAddress]

for p in processes:
    p.wait()
