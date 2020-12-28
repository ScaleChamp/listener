package steps

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"time"
)

type SetupMySQLNode struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *SetupMySQLNode) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	{
		cmds := []string{
			"sudo systemctl start loggly",
			"sudo systemctl start mysql",
		}
		if err := o.executor.MultiRun(node, cmds); err != nil {
			return err
		}
	}
	{
		if err := o.executor.PutString(node, fmt.Sprintf(templates.MySQLAgentConf, instance.Password, node.BcryptSecret()), "/tmp/ssnagent"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo mv /tmp/ssnagent /etc/ssnagent"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo systemctl start ssnagent"); err != nil {
			return err
		}
	}
	{
		if err := o.executor.Run(node, "sudo /usr/local/bin/backup-push.sh"); err != nil {
			return err
		}
	}
	{
		if err := o.executor.Run(node, `sudo mysql -uroot -e "START TRANSACTION; CREATE USER 'ssnuser' IDENTIFIED BY '%s'; CREATE DATABASE ssndb; GRANT ALL PRIVILEGES ON ssndb.* TO 'ssnuser'; FLUSH PRIVILEGES; COMMIT;"`, instance.Password); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo /usr/bin/flock /tmp/wal-g.lock /usr/local/bin/binlog-push.sh"); err != nil {
			return err
		}
	}
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

func (o *SetupMySQLNode) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewSetupMySQLNode(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *SetupMySQLNode {
	return &SetupMySQLNode{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
