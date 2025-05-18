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


type CollegeServiceImpl struct {
	repository 	repository.CollegeRepository
	Lock 		*sync.RWMutex //lock for retrieving the semester info, is Read Write since it depends on if we are reading or writing
	Semester 	domain.SemesterAvailabilityModel
	logWriterChannel chan *domain.AssignationModel
	semesterCreated bool
}

func NewCollegeService(config *domain.ServiceConfig, repository repository.CollegeRepository) (domain.CollegeService, error) {
	//validate if semester already was processed
	sem, err := repository.GetSemester(config.Semester)
	if err != nil{
		return nil, domain.ErrorStartingService
	}
	if sem != nil{
		return nil, domain.ErrorSemesterWasAlreadyProcessed
	}

	//register semester in db
	semester := domain.SemesterAvailabilityModel{
		ID: uuid.New(),
		Semester: config.Semester,
		Classrooms: config.Classrooms,
		Labs: config.Labs,
		MobileLabs: config.MobileLabs,
	}
	//create service
	service := &CollegeServiceImpl{
		repository: repository,
		Semester: semester,
		Lock: &sync.RWMutex{},
		semesterCreated: false,
	}
	//start writers go routines and save the channel to communicate with them
	logWriterChannel := service.startDbLogWriter()
	service.logWriterChannel = logWriterChannel
	return service, nil
}

/*
	The entry method for processing a faculty request, is the main method which many go routines will execute 
	in parallel
*/
func (s *CollegeServiceImpl) ProcessRequest(request domain.DTIRequestDTO, goRoutineId int) (*domain.DTIResponseDTO, error) {
	//validate the request semester
	if request.Semester != s.Semester.Semester{
		return nil, domain.ErrorFacultyInvalidSemester
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

		//create base of program response
		programResponse := domain.DTIProgramResponseDTO{
			ProgramId: program.ProgramId,
			ProgramName: program.ProgramName,
			RequestedClassrooms: program.Classrooms,
			RequestedLabs: program.Labs,
		}

		//if the program isnt a valid faculty program, we add it as a program error
		if _, ok := facultyPrograms[programName]; !ok{
			programResponse.StatusMessage = domain.InvalidProgramMsg
			response.Programs = append(response.Programs, programResponse)
			continue
		}else{
			//we call the method in charge of process the program request, this is where we access the shared resource and all that stuff
			response.Programs = s.processProgramRequest(response.Programs, programResponse, goRoutineId)
		}
	}
	
	return response, nil
}


/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////
//						INTERNAL METHODS								   //
/////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////

func (s *CollegeServiceImpl) processProgramRequest(programs []domain.DTIProgramResponseDTO, programResponse domain.DTIProgramResponseDTO, goRoutineId int) []domain.DTIProgramResponseDTO{
	//check if we already have an assignation for this program for this semester, if we have, then we send the already assignacition
	if ass, _ := s.repository.GetProgramAssignment(programResponse.ProgramId, s.Semester.ID); ass != nil{
		log.Printf("Program %s in semester %s already had assignation", programResponse.ProgramName, s.Semester.Semester)
		programs = append(programs, domain.DTIProgramResponseDTO{
			ProgramId: programResponse.ProgramId,
			Classrooms: ass.Classrooms,
			Labs: ass.Labs,
			MobileLabs: ass.MobileLabs,
			StatusMessage: domain.OkMsg,
		})
		return programs
	}

	//LOCK THE SEMESTER RESOURCES FOR THE PROCESSING
	s.Lock.Lock()
	//if not yet created the semester then we create it
	if !s.semesterCreated{
		s.repository.CreateSemester(&s.Semester)
		s.semesterCreated = true
	}

	//with the logic below we can check all the assignation with or without mobile labs
	mobileLabsNeeded := programResponse.RequestedLabs - s.Semester.Labs	//if with only the availala labs is enough, this will be 0 or negative
	if mobileLabsNeeded < 0{
		mobileLabsNeeded = 0
	}
	remainingClassrooms := s.Semester.Classrooms - mobileLabsNeeded	//if we needed no mobile labs, then this will be equals to the available classrooms

	//by default assign the requested resources (if any was overloaded we will correct later)
	programResponse.Classrooms = programResponse.RequestedClassrooms
	programResponse.Labs = programResponse.RequestedLabs - mobileLabsNeeded
	programResponse.MobileLabs = mobileLabsNeeded

	//if either classrooms or labs are more than available, we generate the error and assign the available ones
	if mobileLabsNeeded > s.Semester.MobileLabs || programResponse.RequestedClassrooms > remainingClassrooms{
		programResponse.StatusMessage = domain.NotEnoughResourcesMsg
		//if requested classrooms were more than available
		if programResponse.RequestedClassrooms > remainingClassrooms{
			programResponse.Classrooms = remainingClassrooms
		}
		//if labs were more than available
		if mobileLabsNeeded > s.Semester.MobileLabs{
			programResponse.Labs = s.Semester.Labs
			programResponse.MobileLabs = s.Semester.MobileLabs
		}
	}
	s.Semester.Classrooms -= programResponse.Classrooms + mobileLabsNeeded //we have less classrooms, the classrooms reserved and the mobile labs used
	s.Semester.Labs -= programResponse.Labs
	s.Semester.MobileLabs -= programResponse.MobileLabs

	//create assignation model (for DB)
	assignation := domain.AssignationModel{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		ProgramId: programResponse.ProgramId,
		SemesterId: s.Semester.ID,
		RequestedClassrooms: programResponse.RequestedClassrooms,
		RequestedLabs: programResponse.RequestedLabs,
		Classrooms: programResponse.Classrooms,
		Labs: programResponse.Labs,
		MobileLabs: programResponse.MobileLabs,
		GoRoutineId: goRoutineId,
		ProgramName: programResponse.ProgramName,
		SemesterName: s.Semester.Semester,
		RemainingCLassrooms: s.Semester.Classrooms,
		RemainingLabs: s.Semester.Labs,
		RemainingMobileLabs: s.Semester.MobileLabs,
	}
	//free lock
	s.Lock.Unlock()

	//if we had less resources than needed, we send the aler
	if programResponse.StatusMessage == domain.NotEnoughResourcesMsg{
		assignation.Alert = true
	}
	s.logWriterChannel <- &assignation
	//now we simply add the program response to the slice and return it, this part is not locking since we already liberated the lock
	programs = append(programs, programResponse)
	return programs
}


/*
	method that checks if a faculty exists, and if it those, then it returns the faculty programs
*/
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


/*
	internal method which will run in a separate go routine, is the writer
	to the db, is so that it isnt a blocking operation, is like a queue type of management
*/
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



func (s *CollegeServiceImpl) convertToBasicString(baseString string) string{
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToLower(baseString), " ", "",
			), "_", "",
		), "-", "",
	)
}

