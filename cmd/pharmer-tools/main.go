package main

import (
	"os"

	"kmodules.xyz/client-go/logs"
	"pharmer.dev/cloud/pkg/cmds"
	_ "pharmer.dev/cloud/pkg/credential/cloud"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
