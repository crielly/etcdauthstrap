package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(strapCmd)
}

var strapCmd = &cobra.Command{
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
		CertFile:      "/home/crielly/etcdapi/client-cert.pem",
		KeyFile:       "/home/crielly/etcdapi/client-key.pem",
		TrustedCAFile: "/home/crielly/etcdapi/root-ca-cert.pem",
	}

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{"https://localhost:2379"},
		DialTimeout: 10 * time.Second,
		TLS:         tlsConfig,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	_, err = cli.Auth.UserAdd(context.TODO(), "kube-apiserver", "Shyrriw00k")

	if err != nil {
		log.Fatal(err)
	}

	userlist, err := cli.Auth.UserList(context.TODO())

	fmt.Println(userlist)
}
