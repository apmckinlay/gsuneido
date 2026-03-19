// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const addr = "127.0.0.1:3248"

func main() {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	exe := fmt.Sprintf("./gs_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	server := exec.Command(exe,
		"Use('stdlib'); Suneido.RunAsStandalone = true; RunSuJSHttpServer(3248)")
	setSysProcAttr(server)
	server.Stdout = logFile
	server.Stderr = logFile
	if err = server.Start(); err != nil {
		log.Fatal(err)
	}

	waitForServer()

	browser := findBrowser()
	if browser == "" {
		log.Fatal("no supported Chromium-based browser found")
	}
	// --app mode gives a minimal window without browser chrome.
	// --user-data-dir forces a new instance even if the browser is already running.
	userDataDir := filepath.Join(os.TempDir(), "sujs-browser")
	browserCmd := exec.Command(browser,
		"--app=http://"+addr,
		"--new-window",
		"--user-data-dir="+userDataDir,
		"--window-size=1510,1024",
	)
	setSysProcAttr(browserCmd)
	browserCmd.Stdout = logFile
	browserCmd.Stderr = logFile
	if err = browserCmd.Start(); err != nil {
		log.Fatal("failed to start browser:", err)
	}

	serverDone := make(chan struct{})
	browserDone := make(chan struct{})

	go func() {
		server.Wait()
		close(serverDone)
	}()

	go func() {
		browserCmd.Wait()
		close(browserDone)
	}()

	select {
	case <-serverDone:
		killProcessGroup(browserCmd)
	case <-browserDone:
		// wait for the server to shut down on its own, then kill as a backup
		select {
		case <-serverDone:
		case <-time.After(5 * time.Second):
			killProcessGroup(server)
		}
	}
}

// waitForServer polls addr until a TCP connection succeeds or 5 seconds elapse.
func waitForServer() {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Fatal("timed out waiting for server to start")
}

func findBrowser() string {
	if b := findBrowserApp(); b != "" {
		return b
	}
	for _, b := range browsers {
		if (runtime.GOOS == "windows") == strings.HasSuffix(b, ".exe") {
			if path, err := exec.LookPath(b); err == nil {
				return path
			}
		}
	}
	return ""
}

// browsers to try in order of preference
var browsers = []string{
	// Chromium / Chrome
	"chrome.exe",
	"ungoogled-chromium",
	"chromium",
	"chromium.exe",
	"chromium-browser",
	"google-chrome",
	"google-chrome-stable",
	"chrome",
	// Edge
	"microsoft-edge",
	"microsoft-edge-stable",
	"msedge",
	"msedge.exe",
	"MicrosoftEdge.exe",
	// Brave
	"brave-browser",
	"brave",
	"brave.exe",
	// Vivaldi
	"vivaldi",
	"vivaldi-stable",
	"vivaldi.exe",
	// Thorium
	"thorium",
	"thorium-browser",
	"thorium.exe",
	// Helium
	"helium",
	"helium-browser",
	"helium.exe",
	// Arc
	"arc",
	"arc.exe",
}
