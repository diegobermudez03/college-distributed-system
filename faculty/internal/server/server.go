package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/faculty/internal/client"
	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/go-zeromq/zmq4"
	"github.com/google/uuid"
)

const (
	internalError = "INTERNAL-ERROR"
	invalidSemester = "INVALID-SEMESTER"
)


//server model
type FacultyServer struct {
	listenPort int
	minPrograms int
	semester string
	semesters map[string]map[uuid.UUID]models.ProgramRequest	//map of semesters, each one with its clients
	client *client.FacultyClient
	socket zmq4.Socket
	closeServerWg *sync.WaitGroup
}

func NewFacultyServer(listenPort, minPrograms int, semester string, client *client.FacultyClient) *FacultyServer {
	return &FacultyServer{
		listenPort: listenPort,
		minPrograms: minPrograms,
		semesters: map[string]map[uuid.UUID]models.ProgramRequest{},
		semester: semester,
		client: client,
		closeServerWg: &sync.WaitGroup{},
	}
}

//method to run the zeromq request reply server and listen for the programs requests
func (s *FacultyServer) Listen() (chan models.SemesterRequest, *sync.WaitGroup,error) {
	channel := make(chan models.SemesterRequest)
	var err error = nil
	wg := sync.WaitGroup{}
	wg.Add(1)
	outerWg := sync.WaitGroup{}
	outerWg.Add(1)
	go func (){
		//start zeromq request reply server
		socket := zmq4.NewRouter(context.Background())
		s.socket = socket
		defer socket.Close()
		if errr := socket.Listen(fmt.Sprintf("tcp://*:%d", s.listenPort)); errr != nil {
			err = errors.New(fmt.Sprint("Unable to start server at port ", s.listenPort))
		}
		log.Printf("Listening at port %d", s.listenPort)
		wg.Done()
		//listen for the program requests
		if errr := s.listenProgramRequests(channel); errr != nil{
			err = errr
		}
		s.closeServerWg.Wait()	//wait till we have sent all the replies to the programs
		//this signals the outer main go routine to stop waiting and end the execution
		outerWg.Done()
	}()
	//this is to wait til the go routine may return an error or not
	wg.Wait()
	return channel,&outerWg,err
}

//method that reads from the reponse channel and then answers to the programs
func (s *FacultyServer) SendReplies(channel chan models.DTIResponse){
	s.closeServerWg.Add(1)
	go func(){
		for response := range channel{
			log.Println("Processing DTI response")
			if response.ErrorFound{
				log.Printf("Error received from DTI: %s", response.ErrorMessage)
			}
			//get semester programs
			semesterPrograms, ok := s.semesters[response.Semester]
			log.Println("Looking for semester: ", response.Semester)
			if !ok{
				log.Println("DIDNT FIND SEMESTER PROGRAMS")
				continue
			}
			//iterate over all responses, get the socket ID for each one, and then send the JSON response
			for _, clientResponse := range response.Programs{
				client, ok := semesterPrograms[clientResponse.ProgramId]
				log.Println("Sending reply to client ", client.ProgramName)
				if !ok{
					continue 
				}
				//transform the response into the valid dto and answer to the program
				var clientDTO models.ProgramResponse
				if response.ErrorFound{
					clientDTO = models.ProgramResponse{
						ClientId: client.ClientId,
						Status: response.ErrorMessage,
						ClassroomsAsigned: 0,
						LabsAsigned: 0,
						ErrorRequest: true,
					}
				}else{
					clientDTO = models.ProgramResponse{
						ClientId: client.ClientId,
						Status: clientResponse.StatusMessage,
						ClassroomsAsigned: clientResponse.Classrooms,
						LabsAsigned: client.Labs,
						ErrorRequest: false,
					}
					if clientDTO.Status != "OK"{
						clientDTO.ErrorRequest = true
					}
				}
				bytes, _ := json.Marshal(clientDTO)
				s.socket.Send(zmq4.NewMsgFrom(client.ClientSocketId, bytes))
			}
			delete(s.semesters, response.Semester)
		}
		s.closeServerWg.Done()
	}()
}


///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////
//						INTERNAL METHODS
///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////


func (s *FacultyServer) listenProgramRequests(channel chan models.SemesterRequest) error{
	//wait for the minimum of programs to communicate with the DTI
	for{
		//extract the id of the client
		message, err := s.socket.Recv()
		clientId := message.Frames[0]
		if err != nil{
			//if theres an error we are going to suppose that it was due to the program client
			continue
		}
		//create program request and unmarshal it from program message
		programRequest := models.ProgramRequest{
			ClientSocketId: clientId,
			ClientId: uuid.New(),
		}
		if err := json.Unmarshal(message.Frames[2], &programRequest); err != nil{
			//if there's an error reading, then we are going to suppose that
			//the program did something wrong, but we wont break, so we simply ignore
			continue
		}
		//check if we have a semester set, if we have, then we only receive programs from that semester
		//if we dont have, then we accept all semester programs (using the min programs number)
		if s.semester != "" && programRequest.Semester != s.semester{
			errorResponse := models.ProgramResponse{
				ClientId: programRequest.ClientId,
				Status: invalidSemester,
				ErrorRequest: true,
			}
			errorBytes, _ := json.Marshal(errorResponse)
			s.socket.Send(zmq4.NewMsgFrom(programRequest.ClientSocketId, errorBytes))
			continue
		}
		//save client in semesters
		semesterPrograms, ok := s.semesters[programRequest.Semester]
		if !ok{
			semesterPrograms = map[uuid.UUID]models.ProgramRequest{}
			log.Println("Saving semester: ", programRequest.Semester)
			s.semesters[programRequest.Semester] = semesterPrograms
		}
		semesterPrograms[programRequest.ClientId] = programRequest
		//if the semester is complete, then we redirect it to the DTI request, we use a new go routine (thread)
		//so that we can still listen for new program requests
		if len(semesterPrograms) == s.minPrograms{
			//send request in channel for the client to manage it
			channel <- models.SemesterRequest{
				Semester: programRequest.Semester,
				Programs: semesterPrograms,
			}
			//if we had a semester configured, then we end listening
			if s.semester != ""{
				return nil
			}
		}
	}
}