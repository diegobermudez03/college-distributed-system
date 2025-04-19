package service

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/repository"
	"github.com/google/uuid"
)

type ServiceCache struct{
	Lock 		*sync.RWMutex //lock for retrieving the semester info, is Read Write since it depends on if we are reading or writing
	Semesters 	map[string]*SemesterAvailability
}

type CollegeServiceImpl struct {
	config     *domain.ServiceConfig
	repository 	repository.CollegeRepository
	Cache 		ServiceCache
	logWriterChannel chan *domain.AssignationModel
	alertWriterChannel chan *domain.AlertModel
}

func NewCollegeService(config *domain.ServiceConfig, repository repository.CollegeRepository) domain.CollegeService {
	service := &CollegeServiceImpl{
		config: config,
		repository: repository,
		Cache: ServiceCache{
			Semesters: make(map[string]*SemesterAvailability),
			Lock: &sync.RWMutex{},
		},
	}
	logWriterChannel := service.startDbLogWriter()
	alertWriterChannel := service.startDbAlertWriter()
	service.logWriterChannel = logWriterChannel
	service.alertWriterChannel = alertWriterChannel
	return service
}

//The entry method for processing a faculty request, is the main method which many go routines will execute 
//in parallel
func (s *CollegeServiceImpl) ProcessRequest(request domain.DTIRequestDTO, goRoutineId int) (*domain.DTIResponseDTO, error) {
	s.Cache.Lock.RLock()	//lock for read, since we are reading the semester
	semester, ok := s.Cache.Semesters[request.Semester]
	//if the semester doesnt already exist, we have to load it into the cache
	s.Cache.Lock.RUnlock()	//unlock the read lock
	if !ok{
		var err error
		semester, err = s.loadSemester(request.Semester)
		if err != nil{
			return nil, err
		}
	}

	//we get the faculty programs, the function will check that the faculty does exist
	facultyPrograms, err := s.getFacultyPrograms(request.FacultyName)
	if err != nil{
		return nil, err
	}
	//base struct for response
	response := new(domain.DTIResponseDTO)
	*response = domain.DTIResponseDTO{
		Semester: request.Semester,
		ErrorFound: false,
		ErrorMessage: "",
		Programs: make([]domain.DTIProgramResponseDTO, 0),
	}
	//now that we are sure that we have the semester info, we iterate over the faculty request programs
	for _, program := range request.Programs{
		programName := s.convertToBasicString(program.ProgramName)

		//if the program isnt a valid faculty program, we add it as a program error
		if programId, ok := facultyPrograms[programName]; !ok{
			response.Programs = append(response.Programs, domain.DTIProgramResponseDTO{
				ProgramId: program.ProgramId,
				Classrooms: 0,
				Labs: 0,
				MobileLabs: 0,
				StatusMessage: domain.InvalidProgramMsg,
			})
			continue
		}else{
			//we call the method in charge of process the program request, this is where we access the shared resource and all that stuff
			response.Programs = s.processProgramRequest(response.Programs, semester, &program, programId, goRoutineId)
		}
	}
	
	return response, nil
}


/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
//						INTERNAL METHODS								   //
/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

