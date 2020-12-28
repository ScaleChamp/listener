package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/clouds"
	"gitlab.com/scalablespace/listener/lib/components"
	"log"
)

type CreateNode struct {
	nodeRepository     components.NodeRepository
	planRepository     components.PlanRepository
	instanceRepository components.InstanceRepository
	cloudAdapter       *clouds.CloudAdapter
}

func (r *CreateNode) Do(instanceId uuid.UUID, planId uuid.UUID) (uuid.UUID, error) {
	plan, err := r.planRepository.FindById(planId)
	if err != nil {
		return uuid.Nil, err
	}
	node := &models.Node{
		State:      1,
		InstanceId: instanceId,
		Cloud:      plan.Cloud,
		Region:     plan.Region,
		Metadata: &models.Metadata{
			PrometheusExporterPassword: uuid.NewV4().String(),
		},
	}
	if err := r.nodeRepository.Insert(node); err != nil {
		return uuid.Nil, err
	}
	if _, err = r.cloudAdapter.CreateNode(node, plan); err != nil {
		return node.Id, err
	}
	if err := r.nodeRepository.UpdateMetadata(node); err != nil {
		return node.Id, err
	}
	return node.Id, nil
}

func (r *CreateNode) Undo(nodeId uuid.UUID, instanceId uuid.UUID, planId uuid.UUID) error {
	log.Println("trying to undo", nodeId.String())
	node, err := r.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	for i := 0; i < 2; i += 1 {
		if err := r.cloudAdapter.DeleteNode(node); err != nil {
			log.Println(err)
			continue
		}
		break
	}
	return r.nodeRepository.Delete(nodeId)
}

func NewCreateNode(nodeRepository components.NodeRepository, ca *clouds.CloudAdapter, instanceRepository components.InstanceRepository, planRepo components.PlanRepository) *CreateNode {
	return &CreateNode{
		nodeRepository,
		planRepo,
		instanceRepository,
		ca,
	}
}

/*
package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/clouds"
	"gitlab.com/scalablespace/listener/lib/components"
)

type CreateVM struct {
	nodeRepository     components.NodeRepository
	planRepository     components.PlanRepository
	instanceRepository components.InstanceRepository
	cloudAdapter       *clouds.CloudAdapter
}

func (r *CreateVM) Do(instanceId uuid.UUID, planId uuid.UUID) (uuid.UUID, error) {
	plan, err := r.planRepository.FindById(planId)
	if err != nil {
		return uuid.Nil, err
	}
	// node id before insert + wait
	node := &models.Node{
		State:      1,
		InstanceId: instanceId,
		Cloud:      plan.Cloud,
		Region:     plan.Region,
		Metadata:   new(models.Metadata),
	}
	if err := r.nodeRepository.Insert(node); err != nil {
		return uuid.Nil, err
	}
	return node.Id, nil
}

func (r *CreateVM) Undo(_, _, _ uuid.UUID) error {
	return nil
}

func NewCreateVM(nodeRepository components.NodeRepository, ca *clouds.CloudAdapter, instanceRepository components.InstanceRepository, planRepo components.PlanRepository) *CreateVM {
	return &CreateVM{
		nodeRepository,
		planRepo,
		instanceRepository,
		ca,
	}
}

*/
