package repository

import (
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollegeRepositoryPostgres struct{
	db *gorm.DB
}

func NewCollegeRepositoryPostgres(db *gorm.DB) CollegeRepository {
	return &CollegeRepositoryPostgres{
		db: db,
	}
}

func (r *CollegeRepositoryPostgres) CreateFaculty(faculty *domain.FacultyModel) error{
	return r.db.Create(&faculty).Error
}

func (r *CollegeRepositoryPostgres) GetFacultiesCount() (int, error){
	var count int64
	err := r.db.Model(&domain.FacultyModel{}).Count(&count).Error
	return int(count), err
}

func (r *CollegeRepositoryPostgres) GetSemester(semester string) (*domain.SemesterAvailabilityModel, error){
	var semesterModel domain.SemesterAvailabilityModel
	err := r.db.Where("semester = ?", semester).First(&semesterModel).Error
	return &semesterModel, err
}

func (r *CollegeRepositoryPostgres) CreateSemester(semester *domain.SemesterAvailabilityModel) error{
	return r.db.Create(&semester).Error
}

func (r *CollegeRepositoryPostgres) GetAssignedResourcesOfSemester(semesterId uuid.UUID) (*domain.AssignedSemesterResources, error){
	assignment := &domain.AssignedSemesterResources{}
	err := r.db.Raw(` SELECT 
		COALESCE(SUM(classrooms),0) AS classrooms, 
		COALESCE(SUM(labs), 0) AS labs, 
		COALESCE(SUM(mobile_labs),0) AS mobile_labs 
		FROM assignation_models 
		WHERE semester_id = ?`, semesterId).Scan(assignment).Error
	return assignment, err
}


func (r *CollegeRepositoryPostgres) GetFullFacultyById(facultyId uuid.UUID) (*domain.FacultyModel, error){
	var faculty domain.FacultyModel
	//preload since we need the programs to be fill
	err := r.db.Preload("Programs").Where("id = ?", facultyId).First(&faculty).Error
	return &faculty, err
}

func (r *CollegeRepositoryPostgres) GetAllFaculties() ([]domain.FacultyModel, error){
	var faculties []domain.FacultyModel
	err := r.db.Find(&faculties).Error
	return faculties, err
}

func (r *CollegeRepositoryPostgres) CreateAssignation(assignation *domain.AssignationModel) error{
	return r.db.Create(&assignation).Error
}

func (r *CollegeRepositoryPostgres) CreateAlert(alert *domain.AlertModel) error{
	return r.db.Create(&alert).Error
}


