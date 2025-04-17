package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)

type ZeroMqServer struct {
	port    int
	service domain.CollegeService
	socket  zmq4.Socket
}

func NewZeroMQServer(port int, service domain.CollegeService) *ZeroMqServer {
	return &ZeroMqServer{
		port: port,
		service: service,
	}
}

//principal method to start server and listen for requests
func (s *ZeroMqServer) Listen() error {
	//first we poblate the DB
	if err:= s.service.PoblateFacultiesAndPrograms(); err != nil{
		return err
	}

	//create zeromq socket and listen in the given port
	socket := zmq4.NewRouter(context.Background())
	s.socket = socket
	defer socket.Close()
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.port)); err != nil {
		return err
	}
	log.Print("Listening on port: ", s.port)
	for{
		message, err := socket.Recv()
		//process in a seaparate new go routine the message to continue listening for new messages
		go s.processMessage(message, err)
	}
}


//internal method to process each request message, it validates the message and communicates with the service
func (s *ZeroMqServer) processMessage(message zmq4.Msg, err error){
	//create a goroutine ID, to identify this go routine
	goRoutineId := uuid.New()
	//if there was an error with the mesage we ignore it then
	if err != nil{
		return 
	}
	clientIdentity := message.Frames[0]
	//read request body
	//if the message is of acceptance, then we ignore
	if string(message.Frames[1]) == "ACCEPT"{
		return 
	}
	clientRequestBytes := message.Frames[1]
	clientRequest := domain.DTIRequestDTO{}
	if err := json.Unmarshal(clientRequestBytes, &clientRequest); err != nil{
		return 
	}

	//process message with the service
	response, err := s.service.ProcessRequest(clientRequest, goRoutineId)
	if err != nil{
		//if there was an error, we send it in the authorized format
		errorResponse := domain.DTIResponseDTO{
			Semester: clientRequest.Semester,
			ErrorFound: true,
			ErrorMessage: err.Error(),
		}
		responseBytes, _ := json.Marshal(errorResponse)
		s.socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes))
		return
	}
	//send response to the spceified client
	responseBytes, _ := json.Marshal(response)
	s.socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes))
}
