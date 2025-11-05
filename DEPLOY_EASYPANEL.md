
## Deploy no Easypanel

Para garantir funcionamento perfeito no Easypanel, siga estas configurações:

### 1. Variáveis de Ambiente Recomendadas:

```env
WUZAPI_ADMIN_TOKEN=seu_token_admin_aqui
WUZAPI_GLOBAL_ENCRYPTION_KEY=sua_chave_aes_256_de_32_bytes
WUZAPI_GLOBAL_HMAC_KEY=sua_chave_hmac_com_minimo_32_caracteres
WUZAPI_GLOBAL_WEBHOOK=sua_url_global_de_webhook
TZ=America/Sao_Paulo
SESSION_DEVICE_NAME=Metrizap
WEBHOOK_FORMAT=json
```

### 2. Configurações de Volume:
- ./dbdata:/root/dbdata (persistência do banco de dados)
- ./files:/root/files (arquivos de mídia)

### 3. Configurações de Segurança:
- Garanta que o diretório dbdata tenha permissões adequadas
- Não compartilhe tokens de produção
- Use chaves fortes para criptografia

### 4. Acesso ao Sistema:
- Cadastro: /user-login.html
- Dashboard: /dashboard/user-dashboard-v4.html
- API: /auth/register, /auth/login

