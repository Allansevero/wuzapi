# Lista de Alterações - Sistema de Planos e Melhorias

## ✅ STATUS GERAL: Backend 100% Implementado | Frontend Precisa de Modernização

---

## Alterações Implementadas

### 1. Sistema de Autenticação de Usuários ✅ COMPLETO
- ✅ Cada usuário tem email e senha para acessar
- ✅ Usuários veem apenas as instâncias relacionadas à sua conta
- ✅ Token de admin gerado automaticamente no cadastro/login
- ✅ Acesso direto ao dashboard após login (sem necessidade de preencher token)
- ✅ JWT authentication implementado
- ✅ Sistema de sessions implementado

### 2. Interface do Dashboard ✅ COMPLETO (Precisa Modernização)
- ✅ Removidas configurações do cabeçalho ao entrar na instância
- ✅ QR Code exibido corretamente ao clicar em "Conectar"
- ✅ Status de conexão atualizado em tempo real
- ✅ Instâncias exibidas em grid de 3 colunas com bordas arredondadas
- ✅ Status "Conectado" aparece apenas quando realmente conectado ao WhatsApp
- 🔨 **PENDENTE**: Aplicar novo design do HTML_FRONTEND_REPLIQUE.md

### 3. Sistema de Webhook Centralizado ✅ COMPLETO
- ✅ Webhook padrão do sistema: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- ✅ Webhook não aparece na configuração das instâncias (é transparente para o usuário)
- ✅ Envio diário automático às 18h (horário de Brasília) de todas as conversas do dia
- ✅ Sistema de cron job implementado (`daily_sender.go`)

### 4. Sistema de Número de Destino ✅ COMPLETO
- ✅ Campo `destination_number` adicionado ao banco de dados
- ✅ Número enviado junto com as mensagens no parâmetro "enviar_para"
- ✅ Parâmetro incluído no envio diário às 18h
- 🔨 **PENDENTE**: Interface para usuário editar o número na página "Seus Dados"

### 5. Histórico de Mensagens ✅ COMPLETO
- ✅ Sistema puxa últimas 100 mensagens por conversa ao conectar
- ✅ Armazena mensagens recebidas e enviadas após o login
- ✅ Histórico disponível para envio no compilado diário
- ✅ Tabela `message_history` criada e funcional
- ✅ Tabela `daily_conversations` para cache de conversas diárias

---

## Sistema de Planos - ✅ BACKEND COMPLETO

### 6. Sistema de Planos de Assinatura ✅ BACKEND COMPLETO

#### Plano Gratuito (Trial) - ✅ Implementado
- **Duração**: 5 dias
- **Limite**: Números ilimitados de WhatsApp (999999)
- **Preço**: R$ 0,00
- **Recursos**: Acesso completo a todas funcionalidades durante o período de trial
- **ID no Banco**: 1
- **Ativação**: Automática no registro do usuário

#### Plano Pro - ✅ Implementado
- **Preço**: R$ 29,00/mês
- **Limite**: 5 números de WhatsApp conectados
- **ID no Banco**: 2
- **Recursos**: 
  - Envio diário de conversas às 18h
  - Webhook centralizado
  - Armazenamento de histórico
  - Suporte por email

#### Plano Analista - ✅ Implementado
- **Preço**: R$ 97,00/mês
- **Limite**: 12 números de WhatsApp conectados
- **ID no Banco**: 3
- **Recursos**:
  - Envio diário de conversas às 18h
  - Webhook centralizado
  - Armazenamento de histórico
  - Suporte prioritário
  - Análises avançadas
  - Mais capacidade de instâncias

### ✅ Backend Implementado

