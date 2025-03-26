#!/bin/bash

export PROJECT_NAME=go-zulip
export DEV_PATH=$(pwd)

docker compose -p ${PROJECT_NAME} \
-f docker-compose-dev.yml \
-f docker-compose-zulip.yml \
build --pull

docker compose -p ${PROJECT_NAME} \
-f docker-compose-dev-env.yml \
-f docker-compose-dev.yml \
-f docker-compose-zulip.yml \
up --detach

docker compose -p ${PROJECT_NAME} \
logs --follow --tail 100

docker compose -p ${PROJECT_NAME} \
down
