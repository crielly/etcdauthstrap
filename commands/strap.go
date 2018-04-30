package commands

import (
	"fmt"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var strapCommand = &cobra.Command{
	Use:   "strap",
	Short: "Bootstrap etcd auth using the provided config",
	Long:  "Configure Users and Roles via both the etcd v2 and v3 APIs and enable Auth",
	Run: func(cmd *cobra.Command, args []string) {
		strap()
	},
}

func strap() {
	fmt.Println("Running strap command, adding users and enabling auth")

	tlsInfo := transport.TLSInfo{
		CertFile:      viper.GetString("tls.certfile"),
		KeyFile:       viper.GetString("tls.keyfile"),
		TrustedCAFile: viper.GetString("tls.cafile"),
	}

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf(
		"%s://%s:%d",
		viper.GetString("connection.scheme"),
		viper.GetString("connection.endpoint"),
		viper.GetInt("connection.port"),
	)

	log.Infof("Connection Endpoint: %s", endpoint)

	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 10 * time.Second,
		TLS:         tlsConfig,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	for _, user := range viper.GetStringSlice("users.users") {
		log.Infof("User: %s", user)
	}

	// _, err = cli.Auth.UserAdd(context.TODO(), "kube-apiserver", "Shyrriw00k")

	if err != nil {
		log.Fatal(err)
	}
}
