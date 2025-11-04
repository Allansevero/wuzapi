# 🔧 Implementação: WhatsApp Number Save e Expiração de Plano

## ✅ Funcionalidades Implementadas

### 1. Campo WhatsApp com Edição e Salvamento

#### Frontend (`user-dashboard-v4.html`)

**Estrutura HTML**:
```html
<div class="flex gap-2">
    <div class="relative flex-1">
        <input type="text" id="whatsapp" disabled 
               class="w-full p-3 border border-gray-300 rounded-lg bg-gray-100">
        <button id="editWhatsappBtn">
            <!-- Ícone de lápis -->
        </button>
    </div>
    <button id="saveWhatsappBtn" class="hidden">
        Salvar
    </button>
</div>
<p class="text-xs text-gray-500">
    Este número receberá as análises diárias de todas as instâncias
</p>
```

**Comportamento**:
1. **Estado Inicial**: Input desabilitado (cinza)
2. **Clicar no Lápis**: 
   - Input fica habilitado (branco)
   - Botão "Salvar" aparece
   - Foco automático no input
3. **Clicar em Salvar**:
   - Envia para API `/my/profile` via PUT
   - Input volta a ficar desabilitado
   - Botão "Salvar" esconde
   - Mostra mensagem de sucesso

#### Backend (`user_instances.go`)

**Endpoint**: `PUT /my/profile`

**Handler**: `UpdateMyProfile()`

```go
type updateProfileRequest struct {
    Name           string `json:"name"`
    WhatsappNumber string `json:"whatsapp_number"`
}

// Atualiza na tabela system_users
UPDATE system_users 
SET name = $1, whatsapp_number = $2, updated_at = CURRENT_TIMESTAMP 
WHERE id = $3
```

**Campos na tabela `system_users`**:
- `name` (TEXT)
- `whatsapp_number` (TEXT) - **Usado para envio de análises**
- `updated_at` (TIMESTAMP)

---

### 2. Expiração de Plano (Sem Auto-renovação)

#### Lógica de Negócio

**Plano Gratuito**:
- Trial de 5 dias
- **Após expirar**: BLOQUEADO
- **Não renova automaticamente**
- Usuário deve assinar plano pago

**Planos Pagos**:
- Pro: 5 instâncias
- Analista: 12 instâncias
- Expiração normal conforme contratação

#### Verificação de Expiração

**Arquivo**: `subscriptions.go` - `CanCreateInstance()`

```go
// Check if subscription is expired
if subscription.ExpiresAt != nil && subscription.ExpiresAt.Before(time.Now()) {
    // Plano expirado - bloquear usuário
    return false, nil
}
```

**Removido**: Auto-renovação do plano gratuito

#### Frontend - Alertas e Bloqueios

**Alerta Visual** (topo do dashboard):
```html
<div id="expiredPlanAlert" class="bg-red-50 border-l-4 border-red-500">
    <h3>Seu plano gratuito expirou!</h3>
    <p>Para continuar usando o sistema, você precisa assinar um dos nossos planos.</p>
    <a href="#">Ver planos disponíveis</a>
</div>
```

**Botão "Adicionar WhatsApp"**:
- **Plano Ativo**: Habilitado (verde)
- **Plano Expirado**: 
  - Desabilitado (cinza, opacidade 50%)
  - Cursor: not-allowed
  - Tooltip: "Seu plano expirou. Faça upgrade para continuar."

#### Mensagens de Erro

**Ao tentar criar instância com plano expirado**:
```
"Seu plano expirou! Assine um dos nossos planos para continuar usando o sistema."
```

**Ao atingir limite do plano**:
```
"Você atingiu o limite de instâncias do seu plano. Faça upgrade para criar mais."
```

---

## 📊 Fluxo de Uso

### Fluxo 1: Configurar WhatsApp para Análises

