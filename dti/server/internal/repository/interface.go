package repository

import "github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"

type CollegeRepository interface {
	CreateFaculty(faculty *domain.FacultyModel) error
	//CreateProgram(program *domain.ProgramModel) error
}