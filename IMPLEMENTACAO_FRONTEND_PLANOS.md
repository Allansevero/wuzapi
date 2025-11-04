# Implementação do Frontend com Sistema de Planos

## Status Atual

### Backend - ✅ COMPLETO
O backend já está 100% implementado e funcional:

1. **Migração de Banco de Dados** (`migrations.go` - Migration 13)
   - Tabela `plans` criada com 3 planos
   - Tabela `user_subscriptions` para assinaturas ativas
   - Tabela `subscription_history` para histórico
   - Índices criados para otimização

2. **Lógica de Negócios** (`subscriptions.go`)
   - `CreateDefaultSubscription()` - Cria trial gratuito de 5 dias
   - `GetActiveSubscription()` - Retorna assinatura ativa
   - `UpdateSubscription()` - Atualiza plano do usuário
   - `CheckSubscriptionExpired()` - Verifica expiração
   - `GetUserInstanceCount()` - Conta instâncias do usuário
   - `CanCreateInstance()` - Valida se pode criar mais instâncias
   - `GetAllPlans()` - Lista todos os planos

3. **Integração com Registro** (`auth.go`)
   - Linha 219: Cria automaticamente subscription trial ao registrar
   - Linha 230: Cria instância padrão automaticamente

4. **APIs REST** (`routes.go`)
   - `GET /user/subscription` - Detalhes da assinatura atual
   - `PUT /user/subscription` - Atualizar plano
   - `GET /user/plans` - Listar planos disponíveis

5. **Planos Configurados no Banco**
   - **Gratuito**: R$ 0,00 - Trial 5 dias - Ilimitado
   - **Pro**: R$ 29,00 - 5 instâncias
   - **Analista**: R$ 97,00 - 12 instâncias

## Frontend - 🔨 PRECISA SER IMPLEMENTADO

### Arquitetura do Frontend
- Stack: **HTML + Tailwind CSS + Vanilla JavaScript**
- Sem frameworks React/Vue/Angular
- Design System: Tailwind CSS via CDN
- Fonte: Inter do Google Fonts

### Estrutura de Arquivos Necessários

```
static/
├── dashboard/
│   ├── user-dashboard-v3.html    (NOVO - Dashboard moderno)
│   ├── js/
│   │   ├── user-dashboard-v3.js  (NOVO - Lógica do dashboard)
│   │   ├── plans.js              (NOVO - Gerenciamento de planos)
│   │   └── api-client.js         (NOVO - Cliente API centralizado)
│   └── css/
│       └── custom.css            (NOVO - Estilos personalizados)
```

### Componentes do Frontend

#### 1. Sidebar (Barra Lateral)
```html
<aside class="w-64 bg-white border-r">
  <!-- Logo -->
  <div class="h-20 px-6">
    <h1 class="text-3xl font-bold">metrizap</h1>
  </div>

  <!-- Navegação -->
  <nav class="px-4 py-4">
    <!-- PRINCIPAL -->
    <a href="#contas" class="sidebar-link">
      <svg>...</svg>
      <span>Contas conectadas</span>
    </a>

    <!-- PERFIL -->
    <a href="#dados" class="sidebar-link">
      <svg>...</svg>
      <span>Seus dados</span>
    </a>
  </nav>

  <!-- Rodapé com indicador de plano -->
  <div class="p-6 border-t">
    <p id="instancesRemaining">4 contas conectadas restantes</p>
    <div class="w-full bg-gray-200 rounded-full h-1.5">
      <div id="instancesProgress" class="bg-mz-green h-1.5 rounded-full" style="width: 60%"></div>
    </div>
  </div>
</aside>
```

**JavaScript para atualizar indicador:**
```javascript
async function updateInstancesIndicator() {
  const subscription = await getActiveSubscription();
  const instances = await getUserInstances();
  
  const used = instances.length;
  const total = subscription.plan.max_instances;
  const remaining = total - used;
  const percentage = (used / total) * 100;
  
  document.getElementById('instancesRemaining').textContent = 
    `${remaining} contas conectadas restantes`;
  document.getElementById('instancesProgress').style.width = 
    `${percentage}%`;
}
```

