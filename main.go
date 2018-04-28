package main

import (
	"os"

	"github.com/crielly/etcdauthstrap/commands"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
