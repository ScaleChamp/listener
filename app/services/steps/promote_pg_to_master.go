package steps

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"time"
)

type PromotePgToMaster struct {
	nodeRepository     components.NodeRepository
	planRepository     components.PlanRepository
	instanceRepository components.InstanceRepository
	executor           *executor.Executor
}

func (p *PromotePgToMaster) Do(nodeId uuid.UUID) error {
	node, err := p.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	instance, err := p.instanceRepository.FindById(node.InstanceId)
	if err != nil {
		return err
	}
	plan, err := p.planRepository.FindById(instance.PlanId)
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

	switch plan.Version {
	case "10", "11":
		if err := p.executor.Run(node, fmt.Sprintf("sudo rm -rf /var/lib/postgresql/%s/main/recovery.conf", plan.Version)); err != nil {
			return err
		}
	case "12":
		if err := p.executor.Run(node, "sudo rm -rf /var/lib/postgresql/12/main/standby.signal"); err != nil {
			return err
		}
	default:
		panic("err: version undefined")
	}
	if err := p.executor.Run(node, `psql postgres://postgres:%s@localhost -c "SELECT pg_create_physical_replication_slot('slot1', true)"`, instance.Secret); err != nil {
		return err
	}
	if err := p.executor.Run(node, `sudo systemctl restart postgresql`); err != nil {
		return err
	}
	hour, min, _ := time.Now().UTC().Add(-1 * time.Minute).Clock()
	if err := p.executor.Run(node, `echo "%d %d * * * /usr/local/bin/backup-push.sh" | sudo crontab -`, min, hour); err != nil {
		return err
	}
	return nil
}

func (p *PromotePgToMaster) Undo(uuid.UUID) error {
	return nil
}

func NewPromotePgToMaster(
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
	executor *executor.Executor,
) *PromotePgToMaster {
	return &PromotePgToMaster{
		nodeRepository,
		planRepository,
		instanceRepository,
		executor,
	}
}
