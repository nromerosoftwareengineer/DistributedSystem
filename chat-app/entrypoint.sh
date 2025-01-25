#!/bin/bash
set -e

# Remove the default config if it exists
if [ -f /etc/nginx/conf.d/default.conf ]; then
    rm /etc/nginx/conf.d/default.conf
fi

# Start NGINX
exec "$@"
