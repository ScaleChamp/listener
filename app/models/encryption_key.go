package models

import (
	uuid "github.com/satori/go.uuid"
)

type EncryptionKey struct {
	Id         uuid.UUID
	PublicKey  string
	PrivateKey string
	InstanceId uuid.UUID
}
