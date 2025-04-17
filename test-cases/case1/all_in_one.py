import subprocess

##################################################################################
#this is the automatized script for the case:
# - El servidor se levanta con un número pequeño de recursos disponibles para todos: 10 Aulas y 5 laboratorios y 2 laboratorios móviles. 
# - Se ejecutan tres Facultades (tres procesos) en forma concurrente, en otra máquina distinta,  con las siguientes peticiones:
# 
# Facultad 1: 6 Aulas, 3 laboratorios
# Facultad 2: 3 Aulas, 3 Laboratorios
# Facultad 3: 3 Aulas,  2 Laboratorios  

#in this case, we start the 3 faculties,  we set the minimum programs for all faculties to 1, so that with only 1 program they send
#we will use the same semester on the 3 faculties, and for each faculty we will run only 1 program

#with this script the faculties and programs are in the same machine, dti server can be in other


dtiServerAddress = "127.0.0.1:6000"
#for this example we need the following flags on the DTI server
# --classrooms=10 --labs=5 --mobile-labs=2

commands = [
    ##start engineering faculty
    f"./faculty/bin/fac --name=ingenieria --semester=2025-10 --dti-server={dtiServerAddress} --min-programs=1 --listen-port=5001", 
    ##start law faculty
    f"./faculty/bin/fac --name=derecho --semester=2025-10 --dti-server={dtiServerAddress} --min-programs=1 --listen-port=5002", 
    ##start arts faculty
    f"./faculty/bin/fac --name=artes --semester=2025-10 --dti-server={dtiServerAddress} --min-programs=1 --listen-port=5003", 

    ##start prograsms, 1 for each faculty
    "./program/bin/program --name=Ingenieria-de-Sistemas --semester=2025-10 --classrooms=6 --labs=3 --faculty-server=127.0.0.1:5001",
    "./program/bin/program --name=Derecho-Penal --semester=2025-10 --classrooms=3 --labs=3 --faculty-server=127.0.0.1:5002",
    "./program/bin/program --name=Teatro --semester=2025-10 --classrooms=3 --labs=2 --faculty-server=127.0.0.1:5003",
]

#we execute all the commands
processes = [subprocess.Popen(cmd, shell=True) for cmd in commands]

#wait to all processes to end
for p in processes:
    p.wait()
