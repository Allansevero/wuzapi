# 🚀 DEPLOY RÁPIDO - Sistema de Planos

## ⚡ Comandos para Deploy em Produção

### 1. Preparar Ambiente
```bash
cd /home/allansevero/wuzapi

# Fazer backup completo
tar -czf backup_$(date +%Y%m%d_%H%M%S).tar.gz \
  wuzapi \
  dbdata/ \
  static/ \
  *.go \
  go.mod \
  go.sum

# Parar serviço atual
sudo systemctl stop wuzapi
# OU
pkill wuzapi
```

### 2. Compilar Nova Versão
```bash
# Compilar com otimizações
go build -ldflags="-s -w" -o wuzapi_v2

# Verificar tamanho
ls -lh wuzapi_v2

# Dar permissões
chmod +x wuzapi_v2
```

### 3. Migrar Banco de Dados
```bash
# Backup do banco
cp -r dbdata/ dbdata.backup.$(date +%Y%m%d_%H%M%S)/

# As migrations rodam automaticamente ao iniciar
# Mas você pode verificar manualmente:
sqlite3 dbdata/users.db "SELECT * FROM migrations ORDER BY id DESC LIMIT 5;"

# Deve mostrar migration ID 13 (add_subscription_plans)
```

### 4. Substituir Binário
```bash
# Renomear atual
mv wuzapi wuzapi.old

# Mover nova versão
mv wuzapi_v2 wuzapi

# Verificar
./wuzapi --help
```

### 5. Iniciar Serviço
```bash
# Se usando systemd
sudo systemctl start wuzapi
sudo systemctl status wuzapi

# OU manual
nohup ./wuzapi > wuzapi.log 2>&1 &

# Verificar logs
tail -f wuzapi.log
```

### 6. Verificar Migrations
```bash
# Conectar ao banco
sqlite3 dbdata/users.db

# Verificar tabelas criadas
.tables
# Deve mostrar: plans, user_subscriptions, subscription_history

# Verificar planos inseridos
SELECT * FROM plans;
# Deve mostrar 3 planos

# Sair
.exit
```

### 7. Teste Rápido
```bash
# Health check
curl http://localhost:8080/health

# Registrar usuário teste
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@wuzapi.com","password":"Admin@123456"}'

# Fazer login
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@wuzapi.com","password":"Admin@123456"}' | \
  jq -r '.data.token')

echo $TOKEN

# Ver subscription
curl -s http://localhost:8080/my/subscription \
  -H "Authorization: Bearer $TOKEN" | jq

# Ver planos
curl -s http://localhost:8080/my/plans \
  -H "Authorization: Bearer $TOKEN" | jq
```

---

## 🔄 Rollback (Se necessário)

```bash
# Parar nova versão
sudo systemctl stop wuzapi
# OU
pkill wuzapi

# Restaurar binário antigo
mv wuzapi wuzapi.failed
mv wuzapi.old wuzapi

# Restaurar banco (se necessário)
rm -rf dbdata/
cp -r dbdata.backup.YYYYMMDD_HHMMSS/ dbdata/

# Reiniciar
sudo systemctl start wuzapi
# OU
nohup ./wuzapi > wuzapi.log 2>&1 &
```

---

## 📊 Monitoramento Pós-Deploy

### 1. Logs em Tempo Real
```bash
# Logs do serviço
tail -f wuzapi.log

# Filtrar por subscription
tail -f wuzapi.log | grep -i subscription

# Filtrar por erros
tail -f wuzapi.log | grep -i error
```

### 2. Verificar Banco de Dados
```bash
# Estatísticas de uso
sqlite3 dbdata/users.db << EOF
-- Total de usuários
SELECT COUNT(*) as total_users FROM system_users;

-- Usuários por plano
SELECT 
  p.name,
  COUNT(us.id) as users_count
FROM plans p
LEFT JOIN user_subscriptions us ON p.id = us.plan_id AND us.is_active = 1
GROUP BY p.id;

-- Subscriptions que expiram em 3 dias
SELECT 
  su.email,
  p.name,
  us.expires_at
FROM user_subscriptions us
JOIN system_users su ON us.system_user_id = su.id
JOIN plans p ON us.plan_id = p.id
WHERE us.is_active = 1
  AND us.expires_at IS NOT NULL
  AND datetime(us.expires_at) <= datetime('now', '+3 days');

-- Total de instâncias por usuário
SELECT 
  su.email,
  COUNT(u.id) as instance_count,
  p.name as plan_name,
  p.max_instances
FROM system_users su
LEFT JOIN users u ON u.system_user_id = su.id
LEFT JOIN user_subscriptions us ON us.system_user_id = su.id AND us.is_active = 1
LEFT JOIN plans p ON us.plan_id = p.id
GROUP BY su.id;
EOF
```

