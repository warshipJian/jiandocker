package main

import (
	"os"
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
		NewParentProcess(context)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `mount proc system busybox`,
	Action: func(context *cli.Context) error {

		// 挂载
		setUpMount()

		// 执行传入的命令
		command := context.Args().Get(0)
		argv := []string{command}
		if err := syscall.Exec(command, argv, os.Environ()); err != nil {
			log.Errorf(err.Error())
		}

		return nil
	},
}
