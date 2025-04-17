package service

import (
	"errors"
	"strings"
	"sync"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/repository"
)

const (
	invalidProgramMsg = "INVALID-PROGRAM" 
	invalidFacultyMsg = "INVALID-FACULTY"
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
func (s *CollegeServiceImpl) ProcessRequest(request domain.DTIRequestDTO) (*domain.DTIResponseDTO, error) {
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
		if _, ok := facultyPrograms[programName]; !ok{
			response.Programs = append(response.Programs, domain.DTIProgramResponseDTO{
				ProgramId: program.ProgramId,
				Classrooms: 0,
				Labs: 0,
				StatusMessage: invalidProgramMsg,
			})
			continue
		}

		//if the program is valid, then we lock the semester labs resources for the assignment
		assignation := domain.AssignationModel{
			SemesterId: semester.Id,
			ProgramId: program.ProgramId,
			Classrooms: 0,
			Labs: 0,
			MobileLabs: 0,
		}
		//we get the lock for the resources
		semester.resourcesLock.Lock()
		//check if with only the labs and classrooms we fulfill
		if semester.Labs >= program.Labs && semester.Classrooms >= program.Classrooms{
			assignation.Classrooms = program.Classrooms
			assignation.Labs = program.Labs
			assignation.MobileLabs = 0
			semester.Labs -= program.Labs
			semester.Classrooms -= program.Classrooms
			semester.resourcesLock.Unlock()
			s.logWriterChannel <- &assignation
		}
		//I LEFT HERE, CONTINUE FROM HERE, I WAS AALZING HOW TO PROCESS THIS, 
	}
	
	return nil, nil
}


/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
//						INTERNAL METHODS								   //
/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

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
func (s *CollegeServiceImpl) getFacultyPrograms(facultyName string) (map[string]bool, error){
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
		return nil, errors.New(invalidFacultyMsg)
	}
	//we now read the full faculty with its programs
	faculty, err = s.repository.GetFullFacultyById(faculty.ID)
	if err != nil{
		return nil, err
	}
	//we store the faculty programs names in a set like map, so that we can easily check if a request program is valid or not
	facultyPrograms := make(map[string]bool)
	for _, program := range faculty.Programs{
		facultyPrograms[s.convertToBasicString(program.Name)] = true
	}
	return facultyPrograms, nil
}


//internal method which will run in a separate go routine, is the writer
//to the db, is so that it isnt a blocking operation, is like a queue type of management
func (s *CollegeServiceImpl) startDbLogWriter()chan *domain.AssignationModel{
	channel := make(chan *domain.AssignationModel)
	go func(){
		for message := range channel{

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

