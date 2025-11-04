# Próximas Implementações - WuzAPI

## Correções Aplicadas Hoje (04/11/2025) ✅

1. **Erro SQLITE_BUSY** - RESOLVIDO
   - Configurado WAL mode
   - Aumentado busy timeout para 10s
   - Pool de conexões otimizado para SQLite

2. **QR Code Não Aparece** - RESOLVIDO
   - Corrigido parsing de múltiplos formatos de resposta
   - Validação de imagem data:image
   - Logs detalhados para debug

---

## Pendências Críticas

### 1. Atualização de Status Pós-Conexão ⚠️ URGENTE

**Problema**:
- QR Code é escaneado com sucesso
- Backend registra conexão
- Frontend continua mostrando "Desconectado"

**Causa Provável**:
- Polling para após detectar conexão, mas UI não atualiza
- `loadInstances()` pode não estar sendo chamado corretamente
- Cache de status pode estar desatualizado

**Solução Proposta**:
```javascript
// Em user-dashboard-v2.js, linha ~275
if (statusData.data?.loggedIn === true && statusData.data?.jid) {
    console.log('✓ WhatsApp connected successfully! JID:', statusData.data.jid);
    clearInterval(qrPollingIntervals[instanceId]);
    delete qrPollingIntervals[instanceId];
    
    // Hide QR section
    $(`#qr-section-${instanceId}`).hide();
    
    showAlert('WhatsApp conectado com sucesso!', 'success');
    
    // FORÇAR atualização imediata
    setTimeout(() => {
        console.log('Reloading instances after successful connection...');
        loadInstances();
    }, 500); // Reduzir delay de 1500ms para 500ms
    return;
}
```

**Testes Necessários**:
1. Conectar via QR Code
2. Verificar se status muda para "Conectado" em até 2 segundos
3. Verificar se botão muda de "Conectar" para "Desconectar"
4. Verificar se número JID aparece corretamente

---

### 2. Sistema de Envio Diário às 18h ⏰ IMPORTANTE

**Status Atual**:
- Código base já existe em `daily_sender.go`
- Endpoint manual de teste criado: `POST /session/test-daily-send`

**O que falta**:
1. Ativar o scheduler automático (cron)
2. Testar envio real para o webhook
3. Validar formato de dados enviados

**Webhook Destino**:
```
https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5
```

**Formato de Dados Esperado**:
```json
{
  "instance_id": "abc123",
  "instance_name": "Minha Instância",
  "user_id": "user456",
  "date": "2025-11-04",
  "enviar_para": "+5551999999999",
  "conversations": [
    {
      "contact": "5551888888888@s.whatsapp.net",
      "contact_name": "João Silva",
      "messages": [
        {
          "timestamp": "2025-11-04T10:30:00-03:00",
          "from_me": false,
          "sender": "5551888888888@s.whatsapp.net",
          "type": "text",
          "content": "Olá, tudo bem?"
        },
        {
          "timestamp": "2025-11-04T10:31:00-03:00",
          "from_me": true,
          "sender": "5551999999999@s.whatsapp.net",
          "type": "text",
          "content": "Oi! Tudo ótimo e você?"
        }
      ]
    }
  ],
  "total_conversations": 1,
  "total_messages": 2
}
```

**Como Testar Manualmente**:
```bash
# Arquivo test_webhook_send.sh já criado
./test_webhook_send.sh SEU_TOKEN_AQUI
```

**Implementação do Scheduler**:
```go
// Em main.go ou arquivo separado
import (
    "github.com/robfig/cron/v3"
)

