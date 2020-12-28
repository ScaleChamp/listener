package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"time"
)

type UpdateMySQLPassword struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *UpdateMySQLPassword) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	if err := o.executor.Run(node, `sudo mysql -uroot -e "START TRANSACTION; ALTER USER 'ssnuser' IDENTIFIED BY '%s'; FLUSH PRIVILEGES; COMMIT;"`, instance.Password); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo /usr/bin/flock /tmp/wal-g.lock /usr/local/bin/binlog-push.sh"); err != nil {
		return err
	}
	time.Sleep(30 * time.Second)
	return nil
}

func (*UpdateMySQLPassword) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewUpdateMySQLPassword(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *UpdateMySQLPassword {
	return &UpdateMySQLPassword{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
