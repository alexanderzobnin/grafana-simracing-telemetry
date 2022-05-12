#!/bin/bash

# Exit script if a statement returns a non-true return value.
set -o errexit
# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

GRAFANA_VERSION='8.5.2'
GRAFANA_DOWNLOAD_URL="https://dl.grafana.com/oss/release/grafana-${GRAFANA_VERSION}.windows-amd64.zip"
GRAFANA_PACKAGE_NAME="grafana-${GRAFANA_VERSION}.windows-amd64.zip"

# Clean up
echo -e "${GREEN}Clean up${NC}"
rm -rf ci/grafana

# Download Grafana package
echo -e "${GREEN}Downloading Grafana ${GRAFANA_VERSION} zip package${NC}"
mkdir -p ci/grafana
curl -L -o "ci/grafana/${GRAFANA_PACKAGE_NAME}" ${GRAFANA_DOWNLOAD_URL}

# Extract files
echo -e "${GREEN}Extracting files${NC}"
cd ci/grafana
unzip ${GRAFANA_PACKAGE_NAME}

# Copy plugin into Grafana package
echo -e "${GREEN}Copy plugin into grafana package${NC}"
mkdir -p "grafana-${GRAFANA_VERSION}/data/plugins"
cp -r ../dist/alexanderzobnin-simracingtelemetry-datasource "./grafana-${GRAFANA_VERSION}/data/plugins"

# Copy configs to the Grafana
echo -e "${GREEN}Copy configs into grafana package${NC}"
cp ../../conf/custom.ini "./grafana-${GRAFANA_VERSION}/conf/"
cp ../../conf/datasources.yaml "./grafana-${GRAFANA_VERSION}/conf/provisioning/datasources/"
cp ../../conf/dashboards.yaml "./grafana-${GRAFANA_VERSION}/conf/provisioning/dashboards/"

# Make zip
echo -e "${GREEN}Packaging Grafana into zip${NC}"
zip -r "grafana-${GRAFANA_VERSION}-bundled.windows-amd64.zip" "grafana-${GRAFANA_VERSION}/"

echo -e "${GREEN}Packaged Grafana located in ${BLUE}ci/grafana${NC}"
ls -lh ./
