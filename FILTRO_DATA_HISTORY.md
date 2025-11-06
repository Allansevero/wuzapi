# Filtro de Data no Endpoint de Histórico

## Modificação Implementada
Adicionado suporte a filtros de data no endpoint `/chat/history`.

## Novos Parâmetros de Query

### 1. `date` - Filtro de Data Única
Filtra mensagens de um dia específico (00:00:00 até 23:59:59).

**Valores aceitos:**
- `today` - Mensagens de hoje
- `YYYY-MM-DD` - Data específica (ex: `2025-11-06`)

### 2. `date_from` - Data Inicial (Intervalo)
Define a data/hora inicial do filtro (00:00:00 do dia).

**Formato:** `YYYY-MM-DD`

### 3. `date_to` - Data Final (Intervalo)
Define a data/hora final do filtro (23:59:59 do dia).

**Formato:** `YYYY-MM-DD`

---

## Exemplos de Uso

### 1️⃣ Mensagens de HOJE (dia atual)
```bash
curl -X GET "http://localhost:8080/chat/history?chat_jid=5511999999999@s.whatsapp.net&date=today" \
  -H "token: SEU_TOKEN_AQUI"
```

### 2️⃣ Mensagens de uma data específica
```bash
# Mensagens do dia 2025-11-06
curl -X GET "http://localhost:8080/chat/history?chat_jid=5511999999999@s.whatsapp.net&date=2025-11-06" \
  -H "token: SEU_TOKEN_AQUI"
```

### 3️⃣ Mensagens a partir de uma data (date_from)
```bash
# Mensagens de 2025-11-01 em diante
curl -X GET "http://localhost:8080/chat/history?chat_jid=5511999999999@s.whatsapp.net&date_from=2025-11-01" \
  -H "token: SEU_TOKEN_AQUI"
```

### 4️⃣ Mensagens até uma data (date_to)
```bash
# Mensagens até 2025-11-06
curl -X GET "http://localhost:8080/chat/history?chat_jid=5511999999999@s.whatsapp.net&date_to=2025-11-06" \
  -H "token: SEU_TOKEN_AQUI"
```

### 5️⃣ Mensagens entre duas datas (intervalo)
```bash
# Mensagens entre 2025-11-01 e 2025-11-06
curl -X GET "http://localhost:8080/chat/history?chat_jid=5511999999999@s.whatsapp.net&date_from=2025-11-01&date_to=2025-11-06" \
  -H "token: SEU_TOKEN_AQUI"
```

### 6️⃣ Mensagens de hoje com limite de 100
```bash
curl -X GET "http://localhost:8080/chat/history?chat_jid=5511999999999@s.whatsapp.net&date=today&limit=100" \
  -H "token: SEU_TOKEN_AQUI"
```

### 7️⃣ Mensagens de grupo de hoje
```bash
curl -X GET "http://localhost:8080/chat/history?chat_jid=120363123456789012@g.us&date=today&limit=200" \
  -H "token: SEU_TOKEN_AQUI"
```

---

## Combinação de Parâmetros

| Parâmetro | Obrigatório | Compatível com | Descrição |
|-----------|-------------|----------------|-----------|
| `chat_jid` | ✅ Sim | Todos | JID do chat (contato ou grupo) |
| `date` | ❌ Não | `limit` | Filtra um dia específico |
| `date_from` | ❌ Não | `date_to`, `limit` | Data inicial do intervalo |
| `date_to` | ❌ Não | `date_from`, `limit` | Data final do intervalo |
| `limit` | ❌ Não | Todos | Máximo de mensagens (padrão: 50) |

⚠️ **Nota:** Não use `date` junto com `date_from` ou `date_to`. Se você usar `date`, os outros filtros de data serão ignorados.

---

## Formato da Resposta

```json
{
  "code": 200,
  "data": [
    {
      "id": 123,
      "user_id": "seu_user_id",
      "chat_jid": "5511999999999@s.whatsapp.net",
      "sender_jid": "5511999999999@s.whatsapp.net",
      "message_id": "3EB0ABC123...",
      "timestamp": "2025-11-06T10:30:00Z",
      "message_type": "text",
      "text_content": "Olá, bom dia!",
      "media_link": "",
      "quoted_message_id": "",
      "datajson": ""
    },
    {
      "id": 124,
      "user_id": "seu_user_id",
      "chat_jid": "5511999999999@s.whatsapp.net",
      "sender_jid": "me",
      "message_id": "3EB0XYZ789...",
      "timestamp": "2025-11-06T10:31:00Z",
      "message_type": "text",
      "text_content": "Bom dia! Como vai?",
      "media_link": "",
      "quoted_message_id": "",
      "datajson": ""
    }
  ],
  "success": true
}
```

---

## Erros Possíveis

### Formato de data inválido
```json
{
  "code": 400,
  "error": "invalid date format. Use YYYY-MM-DD or 'today'",
  "success": false
}
```

### Chat JID obrigatório
```json
{
  "code": 400,
  "error": "chat_jid is required",
  "success": false
}
```

---

## Para Substituir em Produção

- `localhost:8080` → Seu host/IP do servidor
- `SEU_TOKEN_AQUI` → Token da sua instância
- `5511999999999@s.whatsapp.net` → Número do contato
- `120363123456789012@g.us` → JID do grupo

---

## Suporte a Bancos de Dados

✅ **PostgreSQL** - Totalmente suportado  
✅ **SQLite** - Totalmente suportado

A implementação detecta automaticamente o banco de dados e ajusta a query SQL.

---

## Como Testar

1. **Reinicie o servidor wuzapi** para aplicar as mudanças
2. **Use os exemplos acima** substituindo os valores
3. **Verifique os logs** para debug caso necessário

---

## Data de Implementação
**06 de Novembro de 2025**

Modificação realizada em: `/home/allansevero/wuzapi/handlers.go` (linhas 5849-5882)
