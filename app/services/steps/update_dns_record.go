package steps

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type UpdateDNSRecord struct {
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
	cloudflare         *cloudflare.API
	environment        models.Environment
}

func (o *UpdateDNSRecord) Do(instanceId uuid.UUID, previousId uuid.UUID, nextId uuid.UUID) error {
	previous, err := o.nodeRepository.FindById(previousId)
	if err != nil {
		return err
	}
	next, err := o.nodeRepository.FindById(nextId)
	if err != nil {
		return err
	}
	err = o.cloudflare.UpdateDNSRecord(o.environment.CloudflareZone, previous.Metadata.RecordId, cloudflare.DNSRecord{
		Type:      "A",
		Name:      fmt.Sprintf("%s-%s.%s", previous.Metadata.Role, instanceId, o.environment.CloudflareDomain),
		Content:   next.Metadata.IPv4,
		Proxiable: false,
		Proxied:   false,
		TTL:       120,
		Priority:  10,
	})
	if err != nil {
		return err
	}
	next.Metadata.RecordId = previous.Metadata.RecordId
	if err := o.nodeRepository.UpdateMetadata(next); err != nil {
		return err
	}
	return nil
}

func (o *UpdateDNSRecord) Undo(_ uuid.UUID, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

func NewUpdateDNSRecord(
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	cloudflare *cloudflare.API,
	environment models.Environment,
) *UpdateDNSRecord {
	return &UpdateDNSRecord{nodeRepository, instanceRepository, cloudflare, environment}
}
