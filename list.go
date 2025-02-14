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
func ListPools(projectID string, showDeleted bool) ([]string, error) {
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "list",
		"--project", projectID,
		"--location", "global",
		"--format", "value(name)")
	if showDeleted {
		cmd.Args = append(cmd.Args, "--show-deleted")
		cmd.Args = append(cmd.Args, "--filter", "state:DELETED")
	}
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

// ListProviders returns a list of workload identity providers for a given project and pool
func ListProviders(projectID, poolName string, showDeleted bool) ([]string, error) {
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "providers", "list",
		"--project", projectID,
		"--location", "global",
		"--workload-identity-pool", poolName,
		"--format", "value(name)")
	if showDeleted {
		cmd.Args = append(cmd.Args, "--show-deleted")
		cmd.Args = append(cmd.Args, "--filter", "state:DELETED")
	}
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

// SelectFromList presents a numbered list to the user and returns their selection
func SelectFromList(items []string, resourceType string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no %s available", resourceType)
	}

	fmt.Printf("\nAvailable %s:\n", resourceType)
	for i, item := range items {
		fmt.Printf("%d) %s\n", i+1, item)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\nEnter number (1-%d) to select %s: ", len(items), resourceType)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read input: %v", err)
		}
		return "", fmt.Errorf("no input provided")
	}
	input := scanner.Text()

	var num int
	_, err := fmt.Sscanf(input, "%d", &num)
	if err != nil || num < 1 || num > len(items) {
		fmt.Printf("Invalid input. Please enter a number between 1 and %d\n", len(items))
		return SelectFromList(items, resourceType)
	}

	return items[num-1], nil
}
