package repository

import (
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/google/uuid"
)

type DBMock struct {
}

// CreateAssignation implements CollegeRepository.
func (d *DBMock) CreateAssignation(assignation *domain.AssignationModel) error {
	return nil
}

// CreateFaculty implements CollegeRepository.
func (d *DBMock) CreateFaculty(faculty *domain.FacultyModel) error {
	return nil
}

// CreateSemester implements CollegeRepository.
func (d *DBMock) CreateSemester(semester *domain.SemesterAvailabilityModel) error {
	return nil
}

// GetAllFaculties implements CollegeRepository.
func (d *DBMock) GetAllFaculties() ([]domain.FacultyModel, error) {
	return nil, nil
}

// GetAssignedResourcesOfSemester implements CollegeRepository.
func (d *DBMock) GetAssignedResourcesOfSemester(semesterId uuid.UUID) (*domain.AssignedSemesterResources, error) {
	return nil, nil
}

// GetFacultiesCount implements CollegeRepository.
func (d *DBMock) GetFacultiesCount() (int, error) {
	return 10, nil
}

// GetFullFacultyById implements CollegeRepository.
func (d *DBMock) GetFullFacultyById(facultyId uuid.UUID) (*domain.FacultyModel, error) {
	return nil, nil
}

// GetProgramAssignment implements CollegeRepository.
func (d *DBMock) GetProgramAssignment(programId uuid.UUID, semesterId uuid.UUID) (*domain.AssignationModel, error) {
	return nil, nil
}

// GetSemester implements CollegeRepository.
func (d *DBMock) GetSemester(semester string) (*domain.SemesterAvailabilityModel, error) {
	return nil, nil
}

func NewDBMock() CollegeRepository {
	return &DBMock{}
}