#### 2. Header (Cabeçalho)
```html
<header class="flex justify-between items-center mb-8">
  <h1 id="welcomeMessage" class="text-4xl font-semibold">Olá, Allan 👋</h1>
  <img id="userAvatar" class="w-12 h-12 rounded-full" 
       src="https://placehold.co/48x48/E2E8F0/4A5568?text=A" 
       alt="Avatar">
</header>
```

**JavaScript para personalizar:**
```javascript
async function updateHeader() {
  const user = await getCurrentUser();
  const firstName = user.email.split('@')[0];
  document.getElementById('welcomeMessage').textContent = 
    `Olá, ${firstName} 👋`;
}
```

#### 3. Cards de Instância
```html
<!-- Card Conectado -->
<div class="bg-white p-5 rounded-xl shadow-sm border border-mz-green">
  <div class="flex justify-between items-center mb-4">
    <h2 class="font-semibold text-lg">Allan</h2>
    <span class="bg-mz-green-light text-mz-green text-xs font-bold px-3 py-1 rounded-full">
      Conectado
    </span>
  </div>
  <div class="flex justify-between items-center mb-6">
    <div>
      <span class="text-sm text-gray-500">Data da criação</span>
      <p class="font-semibold">03/03/2025</p>
    </div>
    <div>
      <span class="text-sm text-gray-500">Análises concluídas</span>
      <p class="font-semibold">110</p>
    </div>
  </div>
  <div class="flex space-x-2">
    <button class="flex-1 bg-gray-800 text-white font-medium py-2 px-4 rounded-lg">
      Desconectar
    </button>
    <button class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg">
      Excluir
    </button>
  </div>
</div>
```

**JavaScript para renderizar cards:**
```javascript
function createInstanceCard(instance) {
  const isConnected = instance.connected;
  const borderClass = isConnected ? 'border-mz-green' : 'border-gray-200';
  const badgeClass = isConnected ? 
    'bg-mz-green-light text-mz-green' : 
    'bg-gray-100 text-gray-600';
  const statusText = isConnected ? 'Conectado' : 'Desconectado';
  
  return `
    <div class="bg-white p-5 rounded-xl shadow-sm border ${borderClass}">
      <div class="flex justify-between items-center mb-4">
        <h2 class="font-semibold text-lg">${instance.name}</h2>
        <span class="${badgeClass} text-xs font-bold px-3 py-1 rounded-full">
          ${statusText}
        </span>
      </div>
      <div class="flex justify-between items-center mb-6">
        <div>
          <span class="text-sm text-gray-500">Data da criação</span>
          <p class="font-semibold">${formatDate(instance.created_at)}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">Análises concluídas</span>
          <p class="font-semibold">0</p>
        </div>
      </div>
      <div class="flex space-x-2">
        ${isConnected ? `
          <button onclick="disconnectInstance('${instance.id}')" 
                  class="flex-1 bg-gray-800 text-white font-medium py-2 px-4 rounded-lg">
            Desconectar
          </button>
        ` : `
          <button onclick="connectInstance('${instance.id}')" 
                  class="flex-1 bg-mz-green text-white font-medium py-2 px-4 rounded-lg">
            Conectar WhatsApp
          </button>
        `}
        <button onclick="deleteInstance('${instance.id}')" 
                class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg">
          Excluir
        </button>
      </div>
    </div>
  `;
}
```

#### 4. Modal QR Code
```html
<div id="qrModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50 hidden">
  <div class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md">
    <button id="closeQRModal" class="absolute top-4 right-4">
      <svg>...</svg> <!-- X icon -->
    </button>
    
    <p class="text-sm text-gray-600 mb-6">
      Abra seu WhatsApp → Dispositivos conectados → Aponte a câmera
    </p>
    
    <div id="qrCodeContainer" class="w-64 h-64 mx-auto my-4">
      <img id="qrCodeImage" src="" alt="QR Code" class="w-full h-full">
    </div>
    
    <div class="bg-mz-green-light border border-mz-green text-mz-green p-3 rounded-lg">
      <svg>...</svg>
      <p>Diariamente você receberá análises desse número até desconecta-lo.</p>
    </div>
  </div>
</div>
```

