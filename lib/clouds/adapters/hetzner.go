package adapters

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"gitlab.com/scalablespace/listener/app/models"
	"time"
)

type Hetzner struct {
	client *hcloud.Client
}

func (hz *Hetzner) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	result, _, err := hz.client.Server.Create(ctx, hcloud.ServerCreateOpts{
		//Volumes: []*hcloud.Volume{},
		//Automount: hcloud.Bool(false),
		Name: fmt.Sprintf("instance-%.8s", node.Id),
		ServerType: &hcloud.ServerType{
			Name: plan.HetznerCloud.ServerType,
		},
		Image: &hcloud.Image{
			ID: plan.HetznerCloud.ImageId,
		},
		SSHKeys: []*hcloud.SSHKey{
			{
				ID: plan.HetznerCloud.SSHKeyId,
			},
		},
		Location: &hcloud.Location{
			Name: node.Region,
		},
		StartAfterCreate: hcloud.Bool(true),
		Labels: map[string]string{
			"instance_id": node.InstanceId.String(),
			"node_id":     node.Id.String(),
		},
	})
	if err != nil {
		return nil, err
	}
	node.Metadata.IPv4 = result.Server.PublicNet.IPv4.IP.String()
	node.Metadata.IPv6 = result.Server.PublicNet.IPv6.IP.String()
	node.Metadata.ServerId = result.Server.ID
	return node, nil
}

func (hz *Hetzner) DeleteNode(node *models.Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := hz.client.Server.Delete(ctx, &hcloud.Server{
		ID: node.Metadata.ServerId,
	})
	return err
}

func NewHetzner(client *hcloud.Client) *Hetzner {
	return &Hetzner{client}
}
