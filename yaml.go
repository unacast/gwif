package main

import "fmt"

func DumpYAML(cfg *config, projectNumber string) {
	fmt.Println(`
 <-- in your workspace -->
 permissions:
  id-token: write # This is required for requesting the JWT from GCP Workload Identity

 <-- in your job steps -->
      - name: 'Authenticate to Google Cloud'
        uses: 'google-github-actions/auth@v2'
        with:
          project_id: '` + cfg.projectID + `'
          workload_identity_provider: 'projects/` + projectNumber + `/locations/global/workloadIdentityPools/` + cfg.poolName + `/providers/` + cfg.providerName + `'
          service_account: '` + cfg.serviceAccount + `'
          access_token_lifetime: '300s' # optional, default: '3600s' (1 hour)
      - name: 'Set up Cloud SDK'
        uses: 'google-github-actions/setup-gcloud@v2'
`)
}
