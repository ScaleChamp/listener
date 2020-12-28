package steps

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
)

type AllowAccessFromMultipleIP struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
	planRepository     components.PlanRepository
	env                models.Environment
}

func (o *AllowAccessFromMultipleIP) Do(instanceId, destinationId uuid.UUID, sources ...uuid.UUID) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	destination, err := o.nodeRepository.FindById(destinationId)
	if err != nil {
		return err
	}
	ips := make([]string, len(sources))
	for i, src := range sources {
		node, err := o.nodeRepository.FindById(src)
		if err != nil {
			return err
		}
		ips[i] = node.Metadata.IPv4
	}

	destination.Whitelist = ips

	if err := o.nodeRepository.UpdateWhitelist(destination); err != nil {
		return err
	}

	userRules := new(bytes.Buffer)
	params := &templates.Firewall{
		Instance: instance,
		Node:     destination,
	}
	if err := mapKindToUfw[instance.Kind].Execute(userRules, params); err != nil {
		return err
	}
	if err := o.executor.Wait(destination); err != nil {
		return err
	}
	if err := o.executor.Put(destination, userRules, "/tmp/user.rules"); err != nil {
		return err
	}
	if err := o.executor.Run(destination, "sudo mv /tmp/user.rules /etc/ufw/user.rules"); err != nil {
		return err
	}
	if err := o.executor.Run(destination, "sudo ufw reload"); err != nil {
		return err
	}
	// update rules in node whitelist
	return nil
}

func (o *AllowAccessFromMultipleIP) Undo(_, _ uuid.UUID, _ ...uuid.UUID) error {
	return nil
}

func NewAllowAccessFromMultipleIP(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
	env models.Environment,
) *AllowAccessFromMultipleIP {
	return &AllowAccessFromMultipleIP{
		executor,
		nodeRepository,
		instanceRepository,
		planRepository,
		env,
	}
}
