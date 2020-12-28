package models

import (
	uuid "github.com/satori/go.uuid"
)

type AccessKeyPair struct {
	Id              uuid.UUID
	SecretAccessKey string
	AccessKeyId     string
	InstanceId      uuid.UUID
}
