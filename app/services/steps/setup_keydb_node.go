package steps

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type SetupKeyDBNode struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	planRepository          components.PlanRepository
	accessKeyPairRepository components.AccessKeyPairRepository
}

const scaleChampBucket = "scalechamp"

func (o *SetupKeyDBNode) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := o.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	node.Metadata.Role = "master"
	if err := o.nodeRepository.UpdateMetadata(node); err != nil {
		return err
	}
	if err := o.executor.Wait(node); err != nil {
		return err
	}
	plan, err := o.planRepository.FindById(instance.PlanId)
	if err != nil {
		return err
	}
	keydbConf := new(bytes.Buffer)
	if err := templates.KeyDBConfTemplate.Execute(keydbConf, &templates.KeyDBConf{
		RequirePass:     instance.Password,
		Secret:          instance.Secret,
		MaxMemory:       plan.Memory - 1e+8,
		MaxMemoryPolicy: instance.EvictionPolicy,
		ServerThreads:   redisIOThreadsByPlanCpus(plan),
	}); err != nil {
		return err
	}
	if err := o.executor.Put(node, keydbConf, "/tmp/keydb.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/keydb.conf /etc/keydb/keydb.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo chown keydb:keydb /etc/keydb/keydb.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start loggly"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start keydb-server"); err != nil {
		return err
	}
	if err := o.executor.PutString(node, fmt.Sprintf(templates.RedisAgentConf, instance.Password, node.BcryptSecret()), "/tmp/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/ssnagent /etc/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start ssnagent"); err != nil {
		return err
	}
	accessKeyPair, err := o.accessKeyPairRepository.FindByInstanceId(instance.Id)
	if err != nil {
		return err
	}
	b := new(bytes.Buffer)
	if err := templates.RedibuConfTemplate.Execute(b, &templates.RedibuConf{
		Bucket:      scaleChampBucket,
		AwsRegion:   "eu-central-1",
		AccessKeyId: accessKeyPair.AccessKeyId,
		SecretKey:   accessKeyPair.SecretAccessKey,
		InstanceId:  instance.Id,
	}); err != nil {
		return err
	}
	if err := o.executor.Put(node, b, "/tmp/redibu"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/redibu /etc/redibu"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start redibu"); err != nil {
		return err
	}
	return nil
}

func (o *SetupKeyDBNode) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewSetupKeyDBNode(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
) *SetupKeyDBNode {
	return &SetupKeyDBNode{
		executor,
		nodeRepository,
		instanceRepository,
		planRepository,
		accessKeyPairRepository,
	}
}
