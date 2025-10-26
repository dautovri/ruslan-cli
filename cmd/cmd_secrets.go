package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dautovri/ruslan-cli/pkg/vault"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets in Vault",
	Long:  `Commands for reading, writing, and managing secrets in Vault KV v2 engine.`,
}

var secretsListCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List secrets at a path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		secrets, err := client.ListSecrets(path)
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		for _, secret := range secrets {
			fmt.Println(secret)
		}

		return nil
	},
}

var secretsGetCmd = &cobra.Command{
	Use:   "get [path]",
	Short: "Read a secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		format, _ := cmd.Flags().GetString("format")
		field, _ := cmd.Flags().GetString("field")
		
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		secret, err := client.GetSecret(path)
		if err != nil {
			return fmt.Errorf("failed to get secret: %w", err)
		}

		// If specific field requested
		if field != "" {
			if val, ok := secret.Data[field]; ok {
				fmt.Println(val)
				return nil
			}
			return fmt.Errorf("field '%s' not found", field)
		}

		// Output full secret in requested format
		switch format {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(secret.Data)
		case "yaml":
			enc := yaml.NewEncoder(os.Stdout)
			return enc.Encode(secret.Data)
		default: // table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Key", "Value"})
			table.SetBorder(false)
			for k, v := range secret.Data {
				table.Append([]string{k, fmt.Sprintf("%v", v)})
			}
			table.Render()
		}

		return nil
	},
}

var secretsPutCmd = &cobra.Command{
	Use:   "put [path] [key=value ...]",
	Short: "Write a secret",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		dataFile, _ := cmd.Flags().GetString("file")
		
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		var data map[string]interface{}

		// Load from file if specified
		if dataFile != "" {
			fileData, err := os.ReadFile(dataFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			if err := json.Unmarshal(fileData, &data); err != nil {
				return fmt.Errorf("failed to parse JSON: %w", err)
			}
		} else {
			// Parse key=value pairs
			data = make(map[string]interface{})
			for _, arg := range args[1:] {
				parts := strings.SplitN(arg, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid key=value pair: %s", arg)
				}
				data[parts[0]] = parts[1]
			}
		}

		if err := client.PutSecret(path, data); err != nil {
			return fmt.Errorf("failed to write secret: %w", err)
		}

		fmt.Printf("✓ Secret written to %s\n", path)
		return nil
	},
}

var secretsDeleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Short: "Delete a secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		if err := client.DeleteSecret(path); err != nil {
			return fmt.Errorf("failed to delete secret: %w", err)
		}

		fmt.Printf("✓ Secret deleted: %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsGetCmd)
	secretsCmd.AddCommand(secretsPutCmd)
	secretsCmd.AddCommand(secretsDeleteCmd)

	// Flags
	secretsGetCmd.Flags().String("format", "table", "output format (table, json, yaml)")
	secretsGetCmd.Flags().String("field", "", "specific field to retrieve")
	secretsPutCmd.Flags().String("file", "", "JSON file containing secret data")
}
