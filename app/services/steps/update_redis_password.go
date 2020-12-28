package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type UpdateRedisPassword struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *UpdateRedisPassword) Do(instanceId uuid.UUID, nodeId uuid.UUID, previousPassword string) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := o.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	if err := o.executor.Run(node, "redis-cli -u redis://:%[1]s@localhost:6379 %[3]s-CONFIG set requirepass %[2]s", previousPassword, instance.Password, instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, "redis-cli -u redis://:%[1]s@localhost:6379 %[2]s-CONFIG set masterauth %[1]s", instance.Password, instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, "redis-cli -u redis://:%s@localhost:6379 %s-CONFIG rewrite", instance.Password, instance.Secret); err != nil {
		return err
	}
	return nil
}

func (*UpdateRedisPassword) Undo(_ uuid.UUID, _ uuid.UUID, _ string) error {
	return nil
}

func NewUpdateRedisPassword(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *UpdateRedisPassword {
	return &UpdateRedisPassword{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
