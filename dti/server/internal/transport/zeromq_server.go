package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/service"
	"github.com/go-zeromq/zmq4"
)

type ZeroMqServer struct {
	port    int
	service service.CollegeService
	socket  zmq4.Socket
}

func NewZeroMQServer(port int, service service.CollegeService) *ZeroMqServer {
	return &ZeroMqServer{
		port: port,
		service: service,
	}
}

//principal method to start server and listen for requests
func (s *ZeroMqServer) Listen() error {
	//create zeromq socket and listen in the given port
	socket := zmq4.NewRouter(context.Background())
	s.socket = socket
	defer socket.Close()
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.port)); err != nil {
		return err
	}
	for{
		message, err := socket.Recv()
		//process in a seaparate new go routine the message to continue listening for new messages
		go s.processMessage(message, err)
	}
}


//internal method to process each request message, it validates the message and communicates with the service
func (s *ZeroMqServer) processMessage(message zmq4.Msg, err error){
	//if there was an error with the mesage we ignore it then
	if err != nil{
		return 
	}
	clientIdentity := message.Frames[0]
	//read request body
	clientRequestBytes := message.Frames[2]
	clientRequest := service.DTIRequestDTO{}
	if err := json.Unmarshal(clientRequestBytes, &clientRequest); err != nil{
		return 
	}

	//process message with the service
	response, err := s.service.ProcessRequest(clientRequest)
	if err != nil{
		return
	}
	//send response to the spceified client
	responseBytes, _ := json.Marshal(response)
	s.socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes))
}
