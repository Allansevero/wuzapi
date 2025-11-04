# 🚀 Guia Rápido - WuzAPI com Sistema de Usuários

## ✅ Atualização: Token Gerado Automaticamente!

Agora você **NÃO precisa mais fornecer um token** ao criar instâncias!
O sistema gera automaticamente tokens seguros para você.

---

## 📋 Como Usar

### 1️⃣ Cadastrar-se no Sistema

**Via Interface Web:**
1. Acesse: `http://localhost:8080/user-login.html`
2. Clique em "Cadastrar"
3. Preencha email e senha (mínimo 8 caracteres)
4. Clique em "Cadastrar"

**Via API:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com","password":"senha123456"}'
```

---

### 2️⃣ Fazer Login

**Via Interface Web:**
1. Volte para aba "Entrar"
2. Digite email e senha
3. Clique em "Entrar"

**Via API:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com","password":"senha123456"}'

# Retorna JWT token para usar nas APIs
```

---

### 3️⃣ Criar Instância (Token Gerado Automaticamente!)

**Via Dashboard:**
1. Clique em "+ Nova Instância"
2. Preencha **apenas**:
   - Nome da instância
   - Número de destino (opcional)
3. Clique em "Criar"
4. **Popup aparece mostrando o token gerado!**
5. Copie o token e guarde em local seguro

**Via API:**
```bash
curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer SEU_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Atendimento",
    "destination_number": "+5511999999999"
  }'

# Retorna:
# {
#   "data": {
#     "token": "abc123def456..." // TOKEN GERADO!
#   }
# }
```

---

### 4️⃣ Copiar Token da Instância

**No Dashboard:**
- Cada instância tem um botão "📋 Copiar" ao lado do token
- Clique para copiar automaticamente

---

### 5️⃣ Conectar ao WhatsApp

Use o **token da instância** (não o JWT!):

```bash
# Obter QR Code
curl "http://localhost:8080/session/qr?token=TOKEN_DA_INSTANCIA"

# Verificar status
curl "http://localhost:8080/session/status?token=TOKEN_DA_INSTANCIA"
```

---

### 6️⃣ Configurar Número de Destino

**Via Dashboard:**
1. Clique em "📱 Configurar Destino"
2. Digite o número: `+5511999999999`
3. Clique em "Salvar"

**Via API:**
```bash
curl -X POST http://localhost:8080/session/destination-number \
  -H "token: TOKEN_DA_INSTANCIA" \
  -H "Content-Type: application/json" \
  -d '{"number":"+5511999999999"}'
```

---

## 🎯 Diferenças Importantes

### Tipos de Token:

1. **JWT Token** (Login do usuário)
   - Retornado no `/auth/login`
   - Usado em endpoints `/my/*`
   - Header: `Authorization: Bearer {JWT}`
   - Dura 30 dias

2. **Instance Token** (Gerado automaticamente!)
   - Criado automaticamente ao criar instância
   - Usado em endpoints `/session/*` e `/chat/*`
   - Header: `token: {INSTANCE_TOKEN}`
   - Não expira

---

## ⏰ Envio Diário às 18h

- Todas as conversas do dia são enviadas às 18h (Brasília)
- Webhook fixo já configurado
- Payload inclui o campo `enviar_para` com o número configurado

---

## 📱 Endpoints Principais

### Autenticação:
- `POST /auth/register` - Cadastrar
- `POST /auth/login` - Login (retorna JWT)

### Instâncias (usa JWT):
- `POST /my/instances` - Criar (token gerado automaticamente!)
- `GET /my/instances` - Listar minhas instâncias
- `DELETE /my/instances/{id}` - Deletar

### WhatsApp (usa Instance Token):
- `GET /session/qr?token=...` - QR Code
- `GET /session/status?token=...` - Status
- `POST /chat/send/text` - Enviar mensagem
- `POST /session/destination-number` - Configurar número

---

## ✅ Resumo

1. ✅ Cadastre-se no sistema
2. ✅ Faça login (recebe JWT)
3. ✅ Crie instância (token gerado automaticamente!)
4. ✅ Copie e guarde o token da instância
5. ✅ Use token da instância para conectar ao WhatsApp
6. ✅ Configure número de destino
7. ✅ Aguarde envio diário às 18h

**Pronto!** Sistema funcionando! 🎉
