#!/bin/bash
set -e

read -r -p "Please input port (default: 5000): " port
if [[ -z "$port" ]]; then
  port=5000
fi

./import-tools -p $port
