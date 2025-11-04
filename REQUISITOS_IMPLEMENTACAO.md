# Requisitos de Implementação - Sistema WuzAPI

## 1. Sistema de Autenticação Multi-Usuário
- [ ] Cada usuário terá e-mail e senha para acessar
- [ ] Usuários só podem ver instâncias relacionadas à sua conta
- [ ] Token de acesso gerado automaticamente no cadastro/login
- [ ] Fluxo: Cadastro/Login → Dashboard (sem necessidade de preencher token admin)

## 2. Interface do Usuário
- [ ] Remover configurações do cabeçalho ao entrar na instância
- [ ] Instâncias exibidas em grid de 3 colunas com bordas arredondadas
- [ ] Status "Conectado" só aparece quando realmente conectado ao WhatsApp
- [ ] QR Code funcional com botão de conectar visível
- [ ] Atualização automática de status no frontend após conexão
- [ ] Botão para conectar via código (além do QR Code)

## 3. Sistema de Envio Diário de Mensagens
- [ ] Envio consolidado diário às 18h (horário de Brasília)
- [ ] Webhook fixo padrão: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- [ ] Webhook não deve aparecer nas configurações das instâncias
- [ ] Enviar todas as conversas do dia de cada instância
- [ ] Popup para inserir número de destino das mensagens
- [ ] Parâmetro "enviar_para" incluído no envio ao webhook

## 4. Histórico de Mensagens
- [ ] Ao conectar, puxar histórico das últimas 100 mensagens por conversa
- [ ] Armazenar histórico no banco de dados
- [ ] Sincronizar mensagens recebidas e enviadas após login

## 5. Sistema de Planos e Limitações

### Plano Gratuito (5 dias)
- Duração: 5 dias
- WhatsApp conectados: Ilimitados
- Preço: R$ 0,00

### Plano Pro
- Duração: Mensal
- WhatsApp conectados: 5 números
- Preço: R$ 29,00

### Plano Analista
- Duração: Mensal
- WhatsApp conectados: 12 números
- Preço: R$ 97,00

### Funcionalidades de Planos
- [ ] Tabela de planos no banco de dados
- [ ] Associação de usuário com plano
- [ ] Validação de limite de instâncias por plano
- [ ] Controle de expiração de plano gratuito
- [ ] Bloqueio de criação de novas instâncias ao atingir limite
- [ ] Interface para visualização do plano atual
- [ ] Sistema para upgrade/downgrade de planos
- [ ] Registro de histórico de planos por usuário

## 6. Correções Técnicas Realizadas
- [x] QR Code aparecendo corretamente no frontend
- [x] Status de conexão atualizado corretamente
- [x] Erro de "database is locked" resolvido
- [x] Erro de "address already in use" resolvido
- [x] Token admin gerado automaticamente
- [x] Instâncias em grid de 3 colunas

## 7. Próximas Implementações Prioritárias
1. Sistema de planos e limitações
2. Interface de gerenciamento de planos
3. Validação de limites por plano
4. Dashboard administrativo para gerenciar usuários e planos
5. Sistema de notificações de expiração de plano
6. Logs de uso por usuário/instância

## Observações Técnicas
- Frontend: HTML/CSS/JavaScript puro
- Backend: Go
- Banco de dados: SQLite
- Webhook fixo para todos os usuários
- Timezone: America/Sao_Paulo (Brasília)
