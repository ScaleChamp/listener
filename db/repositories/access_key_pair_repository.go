package repositories

import (
	"database/sql"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type accessKeyPairRepository struct {
	db *sql.DB
}

const findAccessKeyByInstanceId = `
SELECT id, access_key_id, secret_access_key FROM access_key_pairs where instance_id = $1
`

func (r *accessKeyPairRepository) FindByInstanceId(id uuid.UUID) (*models.AccessKeyPair, error) {
	a := new(models.AccessKeyPair)
	err := r.db.QueryRow(findAccessKeyByInstanceId, id).Scan(
		&a.Id,
		&a.AccessKeyId,
		&a.SecretAccessKey,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find access key pair by instance id")
	}
	return a, nil
}

func NewAccessKeyPairRepository(db *sql.DB) components.AccessKeyPairRepository {
	return &accessKeyPairRepository{db}
}
