package steps

import (
	"crypto"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"strings"
)

type SetupPGP struct {
	encryptionKeyRepository components.EncryptionKeyRepository
}

func (o *SetupPGP) Do(instanceId uuid.UUID) error {
	if key, _ := o.encryptionKeyRepository.FindByInstanceId(instanceId); key != nil {
		return nil
	}
	config := &packet.Config{
		DefaultHash: crypto.SHA256,
	}
	entity, err := openpgp.NewEntity("ScaleChamp", instanceId.String(), "info@scalechamp.com", config)
	if err != nil {
		return err
	}
	privateKey, err := o.newPrivateKey(entity, config)
	if err != nil {
		return err
	}
	publicKey, err := o.newPublicKey(entity)
	if err != nil {
		return err
	}
	encryptionKey := &models.EncryptionKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		InstanceId: instanceId,
	}
	return o.encryptionKeyRepository.Insert(encryptionKey)
}

func (o *SetupPGP) newPublicKey(e *openpgp.Entity) (string, error) {
	b := new(strings.Builder)
	w, err := armor.Encode(b, openpgp.PublicKeyType, nil)
	if err != nil {
		return "", err
	}
	if err := e.Serialize(w); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (o *SetupPGP) newPrivateKey(e *openpgp.Entity, config *packet.Config) (string, error) {
	b := new(strings.Builder)
	w, err := armor.Encode(b, openpgp.PrivateKeyType, nil)
	if err != nil {
		return "", err
	}
	if err := e.SerializePrivate(w, config); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (o *SetupPGP) Undo(uuid.UUID) error {
	return nil
}

func NewSetupPGP(
	accessKeyPairRepository components.EncryptionKeyRepository,
) *SetupPGP {
	return &SetupPGP{
		accessKeyPairRepository,
	}
}
