#!/bin/bash

source /etc/environment

if [[ -z "${API_KEY}" ]]; then
  echo "Error: API_KEY environment variable is not set."
  exit 1
fi

curl -X GET http://localhost:1323/api/update \
     -H "X-API-Key: ${API_KEY}"