### 3. Endpoints de Saúde
```bash
# Health check geral
curl -s http://localhost:8080/health | jq

# Deve retornar:
{
  "status": "ok",
  "active_connections": N,
  "total_users": N,
  "connected_users": N
}
```

---

## 🔐 Segurança Pós-Deploy

### 1. Configurar JWT Secret
```bash
# Gerar secret aleatório
openssl rand -hex 32

# Editar auth.go (linha 44)
# Substituir: var jwtSecret = []byte("...")
# Por secret gerado acima

# Recompilar
go build -o wuzapi

# Reiniciar
sudo systemctl restart wuzapi
```

### 2. Configurar HTTPS
```bash
# Se usando Nginx como proxy reverso
sudo nano /etc/nginx/sites-available/wuzapi

# Adicionar:
server {
    listen 443 ssl http2;
    server_name seu-dominio.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# Recarregar Nginx
sudo systemctl reload nginx
```

### 3. Firewall
```bash
# Permitir apenas porta necessária
sudo ufw allow 8080/tcp

# Se usando Nginx
sudo ufw allow 'Nginx Full'

# Verificar
sudo ufw status
```

---

## 📧 Configurar Notificações (Opcional)

### Email para Expiração
```bash
# Criar script de verificação
cat > /home/allansevero/wuzapi/check_expiring.sh << 'EOF'
#!/bin/bash
sqlite3 /home/allansevero/wuzapi/dbdata/users.db << SQL
SELECT email FROM system_users su
JOIN user_subscriptions us ON su.id = us.system_user_id
WHERE us.is_active = 1
  AND datetime(us.expires_at) <= datetime('now', '+3 days')
  AND datetime(us.expires_at) > datetime('now')
SQL
EOF

chmod +x /home/allansevero/wuzapi/check_expiring.sh

# Agendar verificação diária (crontab)
crontab -e
# Adicionar:
# 0 9 * * * /home/allansevero/wuzapi/check_expiring.sh | mail -s "Subscriptions Expiring" admin@wuzapi.com
```

---

## 📝 Checklist Final de Deploy

### Pré-Deploy
- [ ] Backup completo realizado
- [ ] Código compilado sem erros
- [ ] Migrations testadas
- [ ] JWT secret configurado

### Durante Deploy
- [ ] Serviço parado
- [ ] Banco backup feito
- [ ] Binário substituído
- [ ] Serviço reiniciado

### Pós-Deploy
- [ ] Health check passou
- [ ] Migrations executadas (ID 13 presente)
- [ ] 3 planos inseridos
- [ ] Teste de registro funcionou
- [ ] Teste de login funcionou
- [ ] API de planos respondendo
- [ ] Interface web carregando
- [ ] Logs sem erros críticos

### Segurança
- [ ] JWT secret alterado
- [ ] HTTPS configurado
- [ ] Firewall configurado
- [ ] Backups agendados

---

## 🎯 Métricas de Sucesso

Após 24h de deploy, verificar:
- [ ] Novos registros criando subscriptions
- [ ] Nenhum erro relacionado a planos nos logs
- [ ] Interface de planos acessível
- [ ] Validações de limite funcionando
- [ ] Zero downtime

---

## 📞 Suporte

Em caso de problemas:

1. **Verificar logs:** `tail -f wuzapi.log`
2. **Verificar banco:** `sqlite3 dbdata/users.db`
3. **Health check:** `curl http://localhost:8080/health`
4. **Rollback se necessário** (comandos acima)

---

## ✅ Deploy Concluído!

Se todos os checkboxes acima estão marcados, o deploy foi **sucesso**!

**Sistema de Planos está em PRODUÇÃO! 🎉**
