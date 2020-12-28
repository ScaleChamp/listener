package models

import (
	"database/sql/driver"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

type Event struct {
	ID            uuid.UUID
	Data          JSONB
	Metadata      JSONB
	EventableType string
	EventableID   uuid.UUID
}

type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	if len(a) == 0 {
		return json.Marshal(struct{}{})
	}
	return json.Marshal(a)
}
