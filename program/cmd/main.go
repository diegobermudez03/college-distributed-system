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

	//already assigned message from dti, so that if is this we dont store file
	AlreadyHasAssignation = "PROGRAM-ALREADY-HAS-RESOURCES-FOR-SEMESTER"
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
	MobileLabsAssigned int `json:"mobile-labs-assigned"`
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


	//write result in file 
	var txtMessage string 

	//if we received an error then we print it and store it in the txt string
	if responseJson.Status != "OK"{
		log.Print("Error received from faculty ", responseJson.Status)
		//if the error is that we already have assignation then we simply return and avoid touching the txt
		if AlreadyHasAssignation == responseJson.Status{
			return
		}
	}
	//print result 
	log.Printf("Received %s, %d classrooms expected: %d received.  %d labs expected: %d received.",
	responseJson.Status, *request.Classrooms, responseJson.ClassroomsAsigned, 
	*request.Labs, responseJson.LabsAsigned)
	txtMessage = fmt.Sprintf(`
	Program %s in SEMESTER: % s
	REQUESTED CLASSROOMS: %d
	REQUESTED LABS: %d
	ASSIGNED CLASSROOMS: %d
	ASSIGNED LABS: %d
	HOW MANY OF THE LABS ARE MOBILE LABS (CLASSROOMS WITH PC'S): %d
	STATUS: %s
	`, request.ProgramName, request.Semester, *request.Classrooms, *request.Labs,
	responseJson.ClassroomsAsigned, responseJson.LabsAsigned, responseJson.MobileLabsAssigned, responseJson.Status)

	//create file and store message
	file, err := os.Create(fmt.Sprintf("%s_%s.txt", request.ProgramName, request.Semester))
	if err != nil{
		log.Print("Unable to write file", err.Error())
	}
	defer file.Close()
	file.Write([]byte(txtMessage))

}
