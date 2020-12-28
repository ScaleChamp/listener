package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"log"
)

type DownloadLatestRedisBackup struct {
	instanceRepository      components.InstanceRepository
	nodeRepository          components.NodeRepository
	encryptionKeyRepository components.EncryptionKeyRepository
	accessKeyPairRepository components.AccessKeyPairRepository
	executor                *executor.Executor
}

func (r *DownloadLatestRedisBackup) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
	node, err := r.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	encryptionKey, err := r.encryptionKeyRepository.FindByInstanceId(node.InstanceId)
	if err != nil {
		return err
	}
	accessKeyPair, err := r.accessKeyPairRepository.FindByInstanceId(node.InstanceId)
	if err != nil {
		return err
	}
	if err := r.executor.PutString(node, encryptionKey.PrivateKey, "/tmp/private.pgp"); err != nil {
		return err
	}
	if err := r.executor.Run(node, "sudo mv /tmp/private.pgp /etc/private.pgp"); err != nil {
		return err
	}
	if err := r.executor.Run(
		node,
		"DUMP_PATH=%s PGP_PUBLIC_KEY_PATH=%s PGP_PRIVATE_KEY_PATH=%s S3_BUCKET=%s S3_PREFIX=%s AWS_ACCESS_KEY_ID=%s AWS_SECRET_KEY=%s AWS_REGION=%s /usr/local/bin/redibu backup-fetch",
		"/var/lib/redis/dump.rdb", "/etc/public.pgp", "/etc/private.pgp", scaleChampBucket, instanceId, accessKeyPair.AccessKeyId, accessKeyPair.SecretAccessKey, "eu-central-1"); err != nil {
		log.Println(err)
		return err
	}
	if err := r.executor.Run(node, "sudo rm -rf /etc/private.pgp"); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *DownloadLatestRedisBackup) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewDownloadLatestRedisBackup(
	instanceRepository components.InstanceRepository,
	nodeRepository components.NodeRepository,
	encryptionKeyRepository components.EncryptionKeyRepository,
	accessKeyPairRepository components.AccessKeyPairRepository,
	executor *executor.Executor,
) *DownloadLatestRedisBackup {
	return &DownloadLatestRedisBackup{
		instanceRepository,
		nodeRepository,
		encryptionKeyRepository,
		accessKeyPairRepository,
		executor,
	}
}
