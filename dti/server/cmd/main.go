package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/repository"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/service"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/transport"
)

const (
	portArg       = "--port"
	classroomsArg = "--classrooms"
	labsArg       = "--labs"
	mobileLabsArg = "--mobile-labs"
)

func main() {
	//input structure
	// --port=<port-number> --classrooms=<number-of-available-classrooms> --labs=<number-of-available-labs> --mobile-labs=<max-number-of-mobile-labs>
	//executable --port=6000 --classrooms=360 --labs=200 --mobile-labs=50

	//initial config with default values
	config := service.ServiceConfig{
		Classrooms: 380,
		Labs:       60,
		MobileLabs: 380,	//this means how many of the classrooms could be used as labs, in this case the 100% of them
	}
	listenPort := 6000 	//default listen port

	//iterate over arguments and read them
	for _, arg := range os.Args {
		if !strings.Contains(arg, "=") {
			continue
		}
		parts := strings.Split(arg, "=")
		//all arguments are numeric, so we convert any argument to number, if it wasnt a nmber, ignore it
		number, err := strconv.Atoi(parts[2])
		if err != nil{
			continue
		}
		switch parts[0]{
		case portArg: listenPort = number
		case classroomsArg: config.Classrooms = number
		case labsArg: config.Labs = number
		case mobileLabsArg: config.MobileLabs = number
		}
	}

	//check that number of mobile labs is at much the same as classrooms
	if config.MobileLabs > config.Classrooms{
		log.Fatal("Mobile labs cant be more than number of classrooms")
	}


	//inject dependencies
	collegeRepository := repository.NewCollegeRepositoryPostgres()
	collegeService := service.NewCollegeService(&config, collegeRepository)
	server := transport.NewZeroMQServer(listenPort, collegeService)

	//start server
	if err := server.Listen(); err != nil{
		log.Fatal("Unable to start server at port: ", listenPort, " error: ", err.Error())
	}
}