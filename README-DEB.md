# Config Pilot .deb Package

## Overview
This is a simple .deb package for Config Pilot - a GitOps configuration monitoring tool.

## Package Structure
```
deb-package/
├── DEBIAN/
│   ├── control          # Package metadata
│   ├── postinst         # Post-installation script
│   └── prerm            # Pre-removal script
├── usr/bin/
│   └── config-pilot     # Main binary
├── etc/config-pilot/
│   └── config.yaml      # Configuration file
└── lib/systemd/system/
    └── config-pilot.service  # Systemd service
```

## Build Instructions

1. **Install dependencies**:
   ```bash
   sudo apt-get update
   sudo apt-get install build-essential golang-go
   ```

2. **Build the package**:
   ```bash
   ./build-deb.sh
   ```

3. **Install the package**:
   ```bash
   sudo dpkg -i config-pilot_1.0.0_amd64.deb
   ```

## Usage

1. **Configure**:
   Edit `/etc/config-pilot/config.yaml` with your GitHub repository details.

2. **Start service**:
   ```bash
   sudo systemctl enable config-pilot
   sudo systemctl start config-pilot
   ```

3. **Check status**:
   ```bash
   sudo systemctl status config-pilot
   ```

## Files Created
- `config-pilot_1.0.0_amd64.deb` - The final .deb package
- `build-deb.sh` - Build script
- `deb-package/` - Package structure directory

## Manual Build (if needed)
```bash
# Build binary
go build -o deb-package/usr/bin/config-pilot cmd/monitor/main.go

# Build .deb
dpkg-deb --build deb-package
```

## Uninstall
```bash
sudo apt-get remove config-pilot
