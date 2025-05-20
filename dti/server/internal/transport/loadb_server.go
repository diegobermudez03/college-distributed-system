package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/go-zeromq/zmq4"
)

type LoadBServer struct {
	port        int
	service     domain.CollegeService
	socket      zmq4.Socket
	counter     int
	endChannel  chan bool
	proxyServer string
	faculties   int
	lock        sync.Mutex
	nWorkers    int
	channels    []chan zmq4.Msg
}

func NewLoadBServer(service domain.CollegeService, config domain.ServerConfig, nWorkers int) *LoadBServer {
	return &LoadBServer{
		port:        config.ListenPort,
		service:     service,
		faculties:   config.NumFaculties,
		endChannel:  config.EndChannel,
		proxyServer: config.ProxyServer,
		lock:        sync.Mutex{},
		nWorkers:    nWorkers,
		channels:    []chan zmq4.Msg{},
	}
}

// principal method to start server and listen for requests
func (s *LoadBServer) Listen() error {
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
	//initialize workers
	for i := 0; i < s.nWorkers; i++ {
		channel := make(chan zmq4.Msg)
		s.channels = append(s.channels, channel)
		go s.worker(channel, i+1)
	}
	go func() {
		defer socket.Close()
		nexInLine := 0
		for {
			message, err := socket.Recv()
			if err == nil {
				s.channels[nexInLine] <- message
				nexInLine = (nexInLine + 1) % s.nWorkers
			}
		}
	}()
	return nil
}

// internal method to process each request message, it validates the message and communicates with the service
func (s *LoadBServer) worker(channel chan zmq4.Msg, goRoutineId int) {
	//if there was an error with the mesage we ignore it then
	for message := range channel {
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
		fmt.Printf("-----------receiving request to process from goroutine %d\n", goRoutineId)
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
			//send response to the spceified client
			responseBytes, _ = json.Marshal(response)
		}
		//send message with client ID (if recived one, means, we are using proxy)
		fmt.Printf("***************GO ROUTINE %d IS GOING TO SEND ANSWER FROM REQUEST\n", goRoutineId)
		if err := s.socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes, clientId)); err != nil {
			log.Printf("ERROR SENDING aANSWERRRRRRRRRRR %v", err.Error())
			continue
		}
		fmt.Printf("++++++++++GO ROUTINE %d FINISHED WITH REQUEST\n", goRoutineId)
		if err != nil {
			//check if we completed all the faculties so we send the end signal
			s.lock.Lock()
			s.counter++
			if s.counter == s.faculties {
				fmt.Printf("SENMDING ENDING SIGNALLLLLLLLLLLL goroutine %d\n", goRoutineId)
				s.lock.Unlock()
				s.endChannel <- true
			}
			s.lock.Unlock()
		}
	}
}
