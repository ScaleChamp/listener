package steps

import (
	"bytes"
	"fmt"
	"github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type SetupRedisNodeAsSlave struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	planRepository     components.PlanRepository
	instanceRepository components.InstanceRepository
}

func (o *SetupRedisNodeAsSlave) Do(instanceId uuid.UUID, slaveId uuid.UUID, masterId uuid.UUID) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	slave, err := o.nodeRepository.FindById(slaveId)
	if err != nil {
		return err
	}
	master, err := o.nodeRepository.FindById(masterId)
	if err != nil {
		return err
	}
	if err := o.executor.Wait(slave); err != nil {
		return err
	}
	plan, err := o.planRepository.FindById(instance.PlanId)
	if err != nil {
		return err
	}
	if err := o.executor.Run(master, "redis-cli -u redis://:%s@localhost:6379 %s-SLAVEOF no one", instance.Password, instance.Secret); err != nil {
		return err
	}
	redisSlaveConf := new(bytes.Buffer)
	if err := templates.RedisSlaveConfTemplate.Execute(redisSlaveConf, &templates.RedisSlaveConf{
		RequirePass:     instance.Password,
		Secret:          instance.Secret,
		MaxMemory:       plan.Memory - 1e+8,
		MaxMemoryPolicy: instance.EvictionPolicy,
		IOThreads:       redisIOThreadsByPlanCpus(plan),
		MasterAuth:      instance.Password,
		ReplicaOf:       master.Metadata.IPv4,
	}); err != nil {
		return err
	}
	if err := o.executor.Put(slave, redisSlaveConf, "/tmp/redis.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo mv /tmp/redis.conf /etc/redis/redis.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo chown redis:redis /etc/redis/redis.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo systemctl start loggly"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo systemctl start redis-server"); err != nil {
		return err
	}
	if err := o.executor.PutString(slave, fmt.Sprintf(templates.RedisAgentConf, instance.Password, slave.BcryptSecret()), "/tmp/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo mv /tmp/ssnagent /etc/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo systemctl start ssnagent"); err != nil {
		return err
	}
	slave.Metadata.Role = "slave"
	if err := o.nodeRepository.UpdateMetadata(slave); err != nil {
		return err
	}
	return nil
}

func (o *SetupRedisNodeAsSlave) Undo(_ uuid.UUID, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

func NewSetupNodeAsSlave(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
) *SetupRedisNodeAsSlave {
	return &SetupRedisNodeAsSlave{
		executor,
		nodeRepository,
		planRepository,
		instanceRepository,
	}
}
