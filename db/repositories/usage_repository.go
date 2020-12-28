package repositories

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
)

type usageRepository struct {
	db *sql.DB
}

func (t *usageRepository) Insert(projectId, instanceId, planId uuid.UUID) error {
	_, err := t.db.Exec("INSERT INTO usages (project_id, instance_id, plan_id, started_at) VALUES($1, $2, $3, current_timestamp)", projectId, instanceId, planId)
	return err
}

func (t *usageRepository) Upsert(projectId, instanceId, planId uuid.UUID) error {
	tx, err := t.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`UPDATE usages SET finished_at = current_timestamp WHERE id = (SELECT id FROM usages WHERE instance_id = $1 and finished_at is null)`, instanceId); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`INSERT INTO usages (project_id, instance_id, plan_id, started_at) VALUES($1, $2, $3, current_timestamp)`, projectId, instanceId, planId); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func NewUsageRepository(db *sql.DB) components.UsageRepository {
	return &usageRepository{db: db}
}
