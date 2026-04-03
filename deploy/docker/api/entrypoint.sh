#!/bin/sh
set -eu

# Move secret to file if GOOGLE_SERVICE_ACCOUNT_KEY_FILE_CONTENT is set (base64 encoded)
if [ -n "${GOOGLE_SERVICE_ACCOUNT_KEY_FILE_CONTENT:-}" ]; then
  echo "$GOOGLE_SERVICE_ACCOUNT_KEY_FILE_CONTENT" | base64 -d > "$GOOGLE_SERVICE_ACCOUNT_KEY_FILE_PATH"
fi

# Start the main container process
exec "$@"
