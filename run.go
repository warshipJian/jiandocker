package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(context *cli.Context) {
	command := context.Args().Get(0)
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	tty := context.Bool("ti")
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Start(); err != nil {
		log.Error(err)
	}

	// 记录容器信息
	containerName := context.String("name")
	containerName, err := recordContainerINfo(cmd.Process.Pid, command, containerName)
	if err != nil {
		log.Error(err)
	}

	limitMemory := context.String("m")
	if limitMemory != "" {
		MemoryLimit(cmd.Process.Pid, limitMemory)
	}
	if tty {
		err := cmd.Wait()
		// 终端模式下退出后，删除信息
		deleteContainerinfo(containerName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