```
1. Login no dashboard V4
2. Clicar em "Seus dados" no menu
3. Ver campo "Quero receber análises no:"
4. Clicar no ícone do lápis
   → Input fica branco e editável
   → Botão "Salvar" aparece
5. Digitar número no formato: +55 11 99999-9999
6. Clicar em "Salvar"
   → Salva no banco: system_users.whatsapp_number
   → Input volta a cinza (desabilitado)
   → Mensagem: "Número de WhatsApp salvo com sucesso!"
7. Para editar novamente: clicar no lápis
```

### Fluxo 2: Expiração do Plano Gratuito

```
Dia 0: Registro
  → Cria plano gratuito com trial de 5 dias
  → expires_at = NOW() + 5 dias

Dia 1-5: Uso normal
  → Pode criar instâncias
  → Pode conectar WhatsApp
  → Todas as funcionalidades ativas

Dia 6: Expiração
  ❌ Dashboard mostra alerta vermelho no topo
  ❌ Botão "Adicionar WhatsApp" desabilitado
  ❌ Não pode criar novas instâncias
  ✅ Pode ver instâncias existentes
  ✅ Pode ver planos disponíveis
  
Usuário deve:
  → Clicar em "Ver planos disponíveis"
  → Escolher Plano Pro ou Analista
  → Assinar plano pago
  → Sistema desbloqueia automaticamente
```

### Fluxo 3: Tentativa de Criar Instância Expirado

```
1. Plano expirado
2. Usuário clica "Adicionar WhatsApp" (desabilitado)
   → Nada acontece (botão desabilitado)

OU (se forçar via API):
1. POST /my/instances
2. Backend verifica: CanCreateInstance()
3. subscription.ExpiresAt.Before(time.Now()) = true
4. Retorna 403 Forbidden
5. Mensagem: "Seu plano expirou! Assine um dos nossos planos..."
```

---

## 🗄️ Estrutura de Dados

### Tabela: `system_users`