#### 1. **Tabelas do Banco de Dados** - COMPLETO
   - ✅ **`plans`** - Armazena os 3 planos disponíveis
     * id, name, price, max_instances, trial_days, is_active, created_at
   
   - ✅ **`user_subscriptions`** - Assinaturas ativas dos usuários
     * id, system_user_id, plan_id, started_at, expires_at, is_active, created_at, updated_at
   
   - ✅ **`subscription_history`** - Histórico de assinaturas
     * id, system_user_id, plan_id, started_at, ended_at, created_at
   
   - ✅ **Índices criados**:
     * `idx_user_subscriptions_user` - Performance em consultas por usuário
     * `idx_user_subscriptions_active` - Filtro rápido de subs ativas
     * `idx_subscription_history_user` - Histórico por usuário

#### 2. **Lógica de Negócios** (`subscriptions.go`) - COMPLETO
   - ✅ `CreateDefaultSubscription()` - Cria trial de 5 dias automaticamente
   - ✅ `GetActiveSubscription()` - Retorna subscription + plan details
   - ✅ `UpdateSubscription()` - Troca de plano (desativa antiga, cria nova)
   - ✅ `CheckSubscriptionExpired()` - Valida se subscription expirou
   - ✅ `GetUserInstanceCount()` - Conta instâncias do usuário
   - ✅ `CanCreateInstance()` - Valida se pode criar mais instâncias
   - ✅ `GetAllPlans()` - Lista todos os planos disponíveis

#### 3. **Integração com Autenticação** (`auth.go`) - COMPLETO
   - ✅ Linha 219: `CreateDefaultSubscription()` chamado no registro
   - ✅ Linha 230: Cria instância padrão "Instância Padrão" automaticamente
   - ✅ Usuário registrado já sai com trial ativo + 1 instância criada

#### 4. **Endpoints API REST** (`routes.go`) - COMPLETO
   - ✅ `GET /user/subscription` - Detalhes da assinatura atual do usuário
   - ✅ `PUT /user/subscription` - Atualizar plano (upgrade/downgrade)
   - ✅ `GET /user/plans` - Listar todos os planos disponíveis
   - ✅ Todas rotas protegidas por JWT authentication

#### 5. **Handlers Implementados** (`handlers.go`)
   - ✅ `GetUserSubscriptionHandler()` - Retorna subscription details
   - ✅ `UpdateUserSubscriptionHandler()` - Processa upgrade de plano
   - ✅ `GetPlansHandler()` - Lista planos com informações completas

### 🔨 Frontend - PENDENTE (Precisa ser Implementado)

#### O que falta fazer:

1. **Modernizar Dashboard** (Ver `HTML_FRONTEND_REPLIQUE.md` para design)
   - [ ] Aplicar novo design com Tailwind CSS
   - [ ] Implementar sidebar com logo "metrizap"
   - [ ] Criar navegação entre "Contas conectadas" e "Seus dados"
   - [ ] Mostrar indicador de progresso de instâncias usadas/disponíveis
   - [ ] Grid de cards 3 colunas responsivo

2. **Página "Seus Dados"**
   - [ ] Formulário com dados do usuário (nome, email, senha)
   - [ ] Campo editável "Quero receber análises no:" (destination_number)
   - [ ] Seção "Plano atual" com cards dos 3 planos
   - [ ] Destacar plano ativo com borda verde
   - [ ] Botão "Fazer upgrade" nos planos não ativos
   - [ ] Indicar features de cada plano

3. **JavaScript - API Client**
   - [ ] Criar `api-client.js` com funções:
     * `getActiveSubscription()` - GET /user/subscription
     * `getAllPlans()` - GET /user/plans
     * `upgradePlan(planId)` - PUT /user/subscription
     * `checkCanCreateInstance()` - Valida antes de criar instância
   - [ ] Implementar tratamento de erros
   - [ ] Implementar loading states

4. **Validação de Limites**
   - [ ] Ao clicar "Adicionar WhatsApp", verificar `checkCanCreateInstance()`
   - [ ] Se limite atingido, mostrar modal de upgrade
   - [ ] Modal deve explicar o motivo e oferecer upgrade
   - [ ] Não permitir criação se subscription expirada

