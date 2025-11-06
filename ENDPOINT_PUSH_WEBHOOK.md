# Endpoint Push History to Webhook

## 📋 Descrição

Novo endpoint que **puxa as conversas** da instância e **envia automaticamente** para um webhook externo (n8n).

**Problema resolvido:** 
- n8n recebe 200 OK mas body vazio
- Agora o servidor WUZAPI envia as conversas diretamente para o n8n

## 🎯 Endpoint

```
POST /chat/history/push
```

## 🔐 Autenticação

Header obrigatório:
```
token: SEU_TOKEN_DA_INSTANCIA
```

## 📝 Parâmetros (Query String)

| Parâmetro | Obrigatório | Descrição | Exemplo |
|-----------|-------------|-----------|---------|
| `chat_jid` | ✅ Sim | JID do chat | `555195611075@s.whatsapp.net` |
| `webhook_url` | ✅ Sim | URL do webhook n8n | `https://metrizap-n8n.dyrluy.easypanel.host/webhook/...` |
| `date` | ❌ Não | Filtro de data única | `today` ou `2025-11-06` |
| `date_from` | ❌ Não | Data inicial | `2025-11-01` |
| `date_to` | ❌ Não | Data final | `2025-11-06` |
| `limit` | ❌ Não | Limite de mensagens | `50` (padrão: 50) |

## 📤 Como Funciona

1. Você faz **POST** para `/chat/history/push`
2. WUZAPI **busca** as mensagens no banco de dados
3. WUZAPI **envia** as mensagens para o webhook do n8n
4. WUZAPI retorna status da operação

```
┌─────────┐    POST     ┌─────────┐   GET DB   ┌──────────┐
│         │──────────>  │         │────────────>│ Database │
│  Você   │             │ WUZAPI  │<────────────│          │
│         │<────────────│         │             └──────────┘
└─────────┘   Response  │         │
                        │         │   POST
                        │         │─────────────>┌─────────┐
                        └─────────┘              │   n8n   │
                                                 │ Webhook │
                                                 └─────────┘
```

## 🚀 Exemplos de Uso

### Exemplo 1: Mensagens de HOJE

```bash
curl -X POST "http://localhost:8080/chat/history/push?chat_jid=555195611075@s.whatsapp.net&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

### Exemplo 2: Mensagens de uma data específica

```bash
curl -X POST "http://localhost:8080/chat/history/push?chat_jid=555195611075@s.whatsapp.net&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date=2025-11-06" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

### Exemplo 3: Mensagens dos últimos 7 dias

```bash
curl -X POST "http://localhost:8080/chat/history/push?chat_jid=555195611075@s.whatsapp.net&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date_from=2025-10-30&date_to=2025-11-06" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

### Exemplo 4: Com limite de mensagens

```bash
curl -X POST "http://localhost:8080/chat/history/push?chat_jid=555195611075@s.whatsapp.net&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date=today&limit=100" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

## 📥 Resposta do Endpoint

### Sucesso (200 OK):

```json
{
  "code": 200,
  "data": {
    "success": true,
    "message_count": 4,
    "webhook_url": "https://metrizap-n8n.dyrluy.easypanel.host/webhook/...",
    "webhook_status": 200
  },
  "success": true
}
```

### Erro (400 Bad Request):

```json
{
  "code": 400,
  "error": "chat_jid is required",
  "success": false
}
```

### Erro (502 Bad Gateway):

```json
{
  "code": 502,
  "error": "webhook returned status 500",
  "success": false
}
```

## 📤 Payload Enviado ao Webhook n8n

O WUZAPI envia este JSON para o webhook:

