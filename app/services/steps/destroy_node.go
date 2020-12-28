package steps

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/clouds"
	"gitlab.com/scalablespace/listener/lib/components"
	"log"
)

type DestroyNode struct {
	nodeRepository components.NodeRepository
	cloudAdapter   *clouds.CloudAdapter
}

func (r *DestroyNode) Do(id uuid.UUID) error {
	node, err := r.nodeRepository.FindById(id)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return nil
	}
	if err := r.cloudAdapter.DeleteNode(node); err != nil {
		log.Println(err)
		return r.nodeRepository.Orphan(id)
	}
	return r.nodeRepository.Delete(id)
}

func (r *DestroyNode) Undo(uuid.UUID) error {
	return nil
}

func NewDestroyNode(
	nodeRepository components.NodeRepository,
	cloudAdapter *clouds.CloudAdapter,
) *DestroyNode {
	return &DestroyNode{
		nodeRepository,
		cloudAdapter,
	}
}
