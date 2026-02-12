// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
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
			"Use('axonlib'); "+
			"Use('Accountinglib'); "+
			"Use('etalib'); "+
			"Use('towlib'); "+
			"Use('pcmiler'); "+
			"Use('ticketlib'); "+
			"Use('joblib'); "+
			"Use('prlib'); "+
			"Use('prcadlib'); "+
			"Use('etaprlib'); "+
			"Use('invenlib'); "+
			"Use('wolib'); "+
			"Use('polib'); "+
			"Use('carslib'); "+
			"Use('configlib'); "+
			"Use('demobookoptions'); "+
			"Use('Test_lib'); "+
			// so closing the WorkSpace exits, requires SuInit change
			"Suneido.RunAsStandalone = true;"+ 
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
	finished := false
	go func() {
		cmd.Wait()
		w.Terminate()
		finished = true
	}()

	w.SetTitle("Suneido")
	w.SetSize(1536, 1024, webview.HintNone)
	w.Navigate("http://127.0.0.1:3248")
	w.Run()

	if !finished {
		// give the server a chance to clean up the client
		time.Sleep(200 * time.Millisecond)
		client := &http.Client{Timeout: 100 * time.Millisecond}
		_, err = client.Post("http://127.0.0.1:3248/shutdown", "text/plain", nil)
		if err != nil {
			cmd.Wait()
		}
	}
}
