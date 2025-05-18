import random
import subprocess

##################################################################################
#this is the automatized script for the case:
# - El servidor tendra 120 salones, 20 laboratorios, y 10 posibles laboratorios moviles
# - Se probara la capacidad de las facultades para procesar varios semestres en simultaneo
# Facultad ingenieria (3 programas):
#           2025-10:
#               -Ingenieria civil 20 salones, 6 labs
#               -Ingenieria industrial 15 salones, 4 labs
#               -Ingenieria Electronica, 12 salones, 8 labs
#           2025-30:
#               -Ingenieria civil 35 salones, 8 labs
#               -Ingenieria industrial 12 salones, 5 labs
#               -Ingenieria Electronica, 13 salones, 5 labs
# Facultad Medicina (3 programas):
#           2025-10:
#               -Terapia fisica 8 salones, 0 labs
#               -Odontologia    10 salones, 0 labs
#               -Enfermeria  7 salones, 0 labs
#           2025-30:
#               -Terapia fisica 12 salones, 0 labs
#               -Odontologia    16 salones, 0 labs
#               -Enfermeria  10 salones, 0 labs
# Facultad Derecho (4 programas):
#           2025-10:
#               -Derecho Internacional 10 salones, 5 labs
#               -Derecho Laboral	    15 salones, 2 labs
#               -Derecho Constitucional  12 salones, 6 labs
#               -Derecho Penal	3 salones, 4 labs
#           2025-30:
#               -Derecho Internacional 12 salones, 6 labs
#               -Derecho Laboral	    17 salones, 0 labs
#               -Derecho Constitucional  14 salones, 0 labs
#               -Derecho Penal	5 salones, 4 labs
# Facultad Ciencias Sociales (2 programas):
#           2025-10:
#               -Antropologia       30 salones, 0 labs
#               -Comunicacion       10 salones, 0 labs
#           2025-30:
#               -Antropologia       35 salones, 0 labs
#               -Comunicacion       11 salones, 0 labs


#with this script the faculties and programs are in the same machine, dti server can be in other


dtiServerAddress = "127.0.0.1:6000"
#for this example we need the following flags on the DTI server
# --classrooms=120 --labs=20 --mobile-labs=10

##in the faculties we dont specify the semester, so they are ready to handle all of them
faculties_commands = [
    ##start engineering faculty
    f"./faculty/bin/fac --name=ingenieria --dti-server={dtiServerAddress} --min-programs=3 --listen-port=5001", 
    ##start medicine faculty
    f"./faculty/bin/fac --name=medicina --dti-server={dtiServerAddress} --min-programs=3 --listen-port=5002", 
    ##start law faculty
    f"./faculty/bin/fac --name=derecho --dti-server={dtiServerAddress} --min-programs=4 --listen-port=5003", 
    ##start medicine social science
    f"./faculty/bin/fac --name=ciencias-sociales --dti-server={dtiServerAddress} --min-programs=2 --listen-port=5004", 
]

##we will handle the program commands in a separate list, is so that we can re order it random, to check the system with random order
program_commands = [
    ##start engineering programs
    "./program/bin/program --name=ingenieria-civil --semester=2025-10 --classrooms=20 --labs=6 --faculty-server=127.0.0.1:5001",
    "./program/bin/program --name=ingenieria-industrial --semester=2025-10 --classrooms=15 --labs=4 --faculty-server=127.0.0.1:5001",
    "./program/bin/program --name=ingenieria-electronica --semester=2025-10 --classrooms=12 --labs=8 --faculty-server=127.0.0.1:5001",
    "./program/bin/program --name=ingenieria-civil --semester=2025-30 --classrooms=35 --labs=8 --faculty-server=127.0.0.1:5001",
    "./program/bin/program --name=ingenieria-industrial --semester=2025-30 --classrooms=12 --labs=5 --faculty-server=127.0.0.1:5001",
    "./program/bin/program --name=ingenieria-electronica --semester=2025-30 --classrooms=13 --labs=5 --faculty-server=127.0.0.1:5001",
    
    ##start  medicine programs
    "./program/bin/program --name=terapia-fisica --semester=2025-10 --classrooms=8 --labs=0 --faculty-server=127.0.0.1:5002",
    "./program/bin/program --name=odontologia --semester=2025-10 --classrooms=10 --labs=0 --faculty-server=127.0.0.1:5002",
    "./program/bin/program --name=enfermeria --semester=2025-10 --classrooms=7 --labs=0 --faculty-server=127.0.0.1:5002",
    "./program/bin/program --name=terapia-fisica --semester=2025-30 --classrooms=12 --labs=0 --faculty-server=127.0.0.1:5002",
    "./program/bin/program --name=odontologia --semester=2025-30 --classrooms=16 --labs=0 --faculty-server=127.0.0.1:5002",
    "./program/bin/program --name=enfermeria --semester=2025-30 --classrooms=10 --labs=0 --faculty-server=127.0.0.1:5002",
    
    ##start  law programs
    "./program/bin/program --name=derecho-internacional --semester=2025-10 --classrooms=10 --labs=5 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-laboral --semester=2025-10 --classrooms=15 --labs=2 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-constitucional --semester=2025-10 --classrooms=12 --labs=6 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-penal --semester=2025-10 --classrooms=3 --labs=4 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-internacional --semester=2025-30 --classrooms=12 --labs=6 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-laboral --semester=2025-30 --classrooms=17 --labs=0 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-constitucional --semester=2025-30 --classrooms=14 --labs=0 --faculty-server=127.0.0.1:5003",
    "./program/bin/program --name=derecho-penal --semester=2025-30 --classrooms=5 --labs=4 --faculty-server=127.0.0.1:5003",
    
    ##start social science programs
    "./program/bin/program --name=antropologia --semester=2025-10 --classrooms=30 --labs=0 --faculty-server=127.0.0.1:5004",
    "./program/bin/program --name=comunicacion --semester=2025-10 --classrooms=10 --labs=0 --faculty-server=127.0.0.1:5004",
    "./program/bin/program --name=antropologia --semester=2025-30 --classrooms=35 --labs=0 --faculty-server=127.0.0.1:5004",
    "./program/bin/program --name=comunicacion --semester=2025-30 --classrooms=11 --labs=0 --faculty-server=127.0.0.1:5004",
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
