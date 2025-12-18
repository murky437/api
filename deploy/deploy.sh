#!/bin/bash
set -eu

execDir=$(pwd)

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
    echo -e "${MAGENTA}[$(date '+%Y-%m-%d %H:%M:%S')][local]${NC} $*"
}

err() {
    echo -e "${MAGENTA}[$(date '+%Y-%m-%d %H:%M:%S')][local]${NC} ${RED}$*${NC}" >&2
}

build() {
  info "Creating a new build"
  docker compose -f ../docker-compose.prod.yml run --rm --build api

  if [ -d new ]; then
    rm -rf new
  fi
  mkdir new
  cp ../build/api new

  info "Adding deploy.json file for info"
  jq -n --arg commit "$(git rev-parse --short HEAD)" \
        --arg timestamp "$(date '+%Y%m%dT%H%M%S')" \
        '{commit: $commit, timestamp: $timestamp}' > "new/deploy.json"

  info "Build complete"
}

generateEnvContent() {
  local envFile="../.env.example"
  local exclude=(SSH_HOST SERVER_APP_DIR)
  local keys=()
  local content=""
  local hasUnsetVars=0

  while IFS='=' read -r k _; do
    [[ "$k" =~ ^#.*$ || -z "$k" ]] && continue
    [[ " ${exclude[*]} " =~ " $k " ]] && continue
    keys+=("$k")
  done < "$envFile"

  for k in "${keys[@]}"; do
    if [ -z "${!k+x}" ]; then
      err "${k} is not set"
      hasUnsetVars=1
    else
      content+="$k=${!k}"$'\n'
    fi
  done
  [ "$hasUnsetVars" -eq 1 ] && exit 1

  echo "$content"
}

deploy() {
  if [ ! -d new ]; then
    err "Missing new deployment, run the build script first"
    exit 1
  fi

  info "Syncing deployment files to server"
  rsync -a --mkpath server_files/ "$SSH_HOST:$SERVER_APP_DIR/api/"
  rsync -a --delete --mkpath new/ "$SSH_HOST:$SERVER_APP_DIR/api/deployments/new/"

  if [ -n "${GOOGLE_SERVICE_ACCOUNT_KEY_FILE_CONTENT-}" ]; then
    info "Creating google secrets file"
    ssh "$SSH_HOST" "mkdir -p $SERVER_APP_DIR/api/secrets && cat > $SERVER_APP_DIR/api/secrets/googleServiceAccountKey.json" <<< "$GOOGLE_SERVICE_ACCOUNT_KEY_FILE_CONTENT"
  elif [ -f ../secrets/googleServiceAccountKey.json ]; then
    info "Copying google secrets file from local path"
    rsync -a --mkpath ../secrets/googleServiceAccountKey.json "$SSH_HOST:$SERVER_APP_DIR/api/secrets/googleServiceAccountKey.json"
  else
    err "Google service account key not provided"
    exit 1
  fi

  info "Creating .env file"
  ssh "$SSH_HOST" "cat > $SERVER_APP_DIR/api/.env" <<< "$(generateEnvContent)"

  info "Switching to new deployment"
  ssh "$SSH_HOST" "bash \"$SERVER_APP_DIR/api/switch_to_new_deployment.sh\""

  info "Deploy complete"
}

loadEnvFromFile() {
  envFile="$1"

  if [ ! -f "$envFile" ]; then
    err "Env file not found: $envFile"
    exit 1
  fi

  set -a
  # shellcheck disable=SC1090
  . "$envFile"
  set +a
}


#
# Main script
#

# Load env vars from first argument (.env file), if it is passed
if [ "${1-}" ]; then
  loadEnvFromFile "$execDir/$1"
fi

build
deploy

