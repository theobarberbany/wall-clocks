#!/usr/bin/env bash
set -euo pipefail

echo "Tag:" $1

make docker-build IMG=controller:$1
docker image tag controller:latest controller:$1
kind load docker-image controller:$1
make deploy IMG=controller:$1
