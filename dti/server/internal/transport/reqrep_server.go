package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/go-zeromq/zmq4"
)

type ReqRepServer struct {
	port        int
	service     domain.CollegeService
	socket      zmq4.Socket
	counter     int
	endChannel  chan bool
	proxyServer string
	faculties   int
	lock        sync.Mutex
}

func NewReqRepServer(service domain.CollegeService, config domain.ServerConfig) *ReqRepServer {
	return &ReqRepServer{
		port:        config.ListenPort,
		service:     service,
		faculties:   config.NumFaculties,
		endChannel:  config.EndChannel,
		proxyServer: config.ProxyServer,
		lock:        sync.Mutex{},
	}
}

// principal method to start server and listen for requests
func (s *ReqRepServer) Listen() error {
	//first we poblate the DB
	if err := s.service.PoblateFacultiesAndPrograms(); err != nil {
		return err
	}

	//if we are using a proxy server, then we suscribe
	if s.proxyServer != "" {
		//establish connection with the faculty
		conn, err := net.Dial("tcp", s.proxyServer)
		if err != nil {
			log.Print(err.Error())
			return errors.New("unable to suscribe with proxy")
		}
		message := strconv.Itoa(s.port) + "\n"
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Print(err.Error())
			return errors.New("unable to suscribe with proxy")
		}
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Print(err.Error())
			return errors.New("unable to suscribe with proxy")
		}
		if response != "OK\n" {
			return errors.New("proxy didnt accept the suscription")
		}
		conn.Close()
	}

	//create zeromq socket and listen in the given port
	socket := zmq4.NewRouter(context.Background())
	socket.SetOption(zmq4.OptionHWM, 10000000)
	s.socket = socket
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.port)); err != nil {
		return err
	}
	log.Print("Listening on port: ", s.port)
	go func() {
		defer socket.Close()
		for {
			message, err := socket.Recv()
			//process in a seaparate new go routine the message to continue listening for new messages
			go s.processMessage(message, err)
		}
	}()
	return nil
}

// internal method to process each request message, it validates the message and communicates with the service
func (s *ReqRepServer) processMessage(message zmq4.Msg, err error) {
	//create a goroutine ID, to identify this go routine
	goRoutineId := rand.Intn(90000) + 10000
	//if there was an error with the mesage we ignore it then
	if err != nil {
		return
	}
	clientIdentity := message.Frames[0]

	//for proxy purposes
	var clientId []byte = []byte{}
	if len(message.Frames) > 2 {
		clientId = message.Frames[2]
	}
	//read request body
	//if the message is of acceptance, then we ignore
	if string(message.Frames[1]) == "ACCEPT" {
		return
	}
	clientRequestBytes := message.Frames[1]
	clientRequest := domain.DTIRequestDTO{}

	////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////
	//if message wasnt a request, we check if it was a HEALTH CHECK
	if err := json.Unmarshal(clientRequestBytes, &clientRequest); err != nil || clientRequest.Semester == "" {
		hCheck := HealthCheckDTO{}
		if err := json.Unmarshal(clientRequestBytes, &hCheck); err != nil {
			return
		}
		//if it was a health check, we answer with a simple 1 byte
		log.Print("ANSWERING HEALTH CHECK")
		s.socket.Send(zmq4.NewMsgFrom(clientIdentity, []byte{1}))
		return
	}
	////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////

	//process message with the service
	response, err := s.service.ProcessRequest(clientRequest, goRoutineId)
	log.Print("IT EXITED THE FUNCTION ProcessRequest")
	var responseBytes []byte
	if err != nil {
		//if there was an error, we send it in the authorized format
		errorResponse := domain.DTIResponseDTO{
			Semester:     clientRequest.Semester,
			ErrorFound:   true,
			ErrorMessage: err.Error(),
		}
		responseBytes, _ = json.Marshal(errorResponse)
	} else {
		//check if we completed all the faculties so we send the end signal
		s.lock.Lock()
		s.counter++
		if s.counter == s.faculties {
			s.endChannel <- true
		}
		s.lock.Unlock()
		//send response to the spceified client
		responseBytes, _ = json.Marshal(response)
	}
	//send message with client ID (if recived one, means, we are using proxy)
	if err := s.socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes, clientId)); err != nil {
		log.Printf("ERROR SENDING aANSWERRRRRRRRRRR %v", err.Error())
	} else {
		log.Print("ANSWEER SENTTTTTTTTTTTTTTTTTT")
	}
}
