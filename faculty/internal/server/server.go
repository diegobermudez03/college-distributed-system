package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/client"
	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)

const (
	internalError = "INTERNAL-ERROR"
	invalidSemester = "INVALID-SEMESTER"
)


//server model
type FacultyServer struct {
	listenPort int
	minPrograms int
	semester string
	clients map[uuid.UUID]models.ProgramRequest
	client *client.FacultyClient
}

func NewFacultyServer(listenPort, minPrograms int, semester string, client *client.FacultyClient) *FacultyServer {
	return &FacultyServer{
		listenPort: listenPort,
		minPrograms: minPrograms,
		clients: map[uuid.UUID]models.ProgramRequest{},
		semester: semester,
		client: client,
	}
}

// method to run the zeromq request reply server
func (s *FacultyServer) Run() error {
	//start zeromq request reply server
	//socket := zmq4.NewRep(context.Background())
	socket := zmq4.NewRouter(context.Background())
	defer socket.Close()
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.listenPort)); err != nil {
		return errors.New(fmt.Sprint("Unable to start server at port ", s.listenPort))
	}
	log.Printf("Listening at port %d", s.listenPort)
	//listen for the program requests
	if err := s.serveServer(socket); err != nil{
		return err
	}
	return nil
}

func (s *FacultyServer) serveServer(socket zmq4.Socket) error{
	//receive minimum program requests
	if err := s.listenProgramRequests(socket); err != nil{
		return err
	}
	//send request to the DTI and receive the responses
	responses, err := s.client.SendFacultyRequest(s.clients)
	if err != nil{
		return err 
	}
	//iterate over all responses, get the socket ID for each one, and then send the JSON response
	for _, clientResponse := range responses{
		client, ok := s.clients[clientResponse.ClientId]
		if !ok{
			continue 
		}
		bytes, _ := json.Marshal(clientResponse)
		socket.Send(zmq4.NewMsgFrom(client.ClientSocketId, bytes))
	}
	return nil
}

func (s *FacultyServer) listenProgramRequests(socket zmq4.Socket) error{
	//wait for the minimum of programs to communicate with the DTI
	for i := 0; i < s.minPrograms; i++ {
		//extract the id of the client
		message, err := socket.Recv()
		clientId := message.Frames[0]
		if err != nil{
			//if theres an error we are going to suppose that it was due to the program client
			i--
			continue
		}
		//create program request and unmarshal it from program message
		programRequest := models.ProgramRequest{
			ClientSocketId: clientId,
			ClientId: uuid.New(),
		}
		if err := json.Unmarshal(message.Frames[2], &programRequest); err != nil{
			//if there's an error reading, then we are going to suppose that
			//the program did something wrong, but we wont break, so we simply ignore
			i--
			continue
		}
		//check if the semester is the correct one, if not, then we answer the program and end
		if programRequest.Semester != s.semester{
			errorResponse := models.ResponseRequest{
				ClientId: programRequest.ClientId,
				Status: invalidSemester,
				ErrorRequest: true,
			}
			errorBytes, _ := json.Marshal(errorResponse)
			socket.Send(zmq4.NewMsgFrom(programRequest.ClientSocketId, errorBytes))
			i--
			continue
		}
		//add client to the message
		s.clients[programRequest.ClientId] = programRequest
	}
	return nil
}