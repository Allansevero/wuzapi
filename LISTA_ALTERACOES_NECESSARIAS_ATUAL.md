# Lista de Alterações Necessárias no Sistema

## 1. Sistema de Autenticação por Usuário
- **Objetivo**: Cada usuário terá e-mail e senha para acessar
- **Funcionalidade**: Usuários só podem ver instâncias relacionadas à sua conta
- **Implementação**:
  - Sistema de login/cadastro com e-mail e senha
  - Token de admin gerado automaticamente no cadastro/login
  - Redirecionamento direto para dashboard após autenticação
  - Isolamento de instâncias por usuário

## 2. Remoção de Configurações no Cabeçalho
- **Objetivo**: Remover opções de configuração ao entrar na instância
- **Implementação**:
  - Ocultar menu de configurações no header
  - Simplificar interface do usuário

## 3. Sistema de Envio Diário de Mensagens
- **Objetivo**: Consolidar envio de mensagens em lote diário
- **Funcionalidade**: 
  - Envio diário às 18h (horário de Brasília)
  - Todas as conversas do dia de cada instância
  - Webhook padrão único para todo o sistema
- **Webhook Padrão**: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- **Implementação**:
  - Armazenar mensagens durante o dia
  - Cron job para envio às 18h BRT
  - Webhook fixo (não aparece nas configurações)
  - Substituir envio individual por envio em lote

## 4. Campo de Número de Destino
- **Objetivo**: Permitir usuário definir número para receber mensagens
- **Funcionalidade**:
  - Botão que abre popup para inserir número
  - Número incluído no parâmetro "enviar_para" no webhook
- **Implementação**:
  - Modal/popup para inserção de número
  - Armazenar número por usuário/instância
  - Incluir no payload do webhook diário

## 5. Sistema de Planos e Assinaturas
- **Plano Gratuito (Trial)**:
  - Duração: 5 dias
  - WhatsApp: Números ilimitados
  - Preço: Gratuito

- **Plano Pro**:
  - Preço: R$ 29,00/mês
  - WhatsApp: 5 números conectados
  - Recursos: Completos

- **Plano Analista**:
  - Preço: R$ 97,00/mês
  - WhatsApp: 12 números conectados
  - Recursos: Completos

- **Implementação**:
  - Tabela de planos no banco de dados
  - Tabela de assinaturas de usuários
  - Validação de limites por plano
  - Sistema de expiração (trial de 5 dias)
  - Bloqueio ao atingir limite de números

## 6. Interface de Usuário (Frontend)
- **Layout**:
  - Instâncias em grid de 3 colunas
  - Cards com bordas arredondadas
  - Status "Conectado" apenas quando realmente conectado
  - Botão de conexão visível e funcional
  - QR Code exibido corretamente no modal

- **Funcionalidades**:
  - Remoção da necessidade de copiar token
  - Atualização em tempo real do status de conexão
  - Interface intuitiva e moderna
  - Histórico das últimas 100 mensagens por conversa ao logar

## 7. Correções Técnicas Necessárias
- **QR Code**:
  - Exibição correta no frontend
  - Sincronização backend-frontend
  - Atualização de status após conexão

- **Banco de Dados**:
  - Resolver problemas de lock do SQLite
  - Otimizar queries de mensagens
  - Implementar armazenamento de histórico

- **Webhook**:
  - Implementar sistema de envio diário
  - Consolidar mensagens do dia
  - Incluir parâmetro "enviar_para"

## Status de Implementação
- [ ] Sistema de autenticação por usuário
- [ ] Geração automática de token admin
- [ ] Remoção de configurações do header
- [ ] Sistema de envio diário às 18h
- [ ] Campo de número de destino
- [ ] Sistema de planos e assinaturas
- [ ] Validação de limites por plano
- [ ] Interface com grid 3 colunas
- [ ] Correção de exibição de QR Code
- [ ] Atualização de status em tempo real
- [ ] Histórico de 100 mensagens por conversa

## Próximos Passos
1. Implementar sistema de planos no banco de dados
2. Criar validações de limite de instâncias
3. Implementar envio diário de mensagens
4. Atualizar interface conforme especificações
5. Testar fluxo completo de usuário
6. Implementar sistema de expiração de trial
