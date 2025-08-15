# GitOps Source Automation

A Go-based GitOps automation tool that monitors GitHub repositories for changes and automatically pulls updates, decrypts encrypted files using SOPS, and executes custom scripts.

## Features

- **Continuous Monitoring**: Monitors GitHub repositories for new commits
- **Automatic Updates**: Automatically pulls latest changes when new commits are detected
- **File Decryption**: Decrypts SOPS-encrypted files using age keys
- **Script Execution**: Runs custom scripts after updates
- **Configurable**: All settings managed through a single YAML configuration file
- **Logging**: Comprehensive logging with structured output

## Prerequisites

- **Go 1.21+** (for building from source)
- **Git** (for repository operations)
- **SOPS** (for file decryption)
- **Age** (for encryption keys)
- **GitHub Personal Access Token** (for repository access)

## Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/aditirvan/config-pilot.git
cd config-pilot

# Install dependencies
go mod download

# Build the application
go build -o gitops-monitor cmd/monitor/main.go
```

### Option 2: Download Binary

Download the latest release from the [releases page](https://github.com/aditirvan/config-pilot/releases).

## Configuration

### Configuration File Location

The application supports two ways to specify the configuration file location:

1. **Environment Variable (Recommended)**: Set the `CONFIG_PATH` environment variable
2. **Default Location**: Uses `config.yaml` in the current working directory

#### Using Environment Variable
```bash
# Set config file path via environment variable
export CONFIG_PATH="/path/to/your/config.yaml"
./gitops-monitor

# Or inline
CONFIG_PATH="/path/to/your/config.yaml" ./gitops-monitor

# Windows
set CONFIG_PATH=C:\path\to\your\config.yaml
gitops-monitor.exe
```

#### Using Default Location
Simply create a `config.yaml` file in the project root directory where you run the application.

### Configuration File Format

Create your configuration file with the following structure:

```yaml
# GitHub Configuration
githubToken: ghp_your_personal_access_token_here
owner: your-github-username
repo: your-repository-name

# Monitoring Configuration
interval: 30  # Check interval in seconds
monitorPath: path/to/monitor  # Optional: specific directory to monitor (empty for entire repo)

# Security Configuration
ageSecret: AGE-SECRET-KEY-your-age-private-key-here

# Script Configuration
script: |
  # Your custom script here
  echo "Running deployment script..."
  kubectl apply -f .
  echo "Deployment completed"
```

### Configuration Parameters

| Parameter | Description | Required | Example |
|-----------|-------------|----------|---------|
| `githubToken` | GitHub Personal Access Token | Yes | `ghp_xxxxxxxxxxxxx` |
| `owner` | GitHub repository owner | Yes | `myusername` |
| `repo` | Repository name | Yes | `my-config-repo` |
| `interval` | Check interval in seconds | No | `30` |
| `monitorPath` | Specific directory to monitor | No | `kubernetes/manifests` |
| `ageSecret` | Age private key for SOPS decryption | Yes | `AGE-SECRET-KEY-...` |
| `script` | Custom script to execute after updates | Yes | See example above |

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `CONFIG_PATH` | Path to configuration file | `config.yaml` | `/etc/gitops/config.yaml` |

### GitHub Token Setup

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Click "Generate new token (classic)"
3. Select scopes:
   - `repo` (Full control of private repositories)
   - `public_repo` (Access public repositories)
4. Copy the generated token to your `config.yaml`

### Age Key Setup

1. Install age: `brew install age` (macOS) or `apt install age` (Ubuntu)
2. Generate age key: `age-keygen -o key.txt`
3. Copy the private key from `key.txt` to your `config.yaml`
4. Share the public key with your team for encryption

## Usage

### Quick Start

1. **Create configuration file**:
   ```bash
   cp config.yaml.sample config.yaml
   # Edit config.yaml with your settings
   ```

2. **Run the application**:
   ```bash
   # If built from source
   ./gitops-monitor

   # Or using go run
   go run cmd/monitor/main.go
   ```

### Using Custom Configuration Path

You can specify a custom configuration file location using the `CONFIG_PATH` environment variable:

```bash
# Use custom config file
export CONFIG_PATH="/etc/gitops/production-config.yaml"
./gitops-monitor

