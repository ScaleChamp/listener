package adapters

import (
	"fmt"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"gitlab.com/scalablespace/listener/app/models"
	"time"
)

type SCW struct {
	instance *instance.API

	env models.Environment
}

func (dc *SCW) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	zone, err := scw.ParseZone(plan.Region)
	if err != nil {
		return nil, err
	}
	createRequest := &instance.CreateServerRequest{
		Organization:      dc.env.SCWOrganisation,
		Zone:              zone,
		Name:              fmt.Sprintf("instance-%s", node.Id),
		DynamicIPRequired: scw.BoolPtr(true),
		CommercialType:    plan.Scaleway.CommercialType,
		Image:             plan.Scaleway.Image,
		BootType:          instance.BootTypeLocal,
		Tags: []string{
			fmt.Sprintf("node-%s-instance-%", node.Id, node.InstanceId),
		},
	}
	createServerResponse, err := dc.instance.CreateServer(createRequest)
	if err != nil {
		return nil, err
	}
	if err := dc.instance.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		Zone:     zone,
		ServerID: createServerResponse.Server.ID,
		Action:   instance.ServerActionPoweron,
		Timeout:  5 * time.Minute,
	}); err != nil {
		return nil, err
	}
	node.Metadata.StringId = createServerResponse.Server.ID
	getServerResponse, err := dc.instance.GetServer(&instance.GetServerRequest{
		Zone:     zone,
		ServerID: createServerResponse.Server.ID,
	})
	if err != nil {
		return nil, err
	}
	node.Metadata.IPv4 = getServerResponse.Server.PublicIP.Address.String()
	for _, v := range getServerResponse.Server.Volumes {
		node.Metadata.ScalewaySnapshots = append(node.Metadata.ScalewaySnapshots, v.ID)
	}
	return node, nil
}

func (dc *SCW) DeleteNode(node *models.Node) error {
	zn, err := scw.ParseZone(node.Region)
	if err != nil {
		return err
	}
	serverActionAndWaitRequest := &instance.ServerActionAndWaitRequest{
		Zone:     zn,
		ServerID: node.Metadata.StringId,
		Action:   instance.ServerActionPoweroff,
		Timeout:  5 * time.Minute,
	}

	if err := dc.instance.ServerActionAndWait(serverActionAndWaitRequest); err != nil {
		return err
	}
	deleteServerRequest := &instance.DeleteServerRequest{
		Zone:     zn,
		ServerID: node.Metadata.StringId,
	}
	if err = dc.instance.DeleteServer(deleteServerRequest); err != nil {
		return err
	}
	for _, v := range node.Metadata.ScalewaySnapshots {
		deleteVolumeRequest := &instance.DeleteVolumeRequest{
			Zone:     zn,
			VolumeID: v,
		}
		if err := dc.instance.DeleteVolume(deleteVolumeRequest); err != nil {
			return nil
		}
	}
	return nil
}

func NewScaleway(client *scw.Client, env models.Environment) *SCW {
	return &SCW{instance.NewAPI(client), env}
}