func (s *CollegeServiceImpl) processProgramRequest(programs []domain.DTIProgramResponseDTO, semester *SemesterAvailability, programRequest *domain.DTIProgramRequestDTO, programId uuid.UUID, goRoutineId int) []domain.DTIProgramResponseDTO{
	//check if we already have an assignation for this program for this semester, if we have, then we send the erroe
	if ass, _ := s.repository.GetProgramAssignment(programId, semester.Id); ass != nil{
		log.Printf("Program %s in semester %s already had assignation", programRequest.ProgramName, semester.Semester)
		programs = append(programs, domain.DTIProgramResponseDTO{
			ProgramId: programRequest.ProgramId,
			Classrooms: ass.Classrooms,
			Labs: ass.Labs,
			MobileLabs: ass.MobileLabs,
			StatusMessage: domain.OkMsg,
		})
		return programs
	}
	
	//create the response DTO struct
	programResponse := domain.DTIProgramResponseDTO{
		ProgramId: programRequest.ProgramId,
	}

	//LOCK THE SEMESTER RESOURCES FOR THE PROCESSING
	semester.resourcesLock.Lock()

	//with the logic below we can check all the assignation with or without mobile labs
	mobileLabsNeeded := programRequest.Labs - semester.Labs	//if with only the availala labs is enough, this will be 0 or negative
	if mobileLabsNeeded < 0{
		mobileLabsNeeded = 0
	}
	remainingClassrooms := semester.Classrooms - mobileLabsNeeded	//if we needed no mobile labs, then this will be equals to the available classrooms
	//if we have enough mobile labs (which implicitely checks if we have enough labs) and we have enough classrooms
	if mobileLabsNeeded <= semester.MobileLabs && remainingClassrooms >= programRequest.Classrooms{
		assignation := domain.AssignationModel{
			ID: uuid.New(),
			SemesterId: semester.Id,
			ProgramId: programId,
			CreatedAt: time.Now(),
			ProgramName: programRequest.ProgramName,
			SemesterName: semester.Semester,
			GoRoutineId: goRoutineId,
		}
		//add classroom asignation and update
		assignation.Classrooms = programRequest.Classrooms
		semester.Classrooms -= (programRequest.Classrooms + mobileLabsNeeded)
		assignation.RemainingCLassrooms = semester.Classrooms
		//add labs
		assignation.Labs = programRequest.Labs
		semester.Labs -= (programRequest.Labs)-mobileLabsNeeded
		assignation.RemainingLabs = semester.Labs
		//add mobile labs
		assignation.MobileLabs = mobileLabsNeeded
		semester.MobileLabs -= mobileLabsNeeded
		//this is to adjust, since if we assigned normal classrooms, its possible that now we have less 
		//classrooms than mobile labs allowed, if thats the case, we need to adjust the mobile labs to the available classrooms
		if semester.MobileLabs > semester.Classrooms{
			semester.MobileLabs = semester.Classrooms
		}
		assignation.RemainingMobileLabs = semester.MobileLabs
		//finally we send the assignation to be logged and written in the channel
		s.logWriterChannel <- &assignation

		//update response DTO
		programResponse.Classrooms = programRequest.Classrooms
		programResponse.Labs = programRequest.Labs
		programResponse.StatusMessage = domain.OkMsg
		programResponse.MobileLabs = mobileLabsNeeded
	}else{
		//this is the case in which we were unable to assign resources, in which case we must return an the error
		alert := domain.AlertModel{
			ID: uuid.New(),
			ProgramId: programId,
			SemesterId: semester.Id,
			Message: domain.NotEnoughResourcesMsg,
			CreatedAt: time.Now(),
			ProgramName: programRequest.ProgramName,
			SemesterName: semester.Semester,
			GoRoutineId: goRoutineId,
			RequestedClassrooms: programRequest.Classrooms,
			RequestedLabs: programRequest.Labs,
			AvailableClassrooms: semester.Classrooms,
			AvailableLabs: semester.Labs,
			AvailableMobileLabs: semester.MobileLabs,
		}
		s.alertWriterChannel <- &alert	//send in the alert channel for queue of alerts and logging

		//update response DTO
		programResponse.Classrooms = 0
		programResponse.Labs = 0
		programResponse.MobileLabs = 0
		programResponse.StatusMessage = domain.NotEnoughResourcesMsg
	}
	//UNLOCK THE LOCK, THE STRATEGY IS TO LOCK FOR EACH PROGRAM, SO WHILE WE DO THE ITERATION OTHER PROGRAM CAN BE PROCESSED
	semester.resourcesLock.Unlock() 

	//now we simply add the program response to the slice and return it, this part is not locking since we already liberated the lock
	programs = append(programs, programResponse)
	return programs
}



