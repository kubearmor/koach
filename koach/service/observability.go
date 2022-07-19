package service

import (
	"github.com/kubearmor/koach/koach/model"
	"github.com/kubearmor/koach/koach/repository"
)

type IObservabilityService interface {
	Get(filter model.ObservabilityFilter) ([]model.Observability, error)
	Save(observability model.Observability) error
	DeleteByAgeSeconds(age int) error
}

type observabilityService struct {
	observabilityRepository repository.IObservabilityRepository
}

func (s *observabilityService) Get(filter model.ObservabilityFilter) ([]model.Observability, error) {
	observabilities, err := s.observabilityRepository.Get(filter)
	if err != nil {
		return nil, err
	}

	if filter.Labels != "" {
		labelsFilter := model.LabelsFilter{}
		labelsFilter.FromString(filter.Labels)

		observabilities = model.FilterObservabilitiesByFilter(observabilities, labelsFilter)
	}

	return observabilities, nil
}

func (s *observabilityService) Save(observability model.Observability) error {
	_, err := s.observabilityRepository.Save(observability)

	return err
}

func (s *observabilityService) DeleteByAgeSeconds(age int) error {
	return s.observabilityRepository.DeleteByAgeSeconds(age)
}

func NewObservabilityService(
	observabilityRepository repository.IObservabilityRepository,
) IObservabilityService {
	return &observabilityService{
		observabilityRepository: observabilityRepository,
	}
}
