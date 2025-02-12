package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ListProjects returns a list of available GCP projects
func ListProjects() ([]string, error) {
	cmd := exec.Command("gcloud", "projects", "list", "--format", "value(projectId)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

// ListPools returns a list of workload identity pools for a given project
func ListPools(projectID string) ([]string, error) {
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "list",
		"--project", projectID,
		"--location", "global",
		"--format", "value(name)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list pools: %v", err)
	}
	if len(output) == 0 {
		return []string{}, nil
	}
	pools := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Extract just the pool name from the full resource name
	for i, pool := range pools {
		parts := strings.Split(pool, "/")
		pools[i] = parts[len(parts)-1]
	}
	return pools, nil
}

// ListServiceAccounts returns a list of service accounts for a given project
func ListServiceAccounts(projectID string) ([]string, error) {
	cmd := exec.Command("gcloud", "iam", "service-accounts", "list",
		"--project", projectID,
		"--format", "value(email)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list service accounts: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

// ListProviders returns a list of workload identity providers for a given project and pool
func ListProviders(projectID, poolName string) ([]string, error) {
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "providers", "list",
		"--project", projectID,
		"--location", "global",
		"--workload-identity-pool", poolName,
		"--format", "value(name)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %v", err)
	}
	if len(output) == 0 {
		return []string{}, nil
	}
	providers := strings.Split(strings.TrimSpace(string(output)), "\n")
	// Extract just the provider name from the full resource name
	for i, provider := range providers {
		parts := strings.Split(provider, "/")
		providers[i] = parts[len(parts)-1]
	}
	return providers, nil
}

// SelectFromList presents a numbered list to the user and returns their selection
func SelectFromList(items []string, resourceType string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no %s available", resourceType)
	}

	fmt.Printf("\nAvailable %s:\n", resourceType)
	for i, item := range items {
		fmt.Printf("%d) %s\n", i+1, item)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nEnter number (1-%d) to select %s: ", len(items), resourceType)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %v", err)
		}

		input = strings.TrimSpace(input)
		var num int
		_, err = fmt.Sscanf(input, "%d", &num)
		if err != nil || num < 1 || num > len(items) {
			fmt.Printf("Invalid input. Please enter a number between 1 and %d\n", len(items))
			continue
		}

		return items[num-1], nil
	}
}

// AssistConfigForPool helps fill in project and pool configuration
func AssistConfigForPool(cfg *config) error {
	if cfg.projectID == "" {
		projects, err := ListProjects()
		if err != nil {
			return err
		}
		cfg.projectID, err = SelectFromList(projects, "projects")
		if err != nil {
			return err
		}
	}

	if cfg.poolName == "" {
		for {
			cfg.poolName = GetInput("Enter new pool name (only letters, numbers, and hyphens allowed):")
			if matched := strings.ContainsFunc(cfg.poolName, func(r rune) bool {
				return !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-", r)
			}); matched {
				fmt.Println("Invalid pool name. Only letters, numbers, and hyphens are allowed.")
				continue
			}
			break
		}
	}

	return nil
}

// AssistConfigForProvider helps fill in project, pool, and provider configuration
func AssistConfigForProvider(cfg *config) error {
	if cfg.projectID == "" {
		projects, err := ListProjects()
		if err != nil {
			return err
		}
		cfg.projectID, err = SelectFromList(projects, "projects")
		if err != nil {
			return err
		}
	}

	if cfg.poolName == "" {
		pools, err := ListPools(cfg.projectID)
		if err != nil {
			return err
		}
		if len(pools) > 0 {
			cfg.poolName, err = SelectFromList(pools, "workload identity pools")
			if err != nil {
				return err
			}
		}
	}

	if cfg.providerName == "" {
		providers, err := ListProviders(cfg.projectID, cfg.poolName)
		if err != nil {
			return err
		}
		if len(providers) > 0 {
			cfg.providerName, err = SelectFromList(providers, "providers")
			if err != nil {
				return err
			}
		}
	}

	if cfg.providerName == "" {
		for {
			cfg.providerName = GetInput("Enter new provider name (only letters, numbers, and hyphens allowed):")
			if matched := strings.ContainsFunc(cfg.providerName, func(r rune) bool {
				return !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-", r)
			}); matched {
				fmt.Println("Invalid provider name. Only letters, numbers, and hyphens are allowed.")
				continue
			}
			break
		}
	}

	return nil
}

// AssistConfigForAuth helps fill in all configuration needed for binding
func AssistConfigForAuth(cfg *config) error {
	if cfg.projectID == "" {
		projects, err := ListProjects()
		if err != nil {
			return err
		}
		cfg.projectID, err = SelectFromList(projects, "projects")
		if err != nil {
			return err
		}
	}

	if cfg.serviceAccount == "" {
		accounts, err := ListServiceAccounts(cfg.projectID)
		if err != nil {
			return err
		}
		cfg.serviceAccount, err = SelectFromList(accounts, "service accounts")
		if err != nil {
			return err
		}
	}

	return AssistConfigForProvider(cfg)
}
