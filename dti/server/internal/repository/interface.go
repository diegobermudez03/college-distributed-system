package repository

import (
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/google/uuid"
)

type CollegeRepository interface {
	CreateFaculty(faculty *domain.FacultyModel) error
	GetFacultiesCount() (int, error)
	GetSemester(semester string) (*domain.SemesterAvailabilityModel, error)
	CreateSemester(semester *domain.SemesterAvailabilityModel) error
	GetAssignedResourcesOfSemester(semesterId uuid.UUID) (*domain.AssignedSemesterResources, error)	//if no resources must return pointer to zero struct
	GetFullFacultyById(facultyId uuid.UUID) (*domain.FacultyModel, error)
	GetAllFaculties() ([]domain.FacultyModel, error)
	CreateAssignation(assignation *domain.AssignationModel) error
	CreateAlert(alert *domain.AlertModel) error
	//CreateProgram(program *domain.ProgramModel) error
}