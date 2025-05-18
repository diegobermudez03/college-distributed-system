package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)



type HealthCheckDTO struct {
	HealthCheck bool `json:"health-check"`
}


const (
	portArg = "--port"
	suscribePortArg = "--sus-port"
	healthTickerArg = "--healt-ticker"
)

var suscribedServers []string = []string{}
var activeServer string = ""
var clientSocket zmq4.Socket = nil
var serverSocket zmq4.Socket = nil
var serverLock sync.Mutex = sync.Mutex{}
var activeClients map[uuid.UUID][]byte = map[uuid.UUID][]byte{}


func main() {
	port := ":6000"
	susPort := ":7000"
	healthTicker := 5
	//read arguments
	for _, arg := range os.Args{
		if !strings.Contains(arg, "="){
			continue
		}
		parts := strings.Split(arg, "=")
		switch parts[0]{
		case portArg: port = ":" + parts[1]
		case suscribePortArg: susPort = ":" + parts[1]
		case healthTickerArg: healthTicker, _ = strconv.Atoi(parts[1])
		}
	}

	//start listener for suscribe servers
	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := listenForSuscriptions(susPort, &wg); err != nil{
		fmt.Println("1")
		log.Fatal(err.Error())
	}
	wg.Wait() //wait for having at least 1 server
	//start the server for ZeroMQ
	messageChannel := make(chan [][]byte)
	if err := listenForZeroMQClients(port, messageChannel); err != nil{
		fmt.Println("2")
		log.Fatal(err.Error())
	}

	//start the ZeroMQ client with the active server
	healthCheckChannel := make(chan bool)
	go startHealthCheck(healthCheckChannel, healthTicker)
	if err := startZeroMQClient(messageChannel, healthCheckChannel); err != nil{
		fmt.Println("3")
		log.Fatal(err.Error())
	}
}


/*
	Method to listen for instances trying to suscribe in the TCP port
*/
func listenForSuscriptions(susPort string, wg *sync.WaitGroup) error{
	//start listener for suscribe servers
	listener, err := net.Listen("tcp", "127.0.0.1" +susPort)
	if err != nil{
		return errors.New("unable to listen for suscribe connections")
	}
	log.Printf("Listening for suscriptions on: %s", listener.Addr().String())
	go func(){
		defer listener.Close()
		for{
			//handle tcp message, server suscription
			conn, err := listener.Accept()
			if err != nil{
				continue
			}
			reader := bufio.NewReader(conn)
			//obtain the server IP and the port which is in the TCP message
			serverAddress := conn.RemoteAddr().String()
			parts := strings.Split(serverAddress, ":")
			ipAddress := parts[0]
			port, err := reader.ReadString('\n')
			port = strings.ReplaceAll(port,"\n", "")
			if err != nil{
				conn.Close()
				continue
			}
			server :=  ipAddress + ":" +port
			suscribedServers = append(suscribedServers, server)
			log.Printf("New suscriber: %s", server)
			//if we dont have yet an active server, then we set this one as it
			if activeServer == ""{
				activeServer = suscribedServers[0]
				suscribedServers = suscribedServers[1:]
				log.Printf("Active server selected: %s", activeServer)
				wg.Done()
			}
			//add suscribed server
			conn.Write([]byte("OK\n"))
			conn.Close()
		}
	}()
	return nil
}


/*
	Method to start the ZeroMQ router to listen to client requests
*/
func listenForZeroMQClients(port string, channel chan [][]byte) error{
	socket := zmq4.NewRouter(context.Background())
	if err := socket.Listen(fmt.Sprintf("tcp://*%s", port)); err != nil {
		return err
	}
	serverSocket = socket
	//listen to all messages received
	go func(){
		defer socket.Close()
		for{
			message, err := socket.Recv()
			if err != nil{
				continue
			}
			channel <- message.Frames
		}
	}()
	return nil
}

/*
	Method to start the ZeroMQ client
*/
func startZeroMQClient(channel chan [][]byte, healthCheck chan bool)error{
	stablishConnection()
	go func(){
		for message := range channel{
			clientId := uuid.New()
			activeClients[clientId] = message[0]
			serverLock.Lock()
			if err := clientSocket.Send(zmq4.NewMsgFrom(message[1], []byte(clientId.String()))); err != nil{
				serverLock.Unlock()
				continue
			}
			serverLock.Unlock()
		}
	}()

	//start listener for server responses
	for{
		log.Print("listening in the clientsocket")
		if clientSocket == nil{
			continue
		}
		response, err := clientSocket.Recv()
		log.Print("message received")
		if err != nil{
			clientSocket = nil
			fmt.Println("in the err != nil")
			log.Print(err.Error())
			continue
		}
		//if it was a healthcheck
		log.Print("before chekcing 1")
		if response.Frames[0][0] == 1{
			healthCheck <- true
			continue 
		}
		log.Print("after checking 1")
		//if it wasnt a health check but an actual response
		clientId, err := uuid.Parse(string(response.Frames[2]))
		if err != nil{
			continue
		}
		clientIdentity, ok := activeClients[clientId]
		if  !ok{
			continue
		}
		serverSocket.Send(zmq4.NewMsgFrom(clientIdentity, response.Frames[1]))
		delete(activeClients, clientId)
	}
	return nil
}


/*
	Method for sending the health check to the server
*/
func startHealthCheck(responseChannel chan bool, healthTicker int){
	ticker := time.NewTicker(time.Second * time.Duration(healthTicker))
	helathCheckMessage := HealthCheckDTO{
		HealthCheck: true,
	}
	message, _ := json.Marshal(helathCheckMessage)
	for range ticker.C{
		if clientSocket == nil{
			stablishConnection()
			continue
		}
		serverLock.Lock()
		log.Printf("sending health check to: %s", activeServer)
		if err := clientSocket.Send(zmq4.NewMsgFrom(message)); err != nil{
			serverLock.Unlock()
			continue
		}
		serverLock.Unlock()
		timeout := time.Tick(time.Second * 5)
		//wait to see if we receive the health check answer, or if not, then we need to find another server
		select{
		case <-responseChannel: log.Print("Health check OK")
		case <- timeout: {
				log.Printf("Health check failed")
				stablishConnection()
			}
		}
	}
}

/*
	Method to try to stabish connection with a new server
*/
func stablishConnection(){
	ctx, _ := context.WithCancel(context.Background())
	clientSocket = zmq4.NewDealer(ctx, zmq4.WithAutomaticReconnect(true))
	//try to connect to servers until one does connect
	serverLock.Lock()
	defer serverLock.Unlock()
	for{
		if len(suscribedServers) > 0{
			activeServer = suscribedServers[0]
			suscribedServers = suscribedServers[1:]
		}
		log.Printf("Sending active server to: %s", activeServer)
		if err := clientSocket.Dial(fmt.Sprintf("tcp://%s", activeServer)); err != nil{
			log.Printf("Active server request failed to %s", activeServer)
			log.Print(err.Error())
			continue
		}
		log.Printf("active server switched to %s", activeServer)
		break
	}
}