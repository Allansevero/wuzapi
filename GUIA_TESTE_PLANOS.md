# 🚀 GUIA RÁPIDO - Sistema de Planos

## ⚡ Início Rápido

### 1. Compilar e Executar

```bash
# Parar processo atual (se estiver rodando)
pkill wuzapi

# Fazer backup do binário atual
cp wuzapi wuzapi.backup.$(date +%Y%m%d)

# Compilar nova versão
go build -o wuzapi

# Executar
./wuzapi
```

### 2. Teste Manual Completo

#### Passo 1: Registrar Usuário
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@teste.com",
    "password": "senha12345"
  }'
```

**Resposta esperada:**
```json
{
  "code": 201,
  "message": "user registered successfully",
  "success": true
}
```

**O que acontece nos bastidores:**
- ✅ Usuário criado na tabela `system_users`
- ✅ Plano Gratuito (5 dias) atribuído automaticamente
- ✅ Instância padrão criada

#### Passo 2: Fazer Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@teste.com",
    "password": "senha12345"
  }'
```

**Resposta esperada:**
```json
{
  "code": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "email": "usuario@teste.com"
  },
  "success": true
}
```

**⚠️ IMPORTANTE:** Copie o token para usar nos próximos comandos!

#### Passo 3: Ver Assinatura Atual
```bash
# Substitua SEU_TOKEN pelo token recebido no login
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/my/subscription \
  -H "Authorization: Bearer $TOKEN"
```

**Resposta esperada:**
```json
{
  "success": true,
  "subscription": {
    "id": 1,
    "system_user_id": 1,
    "plan_id": 1,
    "started_at": "2025-11-04T08:30:00Z",
    "expires_at": "2025-11-09T08:30:00Z",
    "is_active": true,
    "plan": {
      "id": 1,
      "name": "Gratuito",
      "price": 0.00,
      "max_instances": 999999,
      "trial_days": 5
    }
  },
  "instance_count": 1,
  "is_expired": false
}
```

#### Passo 4: Listar Planos Disponíveis
```bash
curl -X GET http://localhost:8080/my/plans \
  -H "Authorization: Bearer $TOKEN"
```

**Resposta esperada:**
```json
{
  "success": true,
  "plans": [
    {
      "id": 1,
      "name": "Gratuito",
      "price": 0.00,
      "max_instances": 999999,
      "trial_days": 5,
      "is_active": true
    },
    {
      "id": 2,
      "name": "Pro",
      "price": 29.00,
      "max_instances": 5,
      "trial_days": 0,
      "is_active": true
    },
    {
      "id": 3,
      "name": "Analista",
      "price": 97.00,
      "max_instances": 12,
      "trial_days": 0,
      "is_active": true
    }
  ]
}
```

#### Passo 5: Fazer Upgrade para Pro
```bash
curl -X PUT http://localhost:8080/my/subscription \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_id": 2
  }'
```

**Resposta esperada:**
```json
{
  "success": true,
  "message": "Subscription updated successfully"
}
```

#### Passo 6: Testar Criação de Instância com Limite

```bash
# Criar instância 1
curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp 1",
    "destination_number": "+5511999999999"
  }'

# Criar instância 2
curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp 2",
    "destination_number": "+5511888888888"
  }'

# ... criar até 5 instâncias (limite do plano Pro)

# Tentar criar a 6ª (deve falhar)
curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp 6",
    "destination_number": "+5511777777777"
  }'
```

**Resposta esperada no 6º:**
```json
{
  "code": 403,
  "error": "You have reached the maximum number of instances for your plan. Please upgrade to create more.",
  "success": false
}
```

---

## 🌐 Teste via Interface Web

### 1. Acessar Login
```
http://localhost:8080/dashboard/login.html
```

### 2. Criar Conta
- Email: usuario@teste.com
- Senha: senha12345
- Clicar em "Registrar"

### 3. Fazer Login
- Email: usuario@teste.com
- Senha: senha12345
- Clicar em "Entrar"

### 4. Ver Dashboard
- Você será redirecionado para o dashboard
- Verá a instância padrão já criada

### 5. Ver Assinatura
- Clicar no botão "📊 Minha Assinatura" no topo
- Verá:
  - Plano atual: Gratuito
  - Dias restantes: 5
  - Instâncias: 1 / ∞
  - 3 cards com os planos disponíveis

