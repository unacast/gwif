package main

import (
	"fmt"
	"os"
	"os/exec"
)

func AuthServiceAccount(cfg *config, projectNumber string) error {
	fmt.Printf(`

|--------------------------------------------------------------------------------------|
|                                 Service Accounts                                     |
|    Using an attribute from the generated JWT we can select which service account     |
|    to associate with the GitHub job auth request, and therefore limiting access      |
|    to the permissions associated with that account.                                  |
|                                                                                      |
|        For best security, create a unique service account for each workflow.         |
|--------------------------------------------------------------------------------------|

`)
	if cfg.serviceAccount == "" {
		cfg.serviceAccount = GetInput("Paste the service account email address (e.g. deploy-sa@project-id.iam.gserviceaccount.com):")
	}
	fmt.Println()
	fmt.Printf(`Select the attribute to use for service account association:
1. workflow [SUGGESTED]
2. repository
3. environment
4. actor
5. ref
`)

	attributeNum := GetInput("Enter number (1-5):")
	var attribute string
	switch attributeNum {
	case "1":
		attribute = "workflow"
	case "2":
		attribute = "repository"
	case "3":
		attribute = "environment"
	case "4":
		attribute = "actor"
	case "5":
		attribute = "ref"
	default:
		return fmt.Errorf("invalid selection: %s", attributeNum)
	}

	fmt.Println()
	switch attribute {
	case "repository":
		fmt.Println("Expected format for [repository]: owner/repo")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("- unacast/actions")
		fmt.Println("- redis/go-redis")
	case "environment":
		fmt.Println("Expected format for [environment]: env")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("- dev")
		fmt.Println("- prod")
	case "workflow":
		fmt.Println("Expected format for [workflow]: workflow-filename (without .yml)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("- build")
		fmt.Println("- deploy")
	case "actor":
		fmt.Println("Expected format for [actor]: username")
	case "ref":
		fmt.Println("Expected format for [ref]: refs/heads/branch-name")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("- refs/heads/main")
		fmt.Println("- refs/heads/feature-branch")
	}

	value := GetInput("Enter value [CASE SENSITIVE]:")

	cmd := exec.Command("gcloud", "iam", "service-accounts", "add-iam-policy-binding",
		cfg.serviceAccount,
		"--project", cfg.projectID,
		"--role", "roles/iam.workloadIdentityUser",
		"--member", fmt.Sprintf("principalSet://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/attribute.%s/%s",
			projectNumber, cfg.poolName, attribute, value))

	cmd.Stderr = os.Stderr
	return cmd.Run()
}
