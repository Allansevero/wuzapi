# Endpoint Push History to Webhook - ATUALIZADO

## 🎯 Mudanças Implementadas

### ✅ Melhorias:

1. **chat_jid agora é OPCIONAL**
   - Se fornecido: busca apenas aquele chat
   - Se omitido: busca TODAS as conversas da instância

2. **destination_number incluído no payload**
   - Campo enviarpara da instância
   - Vai junto com as mensagens

## �� Endpoint

```
POST /chat/history/push
```

## 🔐 Autenticação

```
Header: token: SEU_TOKEN_DA_INSTANCIA
```

## 📝 Parâmetros

| Parâmetro | Obrigatório | Descrição | Exemplo |
|-----------|-------------|-----------|---------|
| `webhook_url` | ✅ Sim | URL do webhook n8n | `https://metrizap-n8n.dyrluy.easypanel.host/webhook/...` |
| `chat_jid` | ❌ Não | JID específico (omitir = todas) | `555195611075@s.whatsapp.net` |
| `date` | ❌ Não | Filtro de data | `today` ou `2025-11-06` |
| `date_from` | ❌ Não | Data inicial | `2025-11-01` |
| `date_to` | ❌ Não | Data final | `2025-11-06` |
| `limit` | ❌ Não | Limite de mensagens | `50` (padrão) |

## 🚀 Exemplos de Uso

### Exemplo 1: TODAS as conversas de HOJE

```bash
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

**Resultado:** Todas as conversas de hoje da instância

### Exemplo 2: Chat específico de hoje

```bash
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&chat_jid=555195611075@s.whatsapp.net&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

**Resultado:** Apenas conversas com esse número de hoje

### Exemplo 3: Todas as conversas dos últimos 7 dias

```bash
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&date_from=2025-10-30&date_to=2025-11-06" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

### Exemplo 4: Todas as conversas (sem filtro de data)

```bash
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.dyrluy.easypanel.host/webhook/9e183064-20c4-4334-a139-f908f684a938&limit=100" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

## 📤 Payload Enviado ao Webhook

### Quando chat_jid é fornecido:

```json
{
  "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
  "destination_number": "5551981936133",
  "chat_jid": "555195611075@s.whatsapp.net",
  "message_count": 4,
  "timestamp": "2025-11-06T12:30:00-03:00",
  "date_from": "2025-11-06",
  "date_to": "2025-11-06",
  "messages": [...]
}
```

### Quando chat_jid NÃO é fornecido (todas as conversas):

```json
{
  "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
  "destination_number": "5551981936133",
  "chat_jid": "all",
  "all_chats": true,
  "message_count": 25,
  "timestamp": "2025-11-06T12:30:00-03:00",
  "date_from": "2025-11-06",
  "date_to": "2025-11-06",
  "messages": [
    {
      "id": 29,
      "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
      "chat_jid": "555195611075@s.whatsapp.net",
      "sender_jid": "555195611075@s.whatsapp.net",
      "message_id": "AC9CF5843805CBB6E88862986EF100E3",
      "timestamp": "2025-11-06T11:31:37-03:00",
      "message_type": "text",
      "text_content": "Parágrafo 6.2",
      "media_link": "",
      "quoted_message_id": "",
      "data_json": "{...}"
    },
    {
      "id": 28,
      "user_id": "2352d29b7b8c13ad75526f7e9b1c6a9f",
      "chat_jid": "555180264083@s.whatsapp.net",
      "sender_jid": "555180264083@s.whatsapp.net",
      "message_id": "ABC123...",
      "timestamp": "2025-11-06T10:45:12-03:00",
      "message_type": "text",
      "text_content": "Outra conversa...",
      "media_link": "",
      "quoted_message_id": "",
      "data_json": "{...}"
    }
  ]
}
```

