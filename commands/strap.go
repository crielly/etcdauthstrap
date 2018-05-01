package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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

	api, err := etcd.New(etcd.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 10 * time.Second,
		TLS:         tlsConfig,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer api.Close()

	ssmclient := getSSMClient()

	for _, user := range viper.GetStringSlice("users.users") {

		parampath := fmt.Sprintf("%s/%s/password", viper.GetString("environment.passpathprefix"), user)

		log.Infof("SSM Parameter Path for %s password: %s", user, parampath)

		getParamInput := &ssm.GetParameterInput{
			Name:           &parampath,
			WithDecryption: aws.Bool(true),
		}

		resp, err := ssmclient.GetParameter(getParamInput)

		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Password for user %s: %s", user, *resp.Parameter.Value)

		_, err = api.Auth.UserAdd(context.TODO(), user, *resp.Parameter.Value)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func getRegion() (region string) {
	sess := session.Must(session.NewSession())

	client := ec2metadata.New(sess, aws.NewConfig())

	region, err := client.Region()

	if err != nil {
		log.Fatal(err)
	}

	return region
}

func getSSMClient() *ssm.SSM {
	region := getRegion()

	sess := session.Must(session.NewSession())

	client := ssm.New(sess, aws.NewConfig().WithRegion(region))

	return client
}
