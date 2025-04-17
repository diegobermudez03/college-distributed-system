package service

import (
	"sync"

	"github.com/google/uuid"
)

type SemesterAvailability struct {
	Id 			uuid.UUID
	resourcesLock      *sync.Mutex
	Semester   string
	Classrooms int
	Labs       int
	MobileLabs int
}