package steps

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type SetupKeydbNodeAsSlave struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
	planRepository     components.PlanRepository
}

func (o *SetupKeydbNodeAsSlave) Do(instanceId uuid.UUID, slaveId uuid.UUID, masterId uuid.UUID) error {
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
	if err := o.executor.Run(master, "keydb-cli -u redis://:%s@localhost:6379 %s-SLAVEOF no one", instance.Password, instance.Secret); err != nil {
		return err
	}
	plan, err := o.planRepository.FindById(instance.PlanId)
	if err != nil {
		return err
	}
	keydbSlaveConf := new(bytes.Buffer)
	if err := templates.KeyDBSlaveConfTemplate.Execute(keydbSlaveConf, &templates.KeyDBSlaveConf{
		ReplicaOf:       master.Metadata.IPv4,
		MasterAuth:      instance.Password,
		RequirePass:     instance.Password,
		Secret:          instance.Secret,
		MaxMemory:       plan.Memory - 1e+8,
		MaxMemoryPolicy: instance.EvictionPolicy,
		ServerThreads:   redisIOThreadsByPlanCpus(plan),
	}); err != nil {
		return err
	}
	if err := o.executor.Put(slave, keydbSlaveConf, "/tmp/keydb.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo mv /tmp/keydb.conf /etc/keydb/keydb.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo chown keydb:keydb /etc/keydb/keydb.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo systemctl start loggly"); err != nil {
		return err
	}
	if err := o.executor.Run(slave, "sudo systemctl start keydb-server"); err != nil {
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

func (o *SetupKeydbNodeAsSlave) Undo(_ uuid.UUID, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

func NewSetupKeydbNodeAsSlave(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
) *SetupKeydbNodeAsSlave {
	return &SetupKeydbNodeAsSlave{
		executor,
		nodeRepository,
		instanceRepository,
		planRepository,
	}
}
