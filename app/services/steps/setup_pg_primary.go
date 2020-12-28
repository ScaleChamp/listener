package steps

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"time"
)

type SetupPgPrimary struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
}

func (o *SetupPgPrimary) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
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
	if err := o.executor.Run(node, "sudo systemctl start loggly"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start postgresql"); err != nil {
		return err
	}
	if err := o.executor.Run(node, `cd /tmp; sudo su - postgres -c "psql -c \"ALTER USER postgres PASSWORD '%s'\""`, instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, `psql postgres://postgres:%s@localhost -c "CREATE USER ssnuser WITH ENCRYPTED PASSWORD '%s';"`, instance.Secret, instance.Password); err != nil {
		return err
	}
	if err := o.executor.Run(node, `psql postgres://postgres:%s@localhost -c "CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD '%s';"`, instance.Secret, instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, "psql postgres://postgres:%s@localhost -c 'CREATE DATABASE ssndb;'", instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, "psql postgres://postgres:%s@localhost -c 'ALTER ROLE ssnuser SUPERUSER;'", instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, "psql postgres://ssnuser:%s@localhost/ssndb -c 'CREATE EXTENSION IF NOT EXISTS pgcrypto;'", instance.Password); err != nil {
		return err
	}
	if err := o.executor.Run(node, "psql postgres://postgres:%s@localhost -c 'ALTER ROLE ssnuser NOSUPERUSER;'", instance.Secret); err != nil {
		return err
	}
	if err := o.executor.Run(node, "psql postgres://postgres:%s@localhost -c 'GRANT ALL PRIVILEGES ON DATABASE ssndb TO ssnuser;'", instance.Secret); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	if err := o.executor.PutString(node, fmt.Sprintf(templates.PgAgentConf, instance.Secret, node.BcryptSecret()), "/tmp/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/ssnagent /etc/ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start ssnagent"); err != nil {
		return err
	}
	if err := o.executor.Run(node, `psql postgres://postgres:%s@localhost -c "SELECT pg_create_physical_replication_slot('slot1', true);"`, instance.Secret); err != nil {
		return err
	}
	// setup cron with pg_base
	hour, min, _ := time.Now().UTC().Add(-1 * time.Minute).Clock()
	if err := o.executor.Run(node, `echo "%d %d * * * /usr/local/bin/backup-push.sh" | sudo crontab -`, min, hour); err != nil {
		return err
	}
	return nil
}

func (*SetupPgPrimary) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewSetupPgPrimary(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
) *SetupPgPrimary {
	return &SetupPgPrimary{
		executor,
		nodeRepository,
		instanceRepository,
	}
}
