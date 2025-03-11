package main

import (
	"log"
	"os"
	"time"

	"github.com/Nigel23932/go-svc/src/installer"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	serviceConfig = &mgr.Config{
		// 1 = SERVICE_SID_TYPE_UNRESTRICTED
		SidType: 1,

		// Start automatically
		StartType: mgr.StartAutomatic,

		// // LocalService account
		ServiceStartName: "NT AUTHORITY\\SYSTEM",

		// Start automatically
		DelayedAutoStart: false,

		DisplayName: "go-svc",
		Description: "My Service is a template for creating Windows services in Go.",
	}
)

func main() {
	cobra.MousetrapHelpText = ""

	defer func() {
		if r := recover(); r != nil {
			var f, err = os.OpenFile("C:\\Users\\PvPBe\\Desktop\\go\\service\\error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to create error log: %v\n", err)
			}
			defer f.Close()

			log.SetOutput(f)
			log.Fatalf("Panic: %v\n", r)
		}
	}()

	var (
		debug    = new(bool)
		instance = &Instance{
			ServiceName: serviceConfig.DisplayName,
			Opts:        serviceConfig,
			Accepted: svc.AcceptStop |
				svc.AcceptShutdown |
				svc.AcceptPauseAndContinue,
		}
		serviceInstallerConfig = &installer.ServiceInstallerConfig{
			Service:         instance,
			EventsSupported: eventlog.Error | eventlog.Warning | eventlog.Info,
			Args:            []string{"run"},
		}
		commandStatus = cobra.Command{
			Use:     "status",
			Short:   "Display the installation status of the service.",
			Aliases: []string{"s"},
			Args:    cobra.NoArgs,
			Run:     InstallationStatus(serviceInstallerConfig),
		}
		commandInstall = cobra.Command{
			Use:     "install",
			Aliases: []string{"i"},
			Short:   "Install the service.",
			Args:    cobra.NoArgs,
			Run:     InstallService(serviceInstallerConfig),
		}
		commandUninstall = cobra.Command{
			Use:     "uninstall",
			Aliases: []string{"u"},
			Short:   "Uninstall the service.",
			Args:    cobra.NoArgs,
			Run:     UninstallService(serviceInstallerConfig),
		}
		commandRun = cobra.Command{
			Use:    "run",
			Short:  "Run the service.",
			Args:   cobra.ArbitraryArgs,
			Run:    RunService(instance, debug),
			Hidden: true,
		}
		rootCmd = cobra.Command{
			Use:   "My Service",
			Short: "Manage the service.",
		}
	)

	rootCmd.PersistentFlags().BoolVar(
		debug, "debug", false, "Run the service in debug mode.",
	)

	rootCmd.AddCommand(
		&commandStatus,
		&commandInstall,
		&commandUninstall,
		&commandRun,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command failed: %v\n", err)
	}

	time.Sleep(5 * time.Second)
}
