package steps

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type PrepareMySQLNode struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	accessKeyPairRepository components.AccessKeyPairRepository
}

func (p *PrepareMySQLNode) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
	instance, err := p.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := p.nodeRepository.FindById(nodeId)
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
	walgParams := &templates.WalGConf{
		Bucket:       scaleChampBucket,
		Region:       "eu-central-1",
		Password:     instance.Secret,
		InstanceId:   instanceId,
		AwsKey:       accessKeyPair.AccessKeyId,
		AwsSecretKey: accessKeyPair.SecretAccessKey,
	}
	{
		backupFetch := new(bytes.Buffer)
		if err := templates.BackupMySQLFetchConfTemplate.Execute(backupFetch, walgParams); err != nil {
			return err
		}
		if err := p.executor.Put(node, backupFetch, "/tmp/backup-fetch.sh"); err != nil {
			return err
		}
		backupFetchCmds := []string{
			"sudo mv /tmp/backup-fetch.sh /usr/local/bin/backup-fetch.sh",
			"sudo chmod +x /usr/local/bin/backup-fetch.sh",
		}
		if err := p.executor.MultiRun(node, backupFetchCmds); err != nil {
			return err
		}
	}
	{
		backupPush := new(bytes.Buffer)
		if err := templates.BackupPushMySQLConfTemplate.Execute(backupPush, walgParams); err != nil {
			return err
		}
		if err := p.executor.Put(node, backupPush, "/tmp/backup-push.sh"); err != nil {
			return err
		}
		backupPushCmds := []string{
			"sudo mv /tmp/backup-push.sh /usr/local/bin/backup-push.sh",
			"sudo chmod +x /usr/local/bin/backup-push.sh",
		}
		if err := p.executor.MultiRun(node, backupPushCmds); err != nil {
			return err
		}
	}
	{
		binlogPush := new(bytes.Buffer)
		if err := templates.BinLogPushConfTemplate.Execute(binlogPush, walgParams); err != nil {
			return err
		}
		if err := p.executor.Put(node, binlogPush, "/tmp/binlog-push.sh"); err != nil {
			return err
		}
		binlogPushCmds := []string{
			"sudo mv /tmp/binlog-push.sh /usr/local/bin/binlog-push.sh",
			"sudo chmod +x /usr/local/bin/binlog-push.sh",
		}
		if err := p.executor.MultiRun(node, binlogPushCmds); err != nil {
			return err
		}
	}
	{
		binlogReplay := new(bytes.Buffer)
		if err := templates.BinLogReplayConfTemplate.Execute(binlogReplay, walgParams); err != nil {
			return err
		}
		if err := p.executor.Put(node, binlogReplay, "/tmp/binlog-replay.sh"); err != nil {
			return err
		}
		binlogReplayCmds := []string{
			"sudo mkdir -p /tmp/binlog",
			"sudo mv /tmp/binlog-replay.sh /usr/local/bin/binlog-replay.sh",
			"sudo chmod +x /usr/local/bin/binlog-replay.sh",
		}
		if err := p.executor.MultiRun(node, binlogReplayCmds); err != nil {
			return err
		}
	}
	return nil
}

func (p *PrepareMySQLNode) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewPrepareMySQLNode(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
) *PrepareMySQLNode {
	return &PrepareMySQLNode{
		executor,
		nodeRepository,
		instanceRepository,
		accessKeyPairRepository,
	}
}
