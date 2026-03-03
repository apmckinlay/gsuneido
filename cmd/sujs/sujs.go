// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/abemedia/go-webview"
	_ "github.com/abemedia/go-webview/embedded" // embed native library
)

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
	cmd := exec.Command(exe,
		"Use('stdlib'); "+
			"Suneido.RunAsStandalone = true;"+ // so closing the WorkSpace exits
			"RunSuJSHttpServer(3248)")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	w := webview.New(true)
	defer w.Destroy()

	// close the window if the server exits
	// e.g. from Shutdown(alsoServer:)
	done := make(chan struct{})
	go func() {
		cmd.Wait()
		close(done)
		w.Terminate()
	}()

	w.SetTitle("Suneido")
	w.SetSize(1510, 1024, webview.HintNone)

	// Bind quit function for macOS Cmd+Q
	w.Bind("quit", func() {
		w.Terminate()
	})
	w.Init(`
		document.addEventListener('keydown', function(e) {
			if (e.metaKey && e.key === 'q') {
				e.preventDefault();
				quit();
			}
		});
	`)

	w.Navigate("http://127.0.0.1:3248")
	w.Run()

	<-done
	time.Sleep(200 * time.Millisecond)
}
