package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/clouds"
	"gitlab.com/scalablespace/listener/lib/components"
	"golang.org/x/sync/errgroup"
	"log"
)

type CreateTwoNodes struct {
	nodeRepository     components.NodeRepository
	planRepository     components.PlanRepository
	instanceRepository components.InstanceRepository
	cloudAdapter       *clouds.CloudAdapter
}

const twoNodes = 2

func (r *CreateTwoNodes) Do(instanceId uuid.UUID, planId uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	plan, err := r.planRepository.FindById(planId)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	var wg errgroup.Group
	results := make(map[int]uuid.UUID, twoNodes)
	for i := 0; i < twoNodes; i += 1 {
		idx := i
		wg.Go(func() error {
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
				return err
			}
			if _, err = r.cloudAdapter.CreateNode(node, plan); err != nil {
				return err
			}
			if err := r.nodeRepository.UpdateMetadata(node); err != nil {
				return err
			}
			results[idx] = node.Id
			return nil
		})
	}
	// node id before insert + wait
	if err := wg.Wait(); err != nil {
		return results[0], results[1], err
	}
	return results[0], results[1], nil
}

func (r *CreateTwoNodes) Undo(nodeId1, nodeId2 uuid.UUID, _ uuid.UUID, _ uuid.UUID) error {
	nodeIds := []uuid.UUID{nodeId1, nodeId2}
	var wg errgroup.Group
	for _, v := range nodeIds {
		if v == uuid.Nil {
			continue
		}
		nodeId := v
		wg.Go(func() error {
			node, err := r.nodeRepository.FindById(nodeId)
			if err != nil {
				log.Println(err)
				return err
			}
			if err := r.nodeRepository.Delete(nodeId); err != nil {
				log.Println(err)
				return err
			}
			return r.cloudAdapter.DeleteNode(node)
		})
	}
	return wg.Wait()
}

func NewCreateTwoNodes(nodeRepository components.NodeRepository, ca *clouds.CloudAdapter, instanceRepository components.InstanceRepository, planRepo components.PlanRepository) *CreateTwoNodes {
	return &CreateTwoNodes{
		nodeRepository,
		planRepo,
		instanceRepository,
		ca,
	}
}
