package service

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/mgr"
)

type EventLog interface {
	Close() error
	Info(eid uint32, msg string) error
	Warning(eid uint32, msg string) error
	Error(eid uint32, msg string) error
}

type Service interface {
	Name() string
	Config() *mgr.Config
	Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32)
	AcceptedCommands() svc.Accepted
	EventLog() EventLog
}

func RunDebug(service Service) error {
	return debug.Run(service.Name(), service)
}

func LoggerDebug(source string) (EventLog, error) {
	return debug.New(source), nil
}
