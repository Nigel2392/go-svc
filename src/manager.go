//go:build !debug
// +build !debug

package service

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

func Run(service Service) error {
	return svc.Run(service.Name(), service)
}

func Logger(source string) (EventLog, error) {
	var eLog, err = eventlog.Open(source)
	return eLog, err
}
