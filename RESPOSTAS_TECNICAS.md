# Respostas Técnicas - WuzAPI

## Pergunta 1: Qual stack estamos usando no frontend?

**Resposta**: **HTML puro + JavaScript vanilla + Fomantic-UI**

- **Não** usamos React, Vue, Angular ou qualquer framework moderno
- **HTML estático** com JavaScript direto
- **Fomantic-UI**: Fork do Semantic UI para componentes visuais (grids, botões, modals, cards)
- **jQuery**: Usado internamente pelo Fomantic-UI

### Arquivos principais:
- `/static/dashboard/user-dashboard-v2.html` - Interface principal
- `/static/dashboard/js/user-dashboard-v2.js` - Lógica do dashboard
- `/static/user-login.html` - Tela de login/cadastro

### Vantagens desta abordagem:
- ✅ Simplicidade (sem build process)
- ✅ Fácil debug no navegador
- ✅ Carregamento rápido
- ✅ Sem dependências npm

### Desvantagens:
- ❌ Sem reatividade automática
- ❌ Precisa manipular DOM manualmente
- ❌ Mais código para manter estado sincronizado

---

## Pergunta 2: Puxamos histórico ao logar ou apenas armazenamos mensagens após login?

**Resposta Original**: Apenas armazenava mensagens após login (sem pull de histórico automático)

### Estado Anterior:
- ❌ Não puxava histórico automaticamente
- ✅ Tinha endpoint `/session/history` para solicitação manual
- ✅ Armazenava mensagens em `message_history` conforme chegavam
- ❌ Usuário perdia contexto de conversas anteriores

### Estado Atual (Após Implementação):
- ✅ **Pull automático de histórico ao conectar**
- ✅ Busca últimas **100 mensagens** por conversa
- ✅ Sincronização em background (não bloqueia conexão)
- ✅ Delay de 5 segundos para garantir estabilidade da conexão

### Implementação Técnica:

**Arquivo**: `wmiau.go`  
**Evento**: `*events.Connected`

```go
// Auto-request history sync when connected
if _, ok := rawEvt.(*events.Connected); ok {
    go func() {
        time.Sleep(5 * time.Second)  // Espera conexão estabilizar
        log.Info().Str("userID", txtid).Msg("Auto-requesting history sync after connection")

        info, found := lastMessageCache.Get(txtid)
        if !found {
            info = &types.MessageInfo{}
        }

        // Solicita últimas 100 mensagens
        historyMsg := mycli.WAClient.BuildHistorySyncRequest(info.(*types.MessageInfo), 100)
        if historyMsg != nil {
            _, err := mycli.WAClient.SendMessage(
                context.Background(), 
                mycli.WAClient.Store.ID.ToNonAD(), 
                historyMsg, 
                whatsmeow.SendRequestExtra{Peer: true}
            )
            if err != nil {
                log.Error().Str("userID", txtid).Err(err).Msg("Failed to auto-request history sync")
            } else {
                log.Info().Str("userID", txtid).Msg("History sync auto-requested successfully")
            }
        }
    }()
}
```

### Como funciona:

1. **Usuário conecta** instância via QR Code ou código
2. **Backend detecta** evento `*events.Connected`
3. **Aguarda 5 segundos** (estabilização)
4. **Solicita histórico** de 100 mensagens ao WhatsApp
5. **WhatsApp envia** mensagens via evento `*events.HistorySync`
6. **Backend armazena** na tabela `message_history`
7. **Mensagens ficam disponíveis** para compilação diária

### Benefícios:

- ✅ Contexto completo de conversas
- ✅ Compilação diária inclui histórico recente
- ✅ Não precisa intervenção manual
- ✅ Execução assíncrona (não trava)

---

## Pergunta 3: Como testar envio ao webhook com dados reais?

**Resposta**: Criado endpoint de teste manual

### Endpoint criado:
```
POST /session/test-daily-send
Authorization: Bearer <seu_token>
```

### Script de teste:
```bash
./test_webhook_send.sh <seu_token>
```

### O que o teste faz:

1. Busca todas as mensagens do **dia atual** da instância
2. Agrupa por **conversa** (chat_jid)
3. Compila no formato esperado pelo webhook
4. Envia para: `https://n8n-webhook.fmy2un.easypanel.host/webhook/0731c270-2870-4bf2-96b1-282ddd0532f5`

### Formato enviado:

```json
{
  "instance_id": "507a6d45c765c6ae5b720e3caa94fca2",
  "date": "2025-11-03",
  "conversations": [
    {
      "contact": "555181936133@s.whatsapp.net",
      "messages": [
        {
          "sender_jid": "555181936133@s.whatsapp.net",
          "message_type": "text",
          "text_content": "Olá, tudo bem?",
          "media_link": "",
          "timestamp": "2025-11-03T14:30:00-03:00",
          "data": { ... }
        }
      ]
    }
  ],
  "enviar_para": "5551999999999"
}
```

### Resposta esperada:

```json
{
  "success": true,
  "message": "Daily messages sent successfully to webhook",
  "date": "2025-11-03"
}
```

### Como usar:

1. **Conecte uma instância** no dashboard
2. **Envie/receba algumas mensagens** de teste
3. **Execute o script**:
   ```bash
   ./test_webhook_send.sh cc52b16537c624c981b5fb54f83e30cc
   ```
4. **Verifique no n8n** se o webhook recebeu os dados

### Nota importante:
- Só envia mensagens do **dia atual** (today)
- Se não houver mensagens hoje, retorna sucesso mas não envia nada
- Para testar com dados reais, precisa ter mensagens na `message_history` de hoje

---

## Problemas Identificados e Status

### 1. QR Code não aparece no frontend ⚠️
**Status**: QR Code é gerado no backend mas frontend não exibe  
**Causa**: Lógica de polling não está capturando o QR Code corretamente  
**Próximo passo**: Revisar `user-dashboard-v2.js` linha 218-313

### 2. Status não atualiza após conexão ⚠️
**Status**: Backend conecta mas frontend continua mostrando "Desconectado"  
**Causa**: Polling de status não detecta mudança  
**Logs mostram**: `"already logged in"` mas frontend não reflete  
**Próximo passo**: Adicionar WebSocket ou melhorar polling

### 3. Erro 500 em /session/qr 🔴
**Erros vistos**:
- `database is locked (SQLITE_BUSY)`
- `not connected`
- `no session`

**Causa**: Concorrência no SQLite  
**Próximo passo**: Implementar retry logic e mutex

---

## Resumo das Implementações

| Feature | Status | Arquivo | Descrição |
|---------|--------|---------|-----------|
| Pull automático de histórico | ✅ | `wmiau.go` | 100 mensagens ao conectar |
| Teste manual webhook | ✅ | `handlers.go`, `routes.go` | POST /session/test-daily-send |
| Script de teste | ✅ | `test_webhook_send.sh` | CLI para testar webhook |
| Documentação técnica | ✅ | `LISTA_ALTERACOES_SISTEMA.md` | Este arquivo |

**Data**: 2025-11-03  
**Próxima ação**: Corrigir exibição de QR Code e sincronização de status
