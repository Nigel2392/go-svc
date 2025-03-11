package main

import (
	"time"

	service "github.com/Nigel23932/go-svc/src"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type InstanceFunc func(i *Instance, r svc.ChangeRequest, status chan<- svc.Status)

type Instance struct {
	ServiceName string
	Opts        *mgr.Config
	ELog        service.EventLog
	Accepted    svc.Accepted

	_t *time.Ticker
}

func (m *Instance) Name() string {
	return m.ServiceName
}

func (m *Instance) Config() *mgr.Config {
	return m.Opts
}

func (m *Instance) EventLog() service.EventLog {
	return m.ELog
}

func (m *Instance) AcceptedCommands() svc.Accepted {
	return m.Accepted
}

func (m *Instance) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {

	// Set the service status to StartPending during initialization.
	m.ELog.Info(1, "Service starting")
	status <- svc.Status{State: svc.StartPending}

	m.ELog.Info(1, "Setting up ticker")
	m._t = time.NewTicker(5 * time.Second)
	defer m._t.Stop()

	m.ELog.Info(1, "Service started")
	// The service is now running.
	// Set the service status to Running.
	status <- svc.Status{State: svc.Running, Accepts: m.Accepted}

mainLoop:
	for {
		select {
		case <-m._t.C:
			m.ELog.Info(1, "I'm alive")

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus

			case svc.Stop, svc.Shutdown:
				m.ELog.Info(1, "Service stopping")
				status <- svc.Status{State: svc.StopPending}
				break mainLoop

			case svc.Pause:
				status <- svc.Status{State: svc.Paused, Accepts: m.Accepted}

			case svc.Continue:
				status <- svc.Status{State: svc.Running, Accepts: m.Accepted}

			default:
				m.ELog.Error(2, "unexpected control request")
			}
		}
	}
	return
}
