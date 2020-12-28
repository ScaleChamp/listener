package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"time"
)

type UpdatePgPassword struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *UpdatePgPassword) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := o.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	if err := o.executor.Wait(node); err != nil {
		return err
	}
	if err := o.executor.Run(node, `psql postgres://postgres:%s@localhost -c "ALTER USER ssnuser WITH PASSWORD '%s';"`, instance.Secret, instance.Password); err != nil {
		return err
	}
	if err := o.executor.Run(node, `psql postgres://postgres:%s@localhost -c "SELECT pg_switch_wal();"`, instance.Secret); err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	return nil
}

func (*UpdatePgPassword) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewUpdatePgPassword(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *UpdatePgPassword {
	return &UpdatePgPassword{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
