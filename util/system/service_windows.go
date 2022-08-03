// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package service is minimal service support for a Windows program.
// It only supports start and stop/shutdown.
// There is minimal error handling.
// It only supports one service per exe.
// The code is based on https://pkg.go.dev/golang.org/x/sys/windows/svc and
// https://www.codeproject.com/Articles/499465/Simple-Windows-Service-in-Cplusplus
//
// It does NOT include support for install or remove.
// Instead, use the Windows sc command:
//
//	sc create <service_name> binpath= "<path_to_exe> <exe_args>" <options>
//	sc start <service_name>
//	sc stop <service_name>
//	sc delete <service_name>
package system

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
)

var (
	initCallbacks       sync.Once
	ctlHandlerCallback  uintptr
	serviceMainCallback uintptr
	serviceName         string
	serviceHandle       windows.Handle
	stopChan            chan bool
	stopFunc            func()
)

// Run handles if started as a service.
// If not started as a service, it does nothing.
// If running as a service it changes the working directory to the exe directory.
// The start function is called if running as a service.
// The stop function is called when the service is stopped.
func Service(name string, start, stop func()) error {
	inService, err := svc.IsWindowsService()
	if err != nil || !inService {
		return err
	}
	os.Chdir(filepath.Dir(os.Args[0]))
	if start != nil {
		start()
	}
	serviceName = name
	stopFunc = stop
	go runService(name)
	return nil
}

func runService(name string) {
	initCallbacks.Do(func() {
		ctlHandlerCallback = windows.NewCallback(ctlHandler)
		serviceMainCallback = windows.NewCallback(serviceMain)
		stopChan = make(chan bool)
	})
	t := []windows.SERVICE_TABLE_ENTRY{
		{ServiceName: windows.StringToUTF16Ptr(name), ServiceProc: serviceMainCallback},
		{ServiceName: nil, ServiceProc: 0},
	}
	err := windows.StartServiceCtrlDispatcher(&t[0])
	if err != nil {
		log.Println("runService", err)
	}
}

func serviceMain(uint32, **uint16) uintptr {
	handle, err := windows.RegisterServiceCtrlHandlerEx(
		windows.StringToUTF16Ptr(serviceName), ctlHandlerCallback,
		uintptr(unsafe.Pointer(&serviceHandle)))
	if sysErr, ok := err.(windows.Errno); ok {
		return uintptr(sysErr)
	} else if err != nil {
		return uintptr(windows.ERROR_UNKNOWN_EXCEPTION)
	}
	serviceHandle = handle
	defer func() {
		serviceHandle = 0
	}()

	updateStatus(windows.SERVICE_START_PENDING, 0)
	updateStatus(windows.SERVICE_RUNNING,
		windows.SERVICE_ACCEPT_STOP|windows.SERVICE_ACCEPT_SHUTDOWN)

	<-stopChan // wait for ctlHandler

	if stopFunc != nil {
		stopFunc()
	}

	updateStatus(windows.SERVICE_STOPPED, 0)
	return windows.NO_ERROR
}

func ctlHandler(ctl, _, _, _ uintptr) uintptr {
	switch ctl {
	case windows.SERVICE_CONTROL_STOP, windows.SERVICE_CONTROL_SHUTDOWN:
		updateStatus(windows.SERVICE_STOP_PENDING, 0)
		stopChan <- true
	}
	return 0
}

func updateStatus(state, accept uint32) {
	var ss windows.SERVICE_STATUS
	ss.ServiceType = windows.SERVICE_WIN32_OWN_PROCESS
	ss.CurrentState = state
	ss.ControlsAccepted = accept
	err := windows.SetServiceStatus(serviceHandle, &ss)
	if err != nil {
		log.Println("updateStatus", err)
	}
}
