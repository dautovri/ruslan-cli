package auth

import (
vaultapi "github.com/hashicorp/vault/api"
)

// LoginWithToken authenticates with a token
func LoginWithToken(client *vaultapi.Client, token string) error {
	client.SetToken(token)
	
	// Verify token is valid
	_, err := client.Auth().Token().LookupSelf()
	return err
}

// LoginWithUserPass authenticates with username/password
func LoginWithUserPass(client *vaultapi.Client, username, password string) (string, error) {
	data := map[string]interface{}{
		"password": password,
	}
	
	secret, err := client.Logical().Write("auth/userpass/login/"+username, data)
	if err != nil {
		return "", err
	}
	
	return secret.Auth.ClientToken, nil
}

// LoginWithAppRole authenticates with AppRole
func LoginWithAppRole(client *vaultapi.Client, roleID, secretID string) (string, error) {
	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}
	
	secret, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		return "", err
	}
	
	return secret.Auth.ClientToken, nil
}
