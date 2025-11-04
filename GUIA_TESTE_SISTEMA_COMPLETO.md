# Guia de Teste do Sistema Completo - Wuzapi

## Sistema de Planos Implementado

### Planos Disponíveis

1. **Plano Gratuito (Trial)**
   - ID: 1
   - Preço: R$ 0,00
   - Números WhatsApp: Ilimitados
   - Duração: 5 dias
   - Criado automaticamente ao cadastrar novo usuário

2. **Plano Pro**
   - ID: 2
   - Preço: R$ 29,00/mês
   - Números WhatsApp: Até 5
   - Duração: Ilimitada (enquanto ativo)

3. **Plano Analista**
   - ID: 3
   - Preço: R$ 97,00/mês
   - Números WhatsApp: Até 12
   - Duração: Ilimitada (enquanto ativo)

## Funcionalidades Implementadas

### 1. Autenticação e Cadastro
- ✅ Registro de novo usuário com e-mail e senha
- ✅ Login com credenciais
- ✅ Token admin gerado automaticamente
- ✅ Criação automática de plano gratuito (5 dias)
- ✅ Redirecionamento direto para dashboard

### 2. Gerenciamento de Instâncias
- ✅ Listagem de instâncias do usuário logado
- ✅ Criação de nova instância
- ✅ Validação de limite por plano
- ✅ Visualização de status de conexão
- ✅ Exclusão de instância

### 3. Sistema de Assinaturas
- ✅ Consulta de plano atual
- ✅ Listagem de planos disponíveis
- ✅ Atualização de plano
- ✅ Verificação de expiração
- ✅ Contagem de instâncias vs limite do plano
- ✅ Histórico de assinaturas

### 4. Webhook Único e Envio Diário
- ✅ Webhook fixo configurado: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- ✅ Envio automático diário às 18:00 (horário de Brasília)
- ✅ Compilação de todas as conversas do dia
- ✅ Parâmetro "enviar_para" com número de destino
- ✅ Endpoint para teste manual de envio

### 5. Número de Destino
- ✅ Configuração de número para receber mensagens compiladas
- ✅ Salvo no banco de dados por instância
- ✅ Incluído no payload do webhook

## APIs Disponíveis

### Autenticação

#### Registrar Novo Usuário
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}

Resposta:
{
  "success": true,
  "message": "User registered successfully",
  "user_id": 1,
  "token": "abc123..." 
}
```

#### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}

Resposta:
{
  "success": true,
  "token": "abc123...",
  "user_id": 1,
  "email": "usuario@exemplo.com"
}
```

### Instâncias (Requer Token)

#### Listar Minhas Instâncias
```bash
GET /my/instances
Authorization: Bearer {token}

Resposta:
{
  "success": true,
  "instances": [...]
}
```

#### Criar Nova Instância
```bash
POST /my/instances
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Meu WhatsApp"
}

Resposta:
{
  "success": true,
  "instance": {...},
  "token": "instance_token_xyz"
}
```

### Planos e Assinaturas (Requer Token)

#### Listar Planos Disponíveis
```bash
GET /my/plans
Authorization: Bearer {token}

Resposta:
{
  "success": true,
  "plans": [
    {
      "id": 1,
      "name": "Gratuito",
      "price": 0,
      "max_instances": 999999,
      "trial_days": 5
    },
    {
      "id": 2,
      "name": "Pro",
      "price": 29,
      "max_instances": 5,
      "trial_days": 0
    },
    {
      "id": 3,
      "name": "Analista",
      "price": 97,
      "max_instances": 12,
      "trial_days": 0
    }
  ]
}
```

#### Ver Assinatura Atual
```bash
GET /my/subscription
Authorization: Bearer {token}

Resposta:
{
  "success": true,
  "subscription": {
    "id": 1,
    "system_user_id": 1,
    "plan_id": 1,
    "started_at": "2025-11-04T10:00:00Z",
    "expires_at": "2025-11-09T10:00:00Z",
    "is_active": true,
    "plan": {
      "id": 1,
      "name": "Gratuito",
      "price": 0,
      "max_instances": 999999,
      "trial_days": 5
    }
  },
  "instance_count": 2,
  "is_expired": false
}
```

#### Atualizar Assinatura
```bash
PUT /my/subscription
Authorization: Bearer {token}
Content-Type: application/json

{
  "plan_id": 2
}

Resposta:
{
  "success": true,
  "message": "Subscription updated successfully"
}
```

### Configuração de Número de Destino (Requer Token de Instância)

#### Definir Número de Destino
```bash
POST /session/destination-number
Authorization: Bearer {instance_token}
Content-Type: application/json

{
  "destination_number": "5511999999999"
}

Resposta:
{
  "success": true,
  "message": "Destination number updated successfully",
  "destination_number": "5511999999999"
}
```

#### Consultar Número de Destino
```bash
GET /session/destination-number
Authorization: Bearer {instance_token}

Resposta:
{
  "success": true,
  "destination_number": "5511999999999"
}
```

### Teste Manual de Envio Diário

