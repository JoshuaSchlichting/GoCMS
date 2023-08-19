#!/bin/bash
# This script is used for starting the GoCMS process via docker compose,
# but with some handy arguments for clean up, teardown, and rebuilding the web
# service only.

# Why would you want to only rebuild the web service? Becuase it's a pain 
# recreating your user in Kratos each time you restart the process.

COMPOSE_ARGS="-f docker-compose.yml -f auth/kratos/quickstart.yml -f auth/kratos/quickstart-standalone.yml -f auth/kratos/quickstart-postgres.yml"

if [ "$1" = "--down" ]; then
  docker compose $COMPOSE_ARGS down
  docker rm gocms-web-1
  exit 0
fi

if [ "$1" = "--clean-up" ]; then
  docker compose $COMPOSE_ARGS down -v
  docker compose $COMPOSE_ARGS rm -fsv
  docker rm gocms-web-1
  exit 0
fi

if [ "$1" = "--rebuild-web" ]; then
  docker compose build web
fi

docker compose $COMPOSE_ARGS up
