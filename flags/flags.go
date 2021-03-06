package flags

import "flag"

var (
	// ConfigPath specifies the path to the configuration file.
	ConfigPath string

	// PIDFilePath specifies the path to the PID file.
	PIDFilePath string

	// LogLevel specifies the log output level.
	LogLevel string

	// Version will cause the agent to show version information.
	Version bool

	// Info will display information about a running agent.
	Info bool

	// CPUProfile specifies the path to output CPU profiling information to.
	// When empty, CPU profiling is disabled.
	CPUProfile string

	// MemProfile specifies the path to output memory profiling information to.
	// When empty, memory profiling is disabled.
	MemProfile string
)

// Win holds a set of flags which will be populated only during the Windows build.
var Win = struct {
	InstallService   bool
	UninstallService bool
	StartService     bool
	StopService      bool
}{}

// AddFlags adds flags for the trace agent.
func AddFlags(flags *flag.FlagSet) {
	flags.StringVar(&ConfigPath, "config", DefaultConfigPath, "Datadog Agent config file location")
	flags.StringVar(&PIDFilePath, "pid", "", "Path to set pidfile for process")
	flags.BoolVar(&Version, "version", false, "Show version information and exit")
	flags.BoolVar(&Info, "info", false, "Show info about running trace agent process and exit")

	// profiling
	flags.StringVar(&CPUProfile, "cpuprofile", "", "Write cpu profile to file")
	flags.StringVar(&MemProfile, "memprofile", "", "Write memory profile to `file`")

	registerOSSpecificFlags(flags)
}
