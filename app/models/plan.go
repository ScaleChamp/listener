package models

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

type Plan struct {
	Id      uuid.UUID
	Kind    string
	Name    string
	Cloud   string
	Region  string
	Version string
	Details
}

type Details struct {
	Rotational          bool                 `json:"rotational"`
	Nodes               int                  `json:"nodes"`
	VirtualCPUs         int                  `json:"vcpu"`
	Memory              uint64               `json:"memory"`
	AmazonWebServices   *AmazonWebServices   `json:"aws,omitempty"`
	GoogleCloudPlatform *GoogleCloudPlatform `json:"gcp,omitempty"`
	DigitalOcean        *DigitalOcean        `json:"do,omitempty"`
	HetznerCloud        *HetznerCloud        `json:"hetzner,omitempty"`
	Scaleway            *Scaleway            `json:"scw,omitempty"`
	Linode              *Linode              `json:"linode,omitempty"`
	Upcloud             *Upcloud             `json:"upcloud,omitempty"`
	Azure               *Azure               `json:"azure,omitempty"`
	Vultr               *Vultr               `json:"vultr,omitempty"`
	IBMCloud            *IBMCloud            `json:"ibm,omitempty"`
	Exoscale            *Exoscale            `json:"exoscale,omitempty"`
	AlibabaCloud        *AlibabaCloud        `json:"alibaba,omitempty"`
}

func (a *Details) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), a)
}

type AlibabaCloud struct {
	ImageId            string `json:"image_id"`
	InstanceType       string `json:"instance_type"`
	InternetChargeType string `json:"internet_charge_type"`
}

type IBMCloud struct {
	Memory int    `json:"memory"`
	Cpus   int    `json:"cpus"`
	OS     string `json:"os"`
	KeyId  int    `json:"key_id"`
}

type Exoscale struct {
	Size              int64  `json:"size"`
	ServiceOfferingID string `json:"service_offering_id"`
	TemplateID        string `json:"template_id"`
	ZoneID            string `json:"zone_id"`
}

type Vultr struct {
	RegionId   int    `json:"region_id"`
	PlanId     int    `json:"plan_id"`
	OsId       int    `json:"os_id"`
	SnapshotId string `json:"snapshot_id"`
}

type Azure struct {
	ResourceGroup string `json:"resource_group"`
	VMSize        string `json:"vm_size"`
	ImageId       string `json:"image_id"`
	SubnetId      string `json:"subnet_id"`
}

type Upcloud struct {
	Plan    string `json:"plan"`
	Storage string `json:"storage"`
	Size    int    `json:"size"`
}

type Linode struct {
	Type  string `json:"type"`
	Image string `json:"image"`
}

type Packet struct {
	Plan      string   `json:"plan"`
	OS        string   `json:"os"`
	ProjectID string   `json:"project_id"`
	Storage   string   `json:"storage"`
	SSHKey    []string `json:"ssh_key"`
}

type GoogleCloudPlatform struct {
	MachineType string `json:"machine_type"`
	SourceImage string `json:"source_image"`
	Zone        string `json:"zone"`
	Subnet      string `json:"subnet"`
}

type Scaleway struct {
	CommercialType string `json:"commercial_type"`
	Image          string `json:"image"`
}

type HetznerCloud struct {
	ImageId    int    `json:"image"`
	ServerType string `json:"server_type"`
	SSHKeyId   int    `json:"ssh_key_id"`
}

type AmazonWebServices struct {
	ImageId      string `json:"image"`
	InstanceType string `json:"instance_type"`
	SubnetId     string `json:"subnet_id"`
	Zone         string `json:"zone"`
	GroupId      string `json:"group_id"`
}

type DigitalOcean struct {
	Size        string `json:"size"`
	Id          int    `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

//type PgTune struct {
//	EffectiveCacheSize string
//	CheckpointCompletionTarget string
//	MaxConnections string
//	DefaultStatisticsTarget string
//	RandomPageCost string
//	SharedBuffers string
//	WalBuffers string
//	EffectiveIOConcurrency string
//	WorkMem string
//	MaintenanceWorkMem string
//	MaxWalSize string
//	MinWalSize string
//	MaxWorkerProcesses string
//	MaxParallelWorkersPerGather string
//	MaxParallelWorkers string
//}
