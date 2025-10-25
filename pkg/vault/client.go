package vault

import (
	"fmt"
	"strings"

	"github.com/dautovri/ruslan-cli/pkg/auth"
	"github.com/dautovri/ruslan-cli/pkg/config"
	vaultapi "github.com/hashicorp/vault/api"
)

type Client struct {
	*vaultapi.Client
	Config *config.Config
}

func NewClient() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	env := cfg.Environments[cfg.CurrentEnvironment]
	if env == nil {
		return nil, fmt.Errorf("environment not found: %s", cfg.CurrentEnvironment)
	}

	// Use configured Vault address
	if env.VaultAddr == "" {
		return nil, fmt.Errorf("vault address not configured for environment: %s", cfg.CurrentEnvironment)
	}

	vaultCfg := vaultapi.DefaultConfig()
	vaultCfg.Address = env.VaultAddr

	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, err
	}

	// Load saved token if exists
	if env.Token != "" {
		client.SetToken(env.Token)
	}

	return &Client{
		Client: client,
		Config: cfg,
	}, nil
}

func (c *Client) SaveToken(token string) error {
	c.Config.Environments[c.Config.CurrentEnvironment].Token = token
	return c.Config.Save()
}

// LoginWithToken authenticates with a token
func (c *Client) LoginWithToken(token string) error {
	err := auth.LoginWithToken(c.Client, token)
	if err != nil {
		return err
	}
	return c.SaveToken(token)
}

// LoginWithUserPass authenticates with username/password
func (c *Client) LoginWithUserPass(username, password string) (string, error) {
	token, err := auth.LoginWithUserPass(c.Client, username, password)
	if err != nil {
		return "", err
	}
	if err := c.SaveToken(token); err != nil {
		return "", err
	}
	return token, nil
}

// LoginWithAppRole authenticates with AppRole
func (c *Client) LoginWithAppRole(roleID, secretID string) (string, error) {
	token, err := auth.LoginWithAppRole(c.Client, roleID, secretID)
	if err != nil {
		return "", err
	}
	if err := c.SaveToken(token); err != nil {
		return "", err
	}
	return token, nil
}

// Logout clears the saved token
func (c *Client) Logout() error {
	c.Config.Environments[c.Config.CurrentEnvironment].Token = ""
	return c.Config.Save()
}

// GetTokenInfo returns information about the current token
func (c *Client) GetTokenInfo() (*vaultapi.Secret, error) {
	return c.Auth().Token().LookupSelf()
}

// ListSecrets lists secrets at a path (handles KV v2)
func (c *Client) ListSecrets(path string) ([]string, error) {
	// For KV v2, we need to add /metadata to list secrets
	metadataPath := path
	if !strings.HasPrefix(path, "secret/metadata/") && strings.HasPrefix(path, "secret/") {
		metadataPath = strings.Replace(path, "secret/", "secret/metadata/", 1)
	}

	secret, err := c.Logical().List(metadataPath)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	result := make([]string, 0, len(keys))
	for _, key := range keys {
		if str, ok := key.(string); ok {
			result = append(result, str)
		}
	}

	return result, nil
}

// GetSecret reads a secret (handles KV v2)
func (c *Client) GetSecret(path string) (*vaultapi.Secret, error) {
	// For KV v2, we need to add /data to the path
	dataPath := path
	if !strings.HasPrefix(path, "secret/data/") && strings.HasPrefix(path, "secret/") {
		dataPath = strings.Replace(path, "secret/", "secret/data/", 1)
	}
	return c.Logical().Read(dataPath)
}

// PutSecret writes a secret (handles KV v2)
func (c *Client) PutSecret(path string, data map[string]interface{}) error {
	// For KV v2, we need to add /data to the path and wrap data in "data" key
	dataPath := path
	if !strings.HasPrefix(path, "secret/data/") && strings.HasPrefix(path, "secret/") {
		dataPath = strings.Replace(path, "secret/", "secret/data/", 1)
	}

	// Wrap the data for KV v2
	wrappedData := map[string]interface{}{
		"data": data,
	}

	_, err := c.Logical().Write(dataPath, wrappedData)
	return err
}

// DeleteSecret deletes a secret (handles KV v2)
func (c *Client) DeleteSecret(path string) error {
	// For KV v2, we need to add /data to the path
	dataPath := path
	if !strings.HasPrefix(path, "secret/data/") && strings.HasPrefix(path, "secret/") {
		dataPath = strings.Replace(path, "secret/", "secret/data/", 1)
	}
	_, err := c.Logical().Delete(dataPath)
	return err
}
