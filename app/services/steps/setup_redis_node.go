package steps

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type SetupRedisNode struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	planRepository          components.PlanRepository
	instanceRepository      components.InstanceRepository
	accessKeyPairRepository components.AccessKeyPairRepository
}

func (o *SetupRedisNode) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	plan, err := o.planRepository.FindById(instance.PlanId)
	if err != nil {
		return err
	}
	redisConf := new(bytes.Buffer)
	if err := templates.RedisConfTemplate.Execute(redisConf, &templates.RedisConf{
		RequirePass:     instance.Password,
		Secret:          instance.Secret,
		MaxMemory:       plan.Memory - 1e+8,
		MaxMemoryPolicy: instance.EvictionPolicy,
		IOThreads:       redisIOThreadsByPlanCpus(plan),
	}); err != nil {
		return err
	}
	if err := o.executor.Put(node, redisConf, "/tmp/redis.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/redis.conf /etc/redis/redis.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo chown redis:redis /etc/redis/redis.conf"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start loggly"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start redis-server"); err != nil {
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
	// in case user cloud profile, get his aws storage keys or any other compatible storage keys
	// on node creation use cloud profile keys as well
	accessKeyPair, err := o.accessKeyPairRepository.FindByInstanceId(instance.Id)
	if err != nil {
		return err
	}
	b := new(bytes.Buffer)
	if err := templates.RedibuConfTemplate.Execute(b, &templates.RedibuConf{
		Bucket:      scaleChampBucket,
		InstanceId:  instance.Id,
		AccessKeyId: accessKeyPair.AccessKeyId,
		SecretKey:   accessKeyPair.SecretAccessKey,
		AwsRegion:   "eu-central-1",
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
	node.Metadata.Role = "master"
	if err := o.nodeRepository.UpdateMetadata(node); err != nil {
		return err
	}
	return nil
}

func (o *SetupRedisNode) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func redisIOThreadsByPlanCpus(plan *models.Plan) int {
	switch {
	case plan.VirtualCPUs >= 4 && plan.VirtualCPUs < 8:
		return 2
	case plan.VirtualCPUs >= 8:
		return 6
	default:
		return 1
	}
}

func NewSetupRedisNode(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	planRepository components.PlanRepository,
	instanceRepository components.InstanceRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
) *SetupRedisNode {
	return &SetupRedisNode{
		executor,
		nodeRepository,
		planRepository,
		instanceRepository,
		accessKeyPairRepository,
	}
}
