package steps

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type PrepagePgWalGWithoutWalPush struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	accessKeyPairRepository components.AccessKeyPairRepository
	planRepository          components.PlanRepository
}

func (p *PrepagePgWalGWithoutWalPush) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	if err := p.executor.Run(node, "sudo chmod +x /tmp/backup-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/backup-fetch.sh /usr/local/bin/backup-fetch.sh"); err != nil {
		return err
	}

	backupPush := new(bytes.Buffer)
	if err := templates.BackupPushConfTemplate.Execute(backupPush, conf); err != nil {
		return err
	}
	if err := p.executor.Put(node, backupPush, "/tmp/backup-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /tmp/backup-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/backup-push.sh /usr/local/bin/backup-push.sh"); err != nil {
		return err
	}

	const walPush = `#!/bin/bash
/bin/true
`
	if err := p.executor.PutString(node, walPush, "/tmp/wal-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /tmp/wal-push.sh && sudo chown postgres:postgres /tmp/wal-push.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/wal-push.sh /usr/local/bin/wal-push.sh"); err != nil {
		return err
	}

	walFetch := new(bytes.Buffer)
	if err := templates.WalFetchConfTemplate.Execute(walFetch, conf); err != nil {
		return err
	}
	if err := p.executor.Put(node, walFetch, "/tmp/wal-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo chmod +x /tmp/wal-fetch.sh && sudo chown postgres:postgres /tmp/wal-fetch.sh"); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo mv /tmp/wal-fetch.sh /usr/local/bin/wal-fetch.sh"); err != nil {
		return err
	}

	return nil
}

func (p *PrepagePgWalGWithoutWalPush) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewPreparePgWalGWithoutWalPush(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
	planRepository components.PlanRepository,
) *PrepagePgWalGWithoutWalPush {
	return &PrepagePgWalGWithoutWalPush{
		executor,
		nodeRepository,
		instanceRepository,
		accessKeyPairRepository,
		planRepository,
	}
}
