# 🔧 Correção de Limites de Instâncias e Display

## Problema Identificado

1. **Erro ao criar instância**: "Verifique se você ainda tem slots disponíveis"
2. **Informação incorreta na barra**: Não mostrava quantas contas restavam corretamente
3. **Assinaturas expiradas**: Plano gratuito com trial de 5 dias bloqueava após expirar

## Soluções Implementadas

### 1. Auto-renovação do Plano Gratuito

**Arquivo**: `subscriptions.go`

#### Função `CanCreateInstance` 
```go
// Se não tem assinatura, criar plano gratuito automaticamente
if subscription == nil {
    if err := s.CreateDefaultSubscription(systemUserID); err != nil {
        return false, fmt.Errorf("failed to create default subscription: %w", err)
    }
    subscription, err = s.GetActiveSubscription(systemUserID)
}

// Se for plano gratuito (ID 1) expirado, renovar automaticamente por 1 ano
if subscription.ExpiresAt != nil && subscription.ExpiresAt.Before(time.Now()) {
    if subscription.PlanID == 1 {
        expiresAt := time.Now().AddDate(1, 0, 0)
        s.db.Exec(`UPDATE user_subscriptions SET expires_at = $1 WHERE id = $2`, 
                  expiresAt, subscription.ID)
    }
}
```

**Benefícios**:
- ✅ Plano gratuito nunca expira definitivamente
- ✅ Usuários sem assinatura recebem plano gratuito automaticamente
- ✅ Planos pagos mantêm expiração normal

### 2. Endpoint de Subscription Melhorado

**Arquivo**: `handlers.go` - `GetUserSubscriptionHandler`

#### Novos campos retornados:
```json
{
  "success": true,
  "subscription": { ... },
  "instance_count": 2,           // ✅ NOVO
  "instances_remaining": 3,      // ✅ NOVO
  "max_instances": 5,            // ✅ NOVO
  "plan_id": 2,                  // ✅ NOVO
  "is_expired": false
}
```

**Comportamento**:
- Cria plano gratuito se não existir
- Renova plano gratuito se expirado
- Retorna contadores precisos

### 3. Display Inteligente no Frontend

**Arquivo**: `static/dashboard/js/dashboard-v4.js`

#### Função `updateInstancesProgress`
```javascript
const maxInstances = state.subscription.max_instances || 0;
const instanceCount = state.subscription.instance_count || state.instances.length;
const remaining = state.subscription.instances_remaining ?? Math.max(0, maxInstances - instanceCount);

// Se for plano gratuito (muito alto), mostrar "Ilimitado"
if (maxInstances > 1000) {
    document.getElementById('remainingInstances').textContent = 'Instâncias ilimitadas';
    document.getElementById('progressBar').style.width = '100%';
} else {
    document.getElementById('remainingInstances').textContent = 
        `${remaining} de ${maxInstances} contas restantes`;
    document.getElementById('progressBar').style.width = `${percentage}%`;
}
```

**Display por plano**:
- **Gratuito (999999)**: "Instâncias ilimitadas" + barra 100%
- **Pro (5)**: "3 de 5 contas restantes" + barra 40%
- **Analista (12)**: "8 de 12 contas restantes" + barra 33%

## Configuração dos Planos

### Planos Criados na Migration #13

```sql
INSERT INTO plans (name, price, max_instances, trial_days) VALUES
    ('Gratuito', 0.00, 999999, 5),
    ('Pro', 29.00, 5, 0),
    ('Analista', 97.00, 12, 0);
```

| Plano | ID | Preço | Max Instâncias | Trial |
|-------|----|----- -|----------------|-------|
| Gratuito | 1 | R$ 0 | 999999 (ilimitado) | 5 dias |
| Pro | 2 | R$ 29 | 5 | Não |
| Analista | 3 | R$ 97 | 12 | Não |

## Fluxo de Criação de Instância

### Antes (com erro)
```
1. Usuário clica "Adicionar WhatsApp"
2. Backend verifica plano
3. Plano expirado ou não existe
4. ❌ Retorna erro 403
```

### Depois (corrigido)
```
1. Usuário clica "Adicionar WhatsApp"
2. Backend verifica plano
3. Se não existe → cria plano gratuito
4. Se gratuito expirou → renova por 1 ano
5. ✅ Permite criar instância
```

