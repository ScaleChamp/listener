package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type PromoteToMasterRedisCompat struct {
	nodeRepository components.NodeRepository
	executor       *executor.Executor
}

func (p *PromoteToMasterRedisCompat) Do(nodeId uuid.UUID) error {
	node, err := p.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	node.Metadata.Role = "master"
	if err := p.nodeRepository.UpdateMetadata(node); err != nil {
		return err
	}
	if err := p.executor.Wait(node); err != nil {
		return err
	}
	if err := p.executor.Run(node, "sudo systemctl start redibu"); err != nil {
		return err
	}
	//setup slave of no one
	return nil
}

func (p *PromoteToMasterRedisCompat) Undo(uuid.UUID) error {
	return nil
}

func NewPromoteToMasterRedisCompat(
	nodeRepository components.NodeRepository,
	executor *executor.Executor,
) *PromoteToMasterRedisCompat {
	return &PromoteToMasterRedisCompat{
		nodeRepository,
		executor,
	}
}
