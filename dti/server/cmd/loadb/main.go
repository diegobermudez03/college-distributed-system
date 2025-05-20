package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain/service"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/repository"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/transport"
	"github.com/joho/godotenv"
)

const (
	portArg            = "--port"
	classroomsArg      = "--classrooms"
	labsArg            = "--labs"
	mobileLabsArg      = "--mobile-labs"
	semesterArg        = "--semester"
	numberFacultiesArg = "--num-faculties"
	failSecondsArg     = "--fail-seconds"
	proxyArg           = "--proxy-server"
	nWorkersArg        = "--n-workers"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	//input structure
	// --port=<port-number> --classrooms=<number-of-available-classrooms> --labs=<number-of-available-labs> --mobile-labs=<max-number-of-mobile-labs> --semester=<semester>
	//executable --port=6000 --classrooms=360 --labs=200 --mobile-labs=50 --semester=2025-10

	//initial config with default values
	config := domain.ServiceConfig{
		Semester:   "",
		Classrooms: 380,
		Labs:       60,
		MobileLabs: 380, //this means how many of the classrooms could be used as labs, in this case the 100% of them
	}
	serverConfig := domain.ServerConfig{
		ListenPort:   6000,
		NumFaculties: 10,
		EndChannel:   make(chan bool),
		ProxyServer:  "",
	}
	failSeconds := 31536000
	nWorkers := 3
	//iterate over arguments and read them
	for _, arg := range os.Args {
		if !strings.Contains(arg, "=") {
			continue
		}
		parts := strings.Split(arg, "=")
		//all arguments are numeric, so we convert any argument to number, if it wasnt a nmber, ignore it
		number, err := strconv.Atoi(parts[1])
		if err != nil {
			switch parts[0] {
			case semesterArg:
				config.Semester = parts[1]
			case proxyArg:
				serverConfig.ProxyServer = parts[1]
			case nWorkersArg:
				nWorkers = number
			}
			continue
		}
		switch parts[0] {
		case portArg:
			serverConfig.ListenPort = number
		case classroomsArg:
			config.Classrooms = number
		case labsArg:
			config.Labs = number
		case mobileLabsArg:
			config.MobileLabs = number
		case numberFacultiesArg:
			serverConfig.NumFaculties = number
		case failSecondsArg:
			failSeconds = number
		}
	}

	//if no semester indicated, we cant initiate
	if config.Semester == "" {
		log.Fatal("Unable to start, Semester is required")
	}

	//check that number of mobile labs is at much the same as classrooms
	if config.MobileLabs > config.Classrooms {
		config.MobileLabs = config.Classrooms
	}

	//get env variables and read them (if they are not already loaded, we load them from the .env file)
	if os.Getenv("POSTGRES_HOST") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	/*dbConfig := repository.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
		SslMode:  os.Getenv("POSTGRES_SSL_MODE"),
		Timezone: os.Getenv("POSTGRES_TIMEZONE"),
	}
	//create, migrate and start db
	db, err := repository.OpenPostgresDb(dbConfig)
	if err != nil{
		log.Fatal("Error connecting to database: ", err.Error())
	}

	//inject dependencies
	collegeRepository := repository.NewCollegeRepositoryPostgres(db)*/
	collegeRepository := repository.NewDBMock()
	collegeService, err := service.NewCollegeService(&config, collegeRepository)
	if err != nil {
		log.Fatalf("Unable to start service: %s", err.Error())
	} else {
		log.Printf("Starting service with %d classrooms, %d labs and %d mobile labs", config.Classrooms, config.Labs, config.MobileLabs)
	}

	endChannel := make(chan bool)
	serverConfig.EndChannel = endChannel
	server := transport.NewLoadBServer(collegeService, serverConfig, nWorkers)
	//start server
	if err := server.Listen(); err != nil {
		log.Fatal("Unable to start server at port: ", serverConfig.ListenPort, " error: ", err.Error())
	}

	//wait for end signal
	select {
	case <-time.Tick(time.Duration(failSeconds) * time.Second):
		log.Printf("Fail simulated after %d seconds", failSeconds)
	case <-endChannel:
		log.Printf("Completed all %d faculties", serverConfig.NumFaculties)
	}
}
