package adapters

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/linode/linodego"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"gitlab.com/scalablespace/listener/app/models"
)

type Linode struct {
	client linodego.Client
}

func (hz *Linode) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	instance, err := hz.client.CreateInstance(context.TODO(), linodego.InstanceCreateOptions{
		Region:         plan.Region,
		Type:           plan.Linode.Type,
		Label:          fmt.Sprintf("node-%.8s", node.Id),
		RootPass:       "",
		AuthorizedKeys: []string{``},
		Image:          plan.Linode.Image,
		Tags:           []string{fmt.Sprintf("instance-%.8s-node-%.8s", node.InstanceId, node.Id)},
		SwapSize:       aws.Int(0),
		Booted:         scw.BoolPtr(true),
	})
	if err != nil {
		return nil, err
	}
	node.Metadata.IPv4 = instance.IPv4[0].String()
	node.Metadata.IPv6 = instance.IPv6
	node.Metadata.ServerId = instance.ID
	return node, nil
}

func (hz *Linode) DeleteNode(node *models.Node) error {
	return hz.client.DeleteInstance(context.TODO(), node.Metadata.ServerId)
}

func NewLinode(client linodego.Client) *Linode {
	return &Linode{client}
}
