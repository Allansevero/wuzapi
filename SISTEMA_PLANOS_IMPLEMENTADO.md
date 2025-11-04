# Sistema de Planos e Limitações - Implementado

## ✅ Implementações Concluídas

### 1. Estrutura de Banco de Dados

#### Tabela `plans`
- `id`: ID do plano
- `name`: Nome do plano (Gratuito, Pro, Analista)
- `price`: Preço do plano
- `max_instances`: Número máximo de instâncias permitidas
- `trial_days`: Dias de trial (5 para gratuito)
- `is_active`: Se o plano está ativo
- `created_at`: Data de criação

**Planos Padrão:**
- Gratuito: R$ 0,00, 999999 instâncias, 5 dias
- Pro: R$ 29,00, 5 instâncias, sem trial
- Analista: R$ 97,00, 12 instâncias, sem trial

#### Tabela `user_subscriptions`
- `id`: ID da subscrição
- `system_user_id`: ID do usuário (FK para system_users)
- `plan_id`: ID do plano (FK para plans)
- `started_at`: Data de início
- `expires_at`: Data de expiração (NULL para planos perpétuos)
- `is_active`: Se a subscrição está ativa
- `created_at`: Data de criação
- `updated_at`: Data de atualização

#### Tabela `subscription_history`
- `id`: ID do histórico
- `system_user_id`: ID do usuário
- `plan_id`: ID do plano
- `started_at`: Data de início
- `ended_at`: Data de término
- `created_at`: Data de criação

### 2. Lógica de Negócio (subscriptions.go)

#### Funções Principais:

**CreateDefaultSubscription(systemUserID int)**
- Cria subscrição gratuita de 5 dias para novos usuários
- Chamada automaticamente no registro

**GetActiveSubscription(systemUserID int)**
- Retorna a subscrição ativa do usuário com detalhes do plano
- Inclui verificação de expiração

**UpdateSubscription(systemUserID, planID int)**
- Atualiza/muda o plano do usuário
- Desativa plano anterior e cria novo
- Registra no histórico

**CheckSubscriptionExpired(systemUserID int)**
- Verifica e desativa subscrições expiradas
- Chamado automaticamente ao verificar limites

**CanCreateInstance(systemUserID int)**
- Verifica se usuário pode criar mais instâncias
- Checa expiração e limite de instâncias do plano
- Retorna true/false

**GetUserInstanceCount(systemUserID int)**
- Conta quantas instâncias o usuário possui

**GetAllPlans()**
- Retorna todos os planos ativos disponíveis

### 3. API Endpoints

#### GET `/my/plans`
Retorna todos os planos disponíveis

**Resposta:**
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
      "is_active": true,
      "created_at": "2025-11-04T..."
    },
    ...
  ]
}
```

#### GET `/my/subscription`
Retorna a subscrição atual do usuário autenticado

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Resposta:**
```json
{
  "success": true,
  "subscription": {
    "id": 1,
    "system_user_id": 1,
    "plan_id": 1,
    "started_at": "2025-11-04T...",
    "expires_at": "2025-11-09T...",
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

#### PUT `/my/subscription`
Atualiza o plano do usuário

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Body:**
```json
{
  "plan_id": 2
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Subscription updated successfully"
}
```

### 4. Validações Implementadas

#### No Registro (auth.go)
- Cria automaticamente subscrição gratuita de 5 dias
- Cria instância padrão para o novo usuário

#### Na Criação de Instância (user_instances.go)
- Verifica se usuário tem subscrição ativa
- Valida se subscrição não expirou
- Checa se não atingiu limite de instâncias do plano
- Retorna mensagem específica se:
  - Subscrição expirada
  - Limite de instâncias atingido
  - Sem subscrição ativa

**Mensagens de Erro:**
```json
{
  "code": 403,
  "error": "Your subscription has expired. Please renew to create more instances.",
  "success": false
}
```

```json
{
  "code": 403,
  "error": "You have reached the maximum number of instances for your plan. Please upgrade to create more.",
  "success": false
}
```

### 5. Migrations

**Migration 13: add_subscription_plans**
- Cria tabela `plans` com planos padrão
- Cria tabela `user_subscriptions`
- Cria tabela `subscription_history`
- Cria índices para performance
- Suporte completo para PostgreSQL e SQLite

### 6. Fluxo de Uso

1. **Novo Usuário:**
   - Registra → Subscrição gratuita (5 dias) é criada automaticamente
   - Pode criar instâncias ilimitadas durante 5 dias
   
2. **Durante Trial:**
   - Usuário conecta WhatsApps
   - Após 5 dias, subscrição expira
   
3. **Após Expiração:**
   - Não pode criar novas instâncias
   - Instâncias existentes continuam funcionando
   - Precisa fazer upgrade para Pro ou Analista
   
4. **Upgrade de Plano:**
   - PUT `/my/subscription` com novo plan_id
   - Plano anterior é desativado
   - Novo plano ativado imediatamente
   - Pode criar até o limite de instâncias do novo plano

### 7. Próximos Passos Sugeridos

1. **Frontend:**
   - Tela de visualização de planos
   - Indicador de uso (X de Y instâncias)
   - Contador de dias restantes do trial
   - Botão de upgrade de plano
   
2. **Integração com Pagamento:**
   - Gateway de pagamento (Stripe, MercadoPago, etc.)
   - Webhook para confirmação de pagamento
   - Renovação automática
   
3. **Notificações:**
   - Email quando trial está próximo do fim
   - Email quando limite de instâncias está próximo
   - Alerta de subscrição expirada
   
4. **Admin Dashboard:**
   - Visualizar todos os usuários e seus planos
   - Histórico de subscrições
   - Métricas de conversão free → paid
   - Gestão manual de planos

## 📝 Notas Técnicas

- Todas as queries suportam PostgreSQL e SQLite
- Timestamps em UTC
- Validações robustas em múltiplas camadas
- Logs detalhados de todas as operações
- Transações para operações críticas
- Índices otimizados para consultas frequentes

## 🧪 Como Testar

### 1. Registrar novo usuário
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"teste@exemplo.com","password":"senha123"}'
```

### 2. Fazer login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teste@exemplo.com","password":"senha123"}'
```

### 3. Ver subscrição atual
```bash
curl -X GET http://localhost:8080/my/subscription \
  -H "Authorization: Bearer <seu_token>"
```

### 4. Listar planos disponíveis
```bash
curl -X GET http://localhost:8080/my/plans \
  -H "Authorization: Bearer <seu_token>"
```

### 5. Criar instância (testando limite)
```bash
curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer <seu_token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Minha Instância","destination_number":"+5511999999999"}'
```

### 6. Fazer upgrade de plano
```bash
curl -X PUT http://localhost:8080/my/subscription \
  -H "Authorization: Bearer <seu_token>" \
  -H "Content-Type: application/json" \
  -d '{"plan_id":2}'
```

## ✅ Checklist de Funcionalidades

- [x] Tabelas de planos no banco
- [x] Inserção de planos padrão
- [x] Associação usuário → plano
- [x] Validação de limite de instâncias
- [x] Controle de expiração (trial 5 dias)
- [x] Bloqueio de criação ao atingir limite
- [x] API para visualizar plano atual
- [x] API para listar planos disponíveis
- [x] API para upgrade/downgrade
- [x] Subscrição automática no registro
- [x] Histórico de subscrições
- [x] Suporte PostgreSQL e SQLite
- [x] Logs completos
- [x] Mensagens de erro específicas
