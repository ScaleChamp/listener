package repositories

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type planRepository struct {
	db *sql.DB
}

const selectPlanById = `
SELECT id, kind, cloud, region, details, version FROM plans where id = $1
`

func (r *planRepository) FindById(id uuid.UUID) (*models.Plan, error) {
	plan := new(models.Plan)
	err := r.db.QueryRow(selectPlanById, id).Scan(
		&plan.Id,
		&plan.Kind,
		&plan.Cloud,
		&plan.Region,
		&plan.Details,
		&plan.Version,
	)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *planRepository) FindByIdTx(tx *sql.Tx, id uuid.UUID) (*models.Plan, error) {
	plan := new(models.Plan)
	err := tx.QueryRow(selectPlanById, id).Scan(
		&plan.Id,
		&plan.Kind,
		&plan.Cloud,
		&plan.Region,
		&plan.Details,
		&plan.Version,
	)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func NewPlanRepository(db *sql.DB) components.PlanRepository {
	return &planRepository{db}
}
