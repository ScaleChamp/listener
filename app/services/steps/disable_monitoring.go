package steps

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
)

type DisableMonitoring struct {
	prometheusRepository components.PrometheusRepository
}

func (example *DisableMonitoring) Do(nodeId uuid.UUID) error {
	err := example.prometheusRepository.Delete(nodeId)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (example *DisableMonitoring) Undo(uuid.UUID) error {
	return nil
}

func NewDisableMonitoring(prometheusRepository components.PrometheusRepository) *DisableMonitoring {
	return &DisableMonitoring{prometheusRepository}
}
