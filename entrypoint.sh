#!/bin/sh
set -e

# Create database if it doesn't exist
/bin/coveritup -migrate create

# Apply migrations
/bin/coveritup -migrate up

# Start the application
# -host 0.0.0.0 is crucial for Docker networking
exec /bin/coveritup -host 0.0.0.0 "$@"
