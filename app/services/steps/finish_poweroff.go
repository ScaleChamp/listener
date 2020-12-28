package steps

import (
	"github.com/pusher/pusher-http-go"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/scalablespace/listener/lib/channels"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/states"
	"log"
)

type FinishPoweroff struct {
	instanceRepository components.InstanceRepository
	client             *pusher.Client
}

func (r *FinishPoweroff) Do(instanceId uuid.UUID) error {
	if err := r.instanceRepository.UpdateState(instanceId, states.InstanceDisabled); err != nil {
		return err
	}
	i := &channels.Instance{
		State: channels.InstanceDisabled,
	}
	if err := r.client.Trigger(channels.NewInstanceChannel(instanceId), channels.InstanceInitialized, i); err != nil {
		log.Println(err)
		return nil
	}
	return nil
}

func (r *FinishPoweroff) Undo(uuid.UUID) error {
	return nil
}

func NewFinishPoweroff(
	instanceRepository components.InstanceRepository,
	client *pusher.Client,
) *FinishPoweroff {
	return &FinishPoweroff{instanceRepository, client}
}
