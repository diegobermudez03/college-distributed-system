package service

import "github.com/google/uuid"

//models
type ServiceConfig struct {
	Classrooms int
	Labs       int
	MobileLabs int
}

//dtos
type DTIResponseDTO struct {
	Semester string               `json:"semester"`
	Programs []struct{
		ProgramId     uuid.UUID `json:"program-id"`
		Classrooms    int       `json:"classrooms"`
		Labs          int       `json:"labs"`
		StatusMessage string    `json:"status-message"`
	} `json:"programs"`
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
	ProcessRequest(request DTIRequestDTO) (*DTIResponseDTO, error)
}