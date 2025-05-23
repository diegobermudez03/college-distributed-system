package models

import "github.com/google/uuid"

//internal models
type SemesterRequest struct{
	Semester string `json:"semester"`
	Programs map[uuid.UUID]ProgramRequest
}

//SERVER DTOS FOR COMMUNICATION WITH PROGRAMS
type ProgramResponse struct {
	ClientId 		  uuid.UUID `json:"-"`
	Status            string 	`json:"status"`
	ClassroomsAsigned int    	`json:"classrooms-assigned"`
	LabsAsigned       int    	`json:"labs-assigned"`
	MobileLabsAssigned int 	 	`json:"mobile-labs-assigned"`
}

type ProgramRequest struct {
	ClientSocketId []byte    `json:"-"`
	ClientId       uuid.UUID `json:"client-id"`
	ProgramName    string    `json:"program-name" validate:"required"`
	Semester       string    `json:"semester" validate:"required"`
	Classrooms     int      `json:"classrooms" validate:"required"`
	Labs           int      `json:"labs" validate:"required"`
}



//REQUEST COMMUNICATION WITH THE DTI
type DTIRequest struct{
	Semester 	string `json:"semester"`
	FacultyName string 	`json:"faculty-name"`
	Programs []DTIProgramRequest `json:"programs"`
}

type DTIProgramRequest struct{
	ProgramId 	uuid.UUID `json:"program-id"`
	ProgramName string `json:"program-name"`
	Classrooms  int    `json:"classrooms"`
	Labs        int    `json:"labs"`
}

//RESPONSE COMMUNICATION WITH THE DTI
type DTIResponse struct{
	Semester 	string `json:"semester"`
	ErrorFound 	bool 				`json:"error-found"`
	ErrorMessage string 			`json:"error-message"`
	Programs []DTIProgramResponse `json:"programs"`
}

type DTIProgramResponse struct{
	ProgramId     uuid.UUID `json:"program-id"`
	ProgramName 	string 	`json:"program-name"`
	RequestedClassrooms 	int 	`json:"requested-classrooms"`
	RequestedLabs 			int 	`json:"requested-labs"`
	Classrooms    int       `json:"classrooms"`
	Labs          int       `json:"labs"`
	MobileLabs    int       `json:"mobile-labs"`
	StatusMessage string    `json:"status-message"`
}