5. **Indicadores Visuais**
   - [ ] Barra de progresso mostrando X/Y instâncias usadas
   - [ ] Badge mostrando plano atual (Trial 5 dias, Pro, Analista)
   - [ ] Contador regressivo para trial (ex: "3 dias restantes")
   - [ ] Alertas de expiração próxima

6. **Fluxos de Usuário**
   - [ ] **Registro → Trial**: Auto-create subscription + instância padrão
   - [ ] **Expiração Trial**: Bloquear novas instâncias, forçar escolha de plano
   - [ ] **Upgrade**: Mostrar confirmação, atualizar UI imediatamente
   - [ ] **Limite Atingido**: Modal explicativo com botão para upgrade

---

## Estrutura de Dados (✅ Banco de Dados)

### Tabela: plans ✅ CRIADA E POPULADA
```sql
CREATE TABLE plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,                    -- Gratuito, Pro, Analista
    price REAL NOT NULL,                   -- 0.00, 29.00, 97.00
    max_instances INTEGER NOT NULL,        -- 999999, 5, 12
    trial_days INTEGER DEFAULT 0,          -- 5, 0, 0
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Dados inseridos automaticamente:
-- 1, 'Gratuito', 0.00, 999999, 5
-- 2, 'Pro', 29.00, 5, 0
-- 3, 'Analista', 97.00, 12, 0
```

### Tabela: user_subscriptions ✅ CRIADA
```sql
CREATE TABLE user_subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_user_id INTEGER NOT NULL REFERENCES system_users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES plans(id),
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,                   -- NULL = mensal recorrente, data = trial
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Índices:
CREATE INDEX idx_user_subscriptions_user ON user_subscriptions (system_user_id);
CREATE INDEX idx_user_subscriptions_active ON user_subscriptions (system_user_id, is_active, expires_at);
```

### Tabela: subscription_history ✅ CRIADA
```sql
CREATE TABLE subscription_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    system_user_id INTEGER NOT NULL REFERENCES system_users(id) ON DELETE CASCADE,
    plan_id INTEGER NOT NULL REFERENCES plans(id),
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscription_history_user ON subscription_history (system_user_id);
```

### Outras Tabelas Relevantes ✅ JÁ EXISTEM

**system_users** - Usuários do sistema
```sql
CREATE TABLE system_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**users** - Instâncias WhatsApp (modificado)
```sql
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    token TEXT NOT NULL,
    webhook TEXT NOT NULL DEFAULT '',
    jid TEXT NOT NULL DEFAULT '',
    qrcode TEXT NOT NULL DEFAULT '',
    connected INTEGER,
    expiration INTEGER,
    events TEXT NOT NULL DEFAULT '',
    proxy_url TEXT DEFAULT '',
    system_user_id INTEGER REFERENCES system_users(id) ON DELETE CASCADE,  -- ✅ NOVO
    destination_number TEXT DEFAULT '',                                      -- ✅ NOVO
    history INTEGER DEFAULT 0,
    -- ... (outros campos S3, HMAC, etc)
);
```

**daily_conversations** - Cache de conversas diárias ✅
```sql
CREATE TABLE daily_conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    date DATE NOT NULL,
    chat_jid TEXT NOT NULL,
    contact TEXT NOT NULL,
    messages TEXT NOT NULL,  -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date, chat_jid)
);
```

---

## Fluxo de Usuário Implementado

### 1. ✅ Registro de Novo Usuário (BACKEND)
```
1. POST /auth/register { email, password }
2. Backend cria system_user
3. Backend cria user_subscription com plan_id=1 (Trial 5 dias)
4. Backend define expires_at = NOW() + 5 dias
5. Backend cria instância padrão "Instância Padrão"
6. Backend retorna JWT token
7. Frontend redireciona para dashboard com token
```

### 2. ✅ Validação de Criação de Instância (BACKEND)
```
1. Frontend chama: POST /user/instances { name }
2. Backend valida com CanCreateInstance(system_user_id)
3. Backend verifica:
   - Subscription está ativa? (is_active = true)
   - Subscription não expirou? (expires_at > NOW() ou NULL)
   - Contagem atual < max_instances do plano?
