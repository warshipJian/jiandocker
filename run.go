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
	args := getCmdArray(context, []string{"init"})
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	tty := context.Bool("ti")
	d := context.Bool("d")
	if tty && d {
		log.Error("ti和d参数不能同时存在")
		return
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		log.Error(err)
	}
	limitMemory := context.String("m")
	if limitMemory != "" {
		MemoryLimit(cmd.Process.Pid, limitMemory)
	}
	if tty {
		err := cmd.Wait()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func getCmdArray(context *cli.Context, cmdArray []string) []string {
	for _, arg := range context.Args() {
		cmdArray = append(cmdArray, arg)
	}
	return cmdArray
}
