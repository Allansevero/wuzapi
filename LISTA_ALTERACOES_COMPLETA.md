# Lista Completa de Alterações do Sistema

## 1. Sistema de Autenticação por Usuário
- [x] Cada usuário terá e-mail e senha para acessar
- [x] Usuários podem ver somente as instâncias relacionadas à sua conta
- [x] Token admin gerado automaticamente no cadastro/login
- [x] Fluxo direto: cadastro/login → dashboard (sem necessidade de preencher token)

## 2. Interface do Dashboard
- [x] Remover configurações do cabeçalho ao entrar na instância
- [x] QR Code deve aparecer corretamente ao clicar em "Conectar"
- [x] Instâncias devem ficar em grid de 3 colunas com bordas arredondadas
- [x] Status "Conectado" só aparece quando realmente conectado ao WhatsApp
- [x] Botão de conectar com código disponível
- [x] Status atualiza automaticamente após conexão

## 3. Sistema de Envio Diário de Mensagens
- [ ] **IMPLEMENTAR**: Envio diário às 18h (horário de Brasília)
- [ ] Webhook padrão do sistema: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- [ ] Webhook não deve aparecer na configuração das instâncias
- [ ] Todas as conversas do dia são compiladas e enviadas juntas
- [ ] Formato: um único envio por instância com todas as mensagens do dia

## 4. Popup de Número de Recebimento
- [ ] **IMPLEMENTAR**: Botão que abre popup para inserir número
- [ ] Número do usuário para receber mensagens
- [ ] Parâmetro "enviar_para" incluído no envio das 18h
- [ ] Número associado à conta do usuário

## 5. Sistema de Planos e Limitações
- [ ] **IMPLEMENTAR**: Plano Gratuito (5 dias)
  - Números ilimitados de WhatsApp conectados
  - Período de teste de 5 dias
  
- [ ] **IMPLEMENTAR**: Plano Pro (R$ 29,00/mês)
  - Até 5 números de WhatsApp conectados
  - Sem limite de tempo
  
- [ ] **IMPLEMENTAR**: Plano Analista (R$ 97,00/mês)
  - Até 12 números de WhatsApp conectados
  - Sem limite de tempo

- [ ] **IMPLEMENTAR**: Tabela no banco de dados para planos
  - Campo: user_id
  - Campo: plan_type (free, pro, analyst)
  - Campo: plan_start_date
  - Campo: plan_end_date
  - Campo: max_instances
  - Campo: status (active, expired, cancelled)

## 6. Interface de Planos
- [ ] **IMPLEMENTAR**: Replicar layout do HTML_FRONTEND_REPLIQUE.md
- [ ] Design moderno com cards de planos
- [ ] Seleção e upgrade de planos
- [ ] Indicador visual do plano atual
- [ ] Contador de instâncias usadas/disponíveis

## 7. Histórico de Mensagens
- [ ] **IMPLEMENTAR**: Puxar histórico ao conectar
- [ ] Últimas 100 mensagens por conversa
- [ ] Armazenar histórico no banco de dados
- [ ] Incluir histórico no envio diário

## 8. Correções Técnicas
- [x] QR Code aparece no frontend após gerar no backend
- [x] Status de conexão sincroniza entre backend e frontend
- [x] Resolver erro "database is locked"
- [x] Resolver erro "address already in use"
- [x] Polling de status funciona corretamente

## Prioridades de Implementação

### Alta Prioridade (Fazer Agora)
1. ✅ Autenticação e isolamento de usuários
2. ✅ Correção do QR Code e status de conexão
3. **Sistema de envio diário às 18h**
4. **Popup de número de recebimento**

### Média Prioridade (Próximos Passos)
5. **Sistema de planos e limitações**
6. **Interface de planos (replicar HTML)**
7. **Histórico de mensagens ao conectar**

### Baixa Prioridade (Melhorias Futuras)
8. Otimizações de performance
9. Logs e monitoramento
10. Documentação adicional

## Notas Técnicas

### Stack Frontend
- HTML puro + JavaScript
- CSS customizado
- Sem frameworks (React, Vue, etc)

### Webhook Padrão
- URL: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- Não configurável pelo usuário
- Envio único diário às 18h BRT

### Estrutura de Dados do Envio Diário
```json
{
  "instance_name": "Nome da Instância",
  "user_id": "ID do Usuário",
  "date": "2025-11-04",
  "enviar_para": "+5551999999999",
  "conversations": [
    {
      "chat_id": "555199999999@s.whatsapp.net",
      "messages": [
        {
          "id": "msg_id",
          "timestamp": "2025-11-04T10:30:00Z",
          "from": "sender",
          "body": "mensagem",
          "type": "text"
        }
      ]
    }
  ]
}
```

## Status Geral
- ✅ Autenticação: **COMPLETO**
- ✅ Dashboard Básico: **COMPLETO**
- ✅ Conexão WhatsApp: **COMPLETO**
- ⏳ Envio Diário: **PENDENTE**
- ⏳ Sistema de Planos: **PENDENTE**
- ⏳ Interface de Planos: **PENDENTE**
- ⏳ Histórico: **PENDENTE**
