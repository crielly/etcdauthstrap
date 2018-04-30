package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of etcdauthstrap",
	Long:  "Print the version number of etcdauthstrap!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("etcdauthstrap v0.1")
	},
}
