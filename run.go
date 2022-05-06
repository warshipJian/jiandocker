package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func NewParentProcess(context *cli.Context) {
	command := context.Args().Get(0)
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Error(err)
	}
	limitMemory := context.String("m")
	if limitMemory != "" {
		MemoryLimit(cmd.Process.Pid, limitMemory)
	}
	err := cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
}
