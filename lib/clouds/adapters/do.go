package adapters

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"gitlab.com/scalablespace/listener/app/models"
	"net/http"
	"time"
)

type DO struct {
	client *godo.Client
}

func (dc *DO) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	createRequest := &godo.DropletCreateRequest{
		Name:   node.Id.String(),
		Region: plan.Region,
		Size:   plan.DigitalOcean.Size,
		Image: godo.DropletCreateImage{
			ID: plan.DigitalOcean.Id,
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{
				Fingerprint: plan.DigitalOcean.Fingerprint,
			},
		},
		Tags: []string{node.Id.String()},
	}
	newDroplet, _, err := dc.client.Droplets.Create(context.TODO(), createRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete droplet")
	}
	for i := 0; i < 100; i += 1 {
		time.Sleep(7 * time.Second)
		droplet, _, err := dc.client.Droplets.Get(context.TODO(), newDroplet.ID)
		if err != nil {
			continue
		}
		if droplet.Status != "active" {
			continue
		}
		node.Metadata.IPv4, err = droplet.PublicIPv4()
		if err != nil {
			return nil, errors.Wrap(err, "instance not have public ip")
		}
		node.Metadata.IPv6, _ = droplet.PublicIPv6()
		return node, nil
	}
	return nil, http.ErrHandlerTimeout
}

func (dc *DO) DeleteNode(node *models.Node) error {
	_, err := dc.client.Droplets.DeleteByTag(context.TODO(), node.Id.String())
	return err
}

func NewDO(client *godo.Client) *DO {
	return &DO{client}
}