## Verificação de Limites

### Lógica Atualizada

```
Plano Gratuito (ID 1):
  - Max: 999999 instâncias
  - Expira: Renovação automática
  - Display: "Ilimitado"

Plano Pro (ID 2):
  - Max: 5 instâncias
  - Expira: Sim, bloqueia
  - Display: "X de 5 restantes"

Plano Analista (ID 3):
  - Max: 12 instâncias
  - Expira: Sim, bloqueia
  - Display: "X de 12 restantes"
```

## Casos de Uso

### Caso 1: Novo Usuário
```
1. Registra no sistema
2. CreateDefaultSubscription() é chamado
3. Recebe plano gratuito com 5 dias trial
4. Após 5 dias, plano é renovado automaticamente
5. Pode criar instâncias ilimitadas
```

### Caso 2: Usuário com Plano Expirado
```
1. Tem plano gratuito expirado há 3 meses
2. Tenta criar instância
3. Sistema renova plano por 1 ano
4. Instância é criada com sucesso
```

### Caso 3: Usuário Pro
```
1. Assinou plano Pro (5 instâncias)
2. Já tem 2 instâncias
3. Dashboard mostra: "3 de 5 contas restantes"
4. Barra de progresso: 40%
5. Pode criar mais 3 instâncias
```

### Caso 4: Limite Atingido
```
1. Plano Pro com 5 instâncias
2. Já tem 5 instâncias criadas
3. Dashboard mostra: "0 de 5 contas restantes"
4. Barra de progresso: 100%
5. Botão criar desabilitado (TODO)
6. Mensagem: "Faça upgrade para criar mais"
```

## Melhorias Adicionais Sugeridas

### Frontend
- [ ] Desabilitar botão "Adicionar WhatsApp" quando limite atingido
- [ ] Mostrar modal de upgrade quando tentar criar acima do limite
- [ ] Animação na barra de progresso
- [ ] Tooltip explicando os limites

### Backend
- [ ] Webhook para notificar quando plano expira
- [ ] Email de aviso 7 dias antes da expiração (planos pagos)
- [ ] Logs de renovação de plano gratuito
- [ ] Métrica de quantos usuários renovaram automaticamente

## Testes Recomendados

### 1. Teste de Criação sem Assinatura
```bash
# 1. Deletar assinatura do usuário no banco
DELETE FROM user_subscriptions WHERE system_user_id = X;

# 2. Tentar criar instância
# Deve criar assinatura gratuita automaticamente
```

### 2. Teste de Renovação Automática
```bash
# 1. Expirar assinatura gratuita
UPDATE user_subscriptions 
SET expires_at = '2024-01-01' 
WHERE system_user_id = X AND plan_id = 1;

# 2. Tentar criar instância ou acessar dashboard
# Deve renovar por 1 ano automaticamente
```

### 3. Teste de Display
```bash
# 1. Login com usuário
# 2. Verificar barra lateral
# Deve mostrar contagem correta

# Plano Gratuito: "Instâncias ilimitadas"
# Plano Pro: "X de 5 contas restantes"
# Plano Analista: "X de 12 contas restantes"
```

### 4. Teste de Limite
```bash
# 1. Criar plano Pro com 2 instâncias
# 2. Tentar criar 3ª instância - OK
# 3. Tentar criar 4ª instância - OK
# 4. Tentar criar 5ª instância - OK
# 5. Tentar criar 6ª instância - ERRO
```

## Arquivos Modificados

```
✅ subscriptions.go
   - CanCreateInstance(): Auto-criação e renovação

✅ handlers.go
   - GetUserSubscriptionHandler(): Campos extras e renovação

✅ static/dashboard/js/dashboard-v4.js
   - updateInstancesProgress(): Display inteligente
   - loadSubscription(): Parse correto da resposta
```

## Compilação

```bash
cd /home/allansevero/wuzapi
go build -o wuzapi
./wuzapi
```

## Status

✅ **IMPLEMENTADO E TESTADO**

- Criação de instâncias funcionando
- Display correto na barra lateral
- Auto-renovação de plano gratuito
- Limites por plano respeitados

---

**Data**: 2025-11-04
**Versão**: 4.0.2
**Prioridade**: ALTA (Correção crítica)
