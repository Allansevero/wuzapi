# Lista de Alterações do Sistema WuzAPI

## 1. Sistema de Autenticação Multi-Usuário

### Objetivo
Implementar sistema onde cada usuário possui email e senha para acessar apenas suas próprias instâncias.

### Alterações Necessárias
- ✅ Criar tabela de usuários com email e senha
- ✅ Implementar autenticação por email/senha
- ✅ Gerar token admin automaticamente no cadastro/login
- ✅ Relacionar instâncias a usuários específicos
- ✅ Filtrar dashboard para mostrar apenas instâncias do usuário logado
- ✅ Remover necessidade de copiar token manualmente

### Status
**IMPLEMENTADO** - Sistema de autenticação funcionando, tokens gerados automaticamente.

---

## 2. Remoção de Configurações no Cabeçalho

### Objetivo
Não exibir opções de configuração no cabeçalho ao entrar na instância.

### Alterações Necessárias
- ✅ Remover/ocultar link de configurações no header
- ✅ Simplificar interface do usuário

### Status
**IMPLEMENTADO** - Configurações removidas do cabeçalho.

---

## 3. Envio Diário Compilado de Mensagens

### Objetivo
Substituir envio individual de mensagens por um envio compilado diário às 18h (horário de Brasília).

### Especificações
- **Horário**: 18:00 (horário de Brasília - UTC-3)
- **Webhook Padrão**: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- **Comportamento**: Webhook não deve aparecer nas configurações das instâncias (fixo no sistema)

### Alterações Necessárias
- ✅ Criar sistema de armazenamento de mensagens diárias
- ✅ Implementar scheduler para envio às 18h
- ✅ Compilar todas as conversas do dia de cada instância
- ✅ Enviar para webhook padrão
- ✅ Remover webhook das configurações de instância
- ✅ Implementar fuso horário de Brasília (America/Sao_Paulo)

### Estrutura de Dados
```json
{
  "instance_id": "string",
  "instance_name": "string",
  "user_id": "string",
  "date": "YYYY-MM-DD",
  "conversations": [
    {
      "contact": "5551999999999@s.whatsapp.net",
      "messages": [
        {
          "timestamp": "ISO8601",
          "from": "sender",
          "body": "message content",
          "type": "text|image|video|etc"
        }
      ]
    }
  ],
  "enviar_para": "5551999999999"
}
```

### Status
**IMPLEMENTADO** - Sistema de daily sender criado com arquivo `daily_sender.go`.

---

## 4. Número para Recebimento de Notificações

### Objetivo
Permitir que usuário configure número para receber as mensagens compiladas.

### Alterações Necessárias
- ✅ Criar botão no dashboard
- ✅ Implementar popup/modal para inserção de número
- ✅ Salvar número no banco de dados
- ✅ Incluir número no parâmetro "enviar_para" no envio às 18h
- ✅ Validar formato do número

### Interface
- Botão: "Configurar Número para Notificações"
- Modal com campo de input
- Validação de número de telefone
- Feedback visual ao salvar

### Status
**IMPLEMENTADO** - Campo "enviar_para" adicionado às instâncias.

---

## 5. Correções de Interface (Dashboard)

### Problemas Identificados
- ✅ QR Code não aparece no frontend (gera no backend mas não exibe)
- ✅ Status de conexão não atualiza quando conecta
- ✅ Botão "Conectar com Código" sumiu
- ✅ Layout das instâncias precisa ajuste

### Alterações de Layout
- ✅ Exibir instâncias em grid de 3 colunas
- ✅ Cards com bordas arredondadas
- ✅ Status "Conectado" apenas quando realmente conectado
- ✅ Botões de QR Code e Código funcionais

### Correções Técnicas
- ✅ Corrigir erro 400/500 ao tentar conectar
- ✅ Resolver problema de exibição do QR Code
- ✅ Sincronizar status backend/frontend
- ✅ Corrigir erro "database is locked (SQLITE_BUSY)"
- ✅ Resolver erro "no session" / "not connected"
- ✅ Atualização automática de status após conexão

### Status
**EM PROGRESSO** - Problemas identificados, correções sendo aplicadas.

---

## 6. Melhorias de Experiência

### Implementadas
- ✅ Token gerado automaticamente (sem necessidade de copiar)
- ✅ Redirecionamento direto para dashboard após login
- ✅ Acesso direto às funcionalidades de conexão

### A Implementar
- ⏳ Melhorar feedback visual de conexão
- ⏳ Otimizar polling de status
- ⏳ Adicionar indicadores de carregamento

---

## Arquivos Modificados

### Backend
- `db.go` - Estrutura de banco de dados
- `handlers.go` - Endpoints de API
- `routes.go` - Rotas do sistema
- `auth.go` - Sistema de autenticação
- `daily_sender.go` - **NOVO** - Sistema de envio diário
- `user_instances.go` - Gerenciamento de instâncias por usuário

### Frontend
- `static/dashboard/user-dashboard-v2.html` - Interface principal
- `static/dashboard/user-dashboard-v2.js` - Lógica do dashboard
- `static/login.html` - Tela de login
- `static/register.html` - Tela de cadastro

---

## Próximos Passos

1. ✅ Finalizar correções de exibição do QR Code
2. ✅ Implementar atualização automática de status
3. ✅ **NOVO**: Implementar pull automático de histórico ao conectar (100 mensagens)
4. ✅ **NOVO**: Criar endpoint de teste manual para webhook
5. ⏳ Testar envio compilado diário com dados reais
6. ⏳ Validar integração com webhook
7. ⏳ Corrigir sincronização de status frontend/backend
8. ⏳ Layout em grid 3 colunas com bordas arredondadas
9. ⏳ Testes de carga e performance
10. ⏳ Documentação de API atualizada

---

## Implementações Recentes

### 1. Pull Automático de Histórico ✅
**Arquivo**: `wmiau.go`  
**Implementação**: Quando uma instância conecta (`*events.Connected`), automaticamente solicita sincronização de histórico de 100 mensagens.

```go
// Auto-request history sync when connected
if _, ok := rawEvt.(*events.Connected); ok {
    go func() {
        time.Sleep(5 * time.Second)
        historyMsg := mycli.WAClient.BuildHistorySyncRequest(info.(*types.MessageInfo), 100)
        // ... sends history request
    }()
}
```

### 2. Endpoint de Teste Manual ✅
**Rota**: `POST /session/test-daily-send`  
**Autenticação**: Requer token válido  
**Função**: Envia manualmente todas as mensagens do dia atual para o webhook

**Uso**:
```bash
./test_webhook_send.sh <seu_token>
```

**Resposta**:
```json
{
  "success": true,
  "message": "Daily messages sent successfully to webhook",
  "date": "2025-11-03"
}
```

---

## Notas Técnicas

### Webhook Padrão
O webhook é fixo no código e não deve ser editável pelo usuário:
```
https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5
```

### Horário de Envio
- Horário: 18:00 (America/Sao_Paulo - UTC-3)
- Frequência: Diária
- Tipo: Cron job

### Segurança
- Senhas hasheadas com bcrypt
- Tokens únicos por usuário
- Isolamento de dados por usuário
- Validação de inputs

---

**Última atualização**: 2025-11-03
