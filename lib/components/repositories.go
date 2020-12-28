package components

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
)

type InstanceRepository interface {
	FindById(id uuid.UUID) (*models.Instance, error)
	UpdatePlan(id, planId uuid.UUID) error
	UpdateState(id uuid.UUID, to int) error
	GetState(id uuid.UUID) (int, error)
	UpdateStateTx(tx *sql.Tx, id uuid.UUID, to int) error
}

type CertificateAuthorityRepository interface {
	FindByProjectId(id uuid.UUID) (*models.CertificateAuthority, error)
}

type NodeRepository interface {
	FindFirstByInstanceId(id uuid.UUID) (*models.Node, error)
	FindByInstanceId(id uuid.UUID) ([]*models.Node, error)
	UpdateState(id uuid.UUID, to int) error
	GetState(id uuid.UUID) (int, error)
	Insert(s *models.Node) error
	UpdateMetadata(s *models.Node) error
	UpdateWhitelist(s *models.Node) error
	Delete(id uuid.UUID) error
	Orphan(id uuid.UUID) error
	FindById(id uuid.UUID) (*models.Node, error)
}

type PlanRepository interface {
	FindByIdTx(*sql.Tx, uuid.UUID) (*models.Plan, error)
	FindById(uuid.UUID) (*models.Plan, error)
}

type AccessKeyPairRepository interface {
	FindByInstanceId(uuid.UUID) (*models.AccessKeyPair, error)
}

type EncryptionKeyRepository interface {
	FindByInstanceId(uuid.UUID) (*models.EncryptionKey, error)
	Insert(keypair *models.EncryptionKey) error
}

type PrometheusRepository interface {
	Get() ([]*models.Prometheus, error)
	Size() (int, error)
	Insert(prometheus *models.Prometheus) error
	Delete(id uuid.UUID) error
}

type TaskRepository interface {
	Get(taskId uuid.UUID) (*models.Task, error) // call before task execution
	Update(*models.Task) error
	Finish(*models.Task) error
}

type UsageRepository interface {
	Insert(projectId, instanceId, planId uuid.UUID) error
	Upsert(projectId, instanceId, planId uuid.UUID) error
}
