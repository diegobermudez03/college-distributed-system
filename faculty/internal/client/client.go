package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/go-zeromq/zmq4"
)

var (
	errUnableToConnectToDTI = errors.New("UNABLE-TO-CONNECT-TO-DTI")
)

const (
	acceptMessage = "ACCEPT"
)

type FacultyClient struct {
	dtiAddress  string
	facultyName string
	ctx         context.Context
	cancel      context.CancelFunc
	socket      zmq4.Socket
	semester    string
}

func NewFacultyClient(dtiAddress, semester, facultyName string) *FacultyClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &FacultyClient{
		dtiAddress:  dtiAddress,
		facultyName: facultyName,
		ctx:         ctx,
		cancel:      cancel,
		semester:    semester,
	}
}

// function that starts the ZMQ4 client with the DTI, and then in a separate go routine (thread) receives requests
// from the channel and sends them to the DTI
func (c *FacultyClient) SendRequests(channel chan models.SemesterRequest) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error = nil
	go func() {
		//connect to the DTI server
		socket := zmq4.NewDealer(c.ctx, zmq4.WithAutomaticReconnect(true))
		if errr := socket.Dial(fmt.Sprintf("tcp://%s", c.dtiAddress)); errr != nil {
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

// method that listens on the ZMQ4 socket with the DTI, for responses, for each response
// received it writes it in the channel for the server to handle the response to the programs
func (c *FacultyClient) ListenResponses(wg *sync.WaitGroup, listenerChannel chan models.SemesterRequest) chan models.DTIResponse {
	channel := make(chan models.DTIResponse)
	wg.Add(1)
	go func() {
		for {
			//receive dti response
			allocation, err := c.socket.Recv()
			log.Print("RECEIVED MESSAGE FROM SERVERRRRRRRRRRRRRRRRRR")
			if err != nil {
				log.Printf("MESSAGE RECEIVED WITH ERRORRRRRRR %v", err.Error())
				continue
			}
			dtiResponse := models.DTIResponse{}
			if errr := json.Unmarshal(allocation.Bytes(), &dtiResponse); errr != nil {
				log.Printf("MESSAGE RECEIVED WITH ERRORRRRRRR %v", err.Error())
				continue
			}
			log.Printf("Received DTI response for semester %s", dtiResponse.Semester)
			//send accept message with the specified semester
			//acceptMsg := fmt.Sprintf("%s-%s", acceptMessage, dtiResponse.Semester)
			//c.socket.Send(zmq4.NewMsgString(acceptMsg))
			//send response in channel for server to handle the response to the programs
			channel <- dtiResponse
			//if we had a specified semester, then after this semester we simply close the channel
			if c.semester != "" {
				close(channel)
				close(listenerChannel)
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

// internal  method to handle the listening on the channel and send to the DTI
func (c *FacultyClient) sendFacultyRequest(channel chan models.SemesterRequest) error {
	//listen to the channel requests to send them to the DTI
	for request := range channel {
		//generate payload to send
		dtiRequest := models.DTIRequest{
			Semester:    request.Semester,
			FacultyName: c.facultyName,
			Programs:    make([]models.DTIProgramRequest, 0, len(request.Programs)),
		}
		//add program requests
		for _, program := range request.Programs {
			dtiRequest.Programs = append(dtiRequest.Programs, models.DTIProgramRequest{
				ProgramId:   program.ClientId,
				ProgramName: program.ProgramName,
				Classrooms:  program.Classrooms,
				Labs:        program.Labs,
			})
		}
		requestBytes, _ := json.Marshal(dtiRequest)

		//send request to the DTI
		log.Printf("Sending semester %s request to DTI", request.Semester)
		if err := c.socket.Send(zmq4.NewMsgFrom(requestBytes)); err != nil {
			return err
		}
	}
	return nil
}
