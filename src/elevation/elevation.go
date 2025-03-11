package elevation

import (
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func Elevate() error {
	if HasAdminRights() {
		return nil
	}

	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var hwnd windows.Handle
	var showCmd int32 = 1 //SW_NORMAL
	err := windows.ShellExecute(
		hwnd,
		verbPtr,
		exePtr,
		argPtr,
		cwdPtr,
		showCmd,
	)
	if err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func HasAdminRights() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}
