package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type MigrateRedis struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *MigrateRedis) Do(instanceId uuid.UUID, sourceId, destinationId uuid.UUID, oldPassword string) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	source, err := o.nodeRepository.FindById(sourceId)
	if err != nil {
		return err
	}
	destination, err := o.nodeRepository.FindById(destinationId)
	if err != nil {
		return err
	}
	if err := o.executor.Wait(destination); err != nil {
		return err
	}
	const migrator = "REDIS_SOURCE=redis://:%[1]s@%[2]s:6379 REDIS_DESTINATION=redis://:%[3]s@%[4]s:6379 CONFIG=%[5]s-CONFIG SLAVEOF=%[5]s-SLAVEOF CLIENT=%[5]s-CLIENT /usr/local/bin/migrator"
	if err := o.executor.Run(destination, migrator, oldPassword, source.Metadata.IPv4, instance.Password, destination.Metadata.IPv4, instance.Secret); err != nil {
		return err
	}
	return nil
}

func (o *MigrateRedis) Undo(_ uuid.UUID, _, _ uuid.UUID, _ string) error {
	return nil
}

func NewMigrateRedis(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *MigrateRedis {
	return &MigrateRedis{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
