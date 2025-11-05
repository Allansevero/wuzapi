# ✅ Checklist para Produção - WuzAPI

## 🔧 Correções Já Aplicadas (04/11/2025)

- [x] Erro SQLITE_BUSY corrigido
- [x] QR Code exibindo corretamente
- [x] Compilação sem erros
- [x] Documentação criada

---

## ⚡ Testes Imediatos (Fazer AGORA)

```bash
# 1. Recompilar (se ainda não fez)
cd /home/allansevero/wuzapi
go build -o wuzapi

# 2. Reiniciar serviço
sudo systemctl restart wuzapi
# OU se rodando manual:
# killall wuzapi
# ./wuzapi &

# 3. Verificar logs
tail -f wuzapi.log
```

### Teste no Navegador:
- [ ] Abrir http://localhost:8080/dashboard/user-dashboard-v2.html
- [ ] Fazer login
- [ ] Clicar "Conectar WhatsApp"
- [ ] **Verificar**: QR Code aparece? ✅ DEVE APARECER
- [ ] Escanear QR Code com WhatsApp
- [ ] **Verificar**: Conectou? (pode demorar alguns segundos)
- [ ] **Verificar**: Status mudou para "Conectado"?

---

## 🐛 Se QR Code NÃO Aparecer

```bash
# 1. Abrir console do navegador (F12)
# 2. Procurar por erros
# 3. Ver logs:
tail -f wuzapi.log | grep -E "QR|qr"

# 4. Verificar resposta do endpoint:
TOKEN="SEU_TOKEN_AQUI"
curl "http://localhost:8080/session/qr?token=$TOKEN"
```

---

## 🐛 Se Status NÃO Atualizar Após Conectar

### Solução Rápida:
```bash
nano static/dashboard/js/user-dashboard-v2.js
# Procurar linha ~275
# Mudar: setTimeout(() => loadInstances(), 1500);
# Para:  setTimeout(() => loadInstances(), 500);
# Salvar (Ctrl+O, Enter, Ctrl+X)
# Recarregar página no navegador (Ctrl+Shift+R)
```

---

## 📤 Teste de Envio para Webhook

```bash
# 1. Pegar token válido
TOKEN=$(sqlite3 dbdata/users.db "SELECT token FROM users LIMIT 1;")
echo "Token: $TOKEN"

# 2. Enviar teste manual
./test_webhook_send.sh $TOKEN

# 3. Verificar logs
tail -f wuzapi.log | grep webhook

# 4. Verificar resposta do N8N
# (deve aparecer no seu workflow do N8N)
```

---

## 🕐 Ativar Envio Diário Automático às 18h

### Verificar se já está ativo:
```bash
grep -A 10 "cron" main.go
```

### Se NÃO estiver, adicionar:
```go
// Em main.go, na função main()
// Depois de inicializar o server:

import (
    "github.com/robfig/cron/v3"
)

// Adicionar antes de router.Run():
c := cron.New(cron.WithLocation(time.FixedZone("BRT", -3*60*60)))
c.AddFunc("0 20 * * *", func() {
    log.Info().Msg("Starting daily message sender...")
    server.sendDailyMessagesToWebhook()
})
c.Start()
log.Info().Msg("Daily sender scheduler started (18:00 BRT)")
```

### Recompilar:
```bash
go build -o wuzapi
sudo systemctl restart wuzapi
```

---

## 📊 Validação de Histórico de Mensagens

```bash
# 1. Conectar uma nova instância
# 2. Esperar 10 segundos após conexão
# 3. Verificar banco de dados:

sqlite3 dbdata/users.db
SELECT COUNT(*) FROM message_history;
SELECT * FROM message_history LIMIT 5;
```

**Esperado**: Deve ter mensagens antigas (histórico puxado)

---

## 🔒 Segurança (Antes de Deploy)

- [ ] Mudar senhas padrão
- [ ] Configurar HTTPS (certificado SSL)
- [ ] Configurar firewall (apenas portas necessárias)
- [ ] Backup do banco de dados
- [ ] Logs em local seguro
- [ ] Monitoramento ativo

---

## 🚀 Deploy em Produção

### 1. Preparação
```bash
# Backup
sudo systemctl stop wuzapi
cp -r dbdata dbdata.backup.$(date +%Y%m%d)
cp wuzapi wuzapi.backup

# Build de produção
go build -ldflags="-s -w" -o wuzapi
```

### 2. Deploy
```bash
# Copiar para servidor
scp wuzapi user@servidor:/path/to/wuzapi/

# No servidor:
sudo systemctl restart wuzapi
sudo systemctl status wuzapi
```

### 3. Validação
```bash
# Verificar se está rodando
curl http://localhost:8080/health

# Ver logs
tail -f /path/to/wuzapi.log
```

---

## 📋 Checklist Final

### Funcionalidades
- [ ] Login funciona
- [ ] Dashboard carrega instâncias
- [ ] QR Code aparece
- [ ] Conexão via QR funciona
- [ ] Conexão via código funciona
- [ ] Status atualiza (pode demorar alguns segundos)
- [ ] Mensagens são recebidas
- [ ] Mensagens são armazenadas no banco
- [ ] Webhook recebe dados (teste manual)

### Performance
- [ ] Sem erros de database locked
- [ ] Sem memory leaks (verificar com `top`)
- [ ] Logs sem erros críticos
- [ ] Tempo de resposta < 2s

### Segurança
- [ ] HTTPS configurado (produção)
- [ ] Firewall ativo
- [ ] Backups automáticos
- [ ] Senhas fortes

---

## 🆘 Em Caso de Problemas

### Logs
```bash
tail -f wuzapi.log
journalctl -u wuzapi -f
```

### Banco de Dados
```bash
# Verificar integridade
sqlite3 dbdata/users.db "PRAGMA integrity_check;"

# Otimizar
sqlite3 dbdata/users.db "VACUUM;"
```

### Reiniciar do Zero
```bash
sudo systemctl stop wuzapi
mv dbdata dbdata.old
./wuzapi  # Criará novo banco
```

---

## 📞 Suporte

**Logs**: `/var/log/wuzapi.log` ou `./wuzapi.log`  
**Banco**: `dbdata/users.db`  
**Porta**: `8080` (padrão)

---

**Última atualização**: 04/11/2025  
**Próxima revisão**: Após testes em produção
