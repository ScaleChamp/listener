package adapters

import (
	"context"
	"fmt"
	"github.com/vultr/govultr"
	"gitlab.com/scalablespace/listener/app/models"
	"time"
)

type Vultr struct {
	service *govultr.Client
}

func (dc *Vultr) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	s, err := dc.service.Server.Create(context.TODO(), plan.Details.Vultr.RegionId, plan.Details.Vultr.PlanId, plan.Details.Vultr.OsId, &govultr.ServerOptions{
		//SSHKeyIDs:            []string{"5e57f06bde86b"},
		SnapshotID:           plan.Details.Vultr.SnapshotId,
		EnableIPV6:           false,
		EnablePrivateNetwork: false,
		Label:                fmt.Sprintf("node-%.8s", node.Id),
		Hostname:             fmt.Sprintf("node-%.8s", node.Id),
		Tag:                  fmt.Sprintf("instance-%.8s-node-%.8s", node.InstanceId, node.Id),
	})
	if err != nil {
		return nil, err
	}
	for i := 0; i < 10; i += 1 {
		time.Sleep(1 * time.Second)
		s, err = dc.service.Server.GetServer(context.TODO(), s.InstanceID)
		if err != nil {
			return nil, err
		}
		if s.MainIP == "" {
			continue
		}
	}
	node.Metadata.StringId = s.InstanceID
	node.Metadata.IPv4 = s.MainIP
	//node.StringToRawMessage.IPv6 = s.V6Networks[0].MainIP
	return node, nil
}

func (dc *Vultr) DeleteNode(node *models.Node) error {
	return dc.service.Server.Delete(context.TODO(), node.Metadata.StringId)
}

func NewVultr(client *govultr.Client) *Vultr {
	return &Vultr{client}
}
