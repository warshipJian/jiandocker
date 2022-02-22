package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `run a basic docker
			jiandocker run [command] -m [limit memory]`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}

		fmt.Println("aaa")

		command := context.Args().Get(0)
		args := []string{"init", command}
		cmd := exec.Command("/proc/self/exe", args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
				syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
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
		cmd.Wait()
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `mount proc system `,
	Action: func(context *cli.Context) error {

		fmt.Println("bbb")

		command := context.Args().Get(0)
		argv := []string{command}
		fmt.Println(argv)

		// unshareM()

		// 挂载
		setUpMount()

		// 执行命令
		if err := syscall.Exec(command, argv, os.Environ()); err != nil {
			log.Errorf(err.Error())
		}
		return nil
	},
}
