# Guia de Teste - WuzAPI

## 🧪 Roteiro de Testes Completo

### Pré-requisitos
- ✅ Servidor rodando em `http://localhost:8080`
- ✅ Navegador aberto (Chrome/Firefox recomendado)
- ✅ WhatsApp instalado no celular para escanear QR code

---

## Teste 1: Verificar Status do Sistema

### Backend
```bash
curl http://localhost:8080/health
```

**Resultado Esperado:**
```json
{
  "status": "ok",
  "active_connections": 1,
  "logged_in_users": 1,
  ...
}
```

### Frontend
1. Abrir navegador em `http://localhost:8080`
2. Deve redirecionar para `/user-login.html`
3. Página de login deve carregar sem erros no console

---

## Teste 2: Cadastro e Login

### 2.1 Criar Nova Conta
1. Ir para `http://localhost:8080/user-login.html`
2. Clicar em "Criar Conta" (se houver) ou usar formulário de registro
3. Preencher:
   - Email: `teste@example.com`
   - Senha: `senha123`
4. Clicar em "Cadastrar"

**Resultado Esperado:**
- ✅ Mensagem de sucesso
- ✅ Redireciona para dashboard
- ✅ Mostra "Instância Padrão" criada automaticamente

### 2.2 Fazer Login
1. Se já tiver conta, fazer login com credenciais
2. Verificar se redireciona para `/dashboard/user-dashboard-v2.html`

**Resultado Esperado:**
- ✅ Dashboard carrega
- ✅ Email aparece no cabeçalho
- ✅ Pelo menos uma instância listada

---

## Teste 3: Conexão via QR Code

### 3.1 Gerar QR Code
1. No dashboard, localizar a instância
2. Verificar se está marcada como "Desconectado"
3. Clicar no botão **"Conectar WhatsApp"**

**Resultado Esperado:**
- ✅ Botão muda ou fica desabilitado
- ✅ Área do QR code aparece (pode levar alguns segundos)
- ✅ QR code é exibido como imagem

**Console do Navegador (F12):**
```
Starting QR polling for instance: [instance-id]
QR Response status: 200
✓ Valid QR code found, displaying...
```

### 3.2 Escanear QR Code
1. Abrir WhatsApp no celular
2. Ir em: **Configurações → Aparelhos Conectados → Conectar um Aparelho**
3. Escanear o QR code exibido

**Resultado Esperado:**
- ✅ WhatsApp no celular confirma pareamento
- ✅ QR code desaparece automaticamente
- ✅ Status muda para "Conectado"
- ✅ Aparece o número conectado

**Console do Navegador:**
```
✓ WhatsApp connected successfully! JID: 5511999999999@s.whatsapp.net
Reloading instances after successful connection...
```

**Tempo Esperado:** 2-5 segundos após escanear

### 3.3 Se QR Code NÃO Aparecer
**Verificar:**
1. Console do navegador tem erros?
2. Backend está gerando QR? (ver logs)
3. Requisição `/session/qr` está retornando 200?

**Comandos de Debug:**
```bash
# Ver logs do backend
tail -f /home/allansevero/wuzapi/wuzapi.log | grep -i qr

# Testar endpoint manualmente
curl "http://localhost:8080/session/qr?token=SEU_TOKEN_AQUI"
```

---

## Teste 4: Conexão via Código de Pareamento

### 4.1 Solicitar Código
1. No dashboard, clicar em **"Código de Pareamento"**
2. Modal deve abrir
3. Digitar número de telefone: `+5511999999999`
4. Clicar em "Solicitar Código"

**Resultado Esperado:**
- ✅ Modal fecha
- ✅ Mensagem exibe código (ex: "Código de pareamento: ABCD-1234")
- ✅ WhatsApp no celular mostra notificação

### 4.2 Inserir Código no WhatsApp
1. Abrir WhatsApp
2. Ir em: **Configurações → Aparelhos Conectados → Conectar com Código**
3. Digitar o código recebido

**Resultado Esperado:**
- ✅ WhatsApp conecta
- ✅ Dashboard atualiza status para "Conectado"

---

## Teste 5: Configurar Número de Destino

### 5.1 Abrir Modal
1. Clicar em **"Config. Destino"** na instância
2. Modal deve abrir

### 5.2 Salvar Número
1. Digitar: `+5511888888888`
2. Clicar em "Salvar"

**Resultado Esperado:**
- ✅ Modal fecha
- ✅ Mensagem de sucesso
- ✅ Número aparece na instância como "Destino: +5511888888888"

**Verificação no Banco:**
```bash
sqlite3 /home/allansevero/wuzapi/dbdata/users.db \
  "SELECT name, destination_number FROM users;"
```

---

## Teste 6: Criar Nova Instância

### 6.1 Criar Instância
1. Clicar no botão **"+ Nova Instância"**
2. Preencher:
   - Nome: `Teste Comercial`
   - Número Destino: `+5511777777777` (opcional)
3. Clicar em "Criar"

**Resultado Esperado:**
- ✅ Modal fecha
- ✅ Nova instância aparece no grid
- ✅ Status inicial: "Desconectado"
- ✅ Grid mantém 3 colunas

---

## Teste 7: Verificar Envio Diário

### 7.1 Verificar Agendamento
```bash
# Ver se cron está configurado
grep -i "cron\|daily" /home/allansevero/wuzapi/wuzapi.log
```

**Resultado Esperado:**
```
Daily message sender cron job initialized
```