**JavaScript para QR Code:**
```javascript
async function showQRCode(instanceId, token) {
  const modal = document.getElementById('qrModal');
  const qrImage = document.getElementById('qrCodeImage');
  
  modal.classList.remove('hidden');
  
  // Inicia polling para QR code
  const pollInterval = setInterval(async () => {
    try {
      const response = await fetch(`/session/qr?token=${token}`);
      const data = await response.json();
      
      if (data.qrcode) {
        qrImage.src = data.qrcode;
      }
      
      if (data.connected) {
        clearInterval(pollInterval);
        modal.classList.add('hidden');
        await refreshInstances();
      }
    } catch (error) {
      console.error('Error polling QR code:', error);
    }
  }, 2000);
}
```

#### 5. Página de Planos (Seus Dados)
```html
<div id="paginaDados" class="hidden">
  <h1 class="text-4xl font-semibold mb-8">Allan, seus dados 👋</h1>

  <!-- Formulário de Dados -->
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6 mb-10">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
      <input id="userName" type="text" disabled 
             class="w-full p-3 border rounded-lg bg-gray-100">
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">E-mail</label>
      <input id="userEmail" type="email" disabled 
             class="w-full p-3 border rounded-lg bg-gray-100">
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Quero receber análises no:</label>
      <input id="destinationNumber" type="text" 
             class="w-full p-3 border rounded-lg">
    </div>
  </div>

  <!-- Plano Atual -->
  <div>
    <h2 class="text-2xl font-semibold mb-6">Plano atual</h2>
    <div id="plansContainer" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Plans will be rendered here -->
    </div>
  </div>
</div>
```

**JavaScript para renderizar planos:**
```javascript
async function renderPlans() {
  const [plans, subscription] = await Promise.all([
    getAllPlans(),
    getActiveSubscription()
  ]);
  
  const container = document.getElementById('plansContainer');
  container.innerHTML = plans.map(plan => {
    const isActive = subscription.plan_id === plan.id;
    const borderClass = isActive ? 'border-2 border-mz-green' : 'border border-gray-200';
    
    return `
      <div class="bg-white p-6 rounded-xl shadow-sm ${borderClass}">
        <h3 class="text-xl font-semibold">${plan.name}</h3>
        <p class="text-3xl font-bold my-4">R$${plan.price.toFixed(2)}</p>
        <ul class="space-y-2 text-gray-600 mb-4">
          <li class="flex items-center space-x-2">
            <svg class="w-5 h-5 text-mz-green">...</svg>
            <span>Análise diária</span>
          </li>
          <li class="flex items-center space-x-2">
            <svg class="w-5 h-5 text-mz-green">...</svg>
            <span>${plan.max_instances} contas conectadas</span>
          </li>
        </ul>
        ${!isActive ? `
          <button onclick="upgradePlan(${plan.id})" 
                  class="w-full bg-mz-green text-white font-semibold py-3 px-5 rounded-lg">
            Fazer upgrade
          </button>
        ` : `
          <div class="w-full bg-mz-green-light text-mz-green font-semibold py-3 px-5 rounded-lg text-center">
            Plano Atual
          </div>
        `}
      </div>
    `;
  }).join('');
}
```

#### 6. Modal de Adicionar Instância com Validação de Plano
```javascript
async function addInstanceClick() {
  const canCreate = await checkCanCreateInstance();
  
  if (!canCreate.allowed) {
    showUpgradeModal(canCreate.reason);
    return;
  }
  
  showAddInstanceModal();
}

function showUpgradeModal(reason) {
  // Show modal explaining limit reached and offering upgrade
  const modal = `
    <div class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center z-50">
      <div class="bg-white p-8 rounded-2xl shadow-xl max-w-md">
        <h3 class="text-xl font-semibold mb-4">Limite de instâncias atingido</h3>
        <p class="text-gray-600 mb-6">${reason}</p>
        <button onclick="navigateToPlans()" 
                class="w-full bg-mz-green text-white font-semibold py-3 rounded-lg">
          Ver planos disponíveis
        </button>
      </div>
    </div>
  `;
  document.body.insertAdjacentHTML('beforeend', modal);
}
```

### APIs JavaScript

