package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
)

func main() {
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if event.Op == 1 {
					goExecutable, _ := exec.LookPath( "go" )

					cmdGoRun := &exec.Cmd {
						Path: goExecutable,
						Args: []string{ goExecutable, "run", "/Users/vitalii/Work/Hakaton/cmd/uplink/main.go", "cp", event.Name, "sj://bucket", },
						Stdout: os.Stdout,
						Stderr: os.Stdout,
					}

					err = cmdGoRun.Run()
					if err != nil {
						fmt.Println("ERROR", err)
					}
				}

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add("/Users/vitalii/Work/Hakaton/bucket"); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
