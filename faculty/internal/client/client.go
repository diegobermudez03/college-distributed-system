package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)

var(
	errUnableToConnectToDTI = errors.New("UNABLE-TO-CONNECT-TO-DTI")
)

const(
	acceptMessage = "ACCEPT"
)


type FacultyClient struct {
	dtiAddress  string
	facultyName string
	semester    string
}

func NewFacultyClient(dtiAddress, semester, facultyName string) *FacultyClient {
	return &FacultyClient{
		dtiAddress:  dtiAddress,
		semester:    semester,
		facultyName: facultyName,
	}
}

func (c *FacultyClient) SendFacultyRequest(programs map[uuid.UUID]models.ProgramRequest) ([]models.ResponseRequest,error) {
	//connect to the DTI server
	socket := zmq4.NewDealer(context.Background())
	if err := socket.Dial(fmt.Sprintf("tcp://%s", c.dtiAddress)); err != nil{
		return nil, errUnableToConnectToDTI
	}
	defer socket.Close()
	
	//generate payload to send
	dtiRequest := models.DTIRequest{
		Semester: c.semester,
		FacultyName: c.facultyName,
		Programs: make([]models.DTIProgramRequest, 0, len(programs)),
	}
	//add program requests
	for _, program := range programs{
		dtiRequest.Programs = append(dtiRequest.Programs, models.DTIProgramRequest{
			ProgramId: program.ClientId,
			ProgramName: program.ProgramName,
			Classrooms: program.Classrooms,
			Labs: program.Labs,
		})
	}
	requestBytes, _ := json.Marshal(dtiRequest)

	//send request to the DTI
	if err := socket.Send(zmq4.NewMsgFrom([]byte{}, requestBytes)); err != nil{
		return nil, err
	}
	//receive dti response
	allocation, err := socket.Recv()
	if err != nil{
		return nil, err 
	}
	socket.Send(zmq4.NewMsgString(acceptMessage))
	


	responses := make([]models.ResponseRequest, 0, len(programs))
	for _, program := range programs {
		response := models.ResponseRequest{
			ClientId:         program.ClientId,
			Status:           "OK",
			ClassroomsAsigned: program.Classrooms,
			LabsAsigned:     program.Labs,
			ErrorRequest: false,
		}
		responses = append(responses, response)
	}
	return responses,nil
}