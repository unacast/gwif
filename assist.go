package main

import (
	"fmt"
	"strings"
)

func AssistConfigForRoot(cfg *config) error {
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

	if err := verifyActiveProject(cfg.projectID); err != nil {
		return err
	}

	return nil
}

func AssistConfigForPoolDelete(cfg *config) error {
	if err := AssistConfigForRoot(cfg); err != nil {
		return err
	}

	if cfg.poolName == "" {
		pools, err := ListPools(cfg.projectID, false)
		if err != nil {
			return err
		}
		if len(pools) == 0 {
			return fmt.Errorf("no workload identity pools found in project %s", cfg.projectID)
		}
		cfg.poolName, err = SelectFromList(pools, "workload identity pools")
		if err != nil {
			return err
		}
	}

	return nil
}

func AssistConfigForPoolCreate(cfg *config) error {
	if err := AssistConfigForRoot(cfg); err != nil {
		return err
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

func AssistConfigForProviderSubcommand(cfg *config) error {
	if err := AssistConfigForRoot(cfg); err != nil {
		return err
	}

	if cfg.poolName == "" {
		pools, err := ListPools(cfg.projectID, false)
		if err != nil {
			return err
		}
		if len(pools) == 0 {
			return fmt.Errorf("no workload identity pools found in project %s", cfg.projectID)
		}
		cfg.poolName, err = SelectFromList(pools, "workload identity pools")
		if err != nil {
			return err
		}
	}
	return nil
}

// AssistConfigForProviderCreate helps fill in project, pool, and provider configuration
func AssistConfigForProviderCreate(cfg *config) error {
	if err := AssistConfigForProviderSubcommand(cfg); err != nil {
		return err
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

	return AssistGithub(cfg)
}

func AssistConfigForProviderDelete(cfg *config) error {
	if err := AssistConfigForProviderSubcommand(cfg); err != nil {
		return err
	}

	if cfg.providerName == "" {
		providers, err := ListProviders(cfg.projectID, cfg.poolName, false)
		if err != nil {
			return err
		}

		if len(providers) == 0 {
			return fmt.Errorf("no workload identity providers found in pool %s", cfg.poolName)
		}
		cfg.providerName, err = SelectFromList(providers, "workload identity providers")
		if err != nil {
			return err
		}
	}
	return nil
}

func AssistConfigForProviderRestore(cfg *config) error {
	if err := AssistConfigForProviderSubcommand(cfg); err != nil {
		return err
	}

	if cfg.providerName == "" {
		providers, err := ListProviders(cfg.projectID, cfg.poolName, true)
		if err != nil {
			return err
		}

		if len(providers) == 0 {
			return fmt.Errorf("no deleted workload identity providers found in pool %s", cfg.poolName)
		}
		cfg.providerName, err = SelectFromList(providers, "workload identity providers")
		if err != nil {
			return err
		}
	}
	return nil
}

func AssistConfigForAuth(cfg *config) error {
	if err := AssistConfigForRoot(cfg); err != nil {
		return err
	}

	if cfg.poolName == "" {
		pools, err := ListPools(cfg.projectID, false)
		if err != nil {
			return err
		}
		if len(pools) == 0 {
			return fmt.Errorf("no workload identity pools found in project %s", cfg.projectID)
		}
		cfg.poolName, err = SelectFromList(pools, "workload identity pools")
		if err != nil {
			return err
		}
	}

	if cfg.providerName == "" {
		providers, err := ListProviders(cfg.projectID, cfg.poolName, false)
		if err != nil {
			return err
		}

		if len(providers) == 0 {
			return fmt.Errorf("no workload identity providers found in pool %s", cfg.poolName)
		}
		cfg.providerName, err = SelectFromList(providers, "workload identity providers")
		if err != nil {
			return err
		}
	}

	if cfg.serviceAccount == "" {
		accounts, err := ListServiceAccounts(cfg.projectID)
		if err != nil {
			return err
		}
		if len(accounts) == 0 {
			return fmt.Errorf("no service accounts found in project %s", cfg.projectID)
		}
		cfg.serviceAccount, err = SelectFromList(accounts, "service accounts")
		if err != nil {
			return err
		}
	}

	return nil
}

func AssistConfigForYaml(cfg *config) error {
	if err := AssistConfigForRoot(cfg); err != nil {
		return err
	}

	if cfg.poolName == "" {
		pools, err := ListPools(cfg.projectID, false)
		if err != nil {
			return err
		}
		if len(pools) == 0 {
			return fmt.Errorf("no workload identity pools found in project %s", cfg.projectID)
		}
		cfg.poolName, err = SelectFromList(pools, "workload identity pools")
		if err != nil {
			return err
		}
	}

	if cfg.providerName == "" {
		providers, err := ListProviders(cfg.projectID, cfg.poolName, false)
		if err != nil {
			return err
		}

		if len(providers) == 0 {
			return fmt.Errorf("no workload identity providers found in pool %s", cfg.poolName)
		}
		cfg.providerName, err = SelectFromList(providers, "workload identity providers")
		if err != nil {
			return err
		}
	}

	if cfg.serviceAccount == "" {
		accounts, err := ListServiceAccounts(cfg.projectID)
		if err != nil {
			return err
		}
		if len(accounts) == 0 {
			return fmt.Errorf("no service accounts found in project %s", cfg.projectID)
		}
		cfg.serviceAccount, err = SelectFromList(accounts, "service accounts")
		if err != nil {
			return err
		}
	}

	return nil
}

func AssistGithub(cfg *config) error {
	if cfg.githubRepositoryOwner == "" {
		for {
			cfg.githubRepositoryOwner = GetInput("Enter GitHub repository owner [CASE SENSITIVE]:")
			if cfg.githubRepositoryOwner == "" {
				fmt.Println("GitHub repository owner is required.")
				continue
			}
			break
		}
	}

	if cfg.githubRepository == "" {
		for {
			cfg.githubRepository = GetInput("Enter GitHub repository name [CASE SENSITIVE]:")
			if cfg.githubRepository == "" {
				fmt.Println("GitHub repository name is required.")
				continue
			}
			break
		}
	}

	return nil
}
