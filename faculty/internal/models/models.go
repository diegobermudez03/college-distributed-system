package models

import "github.com/google/uuid"

type ResponseRequest struct {
	ClientId 		  uuid.UUID  `json:"-"`
	Status            string `json:"status"`
	ClassroomsAsigned int    `json:"classrooms-assigned"`
	LabsAsigned       int    `json:"labs-assigned"`
	ErrorRequest 	  bool 		`json:"error-request"`
}

type ProgramRequest struct {
	ClientSocketId []byte    `json:"-"`
	ClientId       uuid.UUID `json:"client-id"`
	ProgramName    string    `json:"program-name" validate:"required"`
	Semester       string    `json:"semester" validate:"required"`
	Classrooms     int      `json:"classrooms" validate:"required"`
	Labs           int      `json:"labs" validate:"required"`
}

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
