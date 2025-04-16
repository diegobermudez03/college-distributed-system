package service

import "github.com/google/uuid"

type SemesterAvailabilityModel struct {
	Id 			uuid.UUID	`gorm:"primaryKey"`
	Semester 	string 		`gorm:"uniqueIndex;not null"`
	Classrooms 	int			
	Labs 		int
	MobileLabs 	int
}

type FacultyModel struct{
	Id 		uuid.UUID	`gorm:"primaryKey"`
	Name 	string 		`gorm:"uniqueIndex;not null"`
} 

type ProgramModel struct{
	Id 		uuid.UUID	`gorm:"primaryKey"`
	Name 	string 		`gorm:"uniqueIndex;not null"`
}