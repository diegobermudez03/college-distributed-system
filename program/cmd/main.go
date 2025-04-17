package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-zeromq/zmq4"
)

const (
	//flags for arguments
	nameArg = "--name"
	semesterArg = "--semester"
	classroomsArg = "--classrooms"
	labsArg = "--labs"
	facultyServerArg = "--faculty-server"
)

//models for request response processing
type requestStructure struct{
	ProgramName string	`json:"program-name" validate:"required"`
	Semester    string	`json:"semester" validate:"required"`
	Classrooms  *int		`json:"classrooms" validate:"required"`	
	Labs        *int		`json:"labs" validate:"required"`
}

type responseStructure struct{
	Status string `json:"status"`
	ClassroomsAsigned int `json:"classrooms-assigned"`
	LabsAsigned int `json:"labs-assigned"`
	ErrorRequest bool 		`json:"error-request"`
}

func main() {
	//check number of arguments, should be (using labels so that the order dont matter):
	// --name=<program-name> --semester=<semester> --classrooms=<num-classrooms> --labs=<num-labs> --faculty-server=<faculty-address>
	//	executable --name=Ingenieria-Sistemas --semester=2025-10 --classrooms=4 --labs=10 --faculty-server=127.0.0.1:5000  
	if len(os.Args) < 6{
		log.Fatal("invalid number of arguments")
	} 

	request := requestStructure{}
	var facultyServer string
	for _, arg := range os.Args{
		//ignore if atribute doesnt start with --
		if !strings.Contains(arg, "--"){
			continue
		}
		parts := strings.Split(arg, "=")
		if len(parts) != 2{
			continue
		}
		switch parts[0]{
			case nameArg:
				request.ProgramName = parts[1]
			case semesterArg:
				request.Semester = parts[1]
			case classroomsArg:
				if classrooms, err :=strconv.Atoi(parts[1]); err == nil{
					request.Classrooms = &classrooms
				}
			case labsArg:
				if labs, err := strconv.Atoi(parts[1]); err == nil{
					request.Labs = &labs
				}
			case facultyServerArg:
				facultyServer = parts[1]
		}
	}
	
	//check if all the request attributes where obtained from the flags
	if request.ProgramName == "" || request.Semester == "" || request.Classrooms == nil || request.Labs == nil{
		log.Fatal("invalid request")
	}

	//establish connection with the faculty
	socket := zmq4.NewReq(context.Background(), zmq4.WithDialerRetry(time.Second))
	defer socket.Close()
	if err := socket.Dial(fmt.Sprintf("tcp://%s", facultyServer)); err != nil {
		log.Fatal("Couldnt connect with faculty:" , err.Error())
	}

	jsonMessage, err := json.Marshal(request)
	if err != nil{
		log.Fatal("Internal error processing request")
	}
	message := zmq4.NewMsgString(string(jsonMessage))
	socket.Send(message)
	//request sent to the faculty with JSON structure

	//wait for response
	response, err := socket.Recv()
	if err != nil{
		log.Fatal("Error with request ", err.Error())
	}
	//unmarshal response
	responseJson := responseStructure{}
	if err := json.Unmarshal(response.Bytes(), &responseJson); err != nil{
		log.Fatal("Invalid response received ", err.Error())
	}

	//if we received an error then we simply print the error and stop
	if responseJson.ErrorRequest{
		log.Fatal("Error received from faculty ", responseJson.Status)
		return
	}

	//print result 
	log.Printf("Received %s, %d classrooms expected: %d received.  %d labs expected: %d received.",
		 responseJson.Status, *request.Classrooms, responseJson.ClassroomsAsigned, 
		 *request.Labs, responseJson.LabsAsigned,
	)

	//store result in a file
	file, err := os.Create(fmt.Sprintf("%s_%s.txt", request.ProgramName, request.Semester))
	if err != nil{
		log.Print("Unable to write file", err.Error())
	}
	defer file.Close()
	file.Write(response.Bytes())

}
