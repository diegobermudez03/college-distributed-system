package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/zeromq/goczmq"
)

type LoadBServer struct {
	port        int
	service     domain.CollegeService
	counter     int
	endChannel  chan bool
	proxyServer string
	faculties   int
	lock        sync.Mutex
	nWorkers    int
}

func NewLoadBServer(service domain.CollegeService, config domain.ServerConfig, nWorkers int) *LoadBServer {
	return &LoadBServer{
		port:        config.ListenPort,
		service:     service,
		faculties:   config.NumFaculties,
		endChannel:  config.EndChannel,
		proxyServer: config.ProxyServer,
		nWorkers:    nWorkers,
	}
}

func (s *LoadBServer) Listen() error {
	if err := s.service.PoblateFacultiesAndPrograms(); err != nil {
		return err
	}

	//if we receive proxy address, then we suscribe with it
	if s.proxyServer != "" {
		conn, err := net.Dial("tcp", s.proxyServer)
		if err != nil {
			return fmt.Errorf("proxy connection failed: %v", err)
		}
		defer conn.Close()

		if _, err := conn.Write([]byte(strconv.Itoa(s.port) + "\n")); err != nil {
			return fmt.Errorf("proxy registration failed: %v", err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil || strings.TrimSpace(string(buf[:n])) != "OK" {
			return errors.New("proxy registration rejected")
		}
	}

	//CREATING ZEROMQ PROXY
    proxy := goczmq.NewProxy()
    defer proxy.Destroy()

    //SETTING FRONTEND ROUTER FOR OUR ZEROMQ PROXY
    if err := proxy.SetFrontend(goczmq.Router, fmt.Sprintf("tcp://*:%d", s.port)); err != nil {
        return fmt.Errorf("frontend setup failed: %v", err)
    }

    //SETTING THE BACKEND FOR OUR ZEROMQ BACKEND
    if err := proxy.SetBackend(goczmq.Dealer, "inproc://backend"); err != nil {
        return fmt.Errorf("backend setup failed: %v", err)
    }

    //we start the number or workers to work with
    for i := 0; i < s.nWorkers; i++ {
        go s.worker(i + 1)
    }

    log.Printf("Server listening on port %d with %d workers", s.port, s.nWorkers)
    <-s.endChannel
    return nil
}

func (s *LoadBServer) worker(workerID int) {
	//WORKER CONNECTS WITH ITS DEALER SOCKET TO THE BACKEND
    dealer, err := goczmq.NewDealer("inproc://backend")
    if err != nil {
        log.Printf("Worker %d failed to connect: %v", workerID, err)
        return
    }
    defer dealer.Destroy()

    log.Printf("Worker %d initialized and connected", workerID)

    for {
        msg, err := dealer.RecvMessage()
        if err != nil {
            log.Printf("Worker %d receive error: %v", workerID, err)
            continue
        }

        clientIdentity := msg[0]
		body := msg[1]
		//for proxy purposes
		var clientId []byte = []byte{}
		if len(msg) > 2 {
			clientId = msg[2]
		}
		//read request body
		//if the message is of acceptance, then we ignore
		if strings.Contains(string(msg[1]),  "ACCEPT" ){
			continue
		}
		clientRequest := domain.DTIRequestDTO{}

		////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////
		//if message wasnt a request, we check if it was a HEALTH CHECK
		if err := json.Unmarshal(body, &clientRequest); err != nil || clientRequest.Semester == "" {
			hCheck := HealthCheckDTO{}
			if err := json.Unmarshal(body, &hCheck); err != nil {
				continue
			}
			//if it was a health check, we answer with a simple 1 byte
			log.Print("ANSWERING HEALTH CHECK")
			dealer.SendMessage([][]byte{clientIdentity, {1}})
			continue
		}
		////////////  HEALTH CHECK VALIDATION  //////////////////////////////////////////
	

        response, err := s.service.ProcessRequest(clientRequest, workerID)
        var respBytes []byte
        if err != nil {
            respBytes = s.createErrorResponse(clientRequest.Semester, err.Error())
        } else {
            respBytes, _ = json.Marshal(response)
        }

        // Send response
        reply := [][]byte{clientIdentity, respBytes, clientId}
        if err := dealer.SendMessage(reply); err != nil {
            log.Printf("Worker %d response error: %v", workerID, err)
        }

        // Update completion counter
        if err == nil {
            s.lock.Lock()
            s.counter++
            if s.counter >= s.faculties {
               	s.lock.Unlock()
				s.endChannel <- true
				s.endChannel <- true
            }else{
            	s.lock.Unlock()
			}
        }
    }
}

func (s *LoadBServer) createErrorResponse(semester, message string) []byte {
	errResp := domain.DTIResponseDTO{
		Semester:     semester,
		ErrorFound:   true,
		ErrorMessage: message,
	}
	resp, _ := json.Marshal(errResp)
	return resp
}