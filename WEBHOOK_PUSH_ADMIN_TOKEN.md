# Endpoint Push History - Usando ADMIN TOKEN

## ⚠️ MUDANÇA IMPORTANTE

O endpoint `/chat/history/push` agora requer **ADMIN_TOKEN** ao invés do token da instância.

## 🔐 Autenticação

**Antes:** Token da instância  
**Agora:** **ADMIN_TOKEN**

## 🎯 Novo Endpoint

```
POST /admin/chat/history/push
```

## 📝 Parâmetros Obrigatórios

| Parâmetro | Tipo | Descrição | Exemplo |
|-----------|------|-----------|---------|
| `instance_id` | Query | **NOVO!** ID da instância | `2352d29b7b8c13ad75526f7e9b1c6a9f` |
| `webhook_url` | Query | URL do webhook n8n | `https://metrizap-n8n...` |

## 📝 Parâmetros Opcionais

| Parâmetro | Tipo | Descrição | Exemplo |
|-----------|------|-----------|---------|
| `chat_jid` | Query | JID específico (omitir = todas) | `555195611075@s.whatsapp.net` |
| `date` | Query | Filtro de data | `today` ou `2025-11-06` |
| `date_from` | Query | Data inicial | `2025-11-01` |
| `date_to` | Query | Data final | `2025-11-06` |
| `limit` | Query | Limite de mensagens | `50` (padrão) |

## 🚀 Exemplos Atualizados

### Exemplo 1: TODAS as conversas de HOJE

```bash
curl -X POST "http://localhost:8080/admin/chat/history/push?instance_id=2352d29b7b8c13ad75526f7e9b1c6a9f&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date=today" \
  -H "token: SEU_ADMIN_TOKEN_AQUI"
```

### Exemplo 2: Chat específico de HOJE

```bash
curl -X POST "http://localhost:8080/admin/chat/history/push?instance_id=2352d29b7b8c13ad75526f7e9b1c6a9f&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&chat_jid=555195611075@s.whatsapp.net&date=today" \
  -H "token: SEU_ADMIN_TOKEN_AQUI"
```

### Exemplo 3: Últimos 7 dias

```bash
curl -X POST "http://localhost:8080/admin/chat/history/push?instance_id=2352d29b7b8c13ad75526f7e9b1c6a9f&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date_from=2025-10-30&date_to=2025-11-06" \
  -H "token: SEU_ADMIN_TOKEN_AQUI"
```

## 🔍 Como Pegar o Instance ID

### Opção 1: Endpoint Admin

```bash
curl -X GET "http://localhost:8080/admin/instances" \
  -H "token: SEU_ADMIN_TOKEN_AQUI"
```

Retorna lista de todas as instâncias com seus IDs.

### Opção 2: Dashboard

1. Acesse o dashboard da instância
2. Na URL, pegue o ID: `/dashboard/user-dashboard-v4.html?instance=INSTANCE_ID_AQUI`

## 📤 Payload Enviado ao Webhook

```json
{
  "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
  "destination_number": "5551981936133",
  "chat_jid": "all",
  "all_chats": true,
  "message_count": 33,
  "timestamp": "2025-11-06T14:00:00-03:00",
  "date_from": "2025-11-06",
  "date_to": "2025-11-06",
  "messages": [...]
}
```

## ⚠️ Erros Comuns

### Erro 401 Unauthorized

```json
{"code":401,"error":"unauthorized","success":false}
```

**Causa:** Token inválido ou token de instância ao invés de admin  
**Solução:** Use o ADMIN_TOKEN

### Erro 400 Bad Request

```json
{"code":400,"error":"instance_id is required","success":false}
```

**Causa:** Faltou o parâmetro `instance_id`  
**Solução:** Adicione `?instance_id=...` na URL

### Erro 404 Not Found

```json
{"code":404,"error":"instance not found","success":false}
```

**Causa:** instance_id inválido  
**Solução:** Verifique o ID com `GET /admin/instances`

## 💡 Vantagens do ADMIN_TOKEN

✅ **Segurança:** Admin tem controle total  
✅ **Flexibilidade:** Pode acessar qualquer instância  
✅ **Centralização:** Um único token para todas as instâncias  
✅ **Auditoria:** Fácil rastrear quem está puxando dados  

## 🔄 Migração do Código Antigo

**Antes:**
```bash
curl -X POST ".../chat/history/push?..." \
  -H "token: TOKEN_DA_INSTANCIA"
```

**Agora:**
```bash
curl -X POST ".../admin/chat/history/push?instance_id=ID_AQUI&..." \
  -H "token: ADMIN_TOKEN"
```

## 📝 Configuração no n8n

### Webhook Trigger Node:

```
Method: POST
Path: /webhook/9e183064-20c4-4334-a139-f908f684a938
Authentication: None
```

### HTTP Request Node (para chamar o endpoint):

```json
{
  "method": "POST",
  "url": "http://wuzapi:8080/admin/chat/history/push",
  "authentication": "headerAuth",
  "headerAuth": {
    "name": "token",
    "value": "={{ $env.WUZAPI_ADMIN_TOKEN }}"
  },
  "qs": {
    "instance_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
    "webhook_url": "https://metrizap-n8n.dyrluy.easypanel.host/webhook/...",
    "date": "today"
  }
}
```

## 🧪 Teste Rápido

```bash
# 1. Pegar instance_id
curl -X GET "http://localhost:8080/admin/instances" \
  -H "token: SEU_ADMIN_TOKEN" | jq

# 2. Usar o ID para puxar conversas
curl -X POST "http://localhost:8080/admin/chat/history/push?instance_id=ID_COPIADO&webhook_url=https://webhook.site/SEU-UUID&date=today" \
  -H "token: SEU_ADMIN_TOKEN"

# 3. Verificar em webhook.site
```

## ✅ Checklist

- [ ] Tenho o ADMIN_TOKEN
- [ ] Sei o instance_id
- [ ] URL do webhook n8n está correta
- [ ] Endpoint mudou para `/admin/chat/history/push`
- [ ] Adicionei parâmetro `instance_id`
- [ ] Usando header `token: ADMIN_TOKEN`

---

**Data de Atualização:** 06 de Novembro de 2025  
**Mudança:** Endpoint requer ADMIN_TOKEN e instance_id
