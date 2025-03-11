//go:build debug
// +build debug

package service

import (
	"golang.org/x/sys/windows/svc/debug"
)

func Run(service Service) error {
	return debug.Run(service.Name(), service)
}

func Logger(source string) (EventLog, error) {
	return debug.New(source), nil
}
