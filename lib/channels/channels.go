package channels

import (
	"fmt"
)

const (
	InstanceCreated                 = "instance:created"
	InstanceCreatedError            = "instance:created:error"
	InstanceInitialized             = "instance:initialized"
	InstancePlanUpdate              = "instance:plan:update"
	InstancePlanUpdateError         = "instance:plan:update:error"
	InstanceDestroyed               = "instance:destroyed"
	InstanceDestroyedError          = "instance:destroyed:error"
	InstancePasswordUpdate          = "instance:password:update"
	InstanceCertificateUpdate       = "instance:certificate:update"
	InstanceCertificateUpdateError  = "instance:certificate:update:error"
	InstanceCertificateDisable      = "instance:certificate:disable"
	InstanceCertificateDisableError = "instance:certificate:disable:error"
)

type Instance struct {
	State string `json:"state"`
}

const InstanceRunning = "running"
const InstanceDisabled = "disabled"

func NewInstanceChannel(instanceId fmt.Stringer) string {
	return fmt.Sprintf("private-instance-%s", instanceId)
}
