# Alterações no Sistema de Planos - 2025-11-07

## Resumo
Removido o plano "ilimitado" e ajustado o sistema para ter apenas 3 planos:
- **Gratuito** (2 instâncias, sem expiração)
- **Pro** (8 instâncias, R$47/mês)
- **Analista** (20 instâncias, R$97/mês)

## Alterações Realizadas

### 1. migrations.go
- ✅ Ajustado plano Gratuito de 999999 para 2 instâncias
- ✅ Ajustado plano Pro de 5 para 8 instâncias e preço de R$29 para R$47
- ✅ Ajustado plano Analista de 12 para 20 instâncias e preço de R$97
- ✅ Removido trial_days do plano Gratuito (de 5 para 0)
- ✅ Alterações aplicadas tanto para SQLite quanto PostgreSQL

### 2. subscriptions.go
- ✅ Removida lógica de trial com data de expiração
- ✅ Plano gratuito agora não tem expires_at (NULL)
- ✅ Função CreateDefaultSubscription cria subscription permanente
- ✅ Atualizado comentário de "free trial" para "free"

### 3. auth.go
- ✅ Atualizado comentário ao criar subscription padrão
- ✅ Mantida a criação automática de subscription no registro

### 4. static/dashboard/subscription.html
- ✅ Removida referência a "Usuários ilimitados"
- ✅ Alterado para "Múltiplos usuários no sistema"

### 5. static/dashboard/js/dashboard-v4.js
- ✅ Removida lógica de plano Trial com contagem de dias
- ✅ Removida lógica de plano ilimitado (> 1000 instâncias)
- ✅ Simplificada barra de progresso para mostrar apenas WhatsApp restantes
- ✅ Mantida verificação de expiração para planos pagos

### 6. fix_plans.go (NOVO)
- ✅ Função automática para atualizar banco de dados existente
- ✅ Executa ao iniciar o servidor
- ✅ Atualiza limites e preços dos planos
- ✅ Remove expiração de subscriptions no plano gratuito

### 7. main.go
- ✅ Adicionada chamada para FixPlansInDatabase() na inicialização
- ✅ Atualiza banco automaticamente ao reiniciar o servidor

## Comportamento Após as Alterações

### Registro de Novo Usuário
1. Usuário se registra no sistema
2. Automaticamente recebe o plano Gratuito (ID 1)
3. Plano não tem data de expiração (expires_at = NULL)
4. Pode criar até 2 instâncias do WhatsApp
5. Instância padrão é criada automaticamente
6. **Barra de progresso mostra corretamente: "1 WhatsApp restante para conectar"**

### Plano Gratuito
- Limite: 2 instâncias
- Sem data de expiração
- Sem período de trial
- Permanente enquanto o usuário não fazer upgrade

### Planos Pagos
- Pro: 8 instâncias por R$47/mês
- Analista: 20 instâncias por R$97/mês
- Podem ter data de expiração (expires_at)
- Verificação de expiração ativa

## Atualização Automática do Banco de Dados

**IMPORTANTE:** Ao reiniciar o servidor, o sistema automaticamente:
1. Detecta se há planos com valores antigos no banco
2. Atualiza os limites e preços automaticamente
3. Remove expiração de subscriptions no plano gratuito
4. Registra as alterações no log

**Não é necessário executar scripts SQL manualmente!**

O arquivo `fix_plans.go` contém a lógica que:
- Atualiza plano Gratuito para 2 instâncias
- Atualiza plano Pro para 8 instâncias e R$47
- Atualiza plano Analista para 20 instâncias
- Remove expires_at de subscriptions no plano gratuito

## Migração Manual (Opcional)

Se preferir atualizar manualmente, use o arquivo `update_plans.sql`:

```bash
# Para SQLite
sqlite3 wuzapi.db < update_plans.sql

# Para PostgreSQL
psql -U usuario -d banco < update_plans.sql
```

## Arquivos Modificados
1. `/home/allansevero/wuzapi/migrations.go`
2. `/home/allansevero/wuzapi/subscriptions.go`
3. `/home/allansevero/wuzapi/auth.go`
4. `/home/allansevero/wuzapi/static/dashboard/subscription.html`
5. `/home/allansevero/wuzapi/static/dashboard/js/dashboard-v4.js`
6. `/home/allansevero/wuzapi/fix_plans.go` (NOVO)
7. `/home/allansevero/wuzapi/main.go`
8. `/home/allansevero/wuzapi/update_plans.sql` (NOVO)

## Como Aplicar as Alterações

1. **Pare o servidor:**
   ```bash
   # Se estiver rodando como serviço
   sudo systemctl stop wuzapi
   
   # Ou se estiver rodando manualmente, pressione Ctrl+C
   ```

2. **Reinicie o servidor:**
   ```bash
   # Se estiver rodando como serviço
   sudo systemctl start wuzapi
   
   # Ou rode manualmente
   ./wuzapi
   ```

3. **Verifique os logs:**
   ```bash
   # O servidor deve mostrar mensagens como:
   # "✓ Plano Gratuito atualizado: 2 instâncias, R$0,00"
   # "✓ Plano Pro atualizado: 8 instâncias, R$47,00"
   # "✓ Plano Analista atualizado: 20 instâncias, R$97,00"
   ```

4. **Teste:**
   - Cadastre um novo usuário
   - Verifique que a barra de progresso mostra "1 WhatsApp restante para conectar"
   - Tente criar uma segunda instância
   - Verifique que ao criar a 3ª instância é bloqueado

## Testes Recomendados
1. ✅ Verificar que novo usuário recebe plano Gratuito
2. ✅ Verificar que plano Gratuito não expira
3. ✅ Verificar limite de 2 instâncias no plano Gratuito
4. ✅ Verificar que dashboard não mostra "ilimitado"
5. ✅ Verificar que barra de progresso mostra "1 WhatsApp restante" após criar conta
6. ✅ Verificar que ao criar 2ª instância mostra "0 WhatsApp restantes"
7. ✅ Verificar que ao tentar criar 3ª instância é bloqueado

