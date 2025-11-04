# Teste Manual de Envio Diário de Mensagens

## Visão Geral

Foi implementado um endpoint para testar manualmente o envio compilado de mensagens diárias para o webhook.

## Endpoint

```
POST /session/send-daily-test
```

### Parâmetros Query String

- `token` (obrigatório): Token de autenticação do usuário
- `instance_id` (opcional): ID da instância específica. Se não fornecido, usa o ID do usuário autenticado
- `date` (opcional): Data no formato `YYYY-MM-DD`. Se não fornecido, usa a data atual

### Exemplos de Uso

#### 1. Usando cURL diretamente

```bash
# Enviar mensagens de hoje da instância atual
curl -X POST "http://localhost:8080/session/send-daily-test?token=SEU_TOKEN"

# Enviar mensagens de uma data específica
curl -X POST "http://localhost:8080/session/send-daily-test?token=SEU_TOKEN&date=2025-11-03"

# Enviar mensagens de uma instância específica
curl -X POST "http://localhost:8080/session/send-daily-test?token=SEU_TOKEN&instance_id=ID_DA_INSTANCIA"
```

#### 2. Usando o Script de Teste

Foi criado um script bash para facilitar os testes:

```bash
# Enviar mensagens de hoje
./test_daily_send.sh SEU_TOKEN

# Enviar mensagens de uma instância específica
./test_daily_send.sh SEU_TOKEN ID_DA_INSTANCIA

# Enviar mensagens de uma data específica
./test_daily_send.sh SEU_TOKEN ID_DA_INSTANCIA 2025-11-03
```

#### 3. Usando JavaScript (Frontend)

```javascript
async function testDailySend(token, instanceId = null, date = null) {
    let url = `/session/send-daily-test?token=${token}`;
    
    if (instanceId) {
        url += `&instance_id=${instanceId}`;
    }
    
    if (date) {
        url += `&date=${date}`;
    }
    
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        const data = await response.json();
        console.log('Resultado:', data);
        return data;
    } catch (error) {
        console.error('Erro ao enviar:', error);
        throw error;
    }
}

// Exemplo de uso
testDailySend('seu-token-aqui')
    .then(result => console.log('Sucesso:', result))
    .catch(error => console.error('Erro:', error));
```

## Formato da Resposta

### Sucesso

```json
{
    "success": true,
    "message": "Daily messages sent successfully",
    "instance_id": "507a6d45c765c6ae5b720e3caa94fca2",
    "date": "2025-11-03"
}
```

### Erro

```json
{
    "success": false,
    "error": "Descrição do erro"
}
```

## Estrutura dos Dados Enviados ao Webhook

O payload enviado para o webhook tem a seguinte estrutura:

```json
{
    "instance_id": "507a6d45c765c6ae5b720e3caa94fca2",
    "date": "2025-11-03",
    "conversations": [
        {
            "contact": "5551999999999@s.whatsapp.net",
            "messages": [
                {
                    "sender_jid": "5551999999999@s.whatsapp.net",
                    "message_type": "text",
                    "text_content": "Olá, tudo bem?",
                    "media_link": "",
                    "timestamp": "2025-11-03T14:30:00-03:00",
                    "data": {
                        "... dados adicionais da mensagem ..."
                    }
                },
                {
                    "sender_jid": "5551888888888@s.whatsapp.net",
                    "message_type": "text",
                    "text_content": "Tudo ótimo!",
                    "media_link": "",
                    "timestamp": "2025-11-03T14:31:00-03:00"
                }
            ]
        }
    ],
    "enviar_para": "5551999999999"
}
```

### Campos do Payload

- **instance_id**: ID único da instância do WhatsApp
- **date**: Data das mensagens no formato YYYY-MM-DD
- **conversations**: Array de conversas do dia
  - **contact**: JID do contato (número@s.whatsapp.net)
  - **messages**: Array de mensagens da conversa
    - **sender_jid**: JID de quem enviou a mensagem
    - **message_type**: Tipo da mensagem (text, image, video, audio, etc)
    - **text_content**: Conteúdo textual da mensagem
    - **media_link**: Link da mídia (se houver)
    - **timestamp**: Data/hora da mensagem em ISO 8601
    - **data**: Dados JSON adicionais da mensagem original
- **enviar_para**: Número configurado para receber as notificações

## Webhook Padrão

O webhook é fixo no sistema e não pode ser alterado:

```
https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5
```

## Como Obter o Token

1. Faça login no sistema através de `/auth/login`
2. O token será retornado na resposta do login
3. O token também é gerado automaticamente no cadastro

Exemplo:

```bash
# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seu-email@exemplo.com",
    "password": "sua-senha"
  }'
```

Resposta:
```json
{
    "success": true,
    "token": "cc52b16537c624c981b5fb54f83e30cc",
    "user_id": "507a6d45c765c6ae5b720e3caa94fca2",
    "email": "seu-email@exemplo.com"
}
```

## Notas Importantes

1. **Horário Automático**: O sistema envia automaticamente às 18:00 (horário de Brasília)
2. **Apenas Mensagens do Dia**: Só são enviadas mensagens da data especificada
3. **Agrupamento por Conversa**: As mensagens são agrupadas por contato/chat
4. **Ordem Cronológica**: As mensagens dentro de cada conversa são ordenadas por timestamp
5. **Sem Mensagens**: Se não houver mensagens para a data, nenhum envio é feito

## Troubleshooting

### Erro 401 - Unauthorized
- Verifique se o token está correto
- Certifique-se de que está autenticado

### Erro 500 - Internal Server Error
- Verifique os logs do servidor
- Confirme que há mensagens no histórico para a data especificada
- Verifique a conexão com o banco de dados

### Webhook não recebe dados
- Verifique se o webhook está online
- Confirme que o URL do webhook está correto no código
- Verifique os logs do servidor para erros de envio

## Logs

Para visualizar os logs do envio:

```bash
tail -f wuzapi.log | grep "daily"
```

Ou se estiver rodando no console:

```bash
./wuzapi | grep -i daily
```

## Testando o Webhook Localmente

Se quiser testar com um webhook local, você pode usar ferramentas como:

1. **webhook.site** - https://webhook.site
2. **ngrok** - Para expor localhost
3. **RequestBin** - Para capturar requisições

Exemplo com webhook.site:

```bash
# 1. Vá em https://webhook.site e copie sua URL única
# 2. Temporariamente altere FIXED_WEBHOOK_URL em constants.go
# 3. Recompile e teste
```

---

**Data de Criação**: 2025-11-03  
**Versão**: 1.0