func startDailySender(s *server) {
    c := cron.New(cron.WithLocation(time.FixedZone("BRT", -3*60*60)))
    
    // Executar todos os dias às 18:00 horário de Brasília
    c.AddFunc("0 18 * * *", func() {
        log.Info().Msg("Starting daily message sender...")
        s.sendDailyMessagesToWebhook()
    })
    
    c.Start()
    log.Info().Msg("Daily sender scheduler started (18:00 BRT)")
}
```

---

### 3. Histórico de Mensagens ao Conectar 📜 MÉDIA

**Objetivo**:
Ao conectar uma instância, puxar automaticamente as últimas 100 mensagens de cada conversa.

**Status**:
- Código básico existe em `wmiau.go`
- Precisa validação e testes

**Verificar**:
```go
// Em wmiau.go
if _, ok := rawEvt.(*events.Connected); ok {
    go func() {
        time.Sleep(5 * time.Second)
        historyMsg := mycli.WAClient.BuildHistorySyncRequest(info.(*types.MessageInfo), 100)
        // Enviar requisição
    }()
}
```

**Testes Necessários**:
1. Conectar nova instância
2. Verificar logs se histórico foi solicitado
3. Confirmar se mensagens antigas aparecem no banco
4. Validar quantidade (deve ser ~100 mensagens por conversa)

---

### 4. Melhorias de Interface 🎨 BAIXA

**Ajustes Pendentes**:
- Layout responsivo para mobile
- Indicadores de carregamento melhores
- Animações de transição de estado
- Toast notifications ao invés de alerts
- Dark mode (opcional)

**Grid de 3 Colunas**:
Já implementado em `user-dashboard-v2.html`:
```html
<div class="ui three column stackable grid" id="instances-container"></div>
```

**Cards com Bordas Arredondadas**:
Já implementado no CSS:
```css
.instance-card {
    background: white;
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
```

---

## Plano de Ação Recomendado

### Hoje (04/11/2025)
- [x] Corrigir SQLITE_BUSY
- [x] Corrigir exibição de QR Code
- [ ] **PRÓXIMO**: Corrigir atualização de status pós-conexão
- [ ] **PRÓXIMO**: Testar envio manual para webhook

### Esta Semana
- [ ] Ativar scheduler de envio diário às 18h
- [ ] Validar histórico de mensagens ao conectar
- [ ] Testar sistema completo com dados reais
- [ ] Documentar API atualizada

### Próxima Semana
- [ ] Melhorias de interface
- [ ] Testes de carga
- [ ] Otimizações de performance
- [ ] Deploy em produção

---

## Como Continuar o Desenvolvimento

### 1. Corrigir Status Pós-Conexão
```bash
# Editar arquivo
nano static/dashboard/js/user-dashboard-v2.js

# Procurar linha ~275 e ajustar delay
# Recarregar página no navegador (Ctrl+Shift+R)
```

### 2. Testar Envio para Webhook
```bash
# Obter um token válido do banco
sqlite3 dbdata/users.db "SELECT token FROM users LIMIT 1;"

# Executar teste manual
./test_webhook_send.sh TOKEN_AQUI

# Verificar logs
tail -f wuzapi.log | grep webhook
```

### 3. Ativar Scheduler
```bash
# Editar main.go
nano main.go

# Adicionar chamada do scheduler
# Recompilar
go build -o wuzapi

# Reiniciar
sudo systemctl restart wuzapi
```

---

## Ferramentas de Debug

### Logs em Tempo Real
```bash
# Logs gerais
tail -f wuzapi.log

# Apenas conexões
tail -f wuzapi.log | grep -i "connect"

# Apenas QR codes
tail -f wuzapi.log | grep -i "qr"

# Apenas webhook
tail -f wuzapi.log | grep -i "webhook"
```

### Console do Navegador
```javascript
// Verificar instâncias carregadas
console.log(instances);

// Verificar intervalos de polling
console.log(qrPollingIntervals);

// Forçar reload de instâncias
loadInstances();
```

### Banco de Dados
```bash
# Abrir banco
sqlite3 dbdata/users.db

# Ver usuários
SELECT id, email, name FROM users;

# Ver instâncias
SELECT id, name, user_id, jid FROM users WHERE user_id IS NOT NULL;

# Ver mensagens de hoje
SELECT COUNT(*) FROM message_history WHERE DATE(timestamp) = DATE('now');
```

---

## Contato para Suporte

Se encontrar problemas:
1. Verificar logs: `tail -f wuzapi.log`
2. Verificar console do navegador (F12)
3. Testar endpoints manualmente com curl
4. Verificar banco de dados com sqlite3

---

**Última atualização**: 04/11/2025 09:42 UTC
**Próxima ação**: Corrigir atualização de status pós-conexão
