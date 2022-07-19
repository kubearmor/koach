package repository

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/kubearmor/koach/koach/model"
	"gorm.io/gorm"
)

type IObservabilityRepository interface {
	Get(filter model.ObservabilityFilter) ([]model.Observability, error)
	Save(observability model.Observability) (*string, error)
	DeleteByAgeSeconds(age int) error
}

type observabilityRepository struct {
	db *gorm.DB
}

func (r *observabilityRepository) Get(filter model.ObservabilityFilter) ([]model.Observability, error) {
	observabilities := []model.Observability{}

	query := r.db.Table(model.Observability{}.TableName()).
		Select("*").
		Where("deleted_at IS NULL").
		Where("operation = ?", string(filter.OperationType))

	if filter.NamespaceID != "" {
		query = query.Where("namespace_name = ?", filter.NamespaceID)
	}

	if filter.PodID != "" {
		query = query.Where("pod_name = ?", filter.PodID)
	}

	if filter.ContainerID != "" {
		query = query.Where("container_name = ?", filter.ContainerID)
	}

	if filter.SinceTimeSeconds != 0 {
		query = query.Where(fmt.Sprintf("created_at BETWEEN DATETIME('now', '-%s second', 'localtime') AND DATETIME('now', 'localtime')", strconv.Itoa(filter.SinceTimeSeconds)))
	}

	err := query.Scan(&observabilities).Error
	if err != nil {
		return nil, err
	}

	return observabilities, nil
}

func (r *observabilityRepository) Save(observability model.Observability) (*string, error) {
	observability.ID = uuid.New().String()

	err := r.db.Save(&observability).Error
	if err != nil {
		return nil, err
	}

	return &observability.ID, nil
}

func (r *observabilityRepository) DeleteByAgeSeconds(age int) error {
	return r.db.Table(model.Observability{}.TableName()).
		Where("deleted_at IS NULL").
		Delete(&model.Observability{}, fmt.Sprintf("created_at NOT BETWEEN DATETIME('now', '-%s second', 'localtime') AND DATETIME('now', 'localtime')", strconv.Itoa(age))).
		Error
}

func NewObservabilityRepository(db *gorm.DB) IObservabilityRepository {
	return &observabilityRepository{
		db: db,
	}
}
