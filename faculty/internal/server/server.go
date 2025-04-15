package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-zeromq/zmq4"
)

//models
type responseRequest struct{
	Status string `json:"status"`
	ClassroomsAsigned int `json:"classrooms-assigned"`
	LabsAsigned int `json:"labs-assigned"`
}

//server model
type FacultyServer struct {
	listenPort int
	minPrograms int
}

func NewFacultyServer(listenPort, minPrograms int) *FacultyServer {
	return &FacultyServer{
		listenPort: listenPort,
		minPrograms: minPrograms,
	}
}

// method to run the zeromq request reply server
func (s *FacultyServer) Run() error {
	//start zeromq request reply server
	socket := zmq4.NewRep(context.Background())
	defer socket.Close()
	if err := socket.Listen(fmt.Sprintf("tcp://*:%d", s.listenPort)); err != nil {
		return errors.New(fmt.Sprint("Unable to start server at port ", s.listenPort))
	}
	log.Printf("Listening at port %d", s.listenPort)
	//listen for the program requests
	if err := s.listenProgramRequests(socket); err != nil{
		return err
	}
	return nil
}

func (s *FacultyServer) listenProgramRequests(socket zmq4.Socket) error{
	for i := 0; i < s.minPrograms; i++ {
		message, err := socket.Recv()
		if err != nil{
			return err
		}
		log.Print(message.String())

		response := responseRequest{
			Status: "estatus",
			ClassroomsAsigned: 5,
			LabsAsigned: 10,
		}
		jsonBytes, err := json.Marshal(response)
		if err != nil{
			return err
		}
		zmq4.New
		socket.Send(zmq4.NewMsgString(string(jsonBytes)))
	}
	return nil
}