#### api-client.js - Cliente Centralizado
```javascript
const API_BASE = window.location.origin;

// Auth
function getAuthToken() {
  return localStorage.getItem('auth_token');
}

// Generic API call
async function apiCall(endpoint, options = {}) {
  const token = getAuthToken();
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  };
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers
  });
  
  if (!response.ok) {
    throw new Error(`API Error: ${response.statusText}`);
  }
  
  return response.json();
}

// Subscription APIs
async function getActiveSubscription() {
  return apiCall('/user/subscription');
}

async function getAllPlans() {
  return apiCall('/user/plans');
}

async function upgradePlan(planId) {
  return apiCall('/user/subscription', {
    method: 'PUT',
    body: JSON.stringify({ plan_id: planId })
  });
}

// Instance APIs
async function getUserInstances() {
  return apiCall('/user/instances');
}

async function createInstance(name) {
  return apiCall('/user/instances', {
    method: 'POST',
    body: JSON.stringify({ name })
  });
}

async function deleteInstance(instanceId) {
  return apiCall(`/user/instances/${instanceId}`, {
    method: 'DELETE'
  });
}

async function connectInstance(instanceId) {
  return apiCall(`/session/connect`, {
    method: 'POST',
    body: JSON.stringify({ instance_id: instanceId })
  });
}

async function getInstanceQR(token) {
  return apiCall(`/session/qr?token=${token}`);
}

// Check if user can create more instances
async function checkCanCreateInstance() {
  const subscription = await getActiveSubscription();
  const instances = await getUserInstances();
  
  if (!subscription.is_active) {
    return {
      allowed: false,
      reason: 'Sua assinatura expirou. Renove seu plano para continuar.'
    };
  }
  
  if (subscription.expires_at && new Date(subscription.expires_at) < new Date()) {
    return {
      allowed: false,
      reason: 'Seu período de trial expirou. Escolha um plano para continuar.'
    };
  }
  
  if (instances.length >= subscription.plan.max_instances) {
    return {
      allowed: false,
      reason: `Você atingiu o limite de ${subscription.plan.max_instances} instâncias do seu plano ${subscription.plan.name}. Faça upgrade para conectar mais números.`
    };
  }
  
  return { allowed: true };
}
```

### Cores Personalizadas Tailwind
```javascript
tailwind.config = {
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Inter', 'sans-serif'],
      },
      colors: {
        'mz-green': '#28a745',
        'mz-green-light': '#e9f7ec',
        'mz-red': '#dc3545',
        'mz-orange-light': '#fff3e0',
        'mz-orange-dark': '#fd7e14',
      }
    }
  }
}
```

## Fluxo de Integração

### 1. Registro de Novo Usuário
```
1. Usuário preenche email e senha
2. Backend cria system_user
3. Backend cria subscription com plano Gratuito (trial 5 dias)
4. Backend cria instância padrão automaticamente
5. Frontend recebe JWT token
6. Redireciona para dashboard
7. Dashboard carrega:
   - Subscription details (mostra trial 5 dias)
   - Instâncias do usuário (mostra 1 instância padrão)
   - Indicador de limite (mostra: ilimitado durante 5 dias)
```

### 2. Adicionar Nova Instância
```
1. Usuário clica "Adicionar WhatsApp"
2. JavaScript chama checkCanCreateInstance()
3. Se permitido:
   - Mostra modal para nome da instância
   - Cria instância no backend
   - Atualiza lista de instâncias
   - Atualiza indicador de progresso
4. Se não permitido:
   - Mostra modal de upgrade
   - Oferece planos disponíveis
```

### 3. Conectar WhatsApp
```
1. Usuário clica "Conectar WhatsApp"
2. Modal QR Code abre
3. JavaScript faz polling a cada 2s:
   GET /session/qr?token=INSTANCE_TOKEN
4. Mostra QR Code quando disponível
5. Quando conectado:
   - Para polling
   - Fecha modal
   - Atualiza card para status "Conectado"
   - Troca botões (Desconectar/Excluir)
```

