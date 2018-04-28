package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of etcdauthstrap",
	Long:  "Print the version number of etcdauthstrap!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ETCD Auth Strap v0.0 -- HEAD")
	},
}
