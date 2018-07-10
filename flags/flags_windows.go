// +build windows

package flags

import "flag"

const DefaultConfigPath = "c:\\programdata\\datadog\\datadog.yaml"

func registerOSSpecificFlags(flags *flag.FlagSet) {
	flags.BoolVar(&Win.InstallService, "install-service", false, "Install the trace agent to the Service Control Manager")
	flags.BoolVar(&Win.UninstallService, "uninstall-service", false, "Remove the trace agent from the Service Control Manager")
	flags.BoolVar(&Win.StartService, "start-service", false, "Starts the trace agent service")
	flags.BoolVar(&Win.StopService, "stop-service", false, "Stops the trace agent service")
}