### 4. Upgrade de Plano
```
1. Usuário vai em "Seus dados"
2. Vê planos disponíveis
3. Clica "Fazer upgrade" no plano desejado
4. JavaScript chama: PUT /user/subscription { plan_id: X }
5. Backend:
   - Desativa subscription antiga
   - Cria nova subscription
   - Define expires_at = null (mensal recorrente)
6. Frontend atualiza:
   - Badge do plano atual
   - Indicador de limites
   - Lista de instâncias (se exceder limite, avisa)
```

### 5. Expiração de Trial
```
1. Cron job diário verifica subscriptions expiradas
2. Desativa subscriptions onde expires_at < NOW()
3. Próximo login do usuário:
   - checkCanCreateInstance() retorna false
   - Dashboard mostra aviso de expiração
   - Bloqueia criação de novas instâncias
   - Mantém instâncias existentes conectadas (read-only)
4. Usuário deve escolher plano pago para continuar
```

## Checklist de Implementação

### Backend - ✅ COMPLETO
- [x] Migrations criadas
- [x] Models criados (Plan, UserSubscription, etc)
- [x] Business logic implementada
- [x] APIs REST criadas
- [x] Integração com registro
- [x] Validações de limites
- [x] Planos inseridos no banco

### Frontend - 🔨 IMPLEMENTAR

#### Estrutura HTML
- [ ] Criar user-dashboard-v3.html
- [ ] Implementar sidebar com navegação
- [ ] Implementar header personalizado
- [ ] Criar grid de instâncias (3 colunas)
- [ ] Criar página "Seus dados"
- [ ] Adicionar modals (QR, Delete, Add Instance, Upgrade)

#### JavaScript
- [ ] Criar api-client.js
- [ ] Implementar funções de instâncias
- [ ] Implementar funções de subscription
- [ ] Implementar validação de limites
- [ ] Implementar polling de QR code
- [ ] Implementar atualização de UI em tempo real
- [ ] Implementar navegação entre páginas

#### Integrações
- [ ] Conectar com APIs existentes
- [ ] Testar fluxo de registro → trial → dashboard
- [ ] Testar criação de instâncias com limites
- [ ] Testar upgrade de planos
- [ ] Testar expiração de trial

#### UX/UI
- [ ] Aplicar design Tailwind conforme HTML_FRONTEND_REPLIQUE.md
- [ ] Adicionar loading states
- [ ] Adicionar error handling
- [ ] Adicionar confirmações de ações
- [ ] Adicionar toasts/notifications
- [ ] Testar responsividade mobile

## Próximos Passos

1. **Criar api-client.js** com todas as funções de API
2. **Criar user-dashboard-v3.html** baseado no design fornecido
3. **Criar user-dashboard-v3.js** com toda lógica de interação
4. **Testar fluxo completo** de registro até uso
5. **Ajustar backend** se necessário baseado nos testes
6. **Documentar APIs** para referência futura

## Observações Importantes

1. **Webhook Centralizado**: Já está implementado, todas as instâncias usam o webhook padrão do sistema
2. **Envio Diário 18h**: Sistema `daily_sender.go` já implementa isso
3. **Número de Destino**: Campo `destination_number` já existe na tabela users
4. **Trial Automático**: Já criado automaticamente no registro
5. **Limitações**: Backend já valida limites, frontend só precisa chamar as APIs

## Exemplo de Teste Manual

```bash
# 1. Registrar novo usuário
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"12345678"}'

# 2. Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"12345678"}'

# 3. Ver subscription (usar token do login)
curl -X GET http://localhost:8080/user/subscription \
  -H "Authorization: Bearer TOKEN_HERE"

# 4. Listar planos
curl -X GET http://localhost:8080/user/plans \
  -H "Authorization: Bearer TOKEN_HERE"

# 5. Fazer upgrade (ex: Pro - ID 2)
curl -X PUT http://localhost:8080/user/subscription \
  -H "Authorization: Bearer TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"plan_id":2}'
```

## Conclusão

O backend está 100% pronto e funcional. O próximo passo é implementar o frontend modernizado seguindo o design do HTML_FRONTEND_REPLIQUE.md, conectando-o com as APIs já existentes.

Todos os recursos necessários estão disponíveis:
- Sistema de planos e subscriptions funcional
- APIs REST documentadas e testáveis
- Validações de limites implementadas
- Integration com autenticação via JWT
- Design system completo com Tailwind CSS