#### Trigger Manual
```bash
POST /session/send-daily-test
Authorization: Bearer {instance_token}

Resposta:
{
  "success": true,
  "message": "Daily messages sent successfully",
  "instance_id": "abc123",
  "date": "2025-11-04"
}
```

## Estrutura do Payload do Webhook

Quando o sistema envia as mensagens diárias às 18h, o webhook recebe:

```json
{
  "instance_id": "abc123def456",
  "date": "2025-11-04",
  "enviar_para": "5511999999999",
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
          "data": {...}
        },
        {
          "sender_jid": "5511999999999@s.whatsapp.net",
          "message_type": "text",
          "text_content": "Oi, tudo bem?",
          "media_link": "",
          "timestamp": "2025-11-04T10:31:00Z",
          "data": {...}
        }
      ]
    },
    {
      "contact": "5511777777777@s.whatsapp.net",
      "messages": [...]
    }
  ]
}
```

## Fluxo de Teste Completo

### 1. Cadastro e Login
```bash
# 1. Registrar novo usuário
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@exemplo.com",
    "password": "senha123"
  }'

# Resultado: Plano gratuito de 5 dias criado automaticamente

# 2. Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@exemplo.com",
    "password": "senha123"
  }'

# Guardar o token retornado
export AUTH_TOKEN="token_retornado"
```

### 2. Verificar Plano Atual
```bash
curl -X GET http://localhost:8080/my/subscription \
  -H "Authorization: Bearer $AUTH_TOKEN"

# Deve mostrar plano Gratuito ativo por 5 dias
```

### 3. Criar Instância WhatsApp
```bash
curl -X POST http://localhost:8080/my/instances \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp Principal"
  }'

# Guardar o token da instância
export INSTANCE_TOKEN="instance_token_retornado"
```

### 4. Conectar WhatsApp
```bash
# Gerar QR Code
curl -X POST http://localhost:8080/session/connect \
  -H "Authorization: Bearer $INSTANCE_TOKEN"

# Obter QR Code
curl -X GET http://localhost:8080/session/qr \
  -H "Authorization: Bearer $INSTANCE_TOKEN"

# Escanear com WhatsApp
```

### 5. Configurar Número de Destino
```bash
curl -X POST http://localhost:8080/session/destination-number \
  -H "Authorization: Bearer $INSTANCE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "destination_number": "5511999999999"
  }'
```

### 6. Testar Envio Manual
```bash
# Enviar mensagens de teste durante o dia

# Trigger envio manual (simula 18h)
curl -X POST http://localhost:8080/session/send-daily-test \
  -H "Authorization: Bearer $INSTANCE_TOKEN"

# Verificar webhook recebido em:
# https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5
```

### 7. Upgrade de Plano
```bash
# Ver planos disponíveis
curl -X GET http://localhost:8080/my/plans \
  -H "Authorization: Bearer $AUTH_TOKEN"

# Fazer upgrade para Plano Pro
curl -X PUT http://localhost:8080/my/subscription \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_id": 2
  }'

# Verificar nova assinatura
curl -X GET http://localhost:8080/my/subscription \
  -H "Authorization: Bearer $AUTH_TOKEN"
```

## Validações Implementadas

1. **Limite de Instâncias por Plano**
   - Gratuito: Ilimitado
   - Pro: Máximo 5
   - Analista: Máximo 12
   - Tentativa de criar além do limite retorna erro

2. **Expiração do Plano Gratuito**
   - Após 5 dias, plano expira
   - Sistema verifica expiração automaticamente
   - Usuário precisa fazer upgrade para continuar

3. **Segurança**
   - Usuário só vê suas próprias instâncias
   - Token de autenticação obrigatório
   - Validação de permissões em todas as rotas

4. **Webhook Fixo**
   - URL hardcoded no sistema
   - Não aparece nas configurações
   - Todos os envios diários vão para o mesmo endpoint

## Cron Job Configurado

- **Horário**: 18:00 (Brasília - America/Sao_Paulo)
- **Frequência**: Diário
- **Ação**: Envia todas as conversas do dia para o webhook
- **Formato**: Agrupado por instância e por conversa

## Logs e Monitoramento

Os logs mostram:
- Inicialização do cron job
- Execuções diárias às 18h
- Envios bem-sucedidos para webhook
- Erros de envio
- Número de conversas enviadas por instância

## Próximos Passos

1. ✅ Sistema de planos - COMPLETO
2. ✅ Envio diário às 18h - COMPLETO
3. ✅ Webhook fixo - COMPLETO
4. ✅ Número de destino - COMPLETO
5. ⏳ Interface administrativa para gerenciar planos
6. ⏳ Sistema de pagamento integrado
7. ⏳ Notificações de expiração de plano
8. ⏳ Dashboard com métricas de uso

## Notas Importantes

- O webhook URL é **fixo** e **não configurável** pelos usuários
- Mensagens são enviadas **uma vez por dia** às 18h (horário de Brasília)
- O parâmetro `enviar_para` contém o número configurado pelo usuário
- Plano gratuito expira automaticamente após 5 dias
- Usuários podem ter múltiplas instâncias dentro do limite do plano
- Histórico de assinaturas é mantido para auditoria
