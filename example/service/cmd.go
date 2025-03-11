package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	service "github.com/Nigel23932/go-svc/src"
	"github.com/Nigel23932/go-svc/src/elevation"
	"github.com/Nigel23932/go-svc/src/installer"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc"
)

// Keeps the command window open until the user presses Ctrl+C or the program is terminated.
func waitBeforeExit(msg ...any) {
	if len(msg) > 0 {
		log.Println(msg...)
	}

	println("Press Ctrl+C to exit.")
	var exitSignal = make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	os.Exit(0)
}

func InstallService(cnf *installer.ServiceInstallerConfig) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		if err := elevation.Elevate(); err != nil {
			log.Fatalf("Failed to elevate for administrator permissions: %v\n", err)
			return
		}

		var installer, err = installer.NewServiceInstaller(cnf)
		if err != nil {
			waitBeforeExit(err)
			return
		}
		defer installer.Close()

		installed, err := installer.Installed()
		if err != nil {
			waitBeforeExit(err)
			return
		}

		if installed {
			waitBeforeExit("Service already installed.")
			return
		}

		log.Println("Installing service...")
		err = installer.Install()
		log.Println("Service installed.")

		if err != nil {
			log.Println(err)
		}

		waitBeforeExit()
	}
}

func UninstallService(cnf *installer.ServiceInstallerConfig) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		if err := elevation.Elevate(); err != nil {
			log.Fatalf("Failed to elevate for administrator permissions: %v\n", err)
			return
		}

		var installer, err = installer.NewServiceInstaller(cnf)
		if err != nil {
			waitBeforeExit(err)
			return
		}
		defer installer.Close()

		installed, err := installer.Installed()
		if err != nil {
			waitBeforeExit(err)
			return
		}

		if !installed {
			waitBeforeExit("Service not found.")
			return
		}

		log.Println("Uninstalling service...")
		err = installer.Remove()
		if err != nil {
			waitBeforeExit(err)
			return
		}

		log.Println("Service uninstalled.")

		waitBeforeExit()
	}
}

func InstallationStatus(cnf *installer.ServiceInstallerConfig) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		if err := elevation.Elevate(); err != nil {
			log.Fatalf("Failed to elevate for administrator permissions: %v\n", err)
			return
		}

		var installer, err = installer.NewServiceInstaller(cnf)
		if err != nil {
			waitBeforeExit(err)
			return
		}
		defer installer.Close()

		installed, err := installer.Installed()
		if err != nil {
			waitBeforeExit(err)
			return
		}

		if installed {
			log.Println("Service is installed.")

		} else {
			log.Println("Service is not installed.")
			waitBeforeExit()
			return
		}

		status, err := installer.QueryServiceStatus()
		if err != nil {
			waitBeforeExit(err)
			return
		}

		switch status.State {
		case svc.Stopped:
			log.Println("Service is stopped.")
		case svc.StartPending:
			log.Println("Service is starting.")
		case svc.StopPending:
			log.Println("Service is stopping.")
		case svc.Running:
			log.Println("Service is running.")
		case svc.ContinuePending:
			log.Println("Service is continuing.")
		case svc.PausePending:
			log.Println("Service is pausing.")
		case svc.Paused:
			log.Println("Service is paused.")
		default:
			log.Println("Service is in an unknown state.")
		}

		waitBeforeExit()
	}
}

func RunService(instance *Instance, debug *bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		var (
			runService func(service.Service) error
			elog       service.EventLog
			err        error
		)
		if *debug {
			fmt.Println("Running in debug mode.")
			runService = service.RunDebug
			elog, err = service.LoggerDebug(instance.Name())
		} else {
			runService = service.Run
			elog, err = service.Logger(instance.Name())
		}
		if err != nil {
			log.Fatalf("Event Log: %v\n", err)
		}
		defer elog.Close()

		instance.ELog = elog

		elog.Info(1, fmt.Sprintf("Starting %s service...", instance.Name()))

		err = runService(instance)
		if err != nil {
			elog.Error(1, err.Error())
		}
	}
}