```json
{
  "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
  "chat_jid": "555195611075@s.whatsapp.net",
  "message_count": 4,
  "timestamp": "2025-11-06T12:00:00-03:00",
  "date_from": "2025-11-06",
  "date_to": "2025-11-06",
  "messages": [
    {
      "id": 29,
      "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
      "chat_jid": "555195611075@s.whatsapp.net",
      "sender_jid": "555195611075@s.whatsapp.net",
      "message_id": "AC9CF5843805CBB6E88862986EF100E3",
      "timestamp": "2025-11-06T11:31:37.388326399-03:00",
      "message_type": "text",
      "text_content": "Parágrafo 6.2",
      "media_link": "",
      "quoted_message_id": "",
      "data_json": "{...}"
    },
    {
      "id": 28,
      "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
      "chat_jid": "555195611075@s.whatsapp.net",
      "sender_jid": "555195611075@s.whatsapp.net",
      "message_id": "AC244FB16EFE81E8813C15B24226412B",
      "timestamp": "2025-11-06T11:31:30.45548085-03:00",
      "message_type": "text",
      "text_content": "Esta obrigação de sigilo...",
      "media_link": "",
      "quoted_message_id": "",
      "data_json": "{...}"
    }
  ]
}
```

## 🔧 Configuração no n8n

### Passo 1: Criar Webhook no n8n

1. Adicione node **Webhook**
2. Método: **POST**
3. Path: `/webhook/9e183064-20c4-4334-a139-f908f684a938`
4. Response Mode: **Last Node**

### Passo 2: Processar Dados

O webhook receberá automaticamente:
- `body.user_id`
- `body.chat_jid`
- `body.message_count`
- `body.messages[]` - Array de mensagens

### Passo 3: Acessar Mensagens

```javascript
// No Code Node do n8n
const messages = $json.messages;

messages.forEach(msg => {
  console.log('Mensagem:', msg.text_content);
  console.log('Horário:', msg.timestamp);
  console.log('De:', msg.sender_jid);
});
```

## 🧪 Testando

### Teste 1: Verificar se endpoint existe

```bash
curl -X POST "http://localhost:8080/chat/history/push" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

Deve retornar: `chat_jid is required`

### Teste 2: Enviar para webhook de teste

Use [webhook.site](https://webhook.site) para testar:

```bash
curl -X POST "http://localhost:8080/chat/history/push?chat_jid=555195611075@s.whatsapp.net&webhook_url=https://webhook.site/SEU-UUID&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

Veja o payload no webhook.site

### Teste 3: Enviar para n8n

```bash
curl -X POST "http://localhost:8080/chat/history/push?chat_jid=555195611075@s.whatsapp.net&webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

## 💡 Vantagens

✅ **Servidor envia** - Não depende do n8n puxar  
✅ **Formato garantido** - JSON sempre correto  
✅ **Timeout maior** - Servidor aguarda resposta  
✅ **Mais confiável** - Elimina problema de cache  
✅ **Logs melhores** - WUZAPI loga tudo  

## 🔒 Segurança

- Token é validado antes de enviar
- Webhook URL deve ser HTTPS (produção)
- Timeout de 30 segundos
- Logs de todas as operações

## 📊 Logs

O servidor loga:
```
INFO Successfully pushed history to webhook
  user_id=2352d29b7b8c13ad75526f7e9b1c6a9f
  chat_jid=555195611075@s.whatsapp.net
  webhook_url=https://metrizap-n8n...
  message_count=4
```

## 🚨 Troubleshooting

### Erro: "webhook returned status 500"
- Webhook n8n está com erro
- Verifique logs do n8n

### Erro: "failed to send to webhook"
- Webhook URL incorreta
- Firewall bloqueando
- n8n offline

### Erro: "message history is disabled"
- History não habilitado para o usuário
- Configure history no banco de dados

## 📝 Notas

- Limite padrão: 50 mensagens
- Ordenação: DESC (mais recente primeiro)
- Suporta PostgreSQL e SQLite
- Webhook deve aceitar POST com JSON

---

**Data de Implementação:** 06 de Novembro de 2025  
**Arquivos Modificados:**
- `handlers.go` - Nova função `PushHistoryToWebhook()`
- `routes.go` - Nova rota `/chat/history/push`
