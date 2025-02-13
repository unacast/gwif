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

	fmt.Printf(`
|---------------------------------------------------------------------------------------|
|                                 Provider Conditions                                   |
|    You can apply conditions to further restrict the audience of the provider.         |
|    you can later associate service accounts based on some other mapped attribute.     |
|                      e.g. apply a branch condition to provider                        |
|                                        then                                           |
|                        associate service account by workflow                          |
|---------------------------------------------------------------------------------------|

NOTE: Only one attribute can be used for service account assignment

RECOMMENDATION: 
- Be as specific as possible with conditions to improve security e.g. If your CI only runs on main,
  then apply a branch condition to the provider.
- Associate the service account with the workflow name. e.g. deploy -> github-deploy@...com
- Create a service account for each workflow and assign minimum permissions to it to run the workflow.
- Create a provider for each repository sharing the Google Cloud project.

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

	if Ask("[NOT RECOMMENDED] Apply workflow condition to the provider?") {
		workflow := GetInput("Enter your workflow name:")
		attributeCondition = fmt.Sprintf("%s && assertion.workflow=='%s'", attributeCondition, workflow)
	}

	if Ask("[PROBABLY NOT NEEDED] Apply environment condition to the provider?") {
		env := GetInput("Enter your environment name:")
		attributeCondition = fmt.Sprintf("%s && assertion.environment=='%s'", attributeCondition, env)
	}

	if Ask("[PROBABLY NOT NEEDED] Apply branch condition to the provider?") {
		branch := GetInput("Enter your branch name:")
		attributeCondition = fmt.Sprintf("%s && assertion.ref=='refs/heads/%s'", attributeCondition, branch)
	}

	if !Ask("Create provider (" + cfg.providerName + ")?") {
		return fmt.Errorf("cannot continue without a provider")
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
