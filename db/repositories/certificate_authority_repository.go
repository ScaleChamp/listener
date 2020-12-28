package repositories

import (
	"database/sql"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type certificateAuthorityRepository struct {
	db *sql.DB
}

const findCAbyProjectId = `
SELECT id, key, crt, project_id FROM certificate_authorities where project_id = $1
`

func (r *certificateAuthorityRepository) FindByProjectId(id uuid.UUID) (*models.CertificateAuthority, error) {
	ca := new(models.CertificateAuthority)
	err := r.db.QueryRow(findCAbyProjectId, id).Scan(
		&ca.Id,
		&ca.Key,
		&ca.Crt,
		&ca.ProjectId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find certificateAuthorityRepository by project id")
	}
	return ca, nil
}

func NewCertificateAuthorityRepository(db *sql.DB) components.CertificateAuthorityRepository {
	return &certificateAuthorityRepository{db}
}
