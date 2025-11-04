# Status de Implementação - Sistema WuzAPI

## ✅ JÁ IMPLEMENTADO (Backend)

### 1. Sistema de Planos e Assinaturas
- ✅ Tabelas criadas no banco de dados (migrations)
- ✅ Três planos configurados:
  - **Gratuito**: R$ 0,00 - Ilimitado números - 5 dias trial
  - **Pro**: R$ 29,00 - 5 números WhatsApp
  - **Analista**: R$ 97,00 - 12 números WhatsApp
- ✅ Funções de gerenciamento de assinaturas (`subscriptions.go`)
- ✅ Validação de limites por plano (`CanCreateInstance`)
- ✅ Criação automática de assinatura trial no cadastro
- ✅ Verificação de expiração de planos

### 2. Sistema de Autenticação
- ✅ Tabela `system_users` criada
- ✅ Login/Cadastro de usuários com e-mail e senha
- ✅ Geração automática de token admin
- ✅ Isolamento de instâncias por usuário (`system_user_id`)

### 3. Sistema de Envio Diário
- ✅ Tabela `daily_conversations` criada
- ✅ Armazenamento de mensagens do dia
- ✅ Campo `destination_number` para número de destino
- ✅ Webhook fixo configurável (`daily_sender.go`)

### 4. Histórico de Mensagens
- ✅ Tabela `message_history` criada
- ✅ Armazenamento de mensagens com chat_jid
- ✅ Índices para busca eficiente

## ⚠️ PRECISA SER IMPLEMENTADO/CORRIGIDO

### 1. Frontend - Interface Completa
**Status**: Precisa ser totalmente refeito com base no `HTML_FRONTEND_REPLIQUE.md`

**Tarefas**:
- [ ] Criar nova estrutura HTML com Tailwind CSS
- [ ] Sidebar com navegação (Contas conectadas / Seus dados)
- [ ] Grid de 3 colunas para cards de instâncias
- [ ] Cards com bordas arredondadas e design moderno
- [ ] Status real-time (Conectado/Desconectado)
- [ ] Modal QR Code funcional
- [ ] Página "Seus dados" com informações do usuário
- [ ] Campo para número de destino global
- [ ] Exibição de plano atual e opções de upgrade
- [ ] Barra de progresso de contas disponíveis

### 2. Integração Frontend-Backend
**Tarefas**:
- [ ] Conectar API de status com atualização em tempo real
- [ ] Implementar polling de QR Code corretamente
- [ ] Atualizar status após conexão bem-sucedida
- [ ] Implementar validação de limites ao criar instância
- [ ] Exibir mensagem de limite atingido
- [ ] Carregar informações do plano do usuário
- [ ] Salvar/carregar número de destino

### 3. Sistema de Envio Diário
**Tarefas**:
- [ ] Implementar cron job para 18h BRT
- [ ] Consolidar mensagens do dia por instância
- [ ] Enviar para webhook fixo com parâmetro "enviar_para"
- [ ] Limpar tabela `daily_conversations` após envio

### 4. Correções de QR Code
**Tarefas Concluídas Hoje**:
- ✅ QR Code sendo gerado no backend
- ✅ Endpoint `/session/qr` funcionando
- ✅ Conexão WhatsApp funcionando

**Pendente**:
- [ ] Sincronização perfeita frontend-backend
- [ ] Atualização automática de status após scan
- [ ] Fechar modal automaticamente ao conectar

### 5. Histórico de Mensagens ao Logar
**Tarefas**:
- [ ] Buscar últimas 100 mensagens por conversa ao conectar
- [ ] Armazenar no `message_history`
- [ ] Exibir histórico se solicitado

## 🔧 STACK TECNOLÓGICA

### Backend
- **Linguagem**: Go (Golang)
- **Framework Web**: Chi Router
- **Banco de Dados**: SQLite
- **ORM**: sqlx
- **WhatsApp**: whatsmeow library

### Frontend Atual
- **Framework**: HTML + JavaScript Vanilla
- **UI Library**: Fomantic UI (Semantic UI fork)
- **AJAX**: Fetch API

### Frontend Proposto (HTML_FRONTEND_REPLIQUE.md)
- **Framework**: HTML + JavaScript Vanilla
- **CSS Framework**: Tailwind CSS
- **Fonte**: Inter (Google Fonts)
- **Ícones**: SVG inline (Feather Icons style)

## 📋 PRÓXIMOS PASSOS PRIORITÁRIOS

### Prioridade 1: Frontend Funcional
1. Recriar `user-dashboard-v2.html` com novo design
2. Atualizar `user-dashboard-v2.js` para nova estrutura
3. Implementar grid 3 colunas
4. Corrigir exibição de status real-time

### Prioridade 2: Validação de Planos
1. Verificar limites ao criar nova instância
2. Exibir plano atual no frontend
3. Mostrar contador de instâncias restantes
4. Bloquear criação ao atingir limite

### Prioridade 3: Envio Diário
1. Criar job agendado (18h BRT)
2. Consolidar mensagens
3. Enviar para webhook fixo
4. Testar envio manual

### Prioridade 4: Número de Destino
1. Adicionar campo na página "Seus dados"
2. Salvar no banco (tabela `system_users`)
3. Incluir no payload do webhook

## 🐛 BUGS CONHECIDOS

### Resolvidos Hoje
- ✅ Database locked (SQLITE_BUSY) - Resolvido reiniciando serviço
- ✅ QR Code não aparecendo no frontend - Em correção
- ✅ Status não atualiza após conectar - Em correção

### Pendentes
- [ ] Polling de QR Code retorna erro 500 após conexão bem-sucedida
- [ ] Frontend não detecta conexão estabelecida
- [ ] Modal QR Code não fecha automaticamente

## 📊 ESTRUTURA DO BANCO DE DADOS

### Tabelas Principais
1. `system_users` - Usuários do sistema (login/senha)
2. `users` - Instâncias WhatsApp
3. `plans` - Planos disponíveis
4. `user_subscriptions` - Assinaturas ativas
5. `subscription_history` - Histórico de assinaturas
6. `message_history` - Histórico de mensagens
7. `daily_conversations` - Conversas consolidadas do dia

## 🔐 WEBHOOK PADRÃO CONFIGURADO

```
https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5
```

**Payload do envio diário**:
```json
{
  "instance_name": "Allan",
  "instance_id": "abc123",
  "date": "2025-11-04",
  "enviar_para": "+5551981936133",
  "conversations": [
    {
      "contact": "+5551999999999",
      "chat_jid": "5551999999999@s.whatsapp.net",
      "messages": [
        {
          "timestamp": "2025-11-04T10:30:00Z",
          "from_me": false,
          "text": "Olá!"
        }
      ]
    }
  ]
}
```

## 📝 OBSERVAÇÕES IMPORTANTES

1. **Token Admin**: Gerado automaticamente no cadastro/login - usuário não precisa copiar
2. **Configurações no Header**: Removidas da interface (requisito #2)
3. **Webhook nas Configurações**: Não deve aparecer (webhook é fixo e global)
4. **Plano Trial**: 5 dias automáticos ao cadastrar
5. **Instâncias Grid**: 3 colunas em desktop, responsivo em mobile
