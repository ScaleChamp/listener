package adapters

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"gitlab.com/scalablespace/listener/app/models"
	"golang.org/x/sync/errgroup"
)

type Azure struct {
	vmsClient  compute.VirtualMachinesClient
	nicClient  network.InterfacesClient
	pubClient  network.PublicIPAddressesClient
	diskClient compute.DisksClient
}

func (dc *Azure) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	publicIPFuture, err := dc.pubClient.CreateOrUpdate(context.TODO(), plan.Azure.ResourceGroup, fmt.Sprintf("ip-%.8s", node.Id), network.PublicIPAddress{
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: network.Dynamic,
		},
		Location: aws.String(plan.Region),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create public ip")
	}
	if err := publicIPFuture.WaitForCompletionRef(context.TODO(), dc.pubClient.Client); err != nil {
		return nil, errors.Wrap(err, "failed to wait for public ip")
	}
	ipAddr, err := publicIPFuture.Result(dc.pubClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result for public ip")
	}
	nicFuture, err := dc.nicClient.CreateOrUpdate(context.TODO(), plan.Azure.ResourceGroup, fmt.Sprintf("nic-%.8s", node.Id), network.Interface{
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				{
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: network.Dynamic,
						Subnet: &network.Subnet{
							ID: to.StringPtr(plan.Details.Azure.SubnetId),
						},
						Primary:         aws.Bool(true),
						PublicIPAddress: &ipAddr,
					},
					Name: aws.String("unique-name"),
				},
			},
			Primary: aws.Bool(true),
		},
		Location: aws.String(node.Region),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create nic")
	}
	if err := nicFuture.WaitForCompletionRef(context.TODO(), dc.nicClient.Client); err != nil {
		return nil, errors.Wrap(err, "failed to wait for nic")
	}
	nic, err := nicFuture.Result(dc.nicClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result nic")
	}
	vm := compute.VirtualMachine{
		Location: to.StringPtr(plan.Region),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			HardwareProfile: &compute.HardwareProfile{
				VMSize: compute.VirtualMachineSizeTypes(plan.Details.Azure.VMSize),
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: &compute.ImageReference{
					ID: aws.String(plan.Details.Azure.ImageId),
				},
				OsDisk: &compute.OSDisk{
					OsType:       compute.Linux,
					Name:         aws.String(fmt.Sprintf("disk-%.8s", node.Id)),
					CreateOption: compute.DiskCreateOptionTypesFromImage,
					DiskSizeGB:   aws.Int32(30),
					ManagedDisk: &compute.ManagedDiskParameters{
						StorageAccountType: compute.StandardLRS,
					},
				},
			},
			OsProfile: &compute.OSProfile{
				ComputerName:  to.StringPtr(fmt.Sprintf("node-%.8s", node.Id)),
				AdminUsername: to.StringPtr("scalechamp"),
				LinuxConfiguration: &compute.LinuxConfiguration{
					SSH: &compute.SSHConfiguration{
						PublicKeys: &[]compute.SSHPublicKey{
							{
								Path:    to.StringPtr("/home/scalechamp/.ssh/authorized_keys"),
								KeyData: to.StringPtr("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEXlNWESWZ+auYaq/HU7NKvfD57HFe1Sa8YLMGp7BOhtXudL2kxjYQly2iE8lU0EeVjCvuDlT3XMF2Ah+IgtnsT+ySK2V5kNd9hLzDnxwxRulFjoNc+nSrs2nNDDOC7lPnGkWyxu4yLSOIpZwgbwNe9Xe979re3Lg2uxTtg/Yiv89CAiTuf7UjmhHSUKZh/VMT/2PC6ypy0kaYUNE99om/WY5jiF9FxGRjX8BymJjugQ7b8FNOjIb8dphWanvc5JDVONxw+xhm0F3h8YNPGGZQb5XWhQHpZZLhAk87aQpe/AEYX9XFr5OQJpUzXwVQrV8IjXjSVCiCsIXOnjFtc0Dp mikefaraponov@Mikhails-iMac.local"),
							},
						},
					},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: nic.ID,
						NetworkInterfaceReferenceProperties: &compute.NetworkInterfaceReferenceProperties{
							Primary: to.BoolPtr(true),
						},
					},
				},
			},
		},
	}
	future, err := dc.vmsClient.CreateOrUpdate(context.TODO(), plan.Details.Azure.ResourceGroup, fmt.Sprintf("node-%.8s", node.Id), vm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to CreateOrUpdate")
	}
	if err := future.WaitForCompletionRef(context.TODO(), dc.vmsClient.Client); err != nil {
		return nil, errors.Wrap(err, "failed to WaitForCompletionRef vm")
	}
	vm, err = future.Result(dc.vmsClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Result vm")
	}
	ipAddr, err = dc.pubClient.Get(context.TODO(), plan.Azure.ResourceGroup, fmt.Sprintf("ip-%.8s", node.Id), "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get public ip")
	}
	node.Metadata.StringId = plan.Details.Azure.ResourceGroup
	node.Metadata.IPv4 = to.String(ipAddr.IPAddress)
	return node, nil

}

func (dc *Azure) DeleteNode(node *models.Node) error {
	future, err := dc.vmsClient.Delete(context.TODO(), node.Metadata.StringId, fmt.Sprintf("node-%.8s", node.Id))
	if err != nil {
		return errors.Wrap(err, "failed to delete vm")
	}
	if err := future.WaitForCompletionRef(context.TODO(), dc.vmsClient.Client); err != nil {
		return errors.Wrap(err, "failed to wait for vm")
	}
	diskFuture, err := dc.diskClient.Delete(context.TODO(), node.Metadata.StringId, fmt.Sprintf("disk-%.8s", node.Id))
	if err != nil {
		return errors.Wrap(err, "failed to delete disk")
	}
	nicFuture, err := dc.nicClient.Delete(context.TODO(), node.Metadata.StringId, fmt.Sprintf("nic-%.8s", node.Id))
	if err != nil {
		return errors.Wrap(err, "failed to delete nic")
	}

	if err := nicFuture.WaitForCompletionRef(context.TODO(), dc.nicClient.Client); err != nil {
		return errors.Wrap(err, "failed to wait for nic deletion")
	}
	ipFuture, err := dc.pubClient.Delete(context.TODO(), node.Metadata.StringId, fmt.Sprintf("ip-%.8s", node.Id))
	if err != nil {
		return err
	}
	var errGroup errgroup.Group
	errGroup.Go(func() error {
		return diskFuture.WaitForCompletionRef(context.TODO(), dc.diskClient.Client)
	})
	errGroup.Go(func() error {
		return ipFuture.WaitForCompletionRef(context.TODO(), dc.pubClient.Client)
	})
	return errGroup.Wait()
}

func NewAzure(
	vmsClient compute.VirtualMachinesClient,
	nicClient network.InterfacesClient,
	pubClient network.PublicIPAddressesClient,
	diskClient compute.DisksClient,
) *Azure {
	return &Azure{
		vmsClient,
		nicClient,
		pubClient,
		diskClient,
	}
}
