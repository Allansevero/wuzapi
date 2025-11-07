# Alteração: Contagem Baseada em Instâncias Conectadas

**Data:** 2025-11-07  
**Tipo:** Correção de Lógica de Planos

## 📋 Problema Identificado

O sistema estava contando **todas as instâncias criadas** (total), não apenas as **conectadas**, para verificar os limites dos planos. Isso causava problemas como:

- ❌ Usuário cria 5 instâncias mas só conecta 2
- ❌ Sistema bloqueia criação de novas instâncias mesmo tendo 3 desconectadas
- ❌ Limite do plano não refletia a realidade de uso

## ✅ Solução Implementada

**Agora o sistema conta apenas instâncias CONECTADAS ao WhatsApp.**

### Comportamento Atual

- ✅ Plano Gratuito: 2 instâncias **conectadas** simultaneamente
- ✅ Plano Pro: 8 instâncias **conectadas** simultaneamente  
- ✅ Plano Analista: 20 instâncias **conectadas** simultaneamente

### Exemplo Prático

**Plano Gratuito (2 instâncias conectadas):**
- Usuário pode criar quantas instâncias quiser
- Mas só pode ter **2 conectadas** ao mesmo tempo
- Se desconectar 1, pode conectar outra no lugar
- Barra de progresso: "2 de 2 WhatsApp conectados (0 disponíveis)"

## 🔧 Alterações Realizadas

### 1. subscriptions.go
```go
// Nova função para contar apenas conectadas
GetUserConnectedInstanceCount(systemUserID int) (int, error)

// Atualizada para usar contagem de conectadas
CanCreateInstance(systemUserID int) (bool, error)
  - Antes: Verificava total de instâncias criadas
  - Agora: Verifica instâncias conectadas
```

**Removido:**
- ❌ `CanConnectInstance` (duplicada)
- ❌ Lógica de contagem total de instâncias

### 2. handlers.go

**Connect() - Endpoint /session/connect:**
```go
// Verifica limite antes de conectar
canConnect, err := s.CanCreateInstance(systemUserID)
if !canConnect {
    return "connection limit reached"
}
```

**GetUserSubscriptionHandler():**
```go
// Retorna apenas connected_count (não mais instance_count)
{
    "connected_count": 2,
    "instances_remaining": 0,
    "max_instances": 2
}
```

### 3. dashboard-v4.js

**Barra de Progresso:**
```javascript
// Antes
"2 WhatsApp restantes para conectar"

// Agora  
"2 de 2 WhatsApp conectados (0 disponíveis)"
```

**Estado da Subscription:**
- ❌ Removido: `instance_count`
- ✅ Mantido: `connected_count`
- ✅ Mantido: `instances_remaining`

## 🎯 Resultados Esperados

### Cenário 1: Novo Usuário (Plano Gratuito)
1. Cadastra → recebe plano Gratuito (2 conectados)
2. Sistema cria 1 instância automaticamente
3. Instância não está conectada ainda
4. Barra: **"0 de 2 WhatsApp conectados (2 disponíveis)"**
5. Ao conectar a 1ª: **"1 de 2 WhatsApp conectados (1 disponível)"**
6. Ao conectar a 2ª: **"2 de 2 WhatsApp conectados (0 disponíveis)"**
7. Ao tentar conectar 3ª: **BLOQUEADO** ❌

### Cenário 2: Usuário Desconecta Uma Instância
1. Tem 2 instâncias conectadas (limite atingido)
2. Desconecta 1 instância
3. Barra: **"1 de 2 WhatsApp conectados (1 disponível)"**
4. Pode conectar outra instância agora ✅

### Cenário 3: Usuário com 10 Instâncias Criadas
1. Tem 10 instâncias no banco de dados
2. Mas só 2 estão conectadas
3. Barra: **"2 de 2 WhatsApp conectados (0 disponíveis)"**
4. Pode desconectar qualquer uma das 2
5. E conectar qualquer uma das 8 desconectadas

## 📊 Comparação: Antes vs Agora

| Aspecto | Antes | Agora |
|---------|-------|-------|
| **Contagem** | Total criadas | Apenas conectadas |
| **Limite Gratuito** | 2 criadas | 2 conectadas |
| **Flexibilidade** | Baixa | Alta |
| **Bloqueio** | Ao criar | Ao conectar |
| **Desconexão** | Não libera slot | Libera slot ✅ |

## 🔍 Verificações no Banco

### Query para ver instâncias conectadas:
```sql
-- Contar instâncias conectadas de um usuário
SELECT COUNT(*) 
FROM users 
WHERE system_user_id = 1 AND connected = 1;

-- Ver todas as instâncias e status
SELECT id, name, connected, jid, system_user_id 
FROM users 
WHERE system_user_id = 1;
```

## 🚀 Como Testar

1. **Teste com Plano Gratuito:**
   ```bash
   # Crie 3 instâncias
   curl -X POST /user/instance -H "Authorization: Bearer TOKEN" -d '{"name":"Instance 1"}'
   curl -X POST /user/instance -H "Authorization: Bearer TOKEN" -d '{"name":"Instance 2"}'
   curl -X POST /user/instance -H "Authorization: Bearer TOKEN" -d '{"name":"Instance 3"}'
   
   # Conecte a 1ª (deve funcionar)
   curl -X POST /session/connect -H "X-Instance-Token: TOKEN1"
   
   # Conecte a 2ª (deve funcionar)
   curl -X POST /session/connect -H "X-Instance-Token: TOKEN2"
   
   # Conecte a 3ª (deve ser BLOQUEADO)
   curl -X POST /session/connect -H "X-Instance-Token: TOKEN3"
   # Resposta: "connection limit reached. Please upgrade your plan"
   ```

2. **Teste de Desconexão:**
   ```bash
   # Desconecte a 1ª instância
   curl -X POST /session/disconnect -H "X-Instance-Token: TOKEN1"
   
   # Agora pode conectar a 3ª
   curl -X POST /session/connect -H "X-Instance-Token: TOKEN3"
   # Deve funcionar! ✅
   ```

3. **Verificar Dashboard:**
   - Abra o dashboard
   - Verifique a barra de progresso
   - Deve mostrar: "X de Y WhatsApp conectados (Z disponíveis)"

## 📝 Arquivos Modificados

1. `/home/allansevero/wuzapi/subscriptions.go`
   - Adicionada `GetUserConnectedInstanceCount()`
   - Modificada `CanCreateInstance()` para usar contagem de conectadas
   - Removida `CanConnectInstance()` (duplicada)

2. `/home/allansevero/wuzapi/handlers.go`
   - Atualizado `Connect()` para verificar limite de conectadas
   - Atualizado `GetUserSubscriptionHandler()` para retornar apenas `connected_count`

3. `/home/allansevero/wuzapi/static/dashboard/js/dashboard-v4.js`
   - Atualizado `updateInstancesProgress()` para mostrar conectadas
   - Removida referência a `instance_count`

## ⚠️ Importante

**Nenhuma migração de banco necessária!**

A coluna `connected` já existe na tabela `users`:
- `connected = 1` ou `true` → instância conectada
- `connected = 0` ou `false` → instância desconectada

O sistema agora usa essa coluna corretamente para contar limites.

## 🎯 Conclusão

✅ Sistema agora funciona corretamente  
✅ Limites baseados em instâncias **conectadas**  
✅ Usuário tem flexibilidade para desconectar/reconectar  
✅ Barra de progresso mostra informação precisa  
✅ Bloqueios apenas ao tentar **conectar**, não criar  

**Para aplicar: reinicie o servidor e teste!** 🚀
