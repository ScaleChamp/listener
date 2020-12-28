package adapters

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"gitlab.com/scalablespace/listener/app/models"
	"time"
)

type UpCloud struct {
	service *service.Service
}

func (dc *UpCloud) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	details, err := dc.service.CreateServer(&request.CreateServerRequest{
		Zone:             plan.Region,
		Plan:             plan.Upcloud.Plan,
		Title:            fmt.Sprintf("ScaleChamp server %s", node.Id),
		Hostname:         fmt.Sprintf("node-%.8s", node.Id),
		PasswordDelivery: request.PasswordDeliveryNone,
		StorageDevices: []upcloud.CreateServerStorageDevice{
			{
				Action:  upcloud.CreateServerStorageDeviceActionClone,
				Storage: plan.Upcloud.Storage,
				Title:   fmt.Sprintf("disk-%s", node.Id),
				Size:    plan.Upcloud.Size,
				Tier:    upcloud.StorageTierMaxIOPS,
			},
		},
		LoginUser: &request.LoginUser{
			CreatePassword: "no",
			Username:       "root",
			SSHKeys: []string{
				// does not work need to be sshkey from inside custom image or create new image and then install all the staff
				"replaced",
			},
		},
		IPAddresses: []request.CreateServerIPAddress{
			{
				Access: upcloud.IPAddressAccessPrivate,
				Family: upcloud.IPAddressFamilyIPv4,
			},
			{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv4,
			},
			{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv6,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	serverDetails, err := dc.service.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         details.Server.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range serverDetails.IPAddresses {
		if v.Access == upcloud.IPAddressAccessPublic && v.Family == upcloud.IPAddressFamilyIPv4 {
			node.Metadata.IPv4 = v.Address
		}
		if v.Access == upcloud.IPAddressAccessPublic && v.Family == upcloud.IPAddressFamilyIPv6 {
			node.Metadata.IPv6 = v.Address
		}
	}
	node.Metadata.StringId = serverDetails.Server.UUID
	return node, nil
}

func (dc *UpCloud) DeleteNode(node *models.Node) error {
	_, err := dc.service.StopServer(&request.StopServerRequest{
		UUID:     node.Metadata.StringId,
		StopType: request.ServerStopTypeSoft,
		Timeout:  3 * time.Minute,
	})
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	_, err = dc.service.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         node.Metadata.StringId,
		DesiredState: upcloud.ServerStateStopped,
		Timeout:      5 * time.Minute,
	})
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = dc.service.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{
		UUID: node.Metadata.StringId,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewUpcloud(s *service.Service) *UpCloud {
	return &UpCloud{s}
}
