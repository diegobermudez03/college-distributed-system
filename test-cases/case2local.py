import random
import subprocess
import os

##################################################################################
#5 Facultades 
#generando peticiones. 
#5 Programas académicos por Facultad. 
#Cada programa pide los mínimos 7 Aulas y 2 laboratorios (o cómo máximo 2 y 7) 


#with this script the faculties and programs are in the same machine, dti server can be in other


dtiServerAddress = os.getenv('DTI_ADDRESS', '127.0.0.1:6666')
#for this example we need the following flags on the DTI server
# --classrooms=380 --labs=60

##in the faculties we dont specify the semester, so they are ready to handle all of them
faculties_commands = [
    f"./fac --name=ciencias-sociales --dti-server={dtiServerAddress} --listen-port=5001 --semester=2025-10", 
    f"./fac --name=ciencias-naturales --dti-server={dtiServerAddress} --listen-port=5002 --semester=2025-10", 
    f"./fac --name=ingenieria --dti-server={dtiServerAddress} --listen-port=5003 --semester=2025-10", 
    f"./fac --name=medicina --dti-server={dtiServerAddress} --listen-port=5004 --semester=2025-10", 
    f"./fac --name=derecho --dti-server={dtiServerAddress} --listen-port=5005 --semester=2025-10", 
]

##we will handle the program commands in a separate list, is so that we can re order it random, to check the system with random order
program_commands = [
    ##start ciencias sociales programs
    "./program --name=psicologia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5001",
    "./program --name=sociologia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5001",
    "./program --name=trabajo-social --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5001",
    "./program --name=antropologia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5001",
    "./program --name=comunicacion --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5001",

    
    ##start  ciencias naturales programs
    "./program --name=biologia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5002",
    "./program --name=quimica --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5002",
    "./program --name=fisica --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5002",
    "./program --name=geologia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5002",
    "./program --name=ciecias-ambientales --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5002",

    
    ##start  engineering programs
    "./program --name=ingenieria-civil --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-electronica --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-de-sistemas --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-mecanica --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5003",
    "./program --name=ingenieria-industrial --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5003",
   
    
    ##start medicine programs
    "./program --name=medicina-general --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5004",
    "./program --name=enfermeria --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5004",
    "./program --name=odontologia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5004",
    "./program --name=farmacia --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5004",
    "./program --name=terapia-fisica --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5004",

    ##start law programs
    "./program --name=derecho-penal --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5005",
    "./program --name=derecho-civil --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5005",
    "./program --name=derecho-internacional --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5005",
    "./program --name=derecho-laboral --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5005",
    "./program --name=derecho-constitucional --semester=2025-10 --classrooms=10 --labs=4 --faculty-server=127.0.0.1:5005",
   
]


#randomize the order of the program commands, so that its completely random
#random.shuffle(program_commands)

#join the 2 lists of commands in a single command list
commands = faculties_commands + program_commands

#we execute all the commands
processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

#wait to all processes to end
for p in processes:
    p.wait()
