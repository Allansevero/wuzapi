# WuzAPI

WuzAPI é uma API WebSocket para integração com o WhatsApp, com dashboard de gerenciamento de usuários e instâncias.

## Recursos

- API WebSocket para integração com o WhatsApp
- Dashboard completo com múltiplas versões (v2, v3, v4)
- Sistema de autenticação e gerenciamento de usuários
- Sistema de planos e assinaturas
- Envio diário de mensagens
- Webhooks para eventos do WhatsApp
- Histórico de mensagens
- Interface amigável para gerenciamento

## Deploy no EasyPanel

### Opção 1: Usando SQLite (padrão, não requer banco de dados externo)

1. Crie um novo projeto no EasyPanel como "Docker Compose"
2. Cole o conteúdo do arquivo `docker-compose.yml`
3. Defina as variáveis de ambiente (excluindo as variáveis de banco de dados)
4. Mapeie a porta 8080 para 80 (HTTP) ou 443 (HTTPS)
5. Configure volumes para persistência: `./dbdata` → `/root/dbdata`

### Opção 2: Usando PostgreSQL

1. Crie um novo projeto no EasyPanel como "Docker Compose"
2. Crie um banco de dados PostgreSQL separado no EasyPanel
3. Cole o conteúdo do arquivo `docker-compose.yml`
4. Defina todas as variáveis de ambiente, incluindo as do banco de dados PostgreSQL
5. Mapeie a porta 8080 para 80 (HTTP) ou 443 (HTTPS)
6. O volume de dados será gerenciado pelo PostgreSQL

## Variáveis de Ambiente

- `WUZAPI_ADMIN_TOKEN`: Token de administração (gerado automaticamente se não definido)
- `WUZAPI_GLOBAL_ENCRYPTION_KEY`: Chave de criptografia (32 bytes, obrigatória para persistência)
- `WUZAPI_GLOBAL_HMAC_KEY`: Chave HMAC para assinatura de webhooks (mínimo 32 chars)
- `WUZAPI_GLOBAL_WEBHOOK`: URL global para receber eventos do WhatsApp
- `WEBHOOK_FORMAT`: "json" ou "form" (padrão: json)
- `SESSION_DEVICE_NAME`: Nome do dispositivo no WhatsApp (padrão: WuzAPI)
- `TZ`: Fuso horário (padrão: America/Sao_Paulo)

### Variáveis de Banco de Dados (somente para PostgreSQL)
- `DB_USER`: Usuário do banco de dados
- `DB_PASSWORD`: Senha do banco de dados
- `DB_NAME`: Nome do banco de dados
- `DB_HOST`: Host do banco de dados
- `DB_PORT`: Porta do banco de dados
- `DB_SSLMODE`: Modo SSL (padrão: false)

## API Endpoints

### Endpoints Públicos
- `GET /` - Página inicial
- `GET /api/ping` - Verificar status da API

### Endpoints de Admin
- `GET /api/admin/list-users` - Listar todos os usuários
- `POST /api/admin/create-user` - Criar um novo usuário
- `DELETE /api/admin/delete-user` - Excluir um usuário
- `GET /api/admin/get-user` - Obter informações de um usuário

### Endpoints de Usuário
- `GET /api/user/login` - Login do usuário
- `POST /api/user/check-auth` - Verificar autenticação do usuário
- `POST /api/user/register` - Registrar um novo usuário

### Endpoints do WhatsApp
- `POST /api/whatsapp/send-message` - Enviar mensagem
- `POST /api/whatsapp/send-image` - Enviar imagem
- `POST /api/whatsapp/send-file` - Enviar arquivo
- `POST /api/whatsapp/send-sticker` - Enviar sticker
- `POST /api/whatsapp/send-voice` - Enviar áudio
- `GET /api/whatsapp/get-history` - Obter histórico de mensagens
- `GET /api/whatsapp/get-qr` - Obter QR code para conexão
- `POST /api/whatsapp/logout` - Desconectar do WhatsApp

## Dashboard

O dashboard está disponível em `/dashboard/` com as seguintes versões:
- Dashboard principal: `/dashboard/`
- Dashboard do usuário: `/dashboard/user-dashboard.html`
- Dashboard V2: `/dashboard/user-dashboard-v2.html`
- Dashboard V3: `/dashboard/user-dashboard-v3.html`
- Dashboard V4: `/dashboard/user-dashboard-v4.html`

## Segurança

- O sistema implementa proteções contra SSRF (Server-Side Request Forgery)
- Dados sensíveis são criptografados usando a chave `WUZAPI_GLOBAL_ENCRYPTION_KEY`
- Autenticação é exigida para endpoints críticos
- Webhooks podem ser assinados usando HMAC para verificação

## Observações

- Guarde com segurança a chave de criptografia global, pois ela é usada para criptografar dados sensíveis
- Faça backup regular do diretório `dbdata` para preservar dados dos usuários e instâncias
- Ao usar PostgreSQL, os backups devem ser feitos no banco de dados externo