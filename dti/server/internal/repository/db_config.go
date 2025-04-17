package repository

import (
	"fmt"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SslMode  string
	Timezone string
}

func OpenPostgresDb(config PostgresConfig) (*gorm.DB,error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Host, config.Port, config.User, config.Password, config.DbName, config.SslMode, config.Timezone)
	//open db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		return nil, err 
	}
	//clean db just in case, each server execution creates a new db
	db.Exec(`DROP SCHEMA public CASCADE;`)
	db.Exec(`CREATE SCHEMA public;`)

	//migrate db, which means, create tables
	err = db.AutoMigrate(&domain.SemesterAvailabilityModel{}, &domain.FacultyModel{}, 
		&domain.ProgramModel{}, &domain.AssignationModel{}, &domain.AlertModel{})
	if err != nil{
		return nil, err
	}
	return db, nil
}