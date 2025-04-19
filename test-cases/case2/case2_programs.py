import os
import subprocess

facultiesAddress = os.getenv('FAC_ADDRESS')

##we will handle the program commands in a separate list, is so that we can re order it random, to check the system with random order
program_commands = [
    ##start engineering programs
    f"./program --name=ingenieria-civil --semester=2025-10 --classrooms=20 --labs=6 --faculty-server={facultiesAddress}:5001",
    f"./program --name=ingenieria-industrial --semester=2025-10 --classrooms=15 --labs=4 --faculty-server={facultiesAddress}:5001",
    f"./program --name=ingenieria-electronica --semester=2025-10 --classrooms=12 --labs=8 --faculty-server={facultiesAddress}:5001",
    f"./program --name=ingenieria-civil --semester=2025-30 --classrooms=35 --labs=8 --faculty-server={facultiesAddress}:5001",
    f"./program --name=ingenieria-industrial --semester=2025-30 --classrooms=12 --labs=5 --faculty-server={facultiesAddress}:5001",
    f"./program --name=ingenieria-electronica --semester=2025-30 --classrooms=13 --labs=5 --faculty-server={facultiesAddress}:5001",
    
    ##start  medicine programs
    f"./program --name=terapia-fisica --semester=2025-10 --classrooms=8 --labs=0 --faculty-server={facultiesAddress}:5002",
    f"./program --name=odontologia --semester=2025-10 --classrooms=10 --labs=0 --faculty-server={facultiesAddress}:5002",
    f"./program --name=enfermeria --semester=2025-10 --classrooms=7 --labs=0 --faculty-server={facultiesAddress}:5002",
    f"./program --name=terapia-fisica --semester=2025-30 --classrooms=12 --labs=0 --faculty-server={facultiesAddress}:5002",
    f"./program --name=odontologia --semester=2025-30 --classrooms=16 --labs=0 --faculty-server={facultiesAddress}:5002",
    f"./program --name=enfermeria --semester=2025-30 --classrooms=10 --labs=0 --faculty-server={facultiesAddress}:5002",
    
    ##start  law programs
    f"./program --name=derecho-internacional --semester=2025-10 --classrooms=10 --labs=5 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-laboral --semester=2025-10 --classrooms=15 --labs=2 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-constitucional --semester=2025-10 --classrooms=12 --labs=6 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-penal --semester=2025-10 --classrooms=3 --labs=4 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-internacional --semester=2025-30 --classrooms=12 --labs=6 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-laboral --semester=2025-30 --classrooms=17 --labs=0 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-constitucional --semester=2025-30 --classrooms=14 --labs=0 --faculty-server={facultiesAddress}:5003",
    f"./program --name=derecho-penal --semester=2025-30 --classrooms=5 --labs=4 --faculty-server={facultiesAddress}:5003",
    
    ##start social science programs
    f"./program --name=antropologia --semester=2025-10 --classrooms=30 --labs=0 --faculty-server={facultiesAddress}:5004",
    f"./program --name=comunicacion --semester=2025-10 --classrooms=10 --labs=0 --faculty-server={facultiesAddress}:5004",
    f"./program --name=antropologia --semester=2025-30 --classrooms=35 --labs=0 --faculty-server={facultiesAddress}:5004",
    f"./program --name=comunicacion --semester=2025-30 --classrooms=11 --labs=0 --faculty-server={facultiesAddress}:5004",
]


#we execute all the commands
processes = [subprocess.Popen(cmd, shell=True) for cmd in program_commands]

#wait to all processes to end
for p in processes:
    p.wait()
