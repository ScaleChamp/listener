package states

const (
	Pending = iota
	Created
	Running
	Updating
	Down
	Terminating
	Terminated
	Failure
)

const (
	InstancePending = iota
	InstanceRunning
	InstanceTerminated
	InstanceUnhealthy
	InstanceMaintenance
	InstanceDisabled
)
