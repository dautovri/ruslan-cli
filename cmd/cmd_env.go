package cmd

import (
	"fmt"
	"os"

	"github.com/dautovri/ruslan-cli/pkg/config"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
	Long:  `Commands for managing and switching between Vault environments (dev/prod).`,
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Environment", "Current", "Cluster", "Region"})
		table.SetBorder(false)

		for name, env := range cfg.Environments {
			current := ""
			if name == cfg.CurrentEnvironment {
				current = "✓"
			}
			table.Append([]string{name, current, env.ClusterName, env.Region})
		}

		table.Render()
		return nil
	},
}

var envUseCmd = &cobra.Command{
	Use:   "use [environment]",
	Short: "Switch to a different environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := args[0]
		
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if _, exists := cfg.Environments[envName]; !exists {
			return fmt.Errorf("environment '%s' not found", envName)
		}

		cfg.CurrentEnvironment = envName
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Switched to environment: %s\n", envName)
		return nil
	},
}

var envCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println(cfg.CurrentEnvironment)
		return nil
	},
}

var envInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current environment details",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		env, exists := cfg.Environments[cfg.CurrentEnvironment]
		if !exists {
			return fmt.Errorf("current environment not found")
		}

		fmt.Printf("Environment: %s\n", cfg.CurrentEnvironment)
		fmt.Printf("Project ID:  %s\n", env.ProjectID)
		fmt.Printf("Cluster:     %s\n", env.ClusterName)
		fmt.Printf("Region:      %s\n", env.Region)
		fmt.Printf("Namespace:   %s\n", env.Namespace)
		fmt.Printf("Vault Addr:  %s\n", env.VaultAddr)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envUseCmd)
	envCmd.AddCommand(envCurrentCmd)
	envCmd.AddCommand(envInfoCmd)
}
