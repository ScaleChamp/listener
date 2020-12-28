package adapters

import (
	"fmt"
	"github.com/exoscale/egoscale"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"gitlab.com/scalablespace/listener/app/models"
	"net/http"
)

type Exoscale struct {
	client *egoscale.Client
}

func (dc *Exoscale) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	request := &egoscale.DeployVirtualMachine{
		DisplayName:       fmt.Sprintf("%s-%s", plan.Kind, node.Id),
		Name:              fmt.Sprintf("%s-%s", plan.Kind, node.Id),
		Size:              plan.Exoscale.Size,
		ServiceOfferingID: egoscale.MustParseUUID(plan.Exoscale.ServiceOfferingID),
		TemplateID:        egoscale.MustParseUUID(plan.Exoscale.TemplateID),
		ZoneID:            egoscale.MustParseUUID(plan.Exoscale.ZoneID),
		KeyPair:           "main",
		IP4:               scw.BoolPtr(true),
		StartVM:           scw.BoolPtr(true),
	}
	response, err := dc.client.Request(request)
	if err != nil {
		return nil, err
	}
	switch vm := response.(type) {
	case *egoscale.VirtualMachine:
		node.Metadata.IPv4 = vm.IP().String()
		node.Metadata.StringId = vm.ID.String()
		return node, nil
	default:
		return nil, http.ErrNoLocation
	}
}

func (dc *Exoscale) DeleteNode(node *models.Node) error {
	request := &egoscale.DestroyVirtualMachine{
		ID: egoscale.MustParseUUID(node.Metadata.StringId),
	}
	_, err := dc.client.Request(request)
	switch err {
	case egoscale.ErrNotFound:
		return nil
	default:
		return err
	}
}

func NewExoscale(client *egoscale.Client) *Exoscale {
	return &Exoscale{client}
}
