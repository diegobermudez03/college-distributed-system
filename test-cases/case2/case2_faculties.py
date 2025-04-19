import os
import subprocess

dtiServerAddress = os.getenv('DTI_ADDRESS')
dtiServerAddress = dtiServerAddress + ":6000"

##in the faculties we dont specify the semester, so they are ready to handle all of them
faculties_commands = [
    ##start engineering faculty
    f"./fac --name=ingenieria --dti-server={dtiServerAddress} --min-programs=3 --listen-port=5001", 
    ##start medicine faculty
    f"./fac --name=medicina --dti-server={dtiServerAddress} --min-programs=3 --listen-port=5002", 
    ##start law faculty
    f"./fac --name=derecho --dti-server={dtiServerAddress} --min-programs=4 --listen-port=5003", 
    ##start medicine social science
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --min-programs=2 --listen-port=5004", 
]

#we execute all the commands
processes = [subprocess.Popen(cmd, shell=True) for cmd in faculties_commands]

#wait to all processes to end
for p in processes:
    p.wait()
