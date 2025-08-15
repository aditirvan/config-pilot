#!/bin/bash

# Build script for Config Pilot .deb package

set -e

# Clean previous builds
rm -rf deb-package/usr/bin/config-pilot
rm -f config-pilot_*.deb

# Build the Go binary
echo "Building Go binary..."
go build -o deb-package/usr/bin/config-pilot cmd/monitor/main.go

# Set proper permissions
chmod +x deb-package/usr/bin/config-pilot

# Build the .deb package
echo "Building .deb package..."
dpkg-deb --build deb-package

# Rename the package with version
mv deb-package.deb config-pilot_1.0.0_amd64.deb

echo "Build complete!"
echo "Package: config-pilot_1.0.0_amd64.deb"
echo "Install with: sudo dpkg -i config-pilot_1.0.0_amd64.deb"