4. Se OK: Cria instância e retorna sucesso
5. Se NÃO: Retorna erro 403 com motivo
```

### 3. ✅ Upgrade de Plano (BACKEND)
```
1. Frontend chama: PUT /user/subscription { plan_id: 2 }
2. Backend inicia transação
3. Backend desativa subscription atual (is_active = false)
4. Backend cria nova subscription com novo plan_id
5. Backend define expires_at = NULL (mensal recorrente)
6. Backend confirma transação
7. Retorna nova subscription com plan details
```

### 4. ✅ Verificação de Expiração (BACKEND)
```
1. Cron job ou verificação on-demand
2. Backend chama CheckSubscriptionExpired(system_user_id)
3. Backend atualiza: SET is_active = false WHERE expires_at < NOW()
4. CanCreateInstance() automaticamente retorna false se expirado
```

### 5. 🔨 Envio Diário 18h (IMPLEMENTADO mas precisa testar)
```
1. Cron job roda às 18h BRT (daily_sender.go)
2. Para cada usuário ativo:
   - Busca todas conversas do dia na tabela daily_conversations
   - Monta payload JSON com:
     * instanceName
     * conversations (array de conversas)
     * enviar_para (destination_number do usuário)
3. Envia para webhook centralizado:
   POST https://n8n-webhook.fmy2un.easypanel.host/webhook/...
4. Limpa cache diário após envio
```

---

## APIs REST Disponíveis

### Autenticação
```bash
# Registro
POST /auth/register
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "12345678"
}

# Login
POST /auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "12345678"
}
# Retorna: { "token": "JWT_TOKEN", "email": "user@example.com" }
```

### Subscriptions (Protegido por JWT)
```bash
# Ver subscription atual
GET /user/subscription
Authorization: Bearer JWT_TOKEN
# Retorna: { subscription: {...}, plan: {...} }

# Listar planos disponíveis
GET /user/plans
Authorization: Bearer JWT_TOKEN
# Retorna: [{ id: 1, name: "Gratuito", price: 0.00, ... }, ...]

# Fazer upgrade
PUT /user/subscription
Authorization: Bearer JWT_TOKEN
Content-Type: application/json
{
  "plan_id": 2
}
# Retorna: { subscription: {...}, plan: {...} }
```

### Instâncias (Já existentes)
```bash
# Listar instâncias do usuário
GET /user/instances
Authorization: Bearer JWT_TOKEN

# Criar nova instância (valida limite automaticamente)
POST /user/instances
Authorization: Bearer JWT_TOKEN
Content-Type: application/json
{
  "name": "Minha Instância"
}

# Deletar instância
DELETE /user/instances/{instance_id}
Authorization: Bearer JWT_TOKEN
```

---

## Como Testar o Backend

### 1. Testar Registro com Trial Automático
```bash
# 1. Registrar usuário
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"12345678"}'

# 2. Login para pegar token
TOKEN=$(curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"12345678"}' \
  | jq -r '.data.token')

# 3. Ver subscription (deve ser trial 5 dias)
curl -X GET http://localhost:8080/user/subscription \
  -H "Authorization: Bearer $TOKEN" | jq

# Esperado:
# {
#   "subscription": {
#     "plan_id": 1,
#     "started_at": "2025-11-04...",
#     "expires_at": "2025-11-09...",  # 5 dias depois
#     "is_active": true
#   },
#   "plan": {
#     "id": 1,
#     "name": "Gratuito",
#     "price": 0,
#     "max_instances": 999999,
#     "trial_days": 5
#   }
# }
```

### 2. Testar Limite de Instâncias (Pro Plan)
```bash
# 1. Fazer upgrade para Pro (limite 5)
curl -X PUT http://localhost:8080/user/subscription \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plan_id":2}' | jq

# 2. Criar 5 instâncias
for i in {1..5}; do
  curl -X POST http://localhost:8080/user/instances \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Instância $i\"}"
done

