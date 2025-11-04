# Lista de Alterações Necessárias - WuzAPI

## Status das Alterações

### ✅ 1. Sistema de Autenticação por Usuário
**Status:** CONCLUÍDO
- [x] Cada usuário tem e-mail e senha para acessar
- [x] Usuários só veem instâncias relacionadas à sua conta
- [x] Token de admin é gerado automaticamente no cadastro/login
- [x] Usuário vai direto para dashboard após login (sem precisar preencher token)

**Arquivos modificados:**
- `auth.go` - Sistema de login/registro
- `user_instances.go` - Gerenciamento de instâncias por usuário
- `static/user-login.html` - Página de login
- `static/dashboard/user-dashboard-v2.html` - Dashboard do usuário

---

### ✅ 2. Remoção de Configurações no Cabeçalho
**Status:** CONCLUÍDO
- [x] Configurações não aparecem ao entrar na instância
- [x] Interface simplificada apenas com botões de ação essenciais

**Arquivos modificados:**
- `static/dashboard/user-dashboard-v2.html`
- `static/dashboard/js/user-dashboard-v2.js`

---

### ✅ 3. Envio Diário Consolidado de Mensagens
**Status:** CONCLUÍDO
- [x] Webhook padrão configurado: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`
- [x] Mensagens são enviadas consolidadas diariamente às 18h (horário de Brasília)
- [x] Webhook não aparece nas configurações das instâncias (é padrão do sistema)
- [x] Implementado `daily_sender.go` com scheduler

**Arquivos criados/modificados:**
- `daily_sender.go` - Novo arquivo com lógica de envio diário
- `main.go` - Inicialização do scheduler
- `migrations.go` - Tabela para armazenar mensagens do dia

**Detalhes técnicos:**
- Tabela `daily_messages` armazena mensagens recebidas durante o dia
- Scheduler executa às 18:00 BRT todos os dias
- Envia consolidado e limpa tabela após envio

---

### ✅ 4. Configuração de Número de Destino
**Status:** CONCLUÍDO
- [x] Modal para inserir número que receberá mensagens
- [x] Número é enviado no parâmetro `enviar_para` junto com mensagens do dia
- [x] Campo `destination_number` na tabela users

**Arquivos modificados:**
- `user_instances.go` - Endpoint para salvar número
- `handlers.go` - Handler para configurar destination_number
- `static/dashboard/user-dashboard-v2.html` - Modal de configuração
- `static/dashboard/js/user-dashboard-v2.js` - Lógica do modal

---

### 🔧 5. Interface de Usuário Melhorada
**Status:** PARCIALMENTE CONCLUÍDO

**Concluído:**
- [x] Instâncias exibidas em grid de 3 colunas
- [x] Cards com bordas arredondadas
- [x] Status "Conectado" só aparece quando realmente conectado
- [x] Botões de conexão (QR Code e Código de Pareamento)

**Pendente:**
- [ ] QR Code não está aparecendo no frontend (backend gera corretamente)
- [ ] Status não atualiza automaticamente após conexão bem-sucedida
- [ ] Possíveis erros 500 ao tentar obter QR code

**Arquivos envolvidos:**
- `static/dashboard/user-dashboard-v2.html`
- `static/dashboard/js/user-dashboard-v2.js`
- `handlers.go` - GetQR() e GetStatus()
- `wmiau.go` - Geração e armazenamento do QR code

---

## Problemas Conhecidos

### 🐛 Problema 1: QR Code não aparece no Frontend
**Sintoma:** Backend gera QR code corretamente, mas não é exibido no navegador

**Logs do backend mostram:**
```
2025-11-03 20:14:28 -03:00 INFO ... qrcode=data:image/png;base64,...
```

**Logs do frontend mostram:**
```javascript
QR JSON received: {code: 200, data: {...}, success: true}
No QR code in response
```

**Possível causa:**
- O handler `GetQR()` retorna `{"QRCode": "..."}` que é encapsulado em `{code: 200, data: {...}, success: true}`
- JavaScript deve acessar via `qrJson.data.QRCode`
- Pode haver problema de timing (QR ainda não foi gerado quando JS faz polling)

**Status:** EM INVESTIGAÇÃO

---

### 🐛 Problema 2: Status não atualiza após conexão
**Sintoma:** Usuário escaneia QR code, WhatsApp conecta, mas frontend continua mostrando "Desconectado"

**Logs mostram:**
```
2025-11-03 20:14:38 -03:00 INFO Marked self as available
2025-11-03 20:14:38 -03:00 INFO QR pairing ok!
```

**Mas o frontend continua fazendo polling sem detectar a conexão:**
```javascript
Status check: {connected: true, loggedIn: true, jid: "555181936133:64@s.whatsapp.net", ...}
```

**Possível causa:**
- O polling do QR continua mesmo após conexão bem-sucedida
- A lógica de verificação `statusData.loggedIn && statusData.jid` pode não estar sendo executada corretamente
- Pode haver problema com o timeout do polling

**Status:** EM INVESTIGAÇÃO

---

### 🐛 Problema 3: Erros 500 intermitentes
**Sintomas:**
- `database is locked (5) (SQLITE_BUSY)`
- `not connected`
- `no session`

**Causa:** Múltiplas requisições simultâneas ao SQLite

**Possível solução:**
- Implementar connection pooling
- Adicionar retry logic
- Migrar para PostgreSQL (recomendado para produção)

**Status:** CONHECIDO - Limitação do SQLite

---

## Próximos Passos

1. **URGENTE - Corrigir exibição de QR Code:**
   - Verificar formato exato da resposta JSON no frontend
   - Adicionar mais logs para debug
   - Garantir que o polling está acessando o caminho correto

2. **URGENTE - Corrigir atualização de status:**
   - Verificar lógica de detecção de conexão bem-sucedida
   - Garantir que `loadInstances()` é chamado após conexão
   - Adicionar feedback visual imediato

3. **Melhorias de Performance:**
   - Reduzir frequência de polling quando não necessário
   - Implementar WebSocket para updates em tempo real (opcional)
   - Considerar migração para PostgreSQL

4. **Testes:**
   - Testar fluxo completo de cadastro → login → conexão → envio
   - Verificar envio diário às 18h
   - Validar parâmetro `enviar_para` no webhook

---

## Arquivos Principais do Sistema

### Backend (Go)
- `main.go` - Entry point, inicializa servidor e scheduler
- `auth.go` - Autenticação de usuários
- `handlers.go` - Handlers HTTP para API
- `user_instances.go` - Gerenciamento de instâncias
- `daily_sender.go` - Envio diário consolidado
- `wmiau.go` - Integração com WhatsApp (whatsmeow)
- `db.go` - Gerenciamento do banco de dados
- `routes.go` - Definição de rotas

### Frontend (HTML/JS)
- `static/user-login.html` - Página de login
- `static/dashboard/user-dashboard-v2.html` - Dashboard principal
- `static/dashboard/js/user-dashboard-v2.js` - Lógica do dashboard

### Banco de Dados
- `dbdata/` - Arquivos SQLite
- Tabelas principais:
  - `system_users` - Usuários do sistema (email/senha)
  - `users` - Instâncias do WhatsApp
  - `daily_messages` - Mensagens para envio diário

---

## Configurações Importantes

### Webhook Padrão
```
https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5
```

### Horário de Envio Diário
- **Horário:** 18:00
- **Timezone:** America/Sao_Paulo (BRT)

### Formato do Payload Diário
```json
{
  "instance_id": "...",
  "instance_name": "...",
  "date": "2025-11-03",
  "total_messages": 10,
  "messages": [...],
  "enviar_para": "+5511999999999"
}
```

---

**Última atualização:** 2025-11-03 23:30 BRT
