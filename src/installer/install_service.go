package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	service "github.com/Nigel2392/go-svc/src"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

var _ Installer = (*serviceInstaller)(nil)

type ServiceInstallerConfig struct {
	Service         service.Service
	EventsSupported uint32
	Args            []string
}

type serviceInstaller struct {
	cnf ServiceInstallerConfig
	mgr *mgr.Mgr
}

func NewServiceInstaller(cnf *ServiceInstallerConfig) (*serviceInstaller, error) {
	var si = &serviceInstaller{
		cnf: *cnf,
	}

	m, err := mgr.Connect()
	if err != nil {
		return nil, errors.Wrap(
			err, "NewServiceInstaller.mgr.Connect(): failed to connect to service manager",
		)
	}

	si.mgr = m
	return si, nil
}

func (si *serviceInstaller) Close() error {
	if si.mgr != nil {
		var err = si.mgr.Disconnect()
		if err != nil {
			return errors.Wrap(
				err, "Close.mgr.Disconnect(): failed to disconnect from service manager",
			)
		}
		si.mgr = nil
	}
	return nil
}

func (si *serviceInstaller) Install() error {
	var (
		cnf     = si.cnf.Service.Config()
		exepath = cnf.BinaryPathName
		err     error
	)
	exepath = strings.Split(exepath, " ")[0]

	if exepath == "" {
		exepath, err = ExePath()
		if err != nil {
			return errors.Wrap(
				err, "Install.__global__.ExePath(): failed to determine executable path",
			)
		}
	}

	var serviceName = si.cnf.Service.Name()
	s, err := si.mgr.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("Install.mgr.OpenService(): service %s already exists", serviceName)
	}

	s, err = si.mgr.CreateService(serviceName, exepath, *cnf, si.cnf.Args...)
	if err != nil {
		return fmt.Errorf("Install.mgr.CreateService(): failed: %w", err)
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(serviceName, si.cnf.EventsSupported)
	if err != nil {
		s.Delete()
		return fmt.Errorf("Install.eventlog.SetupEventLogSource(): failed: %w", err)
	}

	return nil
}

func (si *serviceInstaller) Remove() error {
	var serviceName = si.cnf.Service.Name()
	var s, err = si.mgr.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("Remove.mgr.OpenService(): service %s is not installed", serviceName)
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("Remove.service.Query(): failed: %s", err)
	}

	if status.State != svc.Stopped && (status.Accepts&svc.AcceptStop != 0 || status.Accepts&svc.AcceptShutdown != 0) {
		var (
			newStatus svc.Status
			err       error
		)

		switch {
		case status.Accepts&svc.AcceptShutdown != 0:
			newStatus, err = s.Control(svc.Shutdown)
		case status.Accepts&svc.AcceptStop != 0:
			newStatus, err = s.Control(svc.Stop)
		}

		if err != nil {
			return fmt.Errorf("Remove.service.Control(): failed to stop service: %w", err)
		}

		if newStatus.State != svc.Stopped {
			return fmt.Errorf("Remove.service.Control(): failed to stop service: state %d", newStatus.State)
		}
	}

	err = s.Delete()
	if err != nil {
		return fmt.Errorf("Remove.service.Delete(): failed: %w", err)
	}

	err = eventlog.Remove(serviceName)
	if err != nil {
		return fmt.Errorf("Remove.eventlog.Remove(): failed: %w", err)
	}

	return nil
}

func (si *serviceInstaller) QueryServiceStatus() (svc.Status, error) {
	var serviceName = si.cnf.Service.Name()
	var s, err = si.mgr.OpenService(serviceName)
	if err != nil {
		return svc.Status{}, fmt.Errorf("QueryServiceStatus.mgr.OpenService(): service %s is not installed", serviceName)
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return svc.Status{}, fmt.Errorf("QueryServiceStatus.service.Query(): failed: %w", err)
	}

	return status, nil
}

func (si *serviceInstaller) Installed() (bool, error) {
	var serviceName = si.cnf.Service.Name()
	s, err := si.mgr.OpenService(serviceName)
	if err != nil {
		return false, nil
	}
	defer s.Close()

	return true, nil
}

func ExePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err = os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}
