#!/bin/bash

docker compose -p go-zulip \
-f docker-compose-dev-env.yml \
-f docker-compose-dev.yml \
-f docker-compose-zulip.yml \
up

docker compose -p go-zulip \
down
