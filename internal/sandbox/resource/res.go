package res

// ResourceManager will be created to limit and monitor the resources' usage of processes.
type ResourceManager interface {
	Apply(pid int) error
	Config(*ResourceConfig) error
	ReadState() (*RunState, error)
	Destroy() error
}

// ResourceConfig records the limits of resource.
type ResourceConfig struct {
	CpuLimit    uint32
	MemoryLimit uint32
	PidsLimit   uint32
}

// ResourceConfig records the resources' usage of processes put into ResourceManager.
type RunState struct {
	CpuTime     uint32
	MemoryUsage uint32
	OOMKill     uint32
}
