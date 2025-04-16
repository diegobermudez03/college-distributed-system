package client

import (
	"github.com/diegobermudez03/college-distributed-system/faculty/internal/models"
	"github.com/google/uuid"
)

type FacultyClient struct {
	dtiAddress  string
	facultyName string
	semester    string
}

func NewFacultyClient(dtiAddress, semester, facultyName string) *FacultyClient {
	return &FacultyClient{
		dtiAddress:  dtiAddress,
		semester:    semester,
		facultyName: facultyName,
	}
}

func (c *FacultyClient) SendFacultyRequest(programs map[uuid.UUID]models.ProgramRequest) ([]models.ResponseRequest,error) {
	responses := make([]models.ResponseRequest, 0, len(programs))
	for _, program := range programs {
		response := models.ResponseRequest{
			ClientId:         program.ClientId,
			Status:           "OK",
			ClassroomsAsigned: program.Classrooms,
			LabsAsigned:     program.Labs,
			ErrorRequest: false,
		}
		responses = append(responses, response)
	}
	return responses,nil
}