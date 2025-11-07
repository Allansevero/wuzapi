# Alteração: Informações da Instância nos Webhooks

**Data:** 2025-11-07  
**Tipo:** Melhoria de Webhooks

## 📋 Alteração Implementada

Adicionadas informações da instância em **todos os eventos** enviados para webhooks:
- ✅ Nome da instância (`instance_name`)
- ✅ Número do WhatsApp da instância (`instance_phone`)
- ✅ JID completo da instância (`instance_jid`)

## 🎯 Objetivo

Permitir que o webhook identifique facilmente de qual instância/número veio cada evento, facilitando:
- Roteamento de mensagens
- Identificação de origem
- Logs e auditoria
- Integração com sistemas externos

## 📦 Dados Adicionados ao Payload

### Antes:
```json
{
  "type": "Message",
  "info": { ... },
  "message": { ... }
}
```

### Agora:
```json
{
  "type": "Message",
  "instance_name": "Instância Padrão",
  "instance_phone": "5511999999999",
  "instance_jid": "5511999999999@s.whatsapp.net",
  "info": { ... },
  "message": { ... }
}
```

## 🔧 Campos Adicionados

| Campo | Tipo | Descrição | Exemplo |
|-------|------|-----------|---------|
| `instance_name` | string | Nome da instância no sistema | "Instância Padrão" |
| `instance_phone` | string | Número do WhatsApp (extraído do JID) | "5511999999999" |
| `instance_jid` | string | JID completo do WhatsApp | "5511999999999@s.whatsapp.net" |

## 📝 Onde é Aplicado

**Todos os eventos enviados para webhook**, incluindo:
- ✅ Message (mensagens recebidas)
- ✅ MessageStatus (status de entrega)
- ✅ Receipt (confirmações de leitura)
- ✅ QR (código QR para conexão)
- ✅ Connected (conexão estabelecida)
- ✅ Disconnected (desconexão)
- ✅ HistorySync (sincronização de histórico)
- ✅ Todos os demais eventos

## 🔍 Como Funciona

### 1. Extração das Informações
```go
// Busca do cache do usuário
userinfo, found := userinfocache.Get(mycli.token)
if found {
    instanceName := userinfo.(Values).Get("Name")
    instanceJid := userinfo.(Values).Get("Jid")
    
    // Extrai número do JID (5511999999999@s.whatsapp.net → 5511999999999)
    instancePhone := strings.Split(instanceJid, "@")[0]
}
```

### 2. Adição ao Payload
```go
postmap["instance_name"] = instanceName
postmap["instance_phone"] = instancePhone
postmap["instance_jid"] = instanceJid
```

### 3. Envio para Webhook
O payload completo com as informações da instância é enviado para:
- Webhook do usuário
- Webhook global (se configurado)
- RabbitMQ (se configurado)

## 💡 Casos de Uso

### 1. Roteamento de Mensagens por Instância
```javascript
// No seu webhook
app.post('/webhook', (req, res) => {
    const { instance_phone, instance_name, message } = req.body;
    
    console.log(`Mensagem recebida de: ${instance_name} (${instance_phone})`);
    
    // Rotear baseado no número
    if (instance_phone === '5511999999999') {
        // Processar para equipe de vendas
    } else if (instance_phone === '5511888888888') {
        // Processar para suporte
    }
});
```

### 2. Identificação em Logs
```javascript
// Logging estruturado
logger.info({
    event: 'message_received',
    instance: req.body.instance_name,
    phone: req.body.instance_phone,
    from: req.body.info.sender,
    message: req.body.message
});
```

### 3. Múltiplas Instâncias
```javascript
// Usuário com múltiplas instâncias
const instances = {
    '5511999999999': { team: 'sales', webhook: 'https://sales.example.com' },
    '5511888888888': { team: 'support', webhook: 'https://support.example.com' }
};

const config = instances[req.body.instance_phone];
// Processar de acordo com a instância
```

## 📊 Exemplo Completo de Payload

### Mensagem de Texto Recebida:
```json
{
  "type": "Message",
  "instance_name": "WhatsApp Vendas",
  "instance_phone": "5511999999999",
  "instance_jid": "5511999999999@s.whatsapp.net",
  "info": {
    "id": "ABC123",
    "timestamp": "2025-11-07T17:30:00.000Z",
    "sender": "5511888888888@s.whatsapp.net",
    "chat": "5511888888888@s.whatsapp.net",
    "pushName": "João Silva",
    "isFromMe": false,
    "isGroup": false
  },
  "message": {
    "conversation": "Olá, gostaria de fazer um pedido"
  }
}
```

### Evento de Status:
```json
{
  "type": "Receipt",
  "instance_name": "WhatsApp Vendas",
  "instance_phone": "5511999999999",
  "instance_jid": "5511999999999@s.whatsapp.net",
  "event": "read",
  "messageId": "ABC123",
  "timestamp": "2025-11-07T17:31:00.000Z",
  "recipient": "5511888888888@s.whatsapp.net"
}
```

## 🔄 Compatibilidade

**Totalmente retrocompatível!**
- ✅ Campos novos não quebram webhooks existentes
- ✅ Webhooks antigos continuam funcionando
- ✅ Apenas ignoram os novos campos se não precisarem

## 📝 Arquivo Modificado

**wmiau.go**
- Função: `sendEventWithWebHook()`
- Linha: ~191-237
- Alteração: Adicionados 3 campos ao postmap antes de enviar

## 🧪 Como Testar

1. **Configure um webhook de teste:**
   ```bash
   # Usando webhook.site para visualizar
   # 1. Acesse https://webhook.site
   # 2. Copie a URL única gerada
   # 3. Configure no seu usuário
   ```

2. **Envie uma mensagem para sua instância**

3. **Veja o payload completo:**
   ```json
   {
     "instance_name": "...",
     "instance_phone": "...",
     "instance_jid": "...",
     ... resto dos dados
   }
   ```

## ⚠️ Observações

1. **Cache de Userinfo:**
   - Informações são buscadas do cache (`userinfocache`)
   - Se não estiver no cache, campos podem estar vazios
   - Cache é populado ao conectar a instância

2. **JID da Instância:**
   - Só estará disponível após a primeira conexão
   - Antes da conexão, `instance_jid` e `instance_phone` podem estar vazios
   - `instance_name` sempre estará disponível

3. **Formato do Número:**
   - `instance_phone` é o número puro: `5511999999999`
   - Sem `@s.whatsapp.net`
   - Ideal para comparações e roteamento

## 🎯 Benefícios

✅ Identificação imediata da origem  
✅ Facilita roteamento de eventos  
✅ Melhora logs e auditoria  
✅ Suporte a múltiplas instâncias  
✅ Integração mais simples  
✅ Retrocompatível  

## 🚀 Para Aplicar

**Apenas reinicie o servidor:**
```bash
sudo systemctl restart wuzapi
```

Todos os eventos novos já incluirão as informações da instância!
