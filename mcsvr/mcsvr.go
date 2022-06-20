// Package mcsvr launches, manages and monitors the running mincraft instance.
package mcsvr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/codeka/mcsvradmin/config"
)

// Instance represents a running instance of Minecraft.
type Instance struct {
	Cmd *exec.Cmd

	// cmdCreator is a function that will create a new instance of exec.Cmd for this command. We use
	// this to create a brand new exec.Cmd whenever it dies.
	cmdCreator func() *exec.Cmd
}

// Launch launches minecraft, starts monitoring it and returns an Instance for controlling it.
func Launch(cfg *config.Config) (*Instance, error) {
	var jarFilePath string
	err := filepath.Walk(cfg.BaseDirectory, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && path.Clean(cfg.BaseDirectory) != path.Clean(p) {
			return filepath.SkipDir
		}
		match, err := path.Match(cfg.JarFilePattern, filepath.Base(filepath.Clean(p)))
		if err != nil {
			return err
		}
		if match {
			jarFilePath = p
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("looking for jar: %v", err)
	}
	if jarFilePath == "" {
		return nil, fmt.Errorf("could not find JAR matching path %s in '%s'", cfg.JarFilePattern, cfg.BaseDirectory)
	}
	log.Printf("Found jar: %s", jarFilePath)

	javaCmd := "java"
	if cfg.JavaPath != "" {
		javaCmd = path.Join(cfg.JavaPath, "bin", "java")
	}

	inst := &Instance{}

	// TODO: use CommandContext so we can kill the process when we die as well.
	inst.cmdCreator = func() *exec.Cmd {
		cmd := exec.Command(javaCmd, "-jar", jarFilePath, "-nogui")
		cmd.Dir = cfg.BaseDirectory
		// TODO: set Stdin, Stdout and Stderr so that we can write to it ourselves.
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd
	}

	err = inst.start()
	if err != nil {
		return nil, fmt.Errorf("error starting java: %v", err)
	}

	return inst, nil
}

func (inst *Instance) start() error {
	if inst.Cmd != nil {
		// Something?
	}
	inst.Cmd = inst.cmdCreator()
	log.Printf("Running: %s", inst.Cmd.String())
	return inst.Cmd.Start()
}

// MonitorProcess monitors the running process and if it dies for whatever reason, waits a minute
// and then restarts it. Expected to run in a separator goroutine.
func (inst *Instance) MonitorProcess() {
	for {
		err := inst.Cmd.Wait()
		if err != nil {
			exitError := err.(*exec.ExitError)
			log.Printf("process exited with error code %d", exitError.ExitCode())
		} else {
			log.Printf("Process exited normally.")
		}

		// TODO: if not exiting...

		log.Printf("Waiting 20 seconds before restarting...")
		time.Sleep(20 * time.Second)

		log.Printf("Restarting the server...")
		inst.start()
	}
}

