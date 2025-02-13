# gwif

A CLI tool to configure GitHub Actions to use Workload Identity Federation.

## Usage

```bash
gwif pools create
gwif providers create
gwif auth
gwif yaml
```

## Installation

```bash
go install github.com/unacast/gwif@latest
```
 
## Dependencies
### Google Cloud SDK
1. Install the Google Cloud SDK
2. Configure gcloud
```bash
gcloud config configurations create <name>
gcloud config set project <project-id>
gcloud auth login
```

### Permissions
- IAM Admin
- Workload Identity Federation Admin
