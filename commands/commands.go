package commands

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config  string
	version bool

	// RootCmd is the root command
	RootCmd = &cobra.Command{
		Use:   "etcdauthstrap",
		Short: "etcd auth config tool",
		Long:  `A utility for configuring and enabling Auth on etcd v2 and v3 APIs`,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			if config != "" {

				absolutepath, err := filepath.Abs(config)

				if err != nil {
					log.Error("Error reading configfile path: ", err)
				}

				base := filepath.Base(absolutepath)

				path := filepath.Dir(absolutepath)

				viper.SetConfigType("toml")
				viper.SetConfigName(strings.Split(base, ".")[0])
				viper.SetConfigFile(config)
				viper.AddConfigPath(path)

				if err := viper.ReadInConfig(); err != nil {
					log.Fatal("Failed to read config file: ", err.Error())
					os.Exit(1)
				}
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

func init() {

	// persistent flags are inherited by subcommands. Setting one on the root
	// command will make it global

	// CONFIG FILE
	RootCmd.PersistentFlags().String("configfile", "/var/cache/etcd/authconfig.toml", "Path to TOML config file")
	viper.BindPFlag("configfile", RootCmd.PersistentFlags().Lookup("configfile"))

	// LOGGING
	RootCmd.PersistentFlags().String("loglevel", "INFO", "Output level of logs TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	RootCmd.PersistentFlags().String("logtype", "stdout", "Log destination (stdout, file)")
	RootCmd.PersistentFlags().String("logfile", "/var/log/etcdauthstrap.log", "If logtype=file, accepts a path to a log file. Otherwise ignored")
	viper.BindPFlag("loglevel", RootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("logtype", RootCmd.PersistentFlags().Lookup("logtype"))
	viper.BindPFlag("logfile", RootCmd.PersistentFlags().Lookup("logfile"))

	// CONNECTION
	RootCmd.PersistentFlags().StringP("endpoint", "e", "localhost", "etcd API endpoint")
	RootCmd.PersistentFlags().StringP("port", "p", "2379", "etcd API port")
	RootCmd.PersistentFlags().StringP("scheme", "s", "https", "Transport Scheme (http, https)")
	viper.BindPFlag("endpoint", RootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("port", RootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("scheme", RootCmd.PersistentFlags().Lookup("scheme"))

	// TLS
	RootCmd.PersistentFlags().StringP("certfile", "c", "/etc/ssl/kubernetes/client-cert.pem", "Client Certificate for connecting to etcd API")
	RootCmd.PersistentFlags().StringP("keyfile", "k", "/etc/ssl/kubernetes/client-key.pem", "Client Key for connecting to etcd API")
	RootCmd.PersistentFlags().StringP("cafile", "a", "/etc/ssl/kubernetes/root-ca-cert.pem", "Trusted CA Certificate for connecting to etcd API")
	viper.BindPFlag("certfile", RootCmd.PersistentFlags().Lookup("certfile"))
	viper.BindPFlag("keyfile", RootCmd.PersistentFlags().Lookup("keyfile"))
	viper.BindPFlag("cafile", RootCmd.PersistentFlags().Lookup("cafile"))
	// local flags apply only to a specific command and are not inherited

	// commands
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(debugCommand)

}
