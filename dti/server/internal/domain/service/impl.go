package service

import (
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/repository"
)

type CollegeServiceImpl struct {
	config     *domain.ServiceConfig
	repository repository.CollegeRepository
}

func NewCollegeService(config *domain.ServiceConfig, repository repository.CollegeRepository) domain.CollegeService {
	return &CollegeServiceImpl{
		config: config,
		repository: repository,
	}
}

func (s *CollegeServiceImpl) ProcessRequest(request domain.DTIRequestDTO) (*domain.DTIResponseDTO, error) {
	return nil, nil
}

