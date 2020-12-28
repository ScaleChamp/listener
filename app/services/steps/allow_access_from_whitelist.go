package steps

import (
	"bytes"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/app/templates"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/executor"
	"text/template"
)

type UpdateAccessFromInstanceWhitelist struct {
	executor           *executor.Executor
	nodeRepository     components.NodeRepository
	instanceRepository components.InstanceRepository
	planRepository     components.PlanRepository
	env                models.Environment
}

var mapKindToUfw = map[string]*template.Template{
	"redis":     templates.UfwRedisTemplate,
	"keydb":     templates.UfwRedisTemplate,
	"keydb-pro": templates.UfwRedisTemplate,
	"pg":        templates.UfwPgTemplate,
	"mysql":     templates.UfwMySQLTemplate,
}

func (o *UpdateAccessFromInstanceWhitelist) Do(instanceId uuid.UUID) error {
	instance, err := o.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	nodes, err := o.nodeRepository.FindByInstanceId(instanceId) // except orphaned or failed
	if err != nil {
		return err
	}
	for _, node := range nodes {
		b := new(bytes.Buffer)
		params := &templates.Firewall{
			Instance: instance,
			Node:     node,
		}
		if err := mapKindToUfw[instance.Kind].Execute(b, params); err != nil {
			return err
		}
		if err := o.executor.Wait(node); err != nil {
			return err
		}
		if err := o.executor.Put(node, b, "/tmp/user.rules"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo mv /tmp/user.rules /etc/ufw/user.rules"); err != nil {
			return err
		}
		if err := o.executor.Run(node, "sudo ufw reload"); err != nil {
			return err
		}
	}
	return nil
}

func (o *UpdateAccessFromInstanceWhitelist) Undo(uuid.UUID) error {
	return nil
}

func NewUpdateAccessFromInstanceWhitelist(
	executor *executor.Executor,
	nodeRepository components.NodeRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
	env models.Environment,
) *UpdateAccessFromInstanceWhitelist {
	return &UpdateAccessFromInstanceWhitelist{
		executor,
		nodeRepository,
		instanceRepository,
		planRepository,
		env,
	}
}
