package steps

import (
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"net"
)

type SetupMonitoring struct {
	prometheusRepository components.PrometheusRepository
	nodeRepository       components.NodeRepository
	instanceRepository   components.InstanceRepository
	planRepository       components.PlanRepository
}

func (example *SetupMonitoring) Do(instanceId uuid.UUID, nodeId uuid.UUID) error {
	instance, err := example.instanceRepository.FindById(instanceId)
	if err != nil {
		return err
	}
	node, err := example.nodeRepository.FindById(nodeId)
	if err != nil {
		return err
	}
	labels := models.Labels{
		NodeId:     nodeId,
		InstanceId: instanceId,
		Kind:       instance.Kind,
		Job:        "node",
		Secret:     node.Metadata.PrometheusExporterPassword,
		Scheme:     "https",
	}
	target := &models.Prometheus{
		Targets: []string{
			net.JoinHostPort(node.Metadata.IPv4, "6666"),
		},
		Labels: labels,
		NodeId: nodeId,
	}
	if err := example.prometheusRepository.Insert(target); err != nil {
		return err
	}
	return nil
}

func (example *SetupMonitoring) Undo(_, _ uuid.UUID) error {
	return nil
}

func NewSetupMonitoring(
	nodeRepository components.NodeRepository,
	prometheusRepository components.PrometheusRepository,
	instanceRepository components.InstanceRepository,
	planRepository components.PlanRepository,
) *SetupMonitoring {
	return &SetupMonitoring{
		prometheusRepository,
		nodeRepository,
		instanceRepository,
		planRepository,
	}
}
