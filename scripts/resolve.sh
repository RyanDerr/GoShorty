#!/bin/bash

# Initialize variables
SHORT=""

# Parse arguments
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --short) SHORT="$2"; shift ;;
    *) echo "Unknown parameter passed: $1"; exit 1 ;;
  esac
  shift
done

# Check if SHORT is provided
if [ -z "$SHORT" ]; then
  echo "Error: --short is required"
  exit 1
fi

API_URL="http://localhost:8080/$SHORT"

curl $API_URL