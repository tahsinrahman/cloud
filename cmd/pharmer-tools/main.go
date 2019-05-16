package main

import (
	"os"

	"github.com/pharmer/cloud/pkg/cmds"
	"kmodules.xyz/client-go/logs"
	_ "github.com/pharmer/cloud/pkg/credential/cloud"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
