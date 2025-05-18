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
	port    int
	service domain.CollegeService
	frontend  zmq4.Socket
	counter int
	endChannel chan bool
	proxyServer string
	faculties int
	lock sync.Mutex
	nWorkers int
	backendAddress string
}

func NewLoadBServer(service domain.CollegeService, config domain.ServerConfig, nWorkers int) *LoadBServer {
	return &LoadBServer{
		port: config.ListenPort,
		service: service,
		faculties: config.NumFaculties,
		endChannel: config.EndChannel,
		proxyServer: config.ProxyServer,
		lock: sync.Mutex{},
		nWorkers: nWorkers,
		backendAddress: "inproc://backend",
	}
}

//principal method to start server and listen for requests
func (s *LoadBServer) Listen() error {
	//first we poblate the DB
	if err:= s.service.PoblateFacultiesAndPrograms(); err != nil{
		return err
	}

	//if we are using a proxy server, then we suscribe
	if s.proxyServer != ""{
		//establish connection with the faculty
		conn, err := net.Dial("tcp", s.proxyServer)
		if err != nil{
			log.Print(err.Error())
			return errors.New("unable to suscribe with proxy")
		}
		message := strconv.Itoa(s.port) + "\n"
		_, err = conn.Write([]byte(message))
		if err != nil{
			log.Print(err.Error())
			return errors.New("unable to suscribe with proxy")
		}
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil{
			log.Print(err.Error())
			return errors.New("unable to suscribe with proxy")
		}
		if response != "OK\n"{
			return errors.New("proxy didnt accept the suscription")
		}
		conn.Close()
	}

	ctx, _ := context.WithCancel(context.Background())
	//create zeromq socket and listen in the given port
	frontend := zmq4.NewRouter(ctx)
	s.frontend = frontend
	if err := frontend.Listen(fmt.Sprintf("tcp://*:%d", s.port)); err != nil {
		return err
	}
	backend := zmq4.NewDealer(ctx)
	if err := backend.Listen(s.backendAddress); err != nil{
		return err 
	}

	for i := 0; i < s.nWorkers; i++{
		go s.worker(ctx, i)
	}

	zmq4.NewProxy(ctx, frontend, backend, nil).Run()
	return nil
}


/*
	Method for the workers
*/
func (s *LoadBServer) worker(ctx context.Context ,goRoutineId int){
	socket := zmq4.NewDealer(ctx)
	defer socket.Close()
	if err := socket.Dial(s.backendAddress); err !=nil{
		return
	}

	for{
		message, err := socket.Recv()
		if err != nil{
			continue 
		}
		clientIdentity := message.Frames[0]

		//for proxy purposes
		var clientId []byte = []byte{} 
		if len(message.Frames) > 2{
			clientId = message.Frames[2]
		}
		//read request body
		//if the message is of acceptance, then we ignore
		if string(message.Frames[1]) == "ACCEPT"{
			continue 
		}
		clientRequestBytes := message.Frames[1]
		clientRequest := domain.DTIRequestDTO{}
		log.Print("before health check")

		////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////
		//if message wasnt a request, we check if it was a HEALTH CHECK
		if err := json.Unmarshal(clientRequestBytes, &clientRequest); err != nil || clientRequest.Semester==""{
			hCheck := HealthCheckDTO{}
			log.Print("in the health check")
			if err := json.Unmarshal(clientRequestBytes, &hCheck); err != nil{
				log.Print(err.Error())
				continue 
			}
			//if it was a health check, we answer with a simple 1 byte
			log.Println("ANSWERING HEALTH CHECK")
			socket.Send(zmq4.NewMsgFrom(clientIdentity, []byte{1}))
			continue 
		}
		log.Print("after healthcheck")
		////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////

		//process message with the service
		response, err := s.service.ProcessRequest(clientRequest, goRoutineId)
		var responseBytes []byte
		if err != nil{
			//if there was an error, we send it in the authorized format
			errorResponse := domain.DTIResponseDTO{
				Semester: clientRequest.Semester,
				ErrorFound: true,
				ErrorMessage: err.Error(),
			}
			responseBytes, _ = json.Marshal(errorResponse)
		}else{
			//check if we completed all the faculties so we send the end signal
			s.lock.Lock()
			s.counter++
			if s.counter == s.faculties{
				s.endChannel <-true
			}
			s.lock.Unlock()
			//send response to the spceified client
			responseBytes, _ = json.Marshal(response)
		}
		//send message with client ID (if recived one, means, we are using proxy)
		socket.Send(zmq4.NewMsgFrom(clientIdentity, responseBytes, clientId))
	}

}
