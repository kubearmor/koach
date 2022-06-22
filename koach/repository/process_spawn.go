package repository

import (
	"fmt"

	"github.com/kubearmor/koach/koach/model"
	"gorm.io/gorm"
)

type IProcessSpawnRepository interface {
	Get(filter model.ObservabilityFilter) ([]model.Observability, error)
}

type processSpawnRepository struct {
	db *gorm.DB
}

func (r *processSpawnRepository) Get(filter model.ObservabilityFilter) ([]model.Observability, error) {
	observabilities := []model.Observability{}

	query := r.db.Table(model.Observability{}.TableName()).
		Select("*").
		Joins(
			fmt.Sprintf("LEFT JOIN %s on %s.id = %s.detail_id",
				model.ProcessSpawn{}.TableName(),
				model.ProcessSpawn{}.TableName(),
				model.Observability{}.TableName(),
			),
		)

	if filter.NamespaceID != "" {
		query = query.Where(fmt.Sprintf("%s.namespace_id = ", model.Observability{}.TableName()), filter.NamespaceID)
	}

	if filter.DeploymentID != "" {
		query = query.Where(fmt.Sprintf("%s.deployment_id = ", model.Observability{}.TableName()), filter.DeploymentID)
	}

	if filter.NodeID != "" {
		query = query.Where(fmt.Sprintf("%s.node_id = ", model.Observability{}.TableName()), filter.NodeID)
	}

	if filter.PodID != "" {
		query = query.Where(fmt.Sprintf("%s.pod_id = ", model.Observability{}.TableName()), filter.PodID)
	}

	if filter.ContainerID != "" {
		query = query.Where(fmt.Sprintf("%s.container_id = ", model.Observability{}.TableName()), filter.ContainerID)
	}

	err := query.Scan(&observabilities).Error
	if err != nil {
		return nil, err
	}

	return observabilities, nil
}

func NewProcessSpawnRepository(db *gorm.DB) IProcessSpawnRepository {
	return &processSpawnRepository{
		db: db,
	}
}