func (s *CollegeServiceImpl) loadSemester(semester string)(*SemesterAvailability, error){
	s.Cache.Lock.Lock()	//compltely lock, Write lock, so that we are able to write and avoid other
	//go gourintes which also notice that they had to load the semester to reload it after already loaded
	defer s.Cache.Lock.Unlock()	//defer so that its executed in any possible case
	//just in case, if perhaps we were ordered to load the semester, but that was while other go routine was already creating it
	if semester, ok := s.Cache.Semesters[semester]; ok{
		return semester, nil
	}
	//check if the semester exists in DB
	semesterModel, err := s.repository.GetSemester(semester)
	if err != nil{
		return nil, err
	}else if semesterModel == nil{	//if the semester wasnt found, then we have to create it
		semesterModel = &domain.SemesterAvailabilityModel{
			ID: uuid.New(),
			Semester: semester,
			Classrooms: s.config.Classrooms,
			Labs: s.config.Labs,
			MobileLabs: s.config.MobileLabs,
		}
		if err := s.repository.CreateSemester(semesterModel); err != nil{
			return nil, err
		}
	}

	//get the already assigned resources from this semester so we have an updated cache
	assignedResources, err := s.repository.GetAssignedResourcesOfSemester(semesterModel.ID)
	if err != nil{
		return nil, err
	}
	//creating the cache model and storing it into the cache
	semesterCache := &SemesterAvailability{
		Id: semesterModel.ID,
		resourcesLock: &sync.Mutex{},
		Semester: semester,
		Classrooms: semesterModel.Classrooms - assignedResources.Classrooms,
		Labs: semesterModel.Labs - assignedResources.Labs,
		MobileLabs: semesterModel.MobileLabs - assignedResources.MobileLabs,
	}
	s.Cache.Semesters[semester] = semesterCache
	return semesterCache, nil
}

//method that checks if a faculty exists, and if it those, then it returns the faculty programs
func (s *CollegeServiceImpl) getFacultyPrograms(facultyName string) (map[string]uuid.UUID, error){
	//we will get the information of the faculty from the request, so that we can valiate the programs
	faculties, err := s.repository.GetAllFaculties()
	if err != nil{
		return nil, err
	}
	var faculty *domain.FacultyModel = nil
	for _, fac := range faculties{
		if s.convertToBasicString(fac.Name) == s.convertToBasicString(facultyName){
			faculty = &fac 
			break;
		}
	}
	//if faculty is still nil, means that the faculty is invalid
	if faculty == nil{
		return nil, errors.New(domain.InvalidFacultyMsg)
	}
	//we now read the full faculty with its programs
	faculty, err = s.repository.GetFullFacultyById(faculty.ID)
	if err != nil{
		return nil, err
	}
	//we store the faculty programs names in a set like map, so that we can easily check if a request program is valid or not
	facultyPrograms := make(map[string]uuid.UUID)
	for _, program := range faculty.Programs{
		facultyPrograms[s.convertToBasicString(program.Name)] = program.ID
	}
	return facultyPrograms, nil
}


//internal method which will run in a separate go routine, is the writer
//to the db, is so that it isnt a blocking operation, is like a queue type of management
func (s *CollegeServiceImpl) startDbLogWriter()chan *domain.AssignationModel{
	channel := make(chan *domain.AssignationModel)
	go func(){
		for message := range channel{
			//log what we just did
			log.Printf("ASSIGNED BY GO ROUTINE: %v PROGRAM: %v SEMESTER: %v CLASSROOMS: %d LABS: %d MOBILE LABS: %d: REMAINING RESOURCES OF SEMESTER C:%d L:%d ML:%d",
				message.GoRoutineId,message.ProgramName,message.SemesterName,message.Classrooms,message.Labs,message.MobileLabs,message.RemainingCLassrooms,
				message.RemainingLabs,message.RemainingMobileLabs,
			)
			//save to db, blocking operation, however this is our purpose, to be the go routine blocked so the main ones are not
			s.repository.CreateAssignation(message)
		}
	}()
	return channel
}

//internal method which will run in a separate go routine, is the writer
//to the db only for alerts, is so that it isnt a blocking operation, is like a queue type of management
func (s *CollegeServiceImpl) startDbAlertWriter()chan *domain.AlertModel{
	channel := make(chan *domain.AlertModel)
	go func(){
		for message := range channel{
			//log what we just did
			log.Printf("XXXXX-ALERT BY GO ROUTINE: %v PROGRAM: %v SEMESTER: %v CLASSROOMS: %d LABS: %d AVILABLE RESOURCES OF SEMESTER C:%d L:%d ML:%d",
				message.GoRoutineId,message.ProgramName,message.SemesterName,message.RequestedClassrooms,message.RequestedLabs,message.AvailableClassrooms,
				message.AvailableLabs,message.AvailableMobileLabs,
			)
			//save to db, blocking operation, however this is our purpose, to be the go routine blocked so the main ones are not
			s.repository.CreateAlert(message)
		}
	}()
	return channel
}


func (s *CollegeServiceImpl) convertToBasicString(baseString string) string{
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToLower(baseString), " ", "",
			), "_", "",
		), "-", "",
	)
}