### 7.2 Trigger Manual (Se Implementado)
Se houver endpoint de envio manual:
```bash
curl -X POST http://localhost:8080/api/trigger-daily-send \
  -H "Authorization: Bearer SEU_TOKEN_ADMIN"
```

### 7.3 Verificar Webhook
1. Acessar n8n: `https://n8n-webhook.fmy2un.easypanel.host`
2. Verificar se webhook foi chamado
3. Verificar payload recebido

**Payload Esperado:**
```json
{
  "instanceName": "Instância Padrão",
  "destination_number": "+5511999999999",
  "date": "2025-11-04",
  "conversations": [
    {
      "chat": "+5511888888888",
      "messages": [
        {
          "from": "...",
          "text": "...",
          "timestamp": "..."
        }
      ]
    }
  ]
}
```

---

## Teste 8: Testar Status em Tempo Real

### 8.1 Enviar Mensagem
1. Com WhatsApp conectado, enviar mensagem para qualquer contato
2. Aguardar processamento

**Verificar Logs:**
```bash
tail -f /home/allansevero/wuzapi/wuzapi.log | grep -i message
```

**Resultado Esperado:**
- ✅ Mensagem é capturada pelo sistema
- ✅ Logs mostram processamento
- ✅ Mensagem é armazenada para envio diário

### 8.2 Verificar Armazenamento
```bash
sqlite3 /home/allansevero/wuzapi/dbdata/users.db \
  "SELECT COUNT(*) FROM message_history;"
```

---

## Teste 9: Desconectar e Deletar

### 9.1 Desconectar
1. Clicar em **"Desconectar"**
2. Confirmar

**Resultado Esperado:**
- ✅ Status muda para "Desconectado"
- ✅ JID desaparece
- ✅ Botões de conexão voltam a aparecer

### 9.2 Deletar Instância
1. Clicar em **"Deletar Instância"**
2. Confirmar

**Resultado Esperado:**
- ✅ Instância removida do grid
- ✅ Grid reorganiza automaticamente
- ✅ Se for última instância, mostra mensagem "Nenhuma instância"

---

## Teste 10: Responsividade

### 10.1 Desktop
- ✅ Grid com 3 colunas
- ✅ Cards com mesmo tamanho
- ✅ Botões visíveis

### 10.2 Tablet (redimensionar navegador)
- ✅ Grid adapta para 2 colunas
- ✅ Layout não quebra

### 10.3 Mobile
- ✅ Grid adapta para 1 coluna
- ✅ Botões empilhados
- ✅ Texto legível

---

## 🐛 Problemas Conhecidos e Soluções

### Problema: QR Code não aparece
**Solução:**
1. Verificar console do navegador
2. Verificar logs: `tail -f wuzapi.log | grep QR`
3. Testar endpoint diretamente: `curl http://localhost:8080/session/qr?token=...`

### Problema: Status não atualiza após conexão
**Solução:**
1. Aguardar 15 segundos (intervalo de refresh)
2. Recarregar página manualmente (F5)
3. Verificar se polling parou: ver console do navegador

### Problema: Database locked
**Solução:**
```bash
# Reiniciar servidor
sudo lsof -ti:8080 | xargs sudo kill -9
./wuzapi
```

### Problema: Porta 8080 ocupada
**Solução:**
```bash
sudo lsof -ti:8080 | xargs sudo kill -9
```

---

## ✅ Checklist Final

Após todos os testes:

- [ ] Login funciona
- [ ] Cadastro funciona
- [ ] QR Code é exibido
- [ ] Conexão via QR funciona
- [ ] Status atualiza (mesmo que demore)
- [ ] Código de pareamento funciona
- [ ] Número de destino pode ser configurado
- [ ] Nova instância pode ser criada
- [ ] Mensagens são capturadas
- [ ] Cron job está ativo
- [ ] Desconectar funciona
- [ ] Deletar funciona
- [ ] Interface é responsiva

---

## 📊 Resultados Esperados

### Todos os Testes Passam (✅)
**Sistema está pronto para uso!**

### Alguns Testes Falham (⚠️)
**Anotar quais falharam e reportar:**
1. Nome do teste
2. Resultado obtido vs esperado
3. Erros no console
4. Logs do backend

### Muitos Testes Falham (🔴)
**Verificar:**
1. Servidor está rodando?
2. Banco de dados está acessível?
3. Houve algum erro na compilação?

---

## 🔧 Ferramentas de Debug

### Console do Navegador
```javascript
// Ver todas as instâncias
console.log(instances);

// Ver intervalos de polling ativos
console.log(qrPollingIntervals);

// Ver token de autenticação
console.log(localStorage.getItem('auth_token'));
```

### SQLite
```bash
# Ver usuários
sqlite3 /home/allansevero/wuzapi/dbdata/users.db \
  "SELECT id, email, name FROM users;"

# Ver mensagens
sqlite3 /home/allansevero/wuzapi/dbdata/users.db \
  "SELECT * FROM message_history ORDER BY timestamp DESC LIMIT 10;"
```

### Logs
```bash
# Ver erros
grep ERROR /home/allansevero/wuzapi/wuzapi.log

# Ver warnings
grep WARN /home/allansevero/wuzapi/wuzapi.log

# Ver tudo em tempo real
tail -f /home/allansevero/wuzapi/wuzapi.log
```

---

**Documento criado em:** 2025-11-04  
**Última atualização:** 2025-11-04 07:35 BRT  
**Versão do Sistema:** 1.0.4
