# Lista de Alterações Necessárias - Wuzapi

## 1. Sistema de Autenticação por Usuário
- [x] Cada usuário terá e-mail e senha para acessar
- [x] Usuários podem ver somente as instâncias relacionadas à sua conta
- [x] Token admin gerado automaticamente ao fazer cadastro/login
- [x] Redirecionamento direto para dashboard após login

## 2. Interface do Dashboard
- [x] Remover configurações do cabeçalho ao entrar na instância
- [x] Instâncias exibidas em grid de 3 colunas com bordas arredondadas
- [x] Status "Conectado" apenas quando realmente conectado ao WhatsApp
- [x] QR Code exibindo corretamente no frontend
- [x] Botão de conectar com código funcional
- [x] Atualização automática do status após conexão

## 3. Sistema de Webhook Único e Envio Diário
- [x] **Webhook padrão global**: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- [x] Webhook não deve aparecer nas configurações das instâncias
- [x] Envio automático diário às 18:00 (horário de Brasília)
- [x] Compilar todas as conversas do dia de cada instância
- [x] Enviar em lote único ao invés de mensagem por mensagem
- [x] Teste manual de envio implementado

## 4. Parâmetro "enviar_para"
- [x] API para inserir número de destino
- [x] Número incluído no parâmetro "enviar_para" do webhook
- [x] Associar número ao usuário/instância
- [ ] Interface (popup) no frontend (pendente)

## 5. Sistema de Planos e Assinaturas
- [x] **Plano Gratuito (5 dias)**
  - Números ilimitados de WhatsApp
  - Período de teste de 5 dias
  - Criação automática ao cadastrar
  
- [x] **Plano Pro (R$ 29/mês)**
  - Até 5 números de WhatsApp conectados
  - Sem limitação de tempo
  
- [x] **Plano Analista (R$ 97/mês)**
  - Até 12 números de WhatsApp conectados
  - Sem limitação de tempo

- [x] Sistema de gerenciamento de planos no banco de dados
- [x] Validação de limites por plano
- [x] Controle de expiração do plano gratuito
- [x] API para consultar/atualizar planos
- [ ] Interface administrativa para gerenciar planos (pendente)

## 6. Histórico de Mensagens
- [x] Armazenar mensagens no banco de dados
- [x] Incluir mensagens no envio diário
- [ ] Puxar histórico das últimas 100 mensagens ao logar (pendente)
- [ ] Interface para visualizar histórico (pendente)

## Problemas Corrigidos
- ✅ QR Code não aparecendo no frontend
- ✅ Botão "Conectar com código" removido (restaurado)
- ✅ Erro 400/500 ao tentar conectar
- ✅ Status não atualizando após conexão
- ✅ Database locked (SQLITE_BUSY)
- ✅ Porta 8080 já em uso

## Stack Técnico
- **Backend**: Go (Golang) com whatsmeow
- **Frontend**: HTML + JavaScript (Vanilla JS)
- **Banco de Dados**: SQLite
- **Servidor Web**: HTTP nativo do Go

## Próximos Passos
1. Implementar sistema de envio diário às 18h
2. Configurar webhook global padrão
3. Criar popup para número "enviar_para"
4. Implementar sistema de planos
5. Adicionar histórico de mensagens
6. Criar painel administrativo para gerenciar planos
