#!/bin/bash

# used for the cron job
printenv >> /etc/environment &

database_path="/backend/config/database.db"
config_json_path="/backend/config/config.json"
config_sample_json_path="/backend/config_sample/config_sample.json"

echo "Waiting for the backend to become available on port 1323..."
while ! curl -s http://localhost:1323/api > /dev/null; do
    sleep 1
done

if [ -e "$config_json_path" ]; then
  echo "Config file already exists, skipping copy"

  curl -X GET http://localhost:1323/api/update \
  -H "X-API-Key: ${API_KEY}"
else
  echo "Config file doesn't exist, creating it"

  cp "$config_sample_json_path" "$config_json_path"

fi


if [[ -z "${API_KEY}" ]]; then
  echo "Error: API_KEY environment variable is not set."
  exit 1
fi

if [ -e "$database_path" ]; then
  echo "Database already exists, skipping user initialization"
else
  echo "Database file doesn't exist, running user initialization"

  sqlite3 "$database_path" "VACUUM;"

  echo "Database created"

  curl -X GET http://localhost:1323/api/initdb \
    -H "X-API-Key: ${API_KEY}"
fi


curl -X GET http://localhost:1323/api/collectStats \
  -H "X-API-Key: ${API_KEY}"