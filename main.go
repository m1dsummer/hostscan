package main

import (
	"flag"
	"fmt"
	"hostscan/core"
	"hostscan/elog"
	"hostscan/utils"
	"hostscan/vars"
	"os"
)

func main() {
	flag.Parse()
	if *vars.Version {
		elog.Info(fmt.Sprintf("Current Hostscan Version: %s", vars.VersionInfo))
		return
	}

	elog.Info("Hostscan Start! Waiting for your good news...")

	if len(*vars.OutFile) > 0 {
		exist, _ := utils.PathExists(*vars.OutFile)
		if exist {
			_ = os.Remove(*vars.OutFile)
		}
	}

	utils.SetUlimitMax()
	err := core.Scan()
	if err != nil {
		elog.Error(fmt.Sprintf("Scan Failed: %v", err))
	}

	if err != nil {
		elog.Error(fmt.Sprintf("ProcessBar Close Failed: %v", err))
		return
	}
}
