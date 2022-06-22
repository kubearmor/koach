package repository

import (
	"fmt"

	"github.com/kubearmor/koach/koach/model"
	"gorm.io/gorm"
)

type INetworkCallRepository interface {
	Get(filter model.ObservabilityFilter) ([]model.Observability, error)
}

type networkCallRepository struct {
	db *gorm.DB
}

func (r *networkCallRepository) Get(filter model.ObservabilityFilter) ([]model.Observability, error) {
	observabilities := []model.Observability{}

	query := r.db.Table(model.Observability{}.TableName()).
		Select("*").
		Joins(
			fmt.Sprintf("LEFT JOIN %s on %s.id = %s.detail_id",
				model.NetworkCall{}.TableName(),
				model.NetworkCall{}.TableName(),
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

func NewNetworkCallRepository(db *gorm.DB) INetworkCallRepository {
	return &networkCallRepository{
		db: db,
	}
}
