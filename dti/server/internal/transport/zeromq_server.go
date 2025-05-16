package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/go-zeromq/zmq4"
)

type ZeroMqServer struct {
	port    int
	service domain.CollegeService
	socket  zmq4.Socket
	counter int
	endChannel chan bool
	proxyServer string
	faculties int
	lock sync.Mutex
}

func NewZeroMQServer(service domain.CollegeService, config domain.ServerConfig) *ZeroMqServer {
	return &ZeroMqServer{
		port: config.ListenPort,
		service: service,
		faculties: config.NumFaculties,
		endChannel: config.EndChannel,
		proxyServer: config.ProxyServer,
		lock: sync.Mutex{},
	}
}

//principal method to start server and listen for requests
func (s *ZeroMqServer) Listen() error {
	//first we poblate the DB
	if err:= s.service.PoblateFacultiesAndPrograms(); err != nil{
		return err
	}

	//if we are using a proxy server, then we suscribe
	if s.proxyServer != ""{
		//establish connection with the faculty
		socket := zmq4.NewReq(context.Background(), zmq4.WithDialerRetry(time.Second))
		defer socket.Close()
		if err := socket.Dial(fmt.Sprintf("tcp://%s", s.proxyServer)); err != nil {
			return errors.New("unable to connect with Proxy")
		}
		suscribeRequest := SuscribeDTO{
			Suscribe: true,
		}
		bytes, _ := json.Marshal(suscribeRequest)
		message := zmq4.NewMsgString(string(bytes))
		socket.Send(message)
		//request sent to the faculty with JSON structure
		//wait for response
		response, err := socket.Recv()
		if err != nil{
			return errors.New("error suscribing with Proxy")
		}

		responseJson := SuscribeResponseDTO{}
		if err := json.Unmarshal(response.Bytes(), &responseJson); err != nil || !responseJson.Suscribed{
			return errors.New("unable to suscribe with the proxy")
		}
	}

	//create zeromq socket and listen in the given port
	socket := zmq4.NewRouter(context.Background())
	s.socket = socket
	defer socket.Close()
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.port)); err != nil {
		return err
	}
	log.Print("Listening on port: ", s.port)
	go func(){
		for{
			message, err := socket.Recv()
			//process in a seaparate new go routine the message to continue listening for new messages
			go s.processMessage(message, err)
		}
	}()
	return nil
}


//internal method to process each request message, it validates the message and communicates with the service
func (s *ZeroMqServer) processMessage(message zmq4.Msg, err error){
	//create a goroutine ID, to identify this go routine
	goRoutineId := rand.Intn(90000) + 10000
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

	////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////
	//if message wasnt a request, we check if it was a HEALTH CHECK
	if err := json.Unmarshal(clientRequestBytes, &clientRequest); err != nil{
		hCheck := HealthCheckDTO{}
		if err := json.Unmarshal(clientRequestBytes, &hCheck); err != nil{
			return 
		}
		//if it was a health check, we answer with a simple 1 byte
		s.socket.Send(zmq4.NewMsgFrom(clientIdentity, []byte{1}))
		return 
	}
	////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////

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
	//check if we completed all the faculties so we send the end signal
	s.lock.Lock()
	s.counter++
	if s.counter == s.faculties{
		s.endChannel <-true
	}
	s.lock.Unlock()
	//send response to the spceified client
	responseBytes, _ := json.Marshal(response)
	s.socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes))
}
