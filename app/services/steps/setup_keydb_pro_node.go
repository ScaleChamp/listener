package steps

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type SetupKeyDBProNode struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	planRepository          components.PlanRepository
	accessKeyPairRepository components.AccessKeyPairRepository
}

func (o *SetupKeyDBProNode) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	keydbProConf := new(bytes.Buffer)
	if err := templates.KeyDBProConfTemplate.Execute(keydbProConf, &templates.KeyDBProConf{
		RequirePass:     instance.Password,
		Secret:          instance.Secret,
		MaxMemory:       plan.Memory - 1e+8,
		MaxMemoryPolicy: instance.EvictionPolicy,
		ServerThreads:   redisIOThreadsByPlanCpus(plan),
		EnablePro:       instance.LicenseKey,
	}); err != nil {
		return err
	}
	if err := o.executor.Put(node, keydbProConf, "/tmp/keydb.conf"); err != nil {
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
	conf := `DB=redis
REDIS_URL=redis://:%s@localhost:6379
`
	if err := o.executor.PutString(node, fmt.Sprintf(conf, instance.Password), "/tmp/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/ssnagent /etc/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start ssnagent"); err != nil {
		return err
	}
	redibuConf := new(bytes.Buffer)
	accessKeyPair, err := o.accessKeyPairRepository.FindByInstanceId(instance.Id)
	if err != nil {
		return err
	}
	if err := templates.RedibuConfTemplate.Execute(redibuConf, &templates.RedibuConf{
		Bucket:      scaleChampBucket,
		InstanceId:  instance.Id,
		AccessKeyId: accessKeyPair.AccessKeyId,
		SecretKey:   accessKeyPair.SecretAccessKey,
		AwsRegion:   "eu-central-1",
	}); err != nil {
		return err
	}
	if err := o.executor.Put(node, redibuConf, "/tmp/redibu"); err != nil {
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

func (o *SetupKeyDBProNode) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewSetupKeyDBProNode(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
) *SetupKeyDBProNode {
	return &SetupKeyDBProNode{
		executor,
		nodeRepository,
		instanceRepository,
		planRepository,
		accessKeyPairRepository,
	}
}
