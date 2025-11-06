#!/bin/bash

# =============================================================================
# WUZAPI - Exemplos de CURL para puxar conversas de HOJE
# Data: 2025-11-06
# =============================================================================

# Configure suas variáveis
TOKEN="seu_token_aqui"
HOST="http://localhost:8080"
NUMERO="5511999999999@s.whatsapp.net"
GRUPO="120363123456789012@g.us"

# =============================================================================
# 1. MENSAGENS DE HOJE - Contato Individual
# =============================================================================
echo "=== 1. Mensagens de HOJE de um contato ==="
curl -X GET "${HOST}/chat/history?chat_jid=${NUMERO}&date=today" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# 2. MENSAGENS DE HOJE - Grupo
# =============================================================================
echo -e "\n=== 2. Mensagens de HOJE de um grupo ==="
curl -X GET "${HOST}/chat/history?chat_jid=${GRUPO}&date=today" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# 3. MENSAGENS DE HOJE com LIMITE de 100
# =============================================================================
echo -e "\n=== 3. Mensagens de HOJE (limite 100) ==="
curl -X GET "${HOST}/chat/history?chat_jid=${NUMERO}&date=today&limit=100" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# 4. MENSAGENS DE UMA DATA ESPECÍFICA (ex: 06/11/2025)
# =============================================================================
echo -e "\n=== 4. Mensagens de 2025-11-06 ==="
curl -X GET "${HOST}/chat/history?chat_jid=${NUMERO}&date=2025-11-06" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# 5. MENSAGENS DOS ÚLTIMOS 7 DIAS
# =============================================================================
DATA_7_DIAS_ATRAS=$(date -d '7 days ago' +%Y-%m-%d)
echo -e "\n=== 5. Mensagens dos últimos 7 dias ==="
curl -X GET "${HOST}/chat/history?chat_jid=${NUMERO}&date_from=${DATA_7_DIAS_ATRAS}" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# 6. MENSAGENS ENTRE DUAS DATAS (ex: 01/11 até 06/11)
# =============================================================================
echo -e "\n=== 6. Mensagens entre 01/11 e 06/11 ==="
curl -X GET "${HOST}/chat/history?chat_jid=${NUMERO}&date_from=2025-11-01&date_to=2025-11-06" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# 7. LISTAR TODOS OS CHATS COM MENSAGENS
# =============================================================================
echo -e "\n=== 7. Índice de todos os chats ==="
curl -X GET "${HOST}/chat/history?chat_jid=index" \
  -H "token: ${TOKEN}" | jq '.'

# =============================================================================
# DICA: Para salvar em arquivo JSON
# =============================================================================
# curl -X GET "${HOST}/chat/history?chat_jid=${NUMERO}&date=today" \
#   -H "token: ${TOKEN}" > conversas_hoje.json

