#!/bin/sh

REDIS_PASSWORD=$(cat /run/secrets/ds_password)
exec redis-server --requirepass "$REDIS_PASSWORD"
