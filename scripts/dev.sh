#!/usr/bin/env bash

PROGNAME="$(basename $0)"
ROOT="$(cd "$(dirname "$0")/.." &>/dev/null; pwd -P)"

DB_HOST="${DB_HOST:-"127.0.0.1"}"
DB_NAME="${DB_NAME:-"eventstoredb"}"
DB_PASSWORD="${DB_PASSWORD:-"changeme"}"
DB_PORT="${DB_PORT:-"5432"}"
DB_USER="${DB_USER:-"postgres"}"
DB_SSLMODE="${DB_SSLMODE:-"disable"}"


function help() {
  echo 1>&2 "Usage: ${PROGNAME} <command>"
  echo 1>&2 ""
  echo 1>&2 "Commands:"
  echo 1>&2 "  start        start the service"
  echo 1>&2 "  migrate		  run migration"
  echo 1>&2 "  up			      pull and start infrastructure images"
  echo 1>&2 "  down			    stop all infrastructure images"
}

function setup_env() {
  export DB_HOST=${DB_HOST}
  export DB_NAME=${DB_NAME}
  export DB_PASSWORD=${DB_PASSWORD}
  export DB_PORT=${DB_PORT}
  export DB_USER=${DB_USER}
  export CONNECTION_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}
}

function start() {
  setup_env
	go run ./cmd/service/main.go
}

function migrate() {
  setup_env
	go run ./cmd/migrate/main.go
}

function up() {
	docker-compose -f ./builders/docker-compose.yml up -d
}
function down() {
	docker-compose -f ./builders/docker-compose.yml down
}

SUBCOMMAND="${1:-}"
case "${SUBCOMMAND}" in
  "" | "help" | "-h" | "--help" )
    help
    ;;

  "start" )
    shift
    start "$@"
    ;;

  "migrate" )
    shift
    migrate "$@"
    ;;

  "up" )
    shift
    up "$@"
    ;;

  "down" )
    shift
    down "$@"
    ;;

  *)
    help
    exit 1
    ;;
esac