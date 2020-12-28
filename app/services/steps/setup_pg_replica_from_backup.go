package steps

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type SetupPgReplicaFromLatestBackup struct {
	executor                *executor.Executor
	nodeRepository          components.NodeRepository
	instanceRepository      components.InstanceRepository
	encryptionKeyRepository components.EncryptionKeyRepository
	planRepository          components.PlanRepository
}

func (o *SetupPgReplicaFromLatestBackup) Do(instanceId, primaryId, replicaId uuid.UUID) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := o.nodeRepository.FindById(replicaId)
	if err != nil {
		return err
	}
	primary, err := o.nodeRepository.FindById(primaryId)
	if err != nil {
		return err
	}
	plan, err := o.planRepository.FindById(instance.PlanId)
	if err != nil {
		return err
	}
	node.Metadata.Role = "slave"

	key, err := o.encryptionKeyRepository.FindByInstanceId(node.InstanceId)
	if err != nil {
		return err
	}
	if err := o.nodeRepository.UpdateMetadata(node); err != nil {
		return err
	}
	if err := o.executor.Wait(node); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl start loggly"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo systemctl stop postgresql"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo rm -rf /var/lib/postgresql/%s/main/", plan.Version); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mkdir -p /var/lib/postgresql/%s/main/", plan.Version); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo chown -R postgres:postgres /var/lib/postgresql/%s/main/", plan.Version); err != nil {
		return err
	}
	if err := o.executor.PutString(node, key.PrivateKey, "/tmp/private.pgp"); err != nil {
		return err
	}
	backupFetch := []string{
		"sudo mv /tmp/private.pgp /etc/private.pgp",
		"sudo /usr/local/bin/backup-fetch.sh",
	}
	if err := o.executor.MultiRun(node, backupFetch); err != nil {
		return err
	}

	const recoveryOptions = `primary_conninfo = 'user=replicator password=%s host=%s port=5432 sslmode=require sslcompression=0 gssencmode=prefer krbsrvname=postgres target_session_attrs=any'
primary_slot_name = 'slot1'
`
	const restoreBefore = `restore_command = '/usr/local/bin/wal-fetch.sh "%f" "%p"'
recovery_target_timeline = 'latest'
`

	switch plan.Version {
	case "10", "11":
		const oldRecoveryOptions = `standby_mode = 'on'
trigger_file = '/tmp/failover.trigger'
`
		if err := o.executor.PutString(node, oldRecoveryOptions+restoreBefore+fmt.Sprintf(recoveryOptions, instance.Secret, primary.Metadata.IPv4), "/tmp/recovery.conf"); err != nil {
			return err
		}
		if err := o.executor.Run(node, fmt.Sprintf("sudo cp /tmp/recovery.conf /var/lib/postgresql/%s/main/", plan.Version)); err != nil {
			return err
		}
	case "12":
		const newRecoveryOptions = `promote_trigger_file = '/tmp/failover.trigger'
`
		if err := o.executor.PutString(node, newRecoveryOptions+restoreBefore+fmt.Sprintf(recoveryOptions, instance.Secret, primary.Metadata.IPv4), "/tmp/recovery.conf"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo cp /etc/postgresql/12/main/postgresql.conf /tmp/postgresql.conf"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo cat /tmp/recovery.conf | sudo tee -a /etc/postgresql/12/main/postgresql.conf"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "touch /var/lib/postgresql/12/main/standby.signal"); err != nil {
			return err
		}
	default:
		panic("err: version undefined")
	}

	backupRestore := []string{
		fmt.Sprintf("sudo chown -R postgres:postgres /var/lib/postgresql/%s/main/", plan.Version),
		"sudo systemctl start postgresql",
	}
	if err := o.executor.MultiRun(node, backupRestore); err != nil {
		return err
	}

	switch plan.Version {
	case "12":
		removeRecovery := []string{
			"sudo mv /tmp/postgresql.conf /etc/postgresql/12/main/",
			"sudo chown -R postgres:postgres /etc/postgresql/12/main/postgresql.conf",
		}
		if err := o.executor.MultiRun(node, removeRecovery); err != nil {
			return err
		}
	default:
	}

	if err := o.executor.Run(node, "sudo rm -rf /etc/private.pgp"); err != nil {
		return err
	}
	if err := o.executor.PutString(node, fmt.Sprintf(templates.PgAgentConf, instance.Secret, node.BcryptSecret()), "/tmp/ssnagent"); err != nil {
		return err
	}
	cmds := []string{
		"sudo mv /tmp/ssnagent /etc/ssnagent",
		"sudo systemctl start ssnagent",
		"sudo systemctl start postgresql",
	}
	if err := o.executor.MultiRun(node, cmds); err != nil {
		return err
	}
	return nil
}

func (*SetupPgReplicaFromLatestBackup) Undo(uuid.UUID, uuid.UUID, uuid.UUID) error {
	return nil
}

func NewSetupPgReplicaFromLatestBackup(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	encryptionKeyRepository components.EncryptionKeyRepository,
	planRepository components.PlanRepository,
) *SetupPgReplicaFromLatestBackup {
	return &SetupPgReplicaFromLatestBackup{
		executor,
		nodeRepository,
		instanceRepository,
		encryptionKeyRepository,
		planRepository,
	}
}
