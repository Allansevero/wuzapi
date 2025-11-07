# Guia de Teste - Histórico com Informações da Instância

## 🧪 Como Testar

### Passo 1: Preparar Webhook de Teste

1. Acesse: https://webhook.site
2. Copie a URL única gerada (ex: `https://webhook.site/abc123-def456`)
3. Deixe a página aberta para ver as requisições

### Passo 2: Reiniciar o Servidor

```bash
sudo systemctl restart wuzapi

# Ou se estiver rodando manualmente
./wuzapi
```

### Passo 3: Verificar Logs (Opcional)

```bash
# Ver logs em tempo real
sudo journalctl -u wuzapi -f

# Procure por:
# "Preparing payload with instance information"
```

### Passo 4: Fazer Requisição de Histórico

Existem **2 formas** de solicitar o histórico:

---

## 📋 OPÇÃO 1: Com Bearer Token (Sistema de Usuário)

```bash
# Substitua:
# - YOUR_BEARER_TOKEN: Token JWT do usuário logado
# - YOUR_WEBHOOK_URL: URL do webhook.site
# - YOUR_INSTANCE_ID: ID da instância

curl -X POST 'http://localhost:8080/user/chat/history/push?webhook_url=YOUR_WEBHOOK_URL&chat_jid=all&date=today' \
  -H 'Authorization: Bearer YOUR_BEARER_TOKEN' \
  -H 'Content-Type: application/json'
```

**Exemplo:**
```bash
curl -X POST 'http://localhost:8080/user/chat/history/push?webhook_url=https://webhook.site/abc123-def456&chat_jid=all&date=today' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIs...' \
  -H 'Content-Type: application/json'
```

---

## 📋 OPÇÃO 2: Com Instance Token

```bash
# Substitua:
# - YOUR_INSTANCE_TOKEN: Token da instância específica
# - YOUR_WEBHOOK_URL: URL do webhook.site

curl -X POST 'http://localhost:8080/chat/history/push?webhook_url=YOUR_WEBHOOK_URL&chat_jid=all&date=today' \
  -H 'X-Instance-Token: YOUR_INSTANCE_TOKEN' \
  -H 'Content-Type: application/json'
```

**Exemplo:**
```bash
curl -X POST 'http://localhost:8080/chat/history/push?webhook_url=https://webhook.site/abc123-def456&chat_jid=all&date=today' \
  -H 'X-Instance-Token: 9a8b7c6d5e4f3g2h1i0j' \
  -H 'Content-Type: application/json'
```

---

## 🔍 Parâmetros da URL

| Parâmetro | Obrigatório | Descrição | Exemplo |
|-----------|-------------|-----------|---------|
| `webhook_url` | ✅ Sim | URL para receber o histórico | `https://webhook.site/abc123` |
| `chat_jid` | ❌ Não | JID do chat ou "all" | `555181936133@s.whatsapp.net` ou `all` |
| `date` | ❌ Não | Filtro de data | `today` ou `2025-11-07` |
| `date_from` | ❌ Não | Data inicial | `2025-11-01` |
| `date_to` | ❌ Não | Data final | `2025-11-07` |
| `limit` | ❌ Não | Limite de mensagens | `50` (padrão) |

---

## ✅ O Que Você Deve Receber no Webhook

```json
{
  "user_id": "1a2df8ab09ccaf2e5fa2a933f5f5cfa2",
  "instance_name": "Instância Padrão",        ⬅️ DEVE APARECER
  "instance_phone": "5511999999999",          ⬅️ DEVE APARECER
  "destination_number": "51995611075",
  "message_count": 50,
  "all_chats": true,
  "chat_jid": "all",
  "date_from": "2025-11-07",
  "date_to": "2025-11-07",
  "timestamp": "2025-11-07T18:00:00Z",
  "messages": [
    {
      "id": 230,
      "user_id": "1a2df8ab09ccaf2e5fa2a933f5f5cfa2",
      "chat_jid": "555181936133@s.whatsapp.net",
      "sender_jid": "555181936133:28@s.whatsapp.net",
      "message_id": "3EB03E10E8188CE3B4CFD1",
      "timestamp": "2025-11-07T14:49:39.944238831-03:00",
      "message_type": "text",
      "text_content": "Olá!",
      "media_link": "",
      "quoted_message_id": "",
      "datajson": ""
    }
  ]
}
```

---

## ❌ Troubleshooting

### Problema: `instance_name` e `instance_phone` estão vazios

**Possíveis causas:**

1. **Instância não está conectada**
   - O JID só é preenchido após a primeira conexão
   - Conecte a instância ao WhatsApp primeiro

2. **Nome da instância não foi definido**
   - Verifique se a instância tem um nome no banco de dados
   - Query: `SELECT id, name, jid FROM users WHERE id = 'seu_id';`

3. **Servidor não foi reiniciado**
   - As alterações só entram em vigor após reiniciar
   - `sudo systemctl restart wuzapi`

### Problema: Erro 404 ou 401

```bash
# Verifique se o endpoint existe
curl http://localhost:8080/chat/history/push

# Verifique se o token é válido
# O Bearer Token deve ser válido e não expirado
```

### Problema: Nenhuma mensagem retornada

```bash
# Verifique se há histórico habilitado
# No banco de dados, coluna 'history' deve ser > 0

SELECT id, name, history FROM users WHERE id = 'seu_id';

# Se history = 0, habilite:
UPDATE users SET history = 100 WHERE id = 'seu_id';
```

---

## 📝 Verificar Logs do Servidor

```bash
# Em tempo real
sudo journalctl -u wuzapi -f | grep "instance"

# Procure por:
# • "Retrieved instance info for webhook push"
# • "Preparing payload with instance information"
# • instance_name=...
# • instance_phone=...
```

---

## 🎯 Exemplo Completo de Teste

```bash
#!/bin/bash

# 1. Configure suas variáveis
WEBHOOK_URL="https://webhook.site/abc123-def456"
INSTANCE_TOKEN="seu_token_aqui"

# 2. Solicite o histórico
echo "🔄 Solicitando histórico..."
curl -X POST \
  "http://localhost:8080/chat/history/push?webhook_url=$WEBHOOK_URL&chat_jid=all&date=today" \
  -H "X-Instance-Token: $INSTANCE_TOKEN" \
  -H "Content-Type: application/json"

echo ""
echo "✅ Requisição enviada!"
echo "📱 Acesse $WEBHOOK_URL para ver o resultado"
```

---

## 🔍 Como Obter os Tokens

### Bearer Token (JWT):
1. Faça login no sistema via `/auth/login`
2. Copie o token retornado no campo `token`

### Instance Token:
1. Acesse o dashboard
2. Vá em "Contas conectadas"
3. Clique nos 3 pontinhos da instância
4. Copie o token mostrado

---

## 📊 Esperado vs Recebido

| Campo | Deve Conter |
|-------|-------------|
| `user_id` | ID da instância (hash) |
| `instance_name` | Nome da instância (ex: "WhatsApp Vendas") |
| `instance_phone` | Número puro (ex: "5511999999999") |
| `destination_number` | Número para enviar análises |
| `message_count` | Quantidade de mensagens |
| `messages` | Array com as mensagens |

**Se algum campo estiver faltando, verifique os logs do servidor!**
