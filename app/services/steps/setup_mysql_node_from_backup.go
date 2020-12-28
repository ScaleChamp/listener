package steps

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"time"
)

type SetupMySQLNodeFromBackup struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	encryptionKeyRepository components.EncryptionKeyRepository
}

func (o *SetupMySQLNodeFromBackup) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	key, err := o.encryptionKeyRepository.FindByInstanceId(node.InstanceId)
	if err != nil {
		return err
	}
	{
		if err := o.executor.PutString(node, key.PrivateKey, "/tmp/private.pgp"); err != nil {
			return err
		}
		cmds := []string{
			"sudo mv /tmp/private.pgp /etc/private.pgp",
			"sudo rm -rf /var/lib/mysql/",
			"sudo mkdir -p /var/lib/mysql/",
			"sudo chown -R mysql:mysql /var/lib/mysql/",
			"sudo /usr/local/bin/backup-fetch.sh",
			"sudo chown -R mysql:mysql /var/lib/mysql/",
			"sudo systemctl start loggly",
			"sudo systemctl start mysql",
			"sudo /usr/local/bin/binlog-replay.sh",
			"sudo rm -rf /etc/private.pgp",
		}
		if err := o.executor.MultiRun(node, cmds); err != nil {
			return err
		}
	}
	{
		if err := o.executor.PutString(node, fmt.Sprintf(templates.MySQLAgentConf, instance.Password, node.BcryptSecret()), "/tmp/ssnagent"); err != nil {
			return err
		}
		cmds := []string{
			"sudo mv /tmp/ssnagent /etc/ssnagent",
			"sudo systemctl start ssnagent",
		}
		if err := o.executor.MultiRun(node, cmds); err != nil {
			return err
		}
	}
	// push current new backup
	{
		hour, min, _ := time.Now().UTC().Add(-1 * time.Minute).Clock()
		crontab := `*/5 * * * * /usr/bin/flock -n /tmp/wal-g.lock /usr/local/bin/binlog-push.sh
%d %d * * * /usr/bin/flock /tmp/wal-g.lock /usr/local/bin/backup-push.sh
`
		if err := o.executor.PutString(node, fmt.Sprintf(crontab, min, hour), "/tmp/crontab"); err != nil {
			return err
		}
		if err := o.executor.Run(node, `cat /tmp/crontab | sudo crontab -`); err != nil {
			return err
		}
		if err := o.executor.Run(node, `rm /tmp/crontab`); err != nil {
			return err
		}
	}
	return nil
}

func (o *SetupMySQLNodeFromBackup) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewSetupMySQLNodeFromBackup(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	encryptionKeyRepository components.EncryptionKeyRepository,
) *SetupMySQLNodeFromBackup {
	return &SetupMySQLNodeFromBackup{
		executor,
		nodeRepository,
		instanceRepository,
		encryptionKeyRepository,
	}
}
