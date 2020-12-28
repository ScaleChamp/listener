package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type UploadPublicPGP struct {
	executor                *executor.Executor
	encryptionKeyRepository components.EncryptionKeyRepository
	nodeRepository          components.NodeRepository
}

func (o *UploadPublicPGP) Do(instanceId, nodeId uuid.UUID) error {
	key, err := o.encryptionKeyRepository.FindByInstanceId(instanceId)
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
	if err := o.executor.PutString(node, key.PublicKey, "/tmp/public.pgp"); err != nil {
		return err
	}
	if err := o.executor.Run(node, "sudo mv /tmp/public.pgp /etc/public.pgp"); err != nil {
		return err
	}
	return nil
}

func (*UploadPublicPGP) Undo(uuid.UUID, uuid.UUID) error {
	return nil
}

func NewUploadPublicPGP(
	executor *executor.Executor,
	encryptionKeyRepository components.EncryptionKeyRepository,
	nodeRepository components.NodeRepository,
) *UploadPublicPGP {
	return &UploadPublicPGP{
		executor,
		encryptionKeyRepository,
		nodeRepository,
	}
}
