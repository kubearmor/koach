package service

import (
	"errors"

	"github.com/kubearmor/koach/koach/model"
	"github.com/kubearmor/koach/koach/repository"
)

type IObservabilityService interface {
	Get(filter model.ObservabilityFilter) ([]model.Observability, error)
}

type observabilityService struct {
	fileAccessRepository   repository.IFileAccessRepository
	networkCallRepository  repository.INetworkCallRepository
	processSpawnRepository repository.IProcessSpawnRepository
}

func (s *observabilityService) Get(filter model.ObservabilityFilter) ([]model.Observability, error) {
	switch {
	case filter.OperationType == model.OperationFileAccess:
		return s.fileAccessRepository.Get(filter)
	case filter.OperationType == model.OperationNetworkCall:
		return s.networkCallRepository.Get(filter)
	case filter.OperationType == model.OperationProcessSpawn:
		return s.processSpawnRepository.Get(filter)
	}

	return nil, errors.New("unidentified operation type")
}

func NewObservabilityService(
	fileAccessRepository repository.IFileAccessRepository,
	networkCallRepository repository.INetworkCallRepository,
	processSpawnRepository repository.IProcessSpawnRepository,
) IObservabilityService {
	return &observabilityService{
		fileAccessRepository:   fileAccessRepository,
		networkCallRepository:  networkCallRepository,
		processSpawnRepository: processSpawnRepository,
	}
}
