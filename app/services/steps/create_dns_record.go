package steps

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"net/http"
)

type CreateDNSRecord struct {
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
	cloudflare         *cloudflare.API
	environment        models.Environment
}

func (o *CreateDNSRecord) Do(nodeId, instanceId uuid.UUID) error {
	node, err := o.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	request := cloudflare.DNSRecord{
		Type:      "A",
		Name:      fmt.Sprintf("%s-%s.%s", node.Metadata.Role, instanceId, o.environment.CloudflareDomain),
		Content:   node.Metadata.IPv4,
		Proxiable: false,
		Proxied:   false,
		TTL:       120,
		Priority:  10,
	}
	response, err := o.cloudflare.CreateDNSRecord(o.environment.CloudflareZone, request)
	if err != nil {
		return err
	}
	if !response.Success {
		responses, err := o.cloudflare.DNSRecords(o.environment.CloudflareZone, request)
		if err != nil {
			return err
		}
		for _, r := range responses {
			node.Metadata.RecordId = r.ID
		}
		if node.Metadata.RecordId == "" {
			return http.ErrNoLocation
		}
	} else {
		node.Metadata.RecordId = response.Result.ID
	}
	return o.nodeRepository.UpdateMetadata(node)
}

func (o *CreateDNSRecord) Undo(_, _ uuid.UUID) error {
	return nil
}

func NewCreateDNSRecord(
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	cloudflare *cloudflare.API,
	environment models.Environment,
) *CreateDNSRecord {
	return &CreateDNSRecord{nodeRepository, instanceRepository, cloudflare, environment}
}
