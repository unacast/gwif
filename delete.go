package main

import (
	"fmt"
	"os"
	"os/exec"
)

func DeletePool(cfg *config) error {
	if Ask(fmt.Sprintf("Are you sure you want to delete the pool [%s > %s]?", cfg.projectID, cfg.poolName)) {
		cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "delete", cfg.poolName, "--project", cfg.projectID, "--location", "global", "--quiet")
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to delete pool: %v", err)
		}
		fmt.Println("Pool deleted successfully - it will be removed after a 30 day grace period and can be restored until then.")
	}
	return nil
}

func DeleteProvider(cfg *config) error {
	if Ask(fmt.Sprintf("Are you sure you want to delete the provider [%s > %s > %s]?", cfg.projectID, cfg.poolName, cfg.providerName)) {
		cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "providers", "delete",
			cfg.providerName,
			"--project", cfg.projectID,
			"--location", "global",
			"--quiet",
			"--workload-identity-pool", cfg.poolName)

		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to delete provider: %v", err)
		}

		fmt.Println("Provider deleted successfully - it will be removed after a 30 day grace period and can be restored until then.")
	}
	return nil
}

func RestoreProvider(cfg *config) error {
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "providers", "undelete",
		cfg.providerName,
		"--project", cfg.projectID,
		"--location", "global",
		"--workload-identity-pool", cfg.poolName)

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore provider: %v", err)
	}
	return nil
}
