package domain

import (
	"errors"

	"github.com/google/uuid"
)

//status messages
const (
	InvalidProgramMsg = "INVALID-PROGRAM" 
	InvalidFacultyMsg = "INVALID-FACULTY"
	NotEnoughResourcesMsg = "NOT-ENOUGH-RESOURCES-FOR-ASSIGNMENT"
	AlreadyHaveAssignation = "PROGRAM-ALREADY-HAS-RESOURCES-FOR-SEMESTER"
	OkMsg = "OK"
)

var (
	ErrorStartingService = errors.New("unable to start service")
	ErrorSemesterWasAlreadyProcessed = errors.New("Semester already processed")
	ErrorFacultyInvalidSemester = errors.New("Faculty semester is not the proccesed one")
)


//models
type ServiceConfig struct {
	Semester string
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
	MobileLabs    int       `json:"mobile-labs"`
	StatusMessage string    `json:"status-message"`
}


type DTIRequestDTO struct{
	Semester 	string `json:"semester"`
	FacultyName string 	`json:"faculty-name"`
	Programs []DTIProgramRequestDTO `json:"programs"`
}

type DTIProgramRequestDTO struct{
	ProgramId 	uuid.UUID `json:"program-id"`
	ProgramName string `json:"program-name"`
	Classrooms  int    `json:"classrooms"`
	Labs        int    `json:"labs"`
}
//interfaces
type CollegeService interface {
	PoblateFacultiesAndPrograms() error
	ProcessRequest(request DTIRequestDTO, goRoutineId int) (*DTIResponseDTO, error)
}