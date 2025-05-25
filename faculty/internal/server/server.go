package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)

const (
	internalError   = "INTERNAL-ERROR"
	invalidSemester = "INVALID-SEMESTER"
)

// server model
type FacultyServer struct {
	listenPort  int
	minPrograms int
	semester    string
	facultyName string
	dtiAddress  string
	clients     map[uuid.UUID]models.ProgramRequest //map of clients
}

func NewFacultyServer(listenPort, minPrograms int, semester string, facultyName, dtiAddress string) *FacultyServer {
	return &FacultyServer{
		listenPort:  listenPort,
		minPrograms: minPrograms,
		clients:     map[uuid.UUID]models.ProgramRequest{},
		semester:    semester,
		facultyName: facultyName,
		dtiAddress:  dtiAddress,
	}
}

func (s *FacultyServer) Listen() error {
	//start zeromq request reply server
	socket := zmq4.NewRouter(context.Background())
	defer socket.Close()
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.listenPort)); err != nil {
		return errors.New(fmt.Sprint("Unable to start server at port ", s.listenPort))
	}
	log.Printf("Listening at port %d", s.listenPort)

	//read for program requests
	s.readPrograms(socket)

	//SEND REQUEST TO THE DTI
	//connect to the DTI server
	dtiSocket := zmq4.NewDealer(context.Background(), zmq4.WithAutomaticReconnect(true))
	if err := dtiSocket.Dial(fmt.Sprintf("tcp://%s", s.dtiAddress)); err != nil {
		return err
	}
	defer dtiSocket.Close()
	//generate payload to send
	dtiRequest := models.DTIRequest{
		Semester:    s.semester,
		FacultyName: s.facultyName,
		Programs:    make([]models.DTIProgramRequest, 0, len(s.clients)),
	}
	//add program requests
	for _, program := range s.clients {
		dtiRequest.Programs = append(dtiRequest.Programs, models.DTIProgramRequest{
			ProgramId:   program.ClientId,
			ProgramName: program.ProgramName,
			Classrooms:  program.Classrooms,
			Labs:        program.Labs,
		})
	}
	requestBytes, _ := json.Marshal(dtiRequest)

	//SEND REQUEST TO DTI
	startTime := time.Now()
	log.Printf("Sending request from Faculty %s Semester %s to DTI", s.facultyName, s.semester)
	if err := dtiSocket.Send(zmq4.NewMsgFrom(requestBytes)); err != nil {
		return err
	}

	allocation, err := dtiSocket.Recv()
	if err != nil {
		return err
	}
	elapsed := time.Since(startTime).Milliseconds()
	//SEND REQUEST TO DTI

	dtiResponse := models.DTIResponse{}
	if err := json.Unmarshal(allocation.Bytes(), &dtiResponse); err != nil {
		return err
	}
	log.Printf("Received DTI response for faculty %s semester %s", dtiRequest.FacultyName, dtiResponse.Semester)

	//send responses to program clients
	if dtiResponse.ErrorFound {
		log.Printf("Error received from DTI: %s", dtiResponse.ErrorMessage)
	}

	//iterate over all responses, get the socket ID for each one, and then send the JSON response
	for _, clientResponse := range dtiResponse.Programs {
		client, ok := s.clients[clientResponse.ProgramId]
		log.Println("Sending reply to client ", client.ProgramName)
		if !ok {
			continue
		}
		//transform the response into the valid dto and answer to the program
		clientDTO := models.ProgramResponse{
			ClientId:           client.ClientId,
			Status:             clientResponse.StatusMessage,
			ClassroomsAsigned:  clientResponse.Classrooms,
			LabsAsigned:        clientResponse.Labs,
			MobileLabsAssigned: clientResponse.MobileLabs,
		}
		bytes, _ := json.Marshal(clientDTO)
		socket.Send(zmq4.NewMsgFrom(client.ClientSocketId, bytes))
	}
	log.Printf("Completed faculty %s in %d ms", s.facultyName, elapsed)
	return nil
}

func (s FacultyServer) readPrograms(socket zmq4.Socket) {
	//read requests from CLIENT PROGRAMS
	for {
		//extract the id of the client
		message, err := socket.Recv()
		clientId := message.Frames[0]
		if err != nil {
			//if theres an error we are going to suppose that it was due to the program client
			continue
		}
		//create program request and unmarshal it from program message
		programRequest := models.ProgramRequest{
			ClientSocketId: clientId,
			ClientId:       uuid.New(),
		}
		if err := json.Unmarshal(message.Frames[2], &programRequest); err != nil {
			//if there's an error reading, then we are going to suppose that
			//the program did something wrong, but we wont break, so we simply ignore
			continue
		}
		//icheck if the program is of our semester
		if programRequest.Semester != s.semester {
			errorResponse := models.ProgramResponse{
				ClientId: programRequest.ClientId,
				Status:   invalidSemester,
			}
			errorBytes, _ := json.Marshal(errorResponse)
			socket.Send(zmq4.NewMsgFrom(programRequest.ClientSocketId, errorBytes))
			continue
		}
		//save client
		s.clients[programRequest.ClientId] = programRequest

		//if the semester is complete, then we continue to the request with the DTI
		if len(s.clients) == s.minPrograms {
			break
		}
	}
}
