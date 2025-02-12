package main

import (
	"fmt"
	"os"
	"os/exec"
)

func CreatePool(cfg *config) error {
	// Check if pool exists
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "describe",
		cfg.poolName,
		"--project", cfg.projectID,
		"--location", "global",
		"--format", "value(name)")

	if err := cmd.Run(); err == nil {
		fmt.Printf("Pool %s already exists... skipping\n", cfg.poolName)
		return nil
	}

	if !Ask("Create pool (" + cfg.poolName + ")?") {
		return fmt.Errorf("cannot continue without a pool")
	}

	cmd = exec.Command("gcloud", "iam", "workload-identity-pools", "create",
		cfg.poolName,
		"--project", cfg.projectID,
		"--location", "global",
		"--display-name", cfg.poolName)

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create pool: %v", err)
	}

	return nil
}

func CreateProvider(cfg *config, projectNumber, githubRepositoryFullName string) error {
	// Check if provider exists
	cmd := exec.Command("gcloud", "iam", "workload-identity-pools", "providers", "describe",
		cfg.providerName,
		"--project", cfg.projectID,
		"--location", "global",
		"--workload-identity-pool", cfg.poolName,
		"--format", "value(name)")

	if err := cmd.Run(); err == nil {
		fmt.Printf("Provider %s already exists... skipping\n", cfg.providerName)
		return nil
	}

	if !Ask("Create provider (" + cfg.providerName + ")?") {
		return fmt.Errorf("cannot continue without a provider")
	}

	fmt.Printf(`
|---------------------------------------------------------------------------------------|
|                                 Provider Conditions                                   |
|    You can apply conditions at the provider level so that you can later associate     |
|                service accounts based on some other mapped attribute.                 |
|                      e.g. apply repository condition to provider                      |
|                                        then                                           |
|                     associate service account by workflow or branch                   |
|---------------------------------------------------------------------------------------|

NOTE: Only one attribute can be used for service account assignment
`)

	audience := fmt.Sprintf("'https://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/providers/%s'",
		projectNumber, cfg.poolName, cfg.providerName)

	attributeMapping := fmt.Sprintf("google.subject=assertion.sub,"+
		"attribute.aud=%s,"+
		"attribute.actor=assertion.actor,"+
		"attribute.repository=assertion.repository,"+
		"attribute.environment=assertion.environment,"+
		"attribute.workflow=assertion.workflow_ref.split('.github/workflows/')[1].split('.')[0].split('@')[0],"+
		"attribute.ref=assertion.ref",
		audience)

	attributeCondition := fmt.Sprintf("assertion.repository_owner=='%s'", cfg.githubRepositoryOwner)

	if cfg.unsafe {
		if Ask("Apply repository condition to the provider?") {
			attributeCondition = fmt.Sprintf("%s && assertion.repository=='%s'", attributeCondition, githubRepositoryFullName)
		} else {
			fmt.Println("WARNING: Not applying repository condition to the provider - MUST use repository full name to associate the service account e.g. owner/repo")
			if !RequiredAsk("Have you read the warning?", "It is critical to use the repository for service account association if not using repository condition") {
				return fmt.Errorf("user declined to acknowledge warning")
			}
		}
	} else {
		attributeCondition = fmt.Sprintf("%s && assertion.repository=='%s'", attributeCondition, githubRepositoryFullName)
	}

	if Ask("Apply workflow condition to the provider?") {
		workflow := GetInput("Enter your workflow name:")
		attributeCondition = fmt.Sprintf("%s && assertion.workflow=='%s'", attributeCondition, workflow)
	}

	if Ask("Apply environment condition to the provider?") {
		env := GetInput("Enter your environment name:")
		attributeCondition = fmt.Sprintf("%s && assertion.environment=='%s'", attributeCondition, env)
	}

	if Ask("Apply branch condition to the provider?") {
		branch := GetInput("Enter your branch name:")
		attributeCondition = fmt.Sprintf("%s && assertion.ref=='refs/heads/%s'", attributeCondition, branch)
	} else {
		fmt.Println("INFO: When associating service account by branch, use the format 'refs/heads/branch-name'")
	}

	cmd = exec.Command("gcloud", "iam", "workload-identity-pools", "providers", "create-oidc",
		cfg.providerName,
		"--project", cfg.projectID,
		"--location", "global",
		"--workload-identity-pool", cfg.poolName,
		"--display-name", cfg.providerName,
		"--attribute-mapping", attributeMapping,
		"--attribute-condition", attributeCondition,
		"--issuer-uri", "https://token.actions.githubusercontent.com")

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create provider: %v", err)
	}

	return nil
}
