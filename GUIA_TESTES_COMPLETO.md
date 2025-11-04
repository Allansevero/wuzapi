# Guia de Testes - WuzAPI

## 🧪 TESTES COMPLETOS DO SISTEMA

### Pré-requisitos
- Sistema compilado: `./wuzapi`
- Banco de dados limpo ou com dados de teste
- Porta 8080 disponível

---

## 1. TESTE DE AUTENTICAÇÃO

### 1.1 Registro de Usuário
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "senha123",
    "name": "Minha Instância"
  }'
```

**Resultado Esperado:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "email": "teste@example.com"
}
```

### 1.2 Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "senha123"
  }'
```

**Resultado Esperado:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "email": "teste@example.com"
}
```

---

## 2. TESTE DE DASHBOARD (Interface Web)

### 2.1 Acessar Dashboard
1. Abra o navegador
2. Acesse: `http://localhost:8080/dashboard/user-dashboard-v2.html`
3. Se não estiver logado, será redirecionado para login
4. Faça login com as credenciais criadas

**Verificar:**
- ✅ Redirecionamento funciona
- ✅ Dashboard carrega
- ✅ E-mail aparece no cabeçalho
- ✅ Instância padrão aparece

### 2.2 Verificar Card de Instância
**Deve exibir:**
- Nome da instância
- Status: "Desconectado" (badge cinza)
- ID da instância (truncado)
- Status: "Não logado" (ícone X vermelho)
- Destino: "Não configurado"
- Botão "Conectar WhatsApp"
- Botão "Código de Pareamento"
- Botão "Config. Destino"

---

## 3. TESTE DE CONEXÃO WHATSAPP

### 3.1 Conectar via QR Code
1. No dashboard, clique em "Conectar WhatsApp"
2. Aguarde alguns segundos
3. QR Code deve aparecer
4. Abra WhatsApp no celular
5. Escaneie o QR Code

**Verificar:**
- ✅ Botão mostra "Conectando..."
- ✅ QR Code aparece em até 10 segundos
- ✅ Após escanear, status muda para "Conectado"
- ✅ Badge fica verde
- ✅ Número do WhatsApp aparece
- ✅ Botões mudam para "Desconectar"

### 3.2 Via API (obter QR Code)
```bash
# Primeiro, obtenha o token da instância no dashboard ou banco
TOKEN="seu-token-aqui"

# Conectar
curl -X POST http://localhost:8080/session/connect \
  -H "token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "Subscribe": ["Message"],
    "Immediate": true
  }'

# Obter QR Code
curl -X GET "http://localhost:8080/session/qr?token=$TOKEN"
```

### 3.3 Verificar Status
```bash
curl -X GET "http://localhost:8080/session/status?token=$TOKEN"
```

**Resultado Esperado (Conectado):**
```json
{
  "code": 200,
  "success": true,
  "data": {
    "connected": true,
    "loggedIn": true,
    "jid": "5511999999999:64@s.whatsapp.net"
  }
}
```

---

## 4. TESTE DE CONFIGURAÇÃO DE NÚMERO

### 4.1 Via Interface
1. No card da instância, clique em "Config. Destino"
2. Digite um número: `+5511999999999`
3. Clique em "Salvar"

**Verificar:**
- ✅ Modal fecha
- ✅ Mensagem de sucesso aparece
- ✅ Número aparece no card: "Destino: +5511999999999"

### 4.2 Via API
```bash
curl -X POST http://localhost:8080/session/destination-number \
  -H "token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "number": "+5511999999999"
  }'
```

**Resultado Esperado:**
```json
{
  "Details": "Destination number configured successfully",
  "Number": "+5511999999999"
}
```

### 4.3 Consultar Número
```bash
curl -X GET "http://localhost:8080/session/destination-number?token=$TOKEN"
```

---

## 5. TESTE DE MENSAGENS

### 5.1 Enviar Mensagem de Teste
```bash
curl -X POST http://localhost:8080/chat/send/text \
  -H "token: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "Phone": "5511988888888",
    "Body": "Mensagem de teste"
  }'
```

### 5.2 Verificar Armazenamento no Banco
```bash
# Conecte no banco
sqlite3 dbdata/users.db

# Execute
SELECT chat_jid, message_type, text_content, timestamp 
FROM message_history 
ORDER BY timestamp DESC 
LIMIT 5;
```

**Deve exibir:**
- Mensagens enviadas e recebidas
- Tipos corretos (text, image, etc)
- Conteúdo armazenado
- Timestamps corretos

---

## 6. TESTE DE ENVIO DIÁRIO

### 6.1 Envio Manual de Teste
```bash
curl -X POST http://localhost:8080/session/send-daily-test \
  -H "token: $TOKEN" \
  -H "Content-Type: application/json"
```

**Resultado Esperado:**
```json
{
  "success": true,
  "message": "Daily messages sent successfully",
  "instance_id": "uuid-da-instancia",
  "date": "2025-11-04"
}
```

### 6.2 Verificar Logs do Backend
```bash
tail -f wuzapi.log | grep -i "daily\|webhook"
```

**Deve exibir:**
```
Manual daily send triggered
Successfully sent daily messages to webhook
```