# 3. Tentar criar a 6ª (deve falhar)
curl -X POST http://localhost:8080/user/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Instância 6"}' | jq

# Esperado: HTTP 403 Forbidden
# {
#   "code": 403,
#   "error": "instance limit reached for your plan",
#   "success": false
# }
```

### 3. Testar Listagem de Planos
```bash
curl -X GET http://localhost:8080/user/plans \
  -H "Authorization: Bearer $TOKEN" | jq

# Esperado:
# [
#   { "id": 1, "name": "Gratuito", "price": 0, "max_instances": 999999, "trial_days": 5 },
#   { "id": 2, "name": "Pro", "price": 29, "max_instances": 5, "trial_days": 0 },
#   { "id": 3, "name": "Analista", "price": 97, "max_instances": 12, "trial_days": 0 }
# ]
```

---

## Próximos Passos

### Prioridade 1: Modernizar Frontend ⭐⭐⭐
1. **Criar novo dashboard** (`user-dashboard-v3.html`)
   - Usar design do `HTML_FRONTEND_REPLIQUE.md`
   - Aplicar Tailwind CSS
   - Implementar sidebar moderna
   - Grid 3 colunas responsivo

2. **Criar API client** (`api-client.js`)
   - Centralizar todas chamadas de API
   - Implementar error handling
   - Adicionar loading states
   - Cache de dados quando apropriado

3. **Página Seus Dados**
   - Formulário de perfil
   - Seleção de planos visual
   - Campo destination_number editável
   - Indicador de uso de instâncias

### Prioridade 2: Validações de Limite no Frontend ⭐⭐
1. Ao clicar "Adicionar WhatsApp":
   - Verificar `checkCanCreateInstance()`
   - Se bloqueado, mostrar modal de upgrade
   - Modal deve ter link direto para planos

2. Indicadores visuais:
   - Barra de progresso (X/Y instâncias)
   - Badge do plano atual
   - Countdown para expiração de trial

### Prioridade 3: Integração de Pagamento ⭐ (Futuro)
1. Integrar Stripe ou Mercado Pago
2. Webhooks para confirmação de pagamento
3. Renovação automática de assinaturas
4. Emissão de notas fiscais

---

## Observações Importantes

### ✅ O que JÁ FUNCIONA
1. **Registro automático** com trial de 5 dias
2. **Validação de limites** no backend
3. **APIs REST** completas e testáveis
4. **Upgrade de planos** funcionando
5. **Webhook centralizado** configurado
6. **Envio diário 18h** implementado (cron job)
7. **Histórico de mensagens** sendo armazenado

### 🔨 O que PRECISA FAZER
1. **Frontend moderno** com Tailwind CSS
2. **Validação visual** de limites
3. **Página de planos** interativa
4. **Indicadores de uso** em tempo real
5. **Experiência de upgrade** fluida

### ⚠️ Limitações Conhecidas
1. **Sistema de pagamento** não integrado (manual por enquanto)
2. **Renovação automática** não implementada
3. **Emails de notificação** não configurados
4. **Dashboard analytics** básico

---

## Conclusão

✅ **Backend**: 100% funcional, todas APIs prontas
🔨 **Frontend**: Funcional mas precisa modernização
📋 **Próximo passo**: Implementar novo frontend seguindo `HTML_FRONTEND_REPLIQUE.md`

O sistema de planos está completamente operacional no backend. Qualquer usuário que se registrar já recebe automaticamente um trial de 5 dias com instâncias ilimitadas. O backend valida corretamente os limites e permite upgrades. 

O que falta é principalmente **melhorar a interface visual** e tornar o processo de upgrade mais intuitivo e atraente para o usuário final.

**Documentação de referência:**
- `IMPLEMENTACAO_FRONTEND_PLANOS.md` - Guia completo de implementação do frontend
- `HTML_FRONTEND_REPLIQUE.md` - Design a ser replicado
- `subscriptions.go` - Lógica de negócios dos planos
- `migrations.go` - Estrutura do banco de dados
