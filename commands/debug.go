package commands

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugCommand = &cobra.Command{
	Use:   "debug",
	Short: "Log debug output",
	Long:  "Log all configuration keys and values",
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range viper.AllSettings() {
			log.Infof("%s: %v", k, v)
		}
	},
}
