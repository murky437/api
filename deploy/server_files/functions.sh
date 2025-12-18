#!/bin/bash
set -eu

cd "$(dirname "$0")"

BLACK='\033[0;30m'
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[0;37m'
NC='\033[0m' # No Color

info() {
    echo -e "${CYAN}[$(date '+%Y-%m-%d %H:%M:%S')][server]${NC} $*"
}

err() {
    echo -e "${CYAN}[$(date '+%Y-%m-%d %H:%M:%S')][server]${NC} ${RED}$*${NC}" >&2
}

deleteOldPreviousDeployments() {
  local BASE="./deployments/previous"

  local dirs
  dirs=$(printf "%s\n" "$BASE"/*/ | xargs --max-args=1 basename | sort --reverse)

  local latest3dirs
  latest3dirs=$(printf "%s\n" "$dirs" | head --lines=3)

  local timestampWeekAgo
  timestampWeekAgo=$(date --date="7 days ago" +"%Y%m%dT%H%M%S")
  local latestWeekDirs
  latestWeekDirs=$(printf "%s\n" "$dirs" | awk --assign c="$timestampWeekAgo" '$0 >= c')

  local keepDirs
  keepDirs=$(printf "%s\n%s\n" "$latest3dirs" "$latestWeekDirs" | sort --unique)

  local d
  for d in $dirs; do
    printf "%s\n" "$keepDirs" | grep --quiet --line-regexp "$d" || rm --recursive --force "${BASE:?}/$d"
  done
}

switchToNewDeployment() {
  if [ ! -f deployments/new/deploy.json ] || [ ! -f deployments/new/api ]; then
    err "ERROR: new deployment is missing required files"
  fi

  info "Copying current database to new deployment"
  mkdir -p deployments/new/db
  if [ -f deployments/current/db/db.sqlite3 ]; then
      sqlite3 deployments/current/db/db.sqlite3 ".backup 'deployments/new/db/db.sqlite3'"
  fi

  info "Moving current deployment to previous deployments directory"
  if [ -f deployments/current/deploy.json ]; then
    timestamp=$(jq -r '.timestamp' deployments/current/deploy.json)
    rsync -a --delete --mkpath deployments/current/ "deployments/previous/${timestamp}/"
  fi

  info "Moving new deployment to current deployment directory"
  rsync -a --delete --mkpath deployments/new/ deployments/current/
  rm -rf deployments/new

  info "Reloading with docker compose"
  docker compose up -d --build --force-recreate

  info "Deleting old previous deployments"
  deleteOldPreviousDeployments

  info "Switch complete"
}

rollBackToPreviousDeployment() {
  latest_prev=$(find deployments/previous -mindepth 1 -maxdepth 1 -type d | sort | tail -n1)

  if [ -z "$latest_prev" ]; then
    info "No previous deployment found."
    exit 1
  fi

  read -r -p "Are you sure you want to roll back to the latest previous deployment '$latest_prev'? (yes/no) " confirm
  if [[ "$confirm" != "yes" ]]; then
    info "Rollback cancelled."
    exit 0
  fi

  info "Rolling back to $latest_prev"

  rsync -a --delete "$latest_prev"/ deployments/current/
  rm -rf "$latest_prev"

  info "Reloading with docker compose"
  docker compose up -d --build

  info "Rollback complete"
}
