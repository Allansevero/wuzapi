# Alterações Necessárias no Sistema WuzAPI

## 1. Sistema de Autenticação por Usuário

**Objetivo:** Implementar sistema de login com email e senha, onde cada usuário visualiza apenas suas próprias instâncias.

### Mudanças necessárias:

- Criar tabela de usuários no banco de dados com campos:
  - `id` (chave primária)
  - `email` (único)
  - `password` (hash)
  - `created_at`
  - `updated_at`

- Adicionar coluna `user_id` na tabela de instâncias para relacionar cada instância a um usuário

- Criar endpoints de autenticação:
  - `POST /login` - autenticação de usuário
  - `POST /register` - registro de novos usuários
  - `POST /logout` - encerrar sessão

- Implementar middleware de autenticação para proteger rotas

- Filtrar instâncias por `user_id` em todas as consultas

- Implementar sistema de sessões ou JWT para gerenciar autenticação

---

## 2. Remover Configurações do Cabeçalho da Instância

**Objetivo:** Ocultar opções de configuração ao acessar uma instância.

### Mudanças necessárias:

- Modificar interface da instância para remover botões/links de configuração do cabeçalho

- Manter configurações acessíveis apenas através de painel administrativo (se necessário)

- Arquivos a modificar:
  - `/static/` (arquivos HTML/JavaScript relacionados à interface da instância)

---

## 3. Envio Diário de Conversas para Webhook

**Objetivo:** Substituir envio imediato de mensagens por envio consolidado diário às 18h (horário de Brasília) para webhook fixo.

### Mudanças necessárias:

- **Remover:** Sistema atual de envio de webhooks por mensagem

- **Implementar:**
  - Sistema de agendamento (cron job) que executa diariamente às 18h (horário de Brasília - UTC-3)
  - Armazenar todas as conversas do dia em banco de dados ou memória
  - Consolidar conversas por instância
  - Enviar payload para webhook fixo: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`

- **Estrutura do payload:**
```json
{
  "instance_id": "string",
  "date": "YYYY-MM-DD",
  "conversations": [
    {
      "contact": "número_contato",
      "messages": [...]
    }
  ],
  "enviar_para": "número_configurado_pelo_usuario"
}
```

- **Remover da configuração:** Campo de webhook personalizado por instância

- **Constante do sistema:** Webhook hardcoded, não editável pelos usuários

- **Timezone:** Configurar para America/Sao_Paulo (horário de Brasília)

---

## 4. Botão para Configurar Número de Recebimento

**Objetivo:** Adicionar interface para usuário configurar número que receberá as mensagens consolidadas.

### Mudanças necessárias:

- Adicionar coluna `numero_destino` na tabela de instâncias ou usuários

- Criar botão na interface com ícone/texto: "Configurar Número para Recebimento"

- Implementar popup/modal com:
  - Campo de input para número de telefone
  - Validação de formato de número
  - Botão "Salvar"
  - Botão "Cancelar"

- Criar endpoint:
  - `PUT /api/instance/:id/numero-destino` - atualizar número de destino

- Incluir campo `enviar_para` no payload enviado ao webhook às 18h

- **Formato do número:** Validar e armazenar no formato internacional (ex: +5511999999999)

---

## Resumo de Arquivos a Modificar

1. **db.go** - Adicionar tabelas e migrações
2. **handlers.go** - Novos endpoints de autenticação e configuração
3. **routes.go** - Novas rotas e middleware de autenticação
4. **main.go** - Configurar cron job para envio diário
5. **constants.go** - Adicionar webhook fixo como constante
6. **static/** - Interface para login, configuração de número e remoção de opções
7. **migrations.go** - Migrações de banco de dados

---

## Considerações Técnicas

- **Segurança:** Usar bcrypt para hash de senhas
- **Timezone:** Garantir uso correto do fuso horário de Brasília (America/Sao_Paulo)
- **Cron:** Biblioteca recomendada: `github.com/robfig/cron`
- **Sessões:** Considerar `gorilla/sessions` ou JWT com `golang-jwt/jwt`
- **Validação:** Validar formato de email e número de telefone

---

## Ordem de Implementação Sugerida

1. Sistema de autenticação e usuários
2. Relacionamento usuário-instâncias
3. Remoção de configurações do cabeçalho
4. Botão e funcionalidade de configuração de número
5. Sistema de coleta de conversas diárias
6. Implementação do cron job de envio às 18h