# Or use with systemd service
CONFIG_PATH="/opt/gitops/configs/staging.yaml" ./gitops-monitor

# Docker example
docker run -e CONFIG_PATH=/app/config/production.yaml -v /host/config:/app/config gitops-monitor
```

### Running with Different Configurations

```bash
# Development
CONFIG_PATH="./configs/dev.yaml" ./gitops-monitor

# Staging
CONFIG_PATH="./configs/staging.yaml" ./gitops-monitor

# Production
CONFIG_PATH="/etc/gitops/production.yaml" ./gitops-monitor
```

### Example Configuration Files

#### Basic Configuration
```yaml
githubToken: ghp_1234567890abcdef
owner: myorg
repo: infrastructure-config
interval: 60
monitorPath: kubernetes
ageSecret: AGE-SECRET-KEY-1U5A44X09L2U4CL2QFPDTX3GQAFXLJCLWT2YRCSXDWASEMSE5T3PS2NGFAW
script: |
  kubectl apply -f kubernetes/
```

#### Advanced Configuration
```yaml
githubToken: ghp_1234567890abcdef
owner: myorg
repo: multi-env-config
interval: 30
monitorPath: ""
ageSecret: AGE-SECRET-KEY-1U5A44X09L2U4CL2QFPDTX3GQAFXLJCLWT2YRCSXDWASEMSE5T3PS2NGFAW
script: |
  #!/bin/bash
  set -e
  
  echo "Starting deployment process..."
  
  # Decrypt and apply staging
  kubectl config use-context staging
  kubectl apply -f staging/
  
  # Decrypt and apply production
  kubectl config use-context production
  kubectl apply -f production/
  
  echo "Deployment completed successfully"
```

## Directory Structure

```
config-pilot/
├── cmd/
│   └── monitor/
│       └── main.go          # Main application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── github/              # GitHub API client
│   ├── monitor/             # Monitoring service
│   └── utils/               # Utility functions
├── data/                    # Working directory for cloned repos
├── config.yaml              # Your configuration file
├── config.yaml.sample       # Sample configuration
└── README.md               # This file
```

## How It Works

1. **Initialization**: The application loads configuration from `config.yaml`
2. **GitHub Connection**: Establishes connection to the specified GitHub repository
3. **Initial Sync**: Pulls the latest commit to establish baseline
4. **Continuous Monitoring**: Checks for new commits at the configured interval
5. **Update Detection**: When new commits are detected:
   - Clones the repository to the local `data/` directory
   - Decrypts any SOPS-encrypted files using the provided age key
   - Executes the custom script defined in the configuration
6. **Logging**: All operations are logged with timestamps and details

## Logging

The application provides comprehensive logging with different log levels:
- **INFO**: General operation information
- **ERROR**: Error messages and stack traces
- **DEBUG**: Detailed debugging information (when enabled)

Logs include:
- Application startup/shutdown
- GitHub API interactions
- File operations
- Script execution output
- Error details

## Troubleshooting

### Common Issues

#### "githubToken must be specified in config.yaml"
**Solution**: Ensure your `config.yaml` contains a valid GitHub token:
```yaml
githubToken: ghp_your_actual_token_here
```

#### "git clone failed"
**Solution**: 
- Verify your GitHub token has the correct permissions
- Check if the repository exists and is accessible
- Ensure the repository owner and name are correct

#### "SOPS decryption failed"
**Solution**:
- Verify your age key is correct
- Ensure files are encrypted with the corresponding public key
- Check file permissions

#### "Script execution failed"
**Solution**:
- Check script syntax
- Ensure required tools (kubectl, helm, etc.) are installed
- Verify file paths in the script

### Debug Mode

To enable debug logging, you can modify the logger configuration in the source code or use environment variables if supported.

## Security Considerations

- **GitHub Token**: Store securely and use minimal required permissions
- **Age Keys**: Keep private keys secure and never commit to version control
- **Configuration**: Ensure `config.yaml` is in `.gitignore` to prevent accidental commits
- **Scripts**: Review scripts for security implications before execution

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests if applicable
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- Create an issue on [GitHub Issues](https://github.com/aditirvan/config-pilot/issues)
- Check existing issues for solutions
- Review the troubleshooting section above
