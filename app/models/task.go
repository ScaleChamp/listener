package models

import (
	"database/sql/driver"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

type Task struct {
	Id         uuid.UUID          `json:"id"`
	State      int                `json:"state"`
	Metadata   StringToRawMessage `json:"metadata"`
	Data       StringToStage      `json:"data"`
	Kind       string             `json:"kind"`
	Action     string             `json:"action"`
	InstanceId uuid.UUID          `json:"instance_id"`
}

type StringToRawMessage map[string]json.RawMessage

func (m StringToRawMessage) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m StringToRawMessage) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &m)
}

type StringToStage map[string]Stage

func (m StringToStage) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m StringToStage) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &m)
}

type Steps []*Step

type Step struct {
	Name     string
	Cmd      interface{}
	Wants    []string
	Provides []string
}

func NewStep(name string, cmd interface{}, wants []string, provides ...[]string) *Step {
	if len(provides) == 0 || len(provides) > 1 {
		return &Step{
			name, cmd, wants, nil,
		}
	}
	return &Step{
		name, cmd, wants, provides[0],
	}
}

type Message struct {
	Id uuid.UUID `json:"id"`
}

type Stage struct {
	Undone       bool `json:"undone"`
	UndoneFailed bool `json:"undone_failed"`
	Done         bool `json:"done"`
	Failed       bool `json:"failed"`
}
