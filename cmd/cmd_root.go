package cmd

import (
	"fmt"
"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	Version string
	Commit  string
	Date    string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "ruslan-cli",
	Short: "CLI tool for managing HashiCorp Vault across multiple environments",
	Long: `ruslan-cli is a command-line tool for managing HashiCorp Vault instances
deployed in GCP/GKE across multiple environments (dev/prod).

It provides easy switching between environments, authentication management,
and common secret operations.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ruslan-cli/config.yaml)")
	rootCmd.PersistentFlags().String("env", "", "environment to use (dev/prod)")
	rootCmd.PersistentFlags().String("format", "table", "output format (table, json, yaml)")
	
	viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("env"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))

	// Version template
	rootCmd.SetVersionTemplate(`ruslan-cli {{.Version}}
Commit: ` + Commit + `
Built: ` + Date + `
`)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Search config in home directory
		viper.AddConfigPath(home + "/.ruslan-cli")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err == nil {
		// Config loaded successfully
	}
}
