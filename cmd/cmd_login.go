package cmd

import (
	"fmt"

	"github.com/dautovri/ruslan-cli/pkg/vault"
	"github.com/spf13/cobra"
)

var (
	loginMethod   string
	loginToken    string
	loginRoleID   string
	loginSecretID string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to Vault",
	Long:  `Authenticate to Vault using various methods (token, userpass, approle).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		switch loginMethod {
		case "token":
			if loginToken == "" {
				return fmt.Errorf("token is required for token auth")
			}
			if err := client.LoginWithToken(loginToken); err != nil {
				return fmt.Errorf("token authentication failed: %w", err)
			}

		case "userpass":
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			if username == "" {
				fmt.Print("Username: ")
				fmt.Scanln(&username)
			}
			if password == "" {
				fmt.Print("Password: ")
				// In production, use terminal.ReadPassword for secure input
				fmt.Scanln(&password)
			}

			_, err := client.LoginWithUserPass(username, password)
			if err != nil {
				return fmt.Errorf("userpass authentication failed: %w", err)
			}

		case "approle":
			if loginRoleID == "" || loginSecretID == "" {
				return fmt.Errorf("role-id and secret-id are required for approle auth")
			}
			_, err := client.LoginWithAppRole(loginRoleID, loginSecretID)
			if err != nil {
				return fmt.Errorf("approle authentication failed: %w", err)
			}

		default:
			return fmt.Errorf("unsupported auth method: %s", loginMethod)
		}

		fmt.Println("✓ Successfully authenticated to Vault")
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove saved authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		if err := client.Logout(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}

		fmt.Println("✓ Logged out successfully")
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Vault client: %w", err)
		}

		tokenInfo, err := client.GetTokenInfo()
		if err != nil {
			fmt.Println("Not authenticated")
			return nil
		}

		fmt.Printf("Authenticated: Yes\n")
		if tokenInfo != nil && tokenInfo.Data != nil {
			if accessor, ok := tokenInfo.Data["accessor"].(string); ok {
				fmt.Printf("Token Accessor: %s\n", accessor)
			}
			if policies, ok := tokenInfo.Data["policies"].([]interface{}); ok {
				fmt.Printf("Policies: %v\n", policies)
			}
			if ttl, ok := tokenInfo.Data["ttl"].(float64); ok {
				fmt.Printf("TTL: %.0fs\n", ttl)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication operations",
	}
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)

	// Login flags
	loginCmd.Flags().StringVar(&loginMethod, "method", "token", "authentication method (token, userpass, approle)")
	loginCmd.Flags().StringVar(&loginToken, "token", "", "vault token")
	loginCmd.Flags().StringVar(&loginRoleID, "role-id", "", "approle role ID")
	loginCmd.Flags().StringVar(&loginSecretID, "secret-id", "", "approle secret ID")
	loginCmd.Flags().String("username", "", "username for userpass auth")
	loginCmd.Flags().String("password", "", "password for userpass auth")
}
