package repositories

import (
	"database/sql"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type instanceRepository struct {
	db *sql.DB
}

const updateInstanceState = `UPDATE instances SET state = $2, updated_at = current_timestamp WHERE id = $1`

func (i *instanceRepository) UpdateState(id uuid.UUID, to int) error {
	_, err := i.db.Exec(updateInstanceState, id, to)
	return err
}

const getState = `SELECT state FROM instances WHERE id = $1`

func (i *instanceRepository) GetState(id uuid.UUID) (state int, err error) {
	return state, i.db.QueryRow(getState, id).Scan(&state)
}

const findById = `SELECT id, name, plan_id, state, plans.kind, password, whitelist, secret, eviction_policy, license_key, project_id FROM instances LEFT JOIN plans on plan_id = plans.id WHERE id = $1`

func (i *instanceRepository) FindById(id uuid.UUID) (*models.Instance, error) {
	instance := new(models.Instance)

	var nullString sql.NullString
	if err := i.db.QueryRow(findById, id).Scan(
		&instance.Id,
		&instance.Name,
		&instance.PlanId,
		&instance.State,
		&instance.Kind,
		&instance.Password,
		pq.Array(&instance.Whitelist),
		&instance.Secret,
		&instance.EvictionPolicy,
		&nullString,
		&instance.ProjectId); err != nil {
		return nil, err
	}
	if nullString.Valid {
		instance.LicenseKey = nullString.String
	}
	return instance, nil
}

func (i *instanceRepository) UpdateStateTx(tx *sql.Tx, id uuid.UUID, to int) error {
	_, err := tx.Exec(updateInstanceState, id, to)
	return err
}

const updateInstancePlan = `UPDATE instances SET state = 2, updated_at = current_timestamp WHERE id = $1`

func (i *instanceRepository) UpdatePlan(id, planId uuid.UUID) error {
	_, err := i.db.Exec(updateInstancePlan, id, planId)
	return err
}

func NewInstanceRepository(db *sql.DB) components.InstanceRepository {
	return &instanceRepository{db}
}
