#!/bin/bash

docker compose -f {{ InstallDir }}/authelia/docker-compose.yml up -d

# Include prompt to take you to the setup page