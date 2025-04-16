package service

import "github.com/diegobermudez03/college-distributed-system/dti/server/internal/repository"

type CollegeServiceImpl struct {
	config     *ServiceConfig
	repository repository.CollegeRepository
}

func NewCollegeService(config *ServiceConfig, repository repository.CollegeRepository) CollegeService {
	return &CollegeServiceImpl{
		config: config,
		repository: repository,
	}
}

func (s *CollegeServiceImpl) ProcessRequest(request DTIRequestDTO) (*DTIResponseDTO, error) {
	return nil, nil
}