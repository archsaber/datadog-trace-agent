// +build !windows

package flags

import "flag"

// DefaultConfigPath specifies the default configuration file path for non-Windows systems.
const DefaultConfigPath = "/opt/datadog-agent/etc/datadog.yaml"

func registerOSSpecificFlags(flags *flag.FlagSet) {}
