# Requisitos e Alterações do Sistema WuzAPI

## Data de Criação
04 de Novembro de 2025

## Visão Geral
Este documento descreve todas as alterações necessárias no sistema WuzAPI para implementar o novo modelo de funcionamento com autenticação por usuário e envio diário consolidado de mensagens.

---

## 1. Sistema de Autenticação Multi-Usuário

### 1.1 Autenticação e Autorização
- ✅ **Implementado**: Cada usuário possui e-mail e senha para login
- ✅ **Implementado**: Token de acesso gerado automaticamente no cadastro/login
- ✅ **Implementado**: Usuários visualizam apenas suas próprias instâncias
- ✅ **Implementado**: Dashboard acessível imediatamente após login sem necessidade de inserir token

### 1.2 Gerenciamento de Sessão
- ✅ **Implementado**: Token JWT armazenado em localStorage
- ✅ **Implementado**: Middleware de autenticação em todas as rotas protegidas
- ✅ **Implementado**: Redirecionamento automático para login se não autenticado

---

## 2. Interface do Usuário

### 2.1 Dashboard de Instâncias
- ✅ **Implementado**: Layout em grade com 3 colunas
- ✅ **Implementado**: Cards de instância com bordas arredondadas
- ✅ **Implementado**: Botão "Conectar" para gerar QR Code
- ✅ **Implementado**: Exibição do QR Code em modal/popup
- ✅ **Implementado**: Status "Conectado" apenas quando WhatsApp está realmente conectado
- ✅ **Implementado**: Atualização automática de status via polling

### 2.2 Remoção de Elementos
- ✅ **Implementado**: Configurações removidas do cabeçalho ao entrar na instância
- ✅ **Implementado**: Campo de webhook removido da configuração de instâncias
- ✅ **Implementado**: Sistema de token manual removido da experiência do usuário

### 2.3 Configuração de Número de Destino
- ⏳ **Pendente**: Botão para abrir popup de configuração
- ⏳ **Pendente**: Modal para inserir número de telefone para receber mensagens
- ⏳ **Pendente**: Validação do formato do número
- ⏳ **Pendente**: Armazenamento do número no banco de dados (campo `destination_number`)
- ⏳ **Pendente**: Exibição do número configurado no card da instância

---

## 3. Sistema de Envio Diário de Mensagens

### 3.1 Armazenamento de Mensagens
- ✅ **Implementado**: Tabela `message_history` para armazenar todas as mensagens
- ✅ **Implementado**: Captura de mensagens recebidas e enviadas
- ✅ **Implementado**: Armazenamento de metadados (timestamp, tipo, conteúdo)
- ⏳ **Pendente**: Buscar histórico das últimas 100 mensagens por conversa ao conectar

### 3.2 Envio Diário Automatizado
- ✅ **Implementado**: Cron job para executar às 18:00 horário de Brasília
- ✅ **Implementado**: Agrupamento de mensagens por conversa
- ✅ **Implementado**: Estrutura do payload com todas as conversas do dia
- ✅ **Implementado**: Envio para webhook fixo
- ✅ **Implementado**: Parâmetro `enviar_para` com número de destino
- ✅ **Implementado**: Endpoint para teste manual de envio

### 3.3 Webhook Fixo
- ✅ **Implementado**: URL fixa do webhook configurada no sistema
- ✅ **Implementado**: Webhook: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- ✅ **Implementado**: Não visível na configuração das instâncias
- ✅ **Implementado**: Timeout de 30 segundos para requisições

---

## 4. Estrutura do Payload de Envio Diário

```json
{
  "instance_id": "uuid-da-instancia",
  "date": "2025-11-04",
  "enviar_para": "+5511999999999",
  "conversations": [
    {
      "contact": "5511888888888@s.whatsapp.net",
      "messages": [
        {
          "sender_jid": "5511888888888@s.whatsapp.net",
          "message_type": "text",
          "text_content": "Olá, como vai?",
          "media_link": "",
          "timestamp": "2025-11-04T10:30:00Z",
          "data": {}
        }
      ]
    }
  ]
}
```

---

## 5. Banco de Dados

### 5.1 Tabelas Principais

#### Tabela: `users`
- `id` - ID único do usuário/instância
- `email` - E-mail de login ✅
- `password` - Senha hash ✅
- `token` - Token de acesso API ✅
- `name` - Nome da instância ✅
- `jid` - JID do WhatsApp conectado ✅
- `destination_number` - Número para receber compilado diário ⏳
- `created_at` - Data de criação ✅