## 📊 Campos do Payload

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `user_id` | string | ID interno do usuário |
| `destination_number` | string | Número para enviar análises (enviarpara) |
| `chat_jid` | string | JID do chat ou "all" |
| `all_chats` | boolean | true quando busca todas as conversas |
| `message_count` | number | Quantidade de mensagens |
| `timestamp` | string | Quando foi gerado o payload |
| `date_from` | string | Data inicial do filtro (se usado) |
| `date_to` | string | Data final do filtro (se usado) |
| `messages` | array | Array de mensagens |

## 🔧 Configuração no n8n

### Webhook Node:

```
Method: POST
Path: /webhook/9e183064-20c4-4334-a139-f908f684a938
Response Mode: Last Node
```

### Acessar dados no n8n:

```javascript
// Verificar se são todas as conversas
const isAllChats = $json.all_chats || false;

// Pegar destination_number
const destinationNumber = $json.destination_number;

// Pegar mensagens
const messages = $json.messages;

// Iterar sobre as mensagens
messages.forEach(msg => {
  console.log('Chat:', msg.chat_jid);
  console.log('Mensagem:', msg.text_content);
  console.log('Enviar para:', destinationNumber);
});
```

## 💡 Casos de Uso

### Caso 1: Relatório Diário de TODAS as Conversas

```bash
# Executar todo dia às 18h (cron job)
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.../webhook/...&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

n8n recebe todas as conversas do dia e envia relatório para `destination_number`

### Caso 2: Monitorar Chat Específico

```bash
# Verificar chat VIP a cada hora
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.../webhook/...&chat_jid=555195611075@s.whatsapp.net&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

### Caso 3: Backup Semanal

```bash
# Todo domingo
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://metrizap-n8n.../webhook/...&date_from=2025-11-01&date_to=2025-11-07&limit=1000" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

## 🎯 Diferenças Entre Modos

| Aspecto | Com chat_jid | Sem chat_jid (todas) |
|---------|--------------|---------------------|
| Filtro | 1 chat específico | Todas as conversas |
| Payload `chat_jid` | JID real | "all" |
| Payload `all_chats` | Não presente | true |
| Performance | Mais rápido | Mais lento (mais dados) |
| Uso | Monitorar cliente específico | Relatórios gerais |

## ✅ Resposta do Endpoint

### Sucesso:

```json
{
  "code": 200,
  "data": {
    "success": true,
    "message_count": 25,
    "webhook_url": "https://metrizap-n8n...",
    "webhook_status": 200
  },
  "success": true
}
```

## 🧪 Testando

### Teste 1: Todas as conversas de hoje

```bash
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://webhook.site/SEU-UUID&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

Veja em webhook.site:
- `all_chats: true`
- `chat_jid: "all"`
- `destination_number: "5551981936133"`
- Array com mensagens de vários chats

### Teste 2: Chat específico

```bash
curl -X POST "http://localhost:8080/chat/history/push?webhook_url=https://webhook.site/SEU-UUID&chat_jid=555195611075@s.whatsapp.net&date=today" \
  -H "token: 975a845ad05c7873423dbcfaf31b6acf"
```

Veja em webhook.site:
- `chat_jid: "555195611075@s.whatsapp.net"`
- `destination_number: "5551981936133"`
- Apenas mensagens desse chat

## 📝 Notas Importantes

1. **destination_number** sempre vem no payload
2. Se não tiver `destination_number` configurado, virá vazio `""`
3. Limite padrão: 50 mensagens
4. Ordenação: DESC (mais recente primeiro)
5. Quando `chat_jid` não é fornecido, pode trazer mensagens de múltiplos chats

## 🚨 Dica de Performance

Para buscar todas as conversas, considere:
- Usar `limit` para não sobrecarregar
- Usar filtro de data (`date=today`)
- Processar em background no n8n

---

**Última Atualização:** 06 de Novembro de 2025  
**Mudanças:**
- `chat_jid` agora é opcional
- `destination_number` incluído no payload
- Suporte para buscar todas as conversas da instância
