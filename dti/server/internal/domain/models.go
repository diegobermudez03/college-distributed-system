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
	Alerts []AlertModel `gorm:"foreignKey:SemesterId"`
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
	Alerts []AlertModel `gorm:"foreignKey:ProgramId"`
}

type AssignationModel struct{
	ID 			uuid.UUID	`gorm:"type:uuid;primaryKey"`
	SemesterId 	uuid.UUID
	ProgramId 	uuid.UUID
	Classrooms 	int
	Labs 		int
	MobileLabs 	int
	CreatedAt 	time.Time
	GoRoutineId	int `gorm:"-"`
	ProgramName string `gorm:"-"`
	SemesterName string `gorm:"-"`
	RemainingCLassrooms int `gorm:"-"`
	RemainingLabs int `gorm:"-"`
	RemainingMobileLabs int `gorm:"-"`
}

type AlertModel struct{
	ID 			uuid.UUID 	`gorm:"type:uuid;primaryKey"`
	ProgramId 	uuid.UUID
	SemesterId 	uuid.UUID
	Message 	string
	RequestedClassrooms 	int
	RequestedLabs 		int
	CreatedAt 	time.Time
	AvailableClassrooms 	int
	AvailableLabs 		int
	AvailableMobileLabs 	int
	GoRoutineId	int `gorm:"-"`
	ProgramName string `gorm:"-"`
	SemesterName string `gorm:"-"`
}

//NON DB TABLES
type AssignedSemesterResources struct{
	Classrooms 	int
	Labs 		int 
	MobileLabs 	int
}