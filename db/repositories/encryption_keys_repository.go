package repositories

import (
	"database/sql"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type encryptionKeyRepository struct {
	db *sql.DB
}

const findEncryptionKeyByInstanceId = `
SELECT id, private_key, public_key, instance_id FROM encryption_keys where instance_id = $1;
`

func (r *encryptionKeyRepository) FindByInstanceId(id uuid.UUID) (*models.EncryptionKey, error) {
	encryptionKey := new(models.EncryptionKey)
	err := r.db.QueryRow(findEncryptionKeyByInstanceId, id).Scan(
		&encryptionKey.Id,
		&encryptionKey.PrivateKey,
		&encryptionKey.PublicKey,
		&encryptionKey.InstanceId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find encryptionKeyRepository by instance id")
	}
	return encryptionKey, nil
}

const insertPgP = `
INSERT INTO encryption_keys (private_key, public_key, instance_id) VALUES ($1, $2, $3) RETURNING id;
`

func (r *encryptionKeyRepository) Insert(encryptionKey *models.EncryptionKey) error {
	return r.db.
		QueryRow(insertPgP, encryptionKey.PrivateKey, encryptionKey.PublicKey, encryptionKey.InstanceId).
		Scan(&encryptionKey.Id)
}

func NewEncryptionKeyRepository(db *sql.DB) components.EncryptionKeyRepository {
	return &encryptionKeyRepository{db}
}
