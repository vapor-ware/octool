package cmd

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/vapor-ware/octool/pkg"
)

var (
	connectAddress    string
	connectTimeout    time.Duration
	connectUsername   string
	connectPassword   string
	connectClientID   string
	connectRootCAs    []string
	connectCert       string
	connectKey        string
	connectSkipVerify bool
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an OpenConfig-enabled server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("connecting to openconfig server...")
		connect()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVarP(&connectAddress, "address", "a", "", "the server address")
	connectCmd.Flags().DurationVarP(&connectTimeout, "timeout", "t", 0, "timeout")
	connectCmd.Flags().StringVarP(&connectUsername, "username", "u", "", "username to login with")
	connectCmd.Flags().StringVarP(&connectPassword, "password", "p", "", "password to login with")
	connectCmd.Flags().StringVarP(&connectClientID, "clientid", "i", "", "unique client ID")
	connectCmd.Flags().StringArrayVar(&connectRootCAs, "ca", []string{}, "server certificate authorities")
	connectCmd.Flags().StringVarP(&connectCert, "cert", "c", "", "path to a tls cert")
	connectCmd.Flags().StringVarP(&connectKey, "key", "k", "", "path to a tls key")
	connectCmd.Flags().BoolVarP(&connectSkipVerify, "skip-verify", "v", false, "skip certificate verification")
}

func connect() {
	if connectAddress != "" {
		cfg.Address = connectAddress
	}
	if connectTimeout != 0 {
		cfg.Timeout = connectTimeout
	}
	if connectUsername != "" {
		cfg.Auth.Username = connectUsername
	}
	if connectPassword != "" {
		cfg.Auth.Password = connectPassword
	}
	if connectClientID != "" {
		cfg.Auth.ClientID = connectClientID
	}
	if len(connectRootCAs) > 0 {
		cfg.TLS.CAs = connectRootCAs
	}
	if connectCert != "" {
		cfg.TLS.Cert = connectCert
	}
	if connectKey != "" {
		cfg.TLS.Key = connectKey
	}
	if connectSkipVerify != false {
		cfg.TLS.SkipVerify = true
	}

	log.Debugf("connect config: %#v", cfg)

	ctx := context.Background()

	client, err := pkg.NewGRPCClient(&cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = client.Connect(ctx, cfg.Address)
	if err != nil {
		log.Fatal(err.Error())
	}
}
