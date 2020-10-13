package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wepala/weos-controller/cmd"
)

var (
	version string
	build   string
)

func main() {
	log.Infof("version=%s", version)
	log.Infof("build=%s", build)
	cmd.Execute()
}
