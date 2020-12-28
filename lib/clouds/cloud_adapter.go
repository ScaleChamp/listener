package clouds

import (
	"github.com/pkg/errors"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/clouds/adapters"
)

type nodeCreatorDeleter interface {
	CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error)
	DeleteNode(node *models.Node) error
}

type CloudAdapter struct {
	adapters map[string]nodeCreatorDeleter
}

func (dc *CloudAdapter) CreateNode(node *models.Node, plan *models.Plan) (*models.Node, error) {
	newNode, err := dc.adapters[node.Cloud].CreateNode(node, plan)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create node for instance %s, node %s, plan %s, cloud %s", node.InstanceId, node.Id, plan.Id, node.Cloud)
	}
	return newNode, nil
}

func (dc *CloudAdapter) DeleteNode(node *models.Node) error {
	if err := dc.adapters[node.Cloud].DeleteNode(node); err != nil {
		return errors.Wrapf(err, "failed to delete node for instance %s, node %s, cloud %s", node.InstanceId, node.Id, node.Cloud)
	}
	return nil
}

func NewCloudAdapter(
	do *adapters.DO,
	hetzner *adapters.Hetzner,
	scw *adapters.SCW,
	aws *adapters.AWS,
	linode *adapters.Linode,
	upcloud *adapters.UpCloud,
	gcp *adapters.GCP,
	azure *adapters.Azure,
	vultr *adapters.Vultr,
	ibm *adapters.IBM,
	alibaba *adapters.Alibaba,
	oracle *adapters.Oracle,
	exoscale *adapters.Exoscale,
	tencent *adapters.TencentCloud,
	selectel *adapters.Selectel,
) *CloudAdapter {
	return &CloudAdapter{
		adapters: map[string]nodeCreatorDeleter{
			"hetzner":  hetzner,
			"scw":      scw,
			"aws":      aws,
			"gcp":      gcp,
			"linode":   linode,
			"azure":    azure,
			"upcloud":  upcloud,
			"vultr":    vultr,
			"do":       do,
			"ibm":      ibm,
			"exoscale": exoscale,
			"alibaba":  alibaba,
			"oracle":   oracle,
			"tencent":  tencent,
			"selectel": selectel,
		},
	}
}
