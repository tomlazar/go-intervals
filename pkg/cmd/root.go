package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "intervals <command> <subcommand> [flags]",
	Short: "Intervals CLI",
	Long:  `Work seamlessly with Intervals from the command line.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	// add flags
	rootCmd.PersistentFlags().String("apikey", "", "sets the api key")
	rootCmd.PersistentFlags().String("baseurl", "https://api.myintervals.com/", "sets the api url")

	// bind to viper
	viper.BindPFlag("intervals.key", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("intervals.url", rootCmd.PersistentFlags().Lookup("baseurl"))

	// add extra commands
	rootCmd.AddCommand(apiCmd)
}

func initConfig() {
	viper.SetConfigName(".intervalsrc")   // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}
}
