package main

import (
	"bitbucket.org/wepala/weos-controller/cmd"
	log "github.com/sirupsen/logrus"
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
