#!/usr/bin/env bash
set -euxo pipefail

sudo -u pronouns git pull
sudo -u pronouns make all

while getopts 'm' OPTION; do
  case "$OPTION" in 
    m)
      sudo systemctl stop pronouns-api pronouns-fe pronouns-exporter
      sudo -u pronouns ./pronouns database migrate
      sudo systemctl start pronouns-api pronouns-fe pronouns-exporter
      ;;
    ?)
      sudo systemctl restart pronouns-api pronouns-fe
      ;;
  esac
done
