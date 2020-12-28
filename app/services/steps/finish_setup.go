package steps

import (
	"github.com/pusher/pusher-http-go"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/channels"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/states"
	"log"
)

type FinishSetup struct {
	instanceRepository components.InstanceRepository
	client             *pusher.Client
}

func (r *FinishSetup) Do(instanceId uuid.UUID) error {
	if err := r.instanceRepository.UpdateState(instanceId, states.InstanceRunning); err != nil {
		return err
	}
	i := &channels.Instance{
		State: channels.InstanceRunning,
	}
	if err := r.client.Trigger(channels.NewInstanceChannel(instanceId), channels.InstanceInitialized, i); err != nil {
		log.Println(err)
		return nil
	}
	return nil
}

func (r *FinishSetup) Undo(uuid.UUID) error {
	return nil
}

func NewFinishSetup(
	instanceRepository components.InstanceRepository,
	client *pusher.Client,
) *FinishSetup {
	return &FinishSetup{instanceRepository, client}
}
