package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

	// 记录容器信息
	containerName := context.String("name")
	command := strings.Join(args[1:], " ")
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

func getCmdArray(context *cli.Context, cmdArray []string) []string {
	for _, arg := range context.Args() {
		cmdArray = append(cmdArray, arg)
	}
	return cmdArray
}
