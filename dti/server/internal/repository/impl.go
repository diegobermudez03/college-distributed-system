package repository

import (
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
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