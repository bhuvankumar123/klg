#!/bin/bash

API_URL="http://localhost:6060/v1.0/logs"

log_levels=("INFO" "WARN" "FATAL" "DEBUG" "ERROR")
messages=(
  "Application started successfully"
  "User authentication failed"
  "Database connection established"
  "Unexpected error occurred"
  "Cache cleared successfully"
  "Disk space running low"
  "New user registered"
  "Payment transaction completed"
  "API request timeout"
  "Service restarted"
)

generate_random_ip() {
  echo "$((RANDOM % 256)).$((RANDOM % 256)).$((RANDOM % 256)).$((RANDOM % 256))"
}

for i in {1..100}
do
  level=${log_levels[$RANDOM % ${#log_levels[@]}]}
  message=${messages[$RANDOM % ${#messages[@]}]}
  user_id=$((1000 + RANDOM % 9000))
  ip_address=$(generate_random_ip)
  session_id=$(openssl rand -hex 16)

  json_data="{\"level\": \"$level\", \"message\": \"$message\", \"metadata\": {\"service\": \"api\", \"version\": \"1.0.0\", \"environment\": \"PROD\", \"user_id\": \"$user_id\", \"ip_address\": \"$ip_address\", \"session_id\": \"$session_id\"}}"

  curl --location "$API_URL" \
    --header "Content-Type: application/json" \
    --data "$json_data"

done