package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/client"
	"github.com/diegobermudez03/college-distributed-system/faculty/internal/server"
)

const (
	//args
	nameArg = "--name"
	semesterArg = "--semester"
	dtiServerArg = "--dti-server"
	minProgramsArg = "--min-programs"
	listenPortArg = "--listen-port"
)
type configuration struct {
	name        string
	semester    string
	dtiServer   string
	minPrograms int
	listenPort  int
}

func main() {
	//input structure (atributes): min programs is optional argument for testing, reduces the number of programs that we wait to process
	// --name=<faculty-name> --semester=<semester> --dti-server=<dti-server-address> --min-programs=<number-of-min-programs> --listen-port=5000
	// executable --name=tecnologia --semester=2025-10 --dti-server=127.0.0.1:5000 --min-programs=2 --listen-port=6000

	//check number of arguments, at least 4 (min programs is optional)
	if len(os.Args) < 4{
		log.Fatal("not enough arguments")
	}
	config := configuration{
		minPrograms: 5,	//by default minPrograms is 5, if the optional argument was added it would change latwer
		listenPort: 5000,
	}
	//read all arguments
	for _, arg := range os.Args{
		//ignore invalid argument
		if !strings.Contains(arg, "="){
			continue
		}
		parts := strings.Split(arg, "=")
		if len(parts) != 2{
			continue
		}
		switch parts[0]{
			case nameArg: config.name = parts[1]
			case semesterArg: config.semester = parts[1]
			case dtiServerArg: config.dtiServer = parts[1]
			case minProgramsArg: {
				if minPrograms, err := strconv.Atoi(parts[1]); err == nil{
					config.minPrograms = minPrograms
				}
			}
			case listenPortArg:{
				if listenPort, err := strconv.Atoi(parts[1]); err == nil{
					config.listenPort = listenPort
				}
			}
		}
	}

	//check if we have all the arguments
	if config.name == "" || config.semester == "" || config.dtiServer == ""{
		log.Fatal("invalid arguments")
	}

	//create faculty client
	client := client.NewFacultyClient(config.dtiServer, config.semester, config.name)
	//start server
	server := server.NewFacultyServer(config.listenPort, config.minPrograms, config.semester, client)
	if err := server.Run(); err != nil{
		log.Fatal("Error ", err.Error())
	}
}

