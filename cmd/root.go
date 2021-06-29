package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vapor-ware/octool/pkg"
)

var cfgFile string
var debug bool
var cfg pkg.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "octool",
	Short: "A simple tool to help debug and verify OpenConfig configurations",
	Long: `This tool is intended to help inspect/debug/diagnose gRPC connections
to OpenConfig servers. To run with additional debug logs for gRPC, run the tool
with the following environment variables set:

  GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info octool ...
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initLogging, initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.octool.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}

func initLogging() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Debug("loading config...")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".octool" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".octool")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	if cfg.Auth.ClientID == "" {
		cfg.Auth.ClientID = pkg.RandStr()
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	log.Debugf("config loaded: %#v", cfg)
}
