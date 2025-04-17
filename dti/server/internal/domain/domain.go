package domain

import "github.com/google/uuid"

//models
type ServiceConfig struct {
	Classrooms int
	Labs       int
	MobileLabs int
}

//dtos
type DTIResponseDTO struct {
	Semester string               	`json:"semester"`
	ErrorFound 	bool 				`json:"error-found"`
	ErrorMessage string 			`json:"error-message"`
	Programs []DTIProgramResponseDTO `json:"programs"`
}

type DTIProgramResponseDTO struct{
	ProgramId     uuid.UUID `json:"program-id"`
	Classrooms    int       `json:"classrooms"`
	Labs          int       `json:"labs"`
	StatusMessage string    `json:"status-message"`
}


type DTIRequestDTO struct{
	Semester 	string `json:"semester"`
	FacultyName string 	`json:"faculty-name"`
	Programs []struct{
		ProgramId 	uuid.UUID `json:"program-id"`
		ProgramName string `json:"program-name"`
		Classrooms  int    `json:"classrooms"`
		Labs        int    `json:"labs"`
	} `json:"programs"`
}

//interfaces
type CollegeService interface {
	PoblateFacultiesAndPrograms() error
	ProcessRequest(request DTIRequestDTO) (*DTIResponseDTO, error)
}