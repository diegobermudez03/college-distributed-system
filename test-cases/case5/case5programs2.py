import subprocess
import os

##################################################################################
#### 30 SALONES 10 LABS
#### SEMESTRE 2025-10


facAddress = os.getenv('DTI_ADDRESS', '127.0.0.1')

program_commands = [
    f"./program --name=ingenieria-civil --semester=2025-10 --classrooms=5 --labs=2 --faculty-server={facAddress}:5003",
    f"./program --name=ingenieria-electronica --semester=2025-10 --classrooms=5 --labs=2 --faculty-server={facAddress}:5003",
    f"./program --name=ingenieria-de-sistemas --semester=2025-10 --classrooms=5 --labs=2 --faculty-server={facAddress}:5003",
]



processes = [subprocess.Popen(cmd, shell=True) for cmd in program_commands]

for p in processes:
    p.wait()
