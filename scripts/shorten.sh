#!/bin/bash

# Initialize variables
URL=""
SHORT=""
EXP=""

# Parse arguments
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --url) URL="$2"; shift ;;
    --short) SHORT="$2"; shift ;;
    --exp) EXP="$2"; shift ;;
    *) echo "Unknown parameter passed: $1"; exit 1 ;;
  esac
  shift
done

# Check if URL is provided
if [ -z "$URL" ]; then
  echo "Error: URL is required"
  exit 1
fi

API_URL="http://localhost:3000/api/v1"
DATA="{\"url\": \"$URL\""

if [ -n "$SHORT" ]; then
  DATA="$DATA, \"short\": \"$SHORT\""
fi

if [ -n "$EXP" ]; then
  DATA="$DATA, \"expiration\": \"$EXP\""
fi

DATA="$DATA}"

curl -X POST -H "Content-Type: application/json" -d "$DATA" $API_URL