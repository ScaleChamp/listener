package models

import uuid "github.com/satori/go.uuid"

type CertificateAuthority struct {
	Id        uuid.UUID
	Key       string
	Crt       string
	ProjectId uuid.UUID
}
