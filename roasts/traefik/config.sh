#!/bin/bash

docker network inspect traefik_public >/dev/null 2>&1 || docker network create --driver bridge traefik_public

docker compose -f {{ InstallDir }}/traefik/docker-compose.yml up -d
