#!/bin/bash

# Script to test manual daily send to webhook
# Usage: ./test_webhook_send.sh <token>

if [ -z "$1" ]; then
    echo "Usage: $0 <token>"
    echo "Example: $0 cc52b16537c624c981b5fb54f83e30cc"
    exit 1
fi

TOKEN=$1
URL="http://localhost:8080/session/test-daily-send"

echo "================================================"
echo "Testing Daily Send to Webhook"
echo "================================================"
echo "Token: $TOKEN"
echo "URL: $URL"
echo ""

response=$(curl -s -X POST "$URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")

echo "Response:"
echo "$response" | jq '.' 2>/dev/null || echo "$response"
echo ""
echo "================================================"
