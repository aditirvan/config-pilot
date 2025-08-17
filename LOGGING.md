# Logging Configuration

This document describes how to configure logging for the config-pilot application.

## Configuration Options

### Environment Variables

You can configure logging using the following environment variables:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `LOG_FILE_PATH` | Path to the log file | (none) | `/var/log/config-pilot/app.log` |
| `LOG_LEVEL` | Logging level | `info` | `debug`, `info`, `warn`, `error` |
| `LOG_TO_FILE` | Enable file logging | `false` | `true`, `1` |

### YAML Configuration

You can also configure logging in your `config.yaml` file:

```yaml
logging:
  logFilePath: "/var/log/config-pilot/app.log"
  logLevel: "info"
  logToFile: false
```

## Usage Examples

### Using Environment Variables

```bash
# Enable file logging
export LOG_TO_FILE=true
export LOG_FILE_PATH="/var/log/config-pilot/app.log"
export LOG_LEVEL=debug

# Run the application
./config-pilot
```

### Using YAML Configuration

```yaml
# config.yaml
githubToken: ghp_xxxxxxxxxxxxx
owner: github-username
repo: your-repository
monitorPath: path-directory-to-monitor
ageSecret: sops age private key
interval: 60
logging:
  logFilePath: "/var/log/config-pilot/app.log"
  logLevel: "debug"
  logToFile: true
```

### Priority Order

Configuration is applied in the following order (later configurations override earlier ones):

1. YAML configuration in `config.yaml`
2. Environment variables

## Log Levels

- `debug`: Detailed debug information
- `info`: General information messages
- `warn`: Warning messages
- `error`: Error messages

## Log Format

Logs are written in structured text format with the following fields:
- Timestamp
- Log level
- Message
- Additional context (when provided)

## File Rotation

The application appends to the log file and does not handle log rotation automatically. For production use, consider using external log rotation tools like `logrotate`.

## Examples

### Basic Usage (stdout only)
```bash
./config-pilot
```

### File Logging
```bash
LOG_TO_FILE=true LOG_FILE_PATH="./logs/app.log" ./config-pilot
```

### Debug Mode with File Logging
```bash
LOG_LEVEL=debug LOG_TO_FILE=true LOG_FILE_PATH="./logs/debug.log" ./config-pilot