### 6.3 Verificar Payload no Webhook
O webhook deve receber:
```json
{
  "instance_id": "uuid",
  "date": "2025-11-04",
  "enviar_para": "+5511999999999",
  "conversations": [
    {
      "contact": "5511888888888@s.whatsapp.net",
      "messages": [
        {
          "sender_jid": "5511888888888@s.whatsapp.net",
          "message_type": "text",
          "text_content": "Olá!",
          "media_link": "",
          "timestamp": "2025-11-04T10:30:00Z",
          "data": {}
        }
      ]
    }
  ]
}
```

---

## 7. TESTE DE HISTÓRICO

### 7.1 Verificar Auto-Request
Após conectar o WhatsApp, verifique os logs:

```bash
tail -f wuzapi.log | grep -i history
```

**Deve exibir:**
```
Auto-requesting history sync after connection
History sync auto-requested successfully
```

### 7.2 Verificar Mensagens no Banco
```bash
sqlite3 dbdata/users.db

SELECT COUNT(*) as total_messages FROM message_history;
```

Deve ter mais de 0 mensagens após alguns minutos conectado.

---

## 8. TESTE DE MÚLTIPLAS INSTÂNCIAS

### 8.1 Criar Nova Instância via Interface
1. No dashboard, clique em "+ Nova Instância"
2. Digite nome: "Instância 2"
3. Digite número (opcional): "+5511888888888"
4. Clique em "Criar"

**Verificar:**
- ✅ Modal fecha
- ✅ Nova instância aparece no grid
- ✅ Cada uma com seu próprio card

### 8.2 Via API
```bash
# Use o token de autenticação JWT (não o token da instância)
AUTH_TOKEN="seu-token-jwt-do-login"

curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Instância Teste 3",
    "destination_number": "+5511777777777"
  }'
```

### 8.3 Listar Instâncias
```bash
curl -X GET http://localhost:8080/my/instances \
  -H "Authorization: Bearer $AUTH_TOKEN"
```

---

## 9. TESTE DE ISOLAMENTO DE USUÁRIOS

### 9.1 Criar Segundo Usuário
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario2@example.com",
    "password": "senha456",
    "name": "Instância User 2"
  }'
```

### 9.2 Verificar Isolamento
1. Faça login com usuário 2
2. Acesse dashboard
3. Deve ver apenas sua instância
4. Não deve ver instâncias do usuário 1

**Teste via API:**
```bash
# Login usuário 2
AUTH_TOKEN_2=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario2@example.com","password":"senha456"}' | jq -r '.token')

# Listar instâncias (deve ser diferente do usuário 1)
curl -X GET http://localhost:8080/my/instances \
  -H "Authorization: Bearer $AUTH_TOKEN_2"
```

---

## 10. TESTE DE DESCONEXÃO

### 10.1 Via Interface
1. No card da instância conectada
2. Clique em "Desconectar"
3. Confirme

**Verificar:**
- ✅ Modal de confirmação
- ✅ Status muda para "Desconectado"
- ✅ Badge fica cinza
- ✅ Botões voltam para "Conectar"

### 10.2 Via API
```bash
curl -X POST http://localhost:8080/session/logout \
  -H "token: $TOKEN"
```

---

## 11. TESTE DE CRON (Envio às 18h)

### 11.1 Ajustar Horário para Teste
Edite `daily_sender.go` temporariamente:

```go
// Linha 45 - mudar para próximo minuto
_, err := c.AddFunc("*/1 * * * *", func() { // A cada minuto para teste
```

### 11.2 Recompilar e Executar
```bash
go build -o wuzapi
./wuzapi
```

### 11.3 Verificar Logs
```bash
tail -f wuzapi.log | grep -i "daily\|18:00"
```

**Deve exibir a cada minuto:**
```
Starting daily message delivery at 18:00 Brasilia time
Successfully sent daily messages to webhook
```

### 11.4 Reverter Alteração
Após teste, volte para `"0 18 * * *"` e recompile.

---

## 12. CHECKLIST FINAL

### Interface
- [ ] Login funciona
- [ ] Cadastro funciona
- [ ] Dashboard carrega
- [ ] Cards aparecem em grid 3 colunas
- [ ] QR Code aparece ao conectar
- [ ] Status atualiza automaticamente
- [ ] Modal de número funciona
- [ ] Criar instância funciona
- [ ] Deletar instância funciona

### Backend
- [ ] Autenticação JWT funciona
- [ ] Isolamento entre usuários
- [ ] WhatsApp conecta via QR
- [ ] WhatsApp conecta via código
- [ ] Mensagens são armazenadas
- [ ] Histórico é buscado ao conectar
- [ ] Número de destino salva
- [ ] Envio manual funciona
- [ ] Cron está configurado

### Segurança
- [ ] Senhas hasheadas
- [ ] Tokens JWT válidos
- [ ] Cada usuário vê só seus dados
- [ ] Prepared statements (SQL)
- [ ] CORS configurado

---

## 🎉 RESULTADO ESPERADO

Se todos os testes passarem, o sistema está **100% FUNCIONAL** e pronto para uso!

### Problemas Comuns

**QR não aparece:**
- Aguarde até 15 segundos
- Verifique logs de erro
- Tente reconectar

**Database locked:**
- Pare processo antigo: `pkill -f wuzapi`
- Aguarde e reinicie

**Status não atualiza:**
- Polling automático leva até 15s
- Recarregue a página

**Webhook não recebe:**
- Verifique URL do webhook
- Teste com curl manualmente
- Veja logs de erro
