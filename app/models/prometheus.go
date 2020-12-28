package models

import (
	"database/sql/driver"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

type Prometheus struct {
	Id      uuid.UUID `json:"-"`
	Labels  Labels    `json:"labels"`
	Targets []string  `json:"targets"`
	NodeId  uuid.UUID `json:"-"`
}

type Labels struct {
	NodeId     uuid.UUID `json:"node_id"`
	InstanceId uuid.UUID `json:"instance_id"`
	Kind       string    `json:"kind"`
	Job        string    `json:"job"`
	Secret     string    `json:"__param_secret,omitempty"`
	Scheme     string    `json:"__scheme__,omitempty"`
}

func (l *Labels) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), l)
}

func (l *Labels) Value() (driver.Value, error) {
	return json.Marshal(l)
}
