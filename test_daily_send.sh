#!/bin/bash

# Script para testar envio manual de mensagens compiladas diárias
# Uso: ./test_daily_send.sh [TOKEN] [INSTANCE_ID] [DATE]

TOKEN="${1:-}"
INSTANCE_ID="${2:-}"
DATE="${3:-}"

if [ -z "$TOKEN" ]; then
    echo "Uso: $0 TOKEN [INSTANCE_ID] [DATE]"
    echo ""
    echo "Exemplos:"
    echo "  $0 seu-token-aqui"
    echo "  $0 seu-token-aqui instance-id-opcional"
    echo "  $0 seu-token-aqui instance-id-opcional 2025-11-03"
    exit 1
fi

URL="http://localhost:8080/session/send-daily-test?token=$TOKEN"

if [ -n "$INSTANCE_ID" ]; then
    URL="${URL}&instance_id=${INSTANCE_ID}"
fi

if [ -n "$DATE" ]; then
    URL="${URL}&date=${DATE}"
fi

echo "============================================"
echo "Testando envio manual de mensagens diárias"
echo "============================================"
echo "URL: $URL"
echo ""

RESPONSE=$(curl -s -X POST "$URL" -H "Content-Type: application/json")

echo "Resposta do servidor:"
echo "$RESPONSE" | jq . 2>/dev/null || echo "$RESPONSE"
echo ""
echo "============================================"
