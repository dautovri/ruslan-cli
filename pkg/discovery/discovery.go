package discovery

import (
"fmt"
"os/exec"
"strings"
)

// DiscoverVaultAddress finds the Vault service LoadBalancer IP
func DiscoverVaultAddress(clusterName, region, namespace, serviceName string) (string, error) {
	// Get kubectl context for the cluster
	cmd := exec.Command("gcloud", "container", "clusters", "get-credentials",
clusterName, "--region", region)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get cluster credentials: %w", err)
	}

	// Get the LoadBalancer IP
	cmd = exec.Command("kubectl", "get", "svc", serviceName,
"-n", namespace,
"-o", "jsonpath={.status.loadBalancer.ingress[0].ip}")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get service IP: %w", err)
	}

	ip := strings.TrimSpace(string(output))
	if ip == "" {
		return "", fmt.Errorf("no external IP found for service %s in namespace %s", serviceName, namespace)
	}

	return fmt.Sprintf("http://%s:8200", ip), nil
}
