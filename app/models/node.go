package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net"
)

type Node struct {
	Id         uuid.UUID
	Name       string
	State      int
	Cloud      string
	Region     string
	Whitelist  []string
	Metadata   *Metadata
	InstanceId uuid.UUID
}

func (n *Node) Names(domain string) []string {
	return []string{
		fmt.Sprintf("master-%s.%s", n.InstanceId, domain),
		fmt.Sprintf("primary-%s.%s", n.InstanceId, domain),
		fmt.Sprintf("slave-%s.%s", n.InstanceId, domain),
		fmt.Sprintf("replica-%s.%s", n.InstanceId, domain),
	}
}

func (n *Node) BcryptSecret() string {
	secret, err := bcrypt.GenerateFromPassword([]byte(n.Metadata.PrometheusExporterPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(secret)
}

type Metadata struct {
	Role                       string   `json:"role,omitempty"`
	RecordId                   string   `json:"record_id,omitempty"`
	IPv4                       string   `json:"ip4,omitempty"`
	IPv6                       string   `json:"ip6,omitempty"`
	ServerId                   int      `json:"id,omitempty"`
	StringId                   string   `json:"string_id,omitempty"`
	ScalewaySnapshots          []string `json:"scaleway_snapshots,omitempty"`
	PrometheusExporterPassword string   `json:"prometheus_exporter_password"`
}

func (m *Metadata) IPs() (ips []net.IP) {
	ips = append(ips, net.ParseIP(m.IPv4))
	if m.IPv6 != "" {
		if ip := net.ParseIP(m.IPv6); ip != nil {
			ips = append(ips, ip)
		}
	}
	return
}

func (m *Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Metadata) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), m)
}