### 6. Fazer Upgrade
- Clicar em "Fazer Upgrade" no card do plano Pro
- Confirmar
- Ver mudança imediata

---

## 🔍 Verificar no Banco de Dados

### SQLite (padrão)
```bash
sqlite3 dbdata/users.db

# Ver planos
SELECT * FROM plans;

# Ver assinaturas
SELECT 
  us.id,
  su.email,
  p.name as plan_name,
  us.started_at,
  us.expires_at,
  us.is_active
FROM user_subscriptions us
JOIN system_users su ON us.system_user_id = su.id
JOIN plans p ON us.plan_id = p.id;

# Ver instâncias por usuário
SELECT 
  su.email,
  COUNT(u.id) as instance_count,
  p.max_instances,
  p.name as plan_name
FROM system_users su
LEFT JOIN users u ON u.system_user_id = su.id
LEFT JOIN user_subscriptions us ON us.system_user_id = su.id AND us.is_active = 1
LEFT JOIN plans p ON us.plan_id = p.id
GROUP BY su.id;

.exit
```

---

## 📊 Cenários de Teste

### Cenário 1: Novo Usuário (Trial)
- ✅ Registro → Trial gratuito 5 dias
- ✅ Pode criar instâncias ilimitadas
- ✅ Após 5 dias → bloqueado
- ✅ Upgrade → desbloqueado

### Cenário 2: Limite de Instâncias
- ✅ Plano Pro → máximo 5 instâncias
- ✅ Criar 5 → OK
- ✅ Tentar criar 6 → Bloqueado
- ✅ Upgrade para Analista → pode criar 7 mais

### Cenário 3: Expiração
- ✅ Trial com 1 dia restante → alerta
- ✅ Trial expirado → bloqueio
- ✅ Upgrade → desbloqueio imediato

### Cenário 4: Downgrade
- ✅ Tem 10 instâncias no plano Analista
- ✅ Fazer downgrade para Pro (máx 5)
- ✅ Mantém as 10 existentes
- ✅ Mas não pode criar novas

---

## 🐛 Troubleshooting

### Erro: "database is locked"
```bash
# Parar o processo
pkill wuzapi

# Verificar se há processos travados
ps aux | grep wuzapi

# Remover locks
rm -f dbdata/*.wal dbdata/*.shm

# Reiniciar
./wuzapi
```

### Erro: "address already in use"
```bash
# Encontrar processo na porta 8080
lsof -i :8080

# Matar processo
kill -9 <PID>

# Ou usar pkill
pkill wuzapi

# Reiniciar
./wuzapi
```

### Erro: "no active subscription"
```bash
# Verificar no banco
sqlite3 dbdata/users.db "SELECT * FROM user_subscriptions WHERE system_user_id = 1;"

# Se não houver, criar manualmente
sqlite3 dbdata/users.db "INSERT INTO user_subscriptions (system_user_id, plan_id, started_at, expires_at, is_active) VALUES (1, 1, datetime('now'), datetime('now', '+5 days'), 1);"
```

### Token expirado
```bash
# Fazer login novamente
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@teste.com","password":"senha12345"}'

# Pegar novo token
```

---

## 📝 Logs Importantes

### Verificar criação de subscription
```bash
# Ao registrar, deve aparecer:
# "Default subscription created for new user"

# Verificar logs
tail -f wuzapi.log | grep subscription
```

### Verificar validação de limites
```bash
# Ao criar instância, deve aparecer:
# "Checking instance limit for user"
# "User can create: true/false"

tail -f wuzapi.log | grep "instance limit"
```

---

## ✅ Checklist de Verificação

Antes de colocar em produção:

- [ ] Compilação sem erros
- [ ] Migrations rodaram com sucesso
- [ ] 3 planos inseridos no banco
- [ ] Registro cria subscription automaticamente
- [ ] Login retorna token válido
- [ ] API `/my/subscription` funciona
- [ ] API `/my/plans` retorna 3 planos
- [ ] Validação de limite funciona
- [ ] Interface web carrega
- [ ] Upgrade de plano funciona
- [ ] Alertas de expiração aparecem
- [ ] Logs estão sendo gerados

---

## 🎉 Pronto!

Se todos os testes passaram, o sistema está **100% funcional** e pronto para uso!

**Próximo passo:** Integrar gateway de pagamento para automatizar cobranças.
