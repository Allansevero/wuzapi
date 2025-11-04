# Lista de Alterações Necessárias - Sistema Wuzapi

## 1. Sistema de Autenticação de Usuários
- Cada usuário terá e-mail e senha para acessar
- Usuários podem ver somente as instâncias relacionadas a sua conta
- Token admin deve ser gerado automaticamente no cadastro/login
- Após login, redirecionar direto para dashboard (sem necessidade de inserir token)

## 2. Interface do Usuário
- Remover configurações do cabeçalho ao entrar na instância
- Instâncias devem ser exibidas em grid de 3 colunas
- Cards de instâncias com bordas arredondadas
- Status "Conectado" deve aparecer apenas quando realmente conectado ao WhatsApp
- Melhorar experiência (remover necessidade de copiar token manualmente)

## 3. Sistema de Envio de Mensagens
- **Webhook Padrão do Sistema**: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- Webhook não deve aparecer na configuração das instâncias (fixo no sistema)
- Envio em lote: todas as conversas do dia enviadas às 18h (horário de Brasília)
- Formato: compilado diário de todas as mensagens da instância

## 4. Número de Destino para Mensagens
- Botão para abrir popup de configuração
- Campo para inserir número que receberá as mensagens
- Número deve ser enviado no parâmetro `enviar_para` junto com o compilado diário

## 5. Sistema de Planos e Assinaturas

### Plano Gratuito (Trial)
- **Duração**: 5 dias
- **WhatsApps**: Ilimitados
- **Valor**: R$ 0,00

### Plano Pro
- **Valor**: R$ 29,00/mês
- **WhatsApps**: Até 5 números conectados
- **Recurso**: Envio diário de conversas

### Plano Analista
- **Valor**: R$ 97,00/mês
- **WhatsApps**: Até 12 números conectados
- **Recurso**: Envio diário de conversas

### Funcionalidades do Sistema de Planos
- Controle de limites de instâncias por plano
- Validação de expiração do plano gratuito
- Sistema de upgrade/downgrade de planos
- Armazenamento de informações de plano no banco de dados
- Interface para gerenciamento de planos

## 6. Correções de Bugs

### QR Code e Conexão
- ✅ QR Code não estava aparecendo no frontend
- ✅ Botão "Conectar" não estava gerando QR Code corretamente
- ✅ Status de conexão não atualizava em tempo real
- ✅ Erro 500 ao tentar conectar instância
- ✅ Problema de "database is locked" (SQLITE_BUSY)

### Histórico de Mensagens
- Implementar pull de histórico ao fazer login
- Buscar últimas 100 mensagens por conversa
- Armazenar mensagens enviadas e recebidas após login

## 7. Layout Frontend (HTML_FRONTEND_REPLIQUE.md)
- Replicar design moderno conforme especificação
- Grid responsivo de 3 colunas para instâncias
- Cards com bordas arredondadas e sombras
- Cores e estilos conforme mockup fornecido

## Status de Implementação

### ✅ Concluído
1. Sistema de autenticação básico
2. Geração automática de token
3. Correção de bugs de conexão QR Code
4. Atualização de status em tempo real

### 🔄 Em Progresso
1. Sistema de planos e limitações
2. Interface de gerenciamento de planos
3. Replicação do layout frontend

### ⏳ Pendente
1. Envio diário compilado às 18h
2. Popup para configurar número de destino
3. Pull de histórico de mensagens ao login
4. Webhook fixo no sistema (não configurável por usuário)

## Notas Técnicas

- **Stack Frontend**: HTML puro + JavaScript (sem React/frameworks)
- **Banco de Dados**: SQLite
- **Horário**: Brasília (America/Sao_Paulo - UTC-3)
- **Webhook**: N8N fixo para todo o sistema
- **Autenticação**: Token-based com geração automática
