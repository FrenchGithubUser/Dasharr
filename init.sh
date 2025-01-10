#!/bin/bash

# this file is ran if no database file is found, considered as a new Dasharr user

if [[ -z "${API_KEY}" ]]; then
  echo "Error: API_KEY environment variable is not set."
  exit 1
fi

database_path="/backend/config/database.db"

if [ -e "$database_path" ]; then
  echo "Database already exists, skipping user initialization"
else
  echo "Database file doesn't exist, running user initialization"

  touch "$database_path"

  curl -X POST http://localhost:1323/initdb \
    -H "X-API-Key: ${API_KEY}" \
fi


curl -X POST http://localhost:1323/initdb \
  -H "X-API-Key: ${API_KEY}" \

