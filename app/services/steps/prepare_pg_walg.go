package steps

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type PreparePgWalG struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	accessKeyPairRepository components.AccessKeyPairRepository
	planRepository          components.PlanRepository
}

func (p *PreparePgWalG) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
	instance, err := p.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := p.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	plan, err := p.planRepository.FindById(instance.PlanId)
	if err != nil {
		return err
	}
	if err := p.executor.Wait(node); err != nil {
		return err
	}
	accessKeyPair, err := p.accessKeyPairRepository.FindByInstanceId(instance.Id)
	if err != nil {
		return err
	}
	conf := &templates.WalGConf{
		Version:      plan.Version,
		Bucket:       scaleChampBucket,
		Region:       "eu-central-1",
		Password:     instance.Secret,
		InstanceId:   instanceId,
		AwsKey:       accessKeyPair.AccessKeyId,
		AwsSecretKey: accessKeyPair.SecretAccessKey,
	}
	backupFetch := new(bytes.Buffer)
	if err := templates.BackupFetchConfTemplate.Execute(backupFetch, conf); err != nil {
		return err
	}
	if err := p.executor.Put(node, backupFetch, "/tmp/backup-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/backup-fetch.sh /usr/local/bin/backup-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /usr/local/bin/backup-fetch.sh"); err != nil {
		return err
	}

	backupPush := new(bytes.Buffer)
	if err := templates.BackupPushConfTemplate.Execute(backupPush, conf); err != nil {
		return err
	}
	if err := p.executor.Put(node, backupPush, "/tmp/backup-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/backup-push.sh /usr/local/bin/backup-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /usr/local/bin/backup-push.sh"); err != nil {
		return err
	}
	walPush := new(bytes.Buffer)
	if err := templates.WalPushConfTemplate.Execute(walPush, conf); err != nil {
		return err
	}
	if err := p.executor.Put(node, walPush, "/tmp/wal-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/wal-push.sh /usr/local/bin/wal-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /usr/local/bin/wal-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chown postgres:postgres /usr/local/bin/wal-push.sh"); err != nil {
		return err
	}

	walFetch := new(bytes.Buffer)
	if err := templates.WalFetchConfTemplate.Execute(walFetch, conf); err != nil {
		return err
	}
	if err := p.executor.Put(node, walFetch, "/tmp/wal-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/wal-fetch.sh /usr/local/bin/wal-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /usr/local/bin/wal-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chown postgres:postgres /usr/local/bin/wal-fetch.sh"); err != nil {
		return err
	}
	return nil
}

func (p *PreparePgWalG) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewPreparePgNode(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
	planRepository components.PlanRepository,
) *PreparePgWalG {
	return &PreparePgWalG{
		executor,
		nodeRepository,
		instanceRepository,
		accessKeyPairRepository,
		planRepository,
	}
}
