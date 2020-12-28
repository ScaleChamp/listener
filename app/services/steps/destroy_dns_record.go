package steps

import (
	"database/sql"
	"github.com/cloudflare/cloudflare-go"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
)

type DestroyDNSRecord struct {
	nodeRepository components.NodeRepository
	cloudflare     *cloudflare.API
	environment    models.Environment
}

func (o *DestroyDNSRecord) Do(nodeId uuid.UUID) error {
	node, err := o.nodeRepository.FindById(nodeId)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	if err := o.cloudflare.DeleteDNSRecord(o.environment.CloudflareZone, node.Metadata.RecordId); err != nil {
		return err
	}
	return nil
}

func (*DestroyDNSRecord) Undo(uuid.UUID) error {
	return nil
}

func NewDestroyDNSRecord(
	nodeRepository components.NodeRepository,
	cloudflare *cloudflare.API,
	environment models.Environment,
) *DestroyDNSRecord {
	return &DestroyDNSRecord{
		nodeRepository,
		cloudflare,
		environment,
	}
}