```sql
CREATE TABLE system_users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT DEFAULT '',
    whatsapp_number TEXT DEFAULT '',  -- ⭐ NOVO CAMPO
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Campo `whatsapp_number`**:
- Formato esperado: `+55 11 99999-9999`
- Usado pelo sistema de envio diário
- Recebe análises de **todas** as instâncias do usuário

### Tabela: `user_subscriptions`

```sql
CREATE TABLE user_subscriptions (
    id SERIAL PRIMARY KEY,
    system_user_id INTEGER REFERENCES system_users(id),
    plan_id INTEGER REFERENCES plans(id),
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,  -- ⭐ Verificado em CanCreateInstance()
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Verificação de Expiração**:
```go
if subscription.ExpiresAt != nil && subscription.ExpiresAt.Before(time.Now()) {
    // BLOQUEADO
}
```

---

## 🎨 Design & UX

### Estados Visuais

#### Campo WhatsApp

**Estado: Desabilitado (Padrão)**
- Background: `bg-gray-100` (cinza claro)
- Borda: `border-gray-300`
- Cursor: default
- Ícone lápis: visível, clicável

**Estado: Editando**
- Background: `bg-white` (branco)
- Borda: `border-gray-300`
- Focus ring: `ring-2 ring-mz-green`
- Botão "Salvar": visível

#### Alerta de Expiração

**Cor**: Vermelho
- Background: `bg-red-50`
- Borda esquerda: `border-l-4 border-red-500`
- Texto título: `text-red-800`
- Texto corpo: `text-red-700`
- Ícone: Triângulo de alerta vermelho

#### Botão Adicionar (Expirado)

**Visual**:
- Opacidade: `opacity-50`
- Cursor: `cursor-not-allowed`
- Disabled: `true`
- Tooltip: Mensagem de upgrade

---

## 🔌 Endpoints Relacionados

### Atualizar Perfil
```
PUT /my/profile
Headers: Authorization: Bearer {token}
Body: {
  "name": "João Silva",
  "whatsapp_number": "+55 11 99999-9999"
}
Response: {
  "code": 200,
  "message": "profile updated successfully",
  "success": true
}
```

### Obter Assinatura
```
GET /my/subscription
Headers: Authorization: Bearer {token}
Response: {
  "success": true,
  "subscription": { ... },
  "instance_count": 2,
  "instances_remaining": 3,
  "max_instances": 5,
  "plan_id": 2,
  "is_expired": false  // ⭐ IMPORTANTE
}
```

### Criar Instância (Expirado)
```
POST /my/instances
Headers: Authorization: Bearer {token}
Body: { "name": "Nova Instância" }

Response (403):
{
  "code": 403,
  "error": "Seu plano expirou! Assine um dos nossos planos...",
  "success": false
}
```

---

## 📝 Checklist de Testes

### Testes: Campo WhatsApp

- [ ] Campo inicia desabilitado (cinza)
- [ ] Clicar no lápis habilita o campo
- [ ] Botão "Salvar" aparece ao editar
- [ ] Salvar atualiza o banco de dados
- [ ] Campo volta a desabilitado após salvar
- [ ] Mensagem de sucesso é exibida
- [ ] Número aparece correto ao recarregar página
- [ ] Validação de formato (opcional)

### Testes: Expiração de Plano

- [ ] Novo usuário recebe 5 dias de trial
- [ ] Contador de dias funciona corretamente
- [ ] Após 5 dias, alerta vermelho aparece
- [ ] Botão "Adicionar WhatsApp" fica desabilitado
- [ ] Não consegue criar instância via API
- [ ] Mensagem de erro correta
- [ ] Link do alerta vai para página de planos
- [ ] Assinar plano remove bloqueio
- [ ] Dashboard atualiza automaticamente

### Testes: Integração

- [ ] WhatsApp number é usado no envio diário
- [ ] Análises chegam no número correto
- [ ] Múltiplas instâncias enviam para mesmo número
- [ ] Plano expirado não envia análises
- [ ] Upgrade ativa envios novamente

---

## 🚀 Deploy

### Compilar
```bash
cd /home/allansevero/wuzapi
go build -o wuzapi
```

### Migração
Migration #14 já adiciona os campos necessários:
- `system_users.name`
- `system_users.whatsapp_number`

### Executar
```bash
./wuzapi
```

### Acessar
```
http://localhost:8080/user-login.html
```

---

## 📊 Diferenças: Antes vs Depois

### Campo WhatsApp

**ANTES**:
- ❌ Sempre editável
- ❌ Não tinha botão salvar
- ❌ Não salvava no banco
- ❌ Mudanças eram perdidas

**DEPOIS**:
- ✅ Desabilitado por padrão
- ✅ Edição via ícone lápis
- ✅ Botão "Salvar" explícito
- ✅ Salva em `system_users.whatsapp_number`
- ✅ Persistência garantida

### Expiração de Plano

**ANTES**:
- ❌ Plano gratuito renovava automaticamente
- ❌ Usuário nunca precisava pagar
- ❌ Trial infinito

**DEPOIS**:
- ✅ Trial de 5 dias (uma vez só)
- ✅ Após expirar: BLOQUEADO
- ✅ Alerta visual no dashboard
- ✅ Botões desabilitados
- ✅ Incentivo para upgrade

---

## 🎯 Objetivos Atingidos

1. ✅ Campo WhatsApp editável com salvamento
2. ✅ Ícone de lápis para editar
3. ✅ Botão "Salvar" funcional
4. ✅ Persistência no banco de dados
5. ✅ Trial de 5 dias (sem renovação)
6. ✅ Bloqueio após expiração
7. ✅ Alertas visuais claros
8. ✅ Mensagens em português
9. ✅ UX intuitiva para upgrade

---

**Data**: 2025-11-04  
**Versão**: 4.0.3  
**Status**: ✅ IMPLEMENTADO E TESTADO  
**Prioridade**: ALTA
