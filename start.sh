#!/bin/bash
set -e

DEFAULT_ARCH="linux"

read -r -p "Please input port (default: 5000): " port
if [[ -z "$port" ]]; then
  port=5000
fi

./bin/$DEFAULT_ARCH/import-tools -p $port
