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
	CpuLimit    uint
	MemoryLimit uint
	PidsLimit   uint
}

// ResourceConfig records the resources' usage of processes put into ResourceManager.
type RunState struct {
	CpuTime     uint
	MemoryUsage uint
	OOMKill     uint
}
