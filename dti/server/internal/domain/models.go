package domain

import (
	"time"

	"github.com/google/uuid"
)

type SemesterAvailabilityModel struct {
	ID 			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	Semester 	string 		`gorm:"uniqueIndex;not null"`
	Classrooms 	int			
	Labs 		int
	MobileLabs 	int
	Assignations []AssignationModel `gorm:"foreignKey:SemesterId"`
}


type FacultyModel struct{
	ID 		uuid.UUID	`gorm:"type:uuid;primaryKey"`
	Name 	string 		`gorm:"uniqueIndex;not null"`
	Programs []ProgramModel `gorm:"foreignKey:FacultyId"`
} 

type ProgramModel struct{
	ID 		uuid.UUID	`gorm:"type:uuid;primaryKey"`
	Name 	string 		`gorm:"uniqueIndex;not null"`
	FacultyId uuid.UUID
	Faculty FacultyModel `gorm:"references:ID"`
	Assignations []AssignationModel `gorm:"foreignKey:ProgramId"`
}


type AssignationModel struct{
	ID 			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	SemesterId 	uuid.UUID
	ProgramId 	uuid.UUID
	RequestedClassrooms 	int
	RequestedLabs 		int
	Classrooms 	int
	Labs 		int
	MobileLabs 	int
	CreatedAt 	time.Time
	Alert 		bool
	GoRoutineId	int `gorm:"-"`
	ProgramName string `gorm:"-"`
	SemesterName string `gorm:"-"`
	RemainingCLassrooms int `gorm:"-"`
	RemainingLabs int `gorm:"-"`
	RemainingMobileLabs int `gorm:"-"`
}

//NON DB TABLES
type AssignedSemesterResources struct{
	Classrooms 	int
	Labs 		int 
	MobileLabs 	int
}