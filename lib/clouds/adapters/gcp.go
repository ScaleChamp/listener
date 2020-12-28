package adapters

import (
	"fmt"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"gitlab.com/scalablespace/listener/app/models"
	"google.golang.org/api/compute/v1"
	"net/http"
	"time"
)

type GCP struct {
	instancesService     *compute.InstancesService
	//networksService     *compute.NetworksService
	//subnetworksService     *compute.SubnetworksService
	zoneOperationService *compute.ZoneOperationsService
}

const projectId = "minewood"

func (dc *GCP) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	//n := compute.Subnetwork{
	//	IPv4Range:             "",
	//	Kind:                  "",
	//	Name:                  "",
	//	Subnetworks: []string,
	//}
	//dc.networksService.Insert(projectId, n)
	//dc.networksService.Insert(projectId, n)
	in := &compute.Instance{
		DeletionProtection: false,
		Description:        "instance for scalablespace",
		DisplayDevice: &compute.DisplayDevice{
			EnableDisplay: false,
		},
		Zone:         fmt.Sprintf("projects/%s/zones/%s", projectId, node.Region),
		MachineType:  fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", projectId, plan.Region, plan.GoogleCloudPlatform.MachineType), // https://compute.googleapis.com/compute/v1/
		Name:         fmt.Sprintf("instance-%s", node.Id),
		CanIpForward: false,
		Disks: []*compute.AttachedDisk{
			{
				InitializeParams: &compute.AttachedDiskInitializeParams{
					Description: "disk for node",
					DiskName:    fmt.Sprintf("instance-%.8s", node.Id),
					DiskSizeGb:  25,
					DiskType:    fmt.Sprintf("projects/%s/zones/%s/diskTypes/pd-standard", projectId, node.Region), // imagePathFromPlan europe-north1-a
					SourceImage: plan.Details.GoogleCloudPlatform.SourceImage,
				},
				Mode:       "READ_WRITE",
				AutoDelete: true,
				Boot:       true,
				Kind:       "compute#attachedDisk",
				Type:       "PERSISTENT",
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Subnetwork: plan.Details.GoogleCloudPlatform.Subnet,
				Kind:       "compute#networkInterface",
				Name:       "Extarnal Nat",
				AccessConfigs: []*compute.AccessConfig{
					{
						Kind:        "compute#accessConfig",
						Name:        "External NAT",
						Type:        "ONE_TO_ONE_NAT",
						NetworkTier: "PREMIUM",
					},
				},
			},
		},
		Scheduling: &compute.Scheduling{
			AutomaticRestart:  scw.BoolPtr(false),
			OnHostMaintenance: "MIGRATE",
		},
	}
	inst := dc.instancesService.Insert(projectId, node.Region, in)
	op, err := inst.Do()
	if err != nil {
		return nil, err
	}
operationLoop:
	for i := 0; i < 50; i += 1 {
		time.Sleep(5 * time.Second)
		op, err = dc.zoneOperationService.Get(projectId, node.Region, op.Name).Do()
		if err != nil {
			panic(err)
		}
		switch op.Status {
		case "RUNNING":
			continue
		case "DONE":
			break operationLoop
		default:
			return nil, http.ErrAbortHandler
		}
	}
	i, err := dc.instancesService.Get(projectId, node.Region, in.Name).Do()
	if err != nil {
		return nil, err
	}
	node.Metadata.IPv4 = i.NetworkInterfaces[0].AccessConfigs[0].NatIP
	return node, err
}

func (dc *GCP) DeleteNode(node *models.Node) error {
	op, err := dc.instancesService.Delete(projectId, node.Region, fmt.Sprintf("instance-%s", node.Id)).Do()
	if err != nil {
		return err
	}
	for i := 0; i < 100; i += 1 {
		time.Sleep(5 * time.Second)
		childOp, err := dc.zoneOperationService.Get(projectId, node.Region, op.Name).Do()
		if err != nil {
			continue
		}
		if childOp.Status == "DONE" {
			break
		}
	}
	return nil
}

func NewGCP(
	instancesService *compute.InstancesService,
	zoneOperationService *compute.ZoneOperationsService,
) *GCP {
	return &GCP{instancesService, zoneOperationService}
}