#### Tabela: `message_history`
- `id` - ID único da mensagem ✅
- `user_id` - ID do usuário/instância ✅
- `chat_jid` - JID do chat/conversa ✅
- `sender_jid` - JID do remetente ✅
- `message_type` - Tipo da mensagem (text, image, etc) ✅
- `text_content` - Conteúdo texto ✅
- `media_link` - Link para mídia ✅
- `timestamp` - Data/hora da mensagem ✅
- `datajson` - JSON completo da mensagem ✅

---

## 6. Melhorias Técnicas

### 6.1 Frontend
- Stack: HTML puro + JavaScript Vanilla (sem frameworks)
- ✅ Polling automático de status
- ✅ Modais para QR Code
- ⏳ Modal para configuração de número de destino
- ⏳ Validação de formulários
- ⏳ Feedback visual de ações

### 6.2 Backend
- ✅ Conexão SQLite com pool de conexões
- ✅ Middleware de autenticação
- ✅ Handlers separados por funcionalidade
- ✅ Sistema de cron para tarefas agendadas
- ✅ Logs estruturados com zerolog

### 6.3 Segurança
- ✅ Senhas hasheadas com bcrypt
- ✅ Tokens JWT para sessões
- ✅ Validação de permissões por usuário
- ✅ CORS configurado
- ✅ SQL injection prevention

---

## 7. Funcionalidades Pendentes de Implementação

### 7.1 Configuração de Número de Destino
**Prioridade: ALTA**

Frontend:
- [ ] Botão "Configurar Número" no card da instância
- [ ] Modal com formulário para inserir número
- [ ] Máscara de entrada para formato de telefone
- [ ] Validação de formato (internacional)
- [ ] Exibição do número atual (se configurado)

Backend:
- [ ] Endpoint POST `/api/user/destination-number`
- [ ] Validação do número no backend
- [ ] Atualização no banco de dados
- [ ] Retorno de confirmação

### 7.2 Busca de Histórico ao Conectar
**Prioridade: MÉDIA**

- [ ] Ao conectar WhatsApp, buscar últimas 100 mensagens por conversa
- [ ] Armazenar no `message_history`
- [ ] Evitar duplicatas (verificar por message ID)
- [ ] Processar em background para não travar conexão

### 7.3 Melhorias na Interface
**Prioridade: MÉDIA**

- [ ] Loading states durante ações
- [ ] Mensagens de erro mais descritivas
- [ ] Confirmação antes de desconectar instância
- [ ] Indicador visual de sincronização
- [ ] Botão de atualizar status manual

---

## 8. Testes Necessários

### 8.1 Testes de Integração
- [ ] Teste completo de fluxo de cadastro
- [ ] Teste de login e geração de token
- [ ] Teste de conexão WhatsApp via QR Code
- [ ] Teste de captura de mensagens
- [ ] Teste de envio diário às 18h
- [ ] Teste de envio manual

### 8.2 Testes de Carga
- [ ] Múltiplos usuários simultâneos
- [ ] Múltiplas instâncias por usuário
- [ ] Grande volume de mensagens
- [ ] Comportamento do cron com muitos usuários

---

## 9. Documentação

### 9.1 Documentação Técnica
- ✅ API endpoints documentados
- ✅ Estrutura do banco de dados
- ✅ Fluxo de autenticação
- ⏳ Guia de deploy em produção
- ⏳ Configuração de variáveis de ambiente

### 9.2 Documentação do Usuário
- ⏳ Como criar conta
- ⏳ Como conectar WhatsApp
- ⏳ Como configurar número de destino
- ⏳ Entendendo o envio diário
- ⏳ Troubleshooting

---

## 10. Checklist de Produção

- [ ] Variáveis de ambiente configuradas
- [ ] Banco de dados com backup automático
- [ ] HTTPS configurado
- [ ] Logs centralizados
- [ ] Monitoramento de saúde
- [ ] Alertas configurados
- [ ] Documentação completa
- [ ] Testes executados com sucesso
- [ ] Plano de rollback definido

---

## Status Geral

| Categoria | Status |
|-----------|--------|
| Autenticação | ✅ Completo |
| Dashboard Básico | ✅ Completo |
| Conexão WhatsApp | ✅ Completo |
| Envio Diário | ✅ Completo |
| Config. Número Destino | ⏳ Pendente |
| Histórico ao Conectar | ⏳ Pendente |
| Testes | ⏳ Parcial |
| Documentação | ⏳ Parcial |

**Legenda:**
- ✅ Implementado e testado
- ⏳ Pendente ou em progresso
- ❌ Bloqueado ou com problema

---

## Próximos Passos Imediatos

1. Implementar modal de configuração de número de destino
2. Criar endpoint para salvar número de destino
3. Implementar busca de histórico ao conectar WhatsApp
4. Realizar testes completos de fluxo
5. Criar documentação de usuário
6. Preparar para deploy em produção
