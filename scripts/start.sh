#!/bin/bash
set -e

DEFAULT_ARCH="${DEFAULT_ARCH:-linux}"
WORKSPACE=$(pwd -P)

function init_workspace() {
  if [[ "$WORKSPACE" = */scripts ]]; then
    WORKSPACE="$(dirname "$WORKSPACE")"
  fi
}

function read_params() {
  read -r -p "Please input port (default: 5000): " port
  if [[ -z "$port" ]]; then
    port=5000
  fi

  read -r -p "Whether to use a shared disk, yes or no (default: no): " use_shared_disk
  if [[ "$use_shared_disk" = "yes" ]]; then
    while [[ -z "$shared_disk_path" ]]; do
      read -r -p "Please input the path of the shared disk: " shared_disk_path
      if [[ -z "$shared_disk_path" ]]; then
        echo "The path of the shared disk cannot be empty!"
      fi
    done
  fi
}

function start() {
  if [[ -z "$shared_disk_path" ]]; then
    $WORKSPACE/bin/$DEFAULT_ARCH/import-tools \
      -port $port
  else
    $WORKSPACE/bin/$DEFAULT_ARCH/import-tools \
      -port $port \
      -shared-disk-path $shared_disk_path
  fi
}

function main() {
  init_workspace
  read_params
  start
}

main "$@"
