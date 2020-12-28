package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type UpdateRedisCompatEviction struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *UpdateRedisCompatEviction) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	if err := o.executor.Run(node, "redis-cli -u redis://:%s@localhost:6379 %s-CONFIG SET maxmemory-policy %s", instance.Password, instance.Secret, instance.EvictionPolicy); err != nil {
		return err
	}
	return nil
}

func (*UpdateRedisCompatEviction) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewUpdateRedisCompatEviction(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *UpdateRedisCompatEviction {
	return &UpdateRedisCompatEviction{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
