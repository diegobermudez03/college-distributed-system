import os
import subprocess

dtiServerAddress = os.getenv('DTI_ADDRESS')
dtiServerAddress = dtiServerAddress + ":6000"

commands = [
    ##start engineering faculty
    f"./fac --name=ingenieria --semester=2025-10 --dti-server={dtiServerAddress} --min-programs=1 --listen-port=5001", 
    ##start law faculty
    f"./fac --name=derecho --semester=2025-10 --dti-server={dtiServerAddress} --min-programs=1 --listen-port=5002", 
    ##start arts faculty
    f"./fac --name=artes --semester=2025-10 --dti-server={dtiServerAddress} --min-programs=1 --listen-port=5003", 
]

#we execute all the commands
processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

#wait to all processes to end
for p in processes:
    p.wait()
