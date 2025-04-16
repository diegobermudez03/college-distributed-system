package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/go-zeromq/zmq4"
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
	socket 		zmq4.Socket
	semester 	string
}

func NewFacultyClient(dtiAddress, semester, facultyName string) *FacultyClient {
	return &FacultyClient{
		dtiAddress:  dtiAddress,
		facultyName: facultyName,
		semester: semester,
	}
}

//function that starts the ZMQ4 client with the DTI, and then in a separate go routine (thread) receives requests
//from the channel and sends them to the DTI
func (c *FacultyClient) SendRequests(channel chan  models.SemesterRequest) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error = nil
	go func(){
		//connect to the DTI server
		socket := zmq4.NewDealer(context.Background())
		if errr := socket.Dial(fmt.Sprintf("tcp://%s", c.dtiAddress)); errr != nil{
			err = errUnableToConnectToDTI
		}
		defer socket.Close()
		c.socket = socket
		wg.Done()
		//we go listen for the requests
		c.sendFacultyRequest(channel)
	}()
	//this is to wait til the go routine may return an error or not
	wg.Wait()
	return err
}

//method that listens on the ZMQ4 socket with the DTI, for responses, for each response
//received it writes it in the channel for the server to handle the response to the programs
func (c *FacultyClient) ListenResponses(wg *sync.WaitGroup) (chan models.DTIResponse){
	channel := make(chan models.DTIResponse)
	wg.Add(1)
	go func(){
		for{
			//receive dti response
			allocation, err := c.socket.Recv()
			if err != nil{
				continue
			}
			dtiResponse := models.DTIResponse{}
			if errr := json.Unmarshal(allocation.Bytes(), &dtiResponse); errr != nil{
				continue
			}
			//send accept message with the specified semester
			acceptMsg := fmt.Sprintf("%s-%s", acceptMessage, dtiResponse.Semester)
			c.socket.Send(zmq4.NewMsgString(acceptMsg))
			//send response in channel for server to handle the response to the programs
			channel <- dtiResponse
			//if we had a specified semester, then after this semester we simply close the channel
			if c.semester != ""{
				close(channel)
				wg.Done()
				break
			}
		}
	}()
	return channel
}


///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////
//						INTERNAL METHODS
///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////


//internal  method to handle the listening on the channel and send to the DTI
func (c *FacultyClient) sendFacultyRequest(channel chan models.SemesterRequest) error{
	//listen to the channel requests to send them to the DTI
	for request := range channel{
		//generate payload to send
		dtiRequest := models.DTIRequest{
			Semester: request.Semester,
			FacultyName: c.facultyName,
			Programs: make([]models.DTIProgramRequest, 0, len(request.Programs)),
		}
		//add program requests
		for _, program := range request.Programs{
			dtiRequest.Programs = append(dtiRequest.Programs, models.DTIProgramRequest{
				ProgramId: program.ClientId,
				ProgramName: program.ProgramName,
				Classrooms: program.Classrooms,
				Labs: program.Labs,
			})
		}
		requestBytes, _ := json.Marshal(dtiRequest)

		//send request to the DTI
		if err := c.socket.Send(zmq4.NewMsgFrom([]byte{}, requestBytes)); err != nil{
			return err
		}
	}
	return nil
}