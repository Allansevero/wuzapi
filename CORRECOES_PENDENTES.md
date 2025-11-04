# Correções Pendentes - Sistema WuzAPI

## Data: 2025-11-04

## Problemas Identificados e Correções Necessárias

### 1. **Status de Conexão não Atualiza no Frontend** ✅ CRÍTICO
**Problema:**
- Backend gera QR code e conecta corretamente
- Frontend não reflete a mudança de status de "Desconectado" para "Conectado"
- Erro 500 no console: `database is locked (5) (SQLITE_BUSY)`
- Mensagem de log: `already logged in` mas frontend continua mostrando desconectado

**Causa Raiz:**
- Resposta da API `/session/status` retorna dados dentro de um objeto `data`
- Frontend está lendo diretamente `statusData.connected` quando deveria ser `statusData.data.connected`
- Polling do QR code continua mesmo após conexão bem-sucedida
- Múltiplas requisições simultâneas causando lock no SQLite

**Correção Necessária:**
- ✅ Atualizar o frontend para ler corretamente `statusData.data.connected`, `statusData.data.loggedIn` e `statusData.data.jid`
- ✅ Garantir que o polling de QR code pare corretamente quando a instância conectar
- ✅ Adicionar verificação de status antes de mostrar QR code
- ⚠️ Implementar mecanismo de retry com backoff exponencial no frontend
- ⚠️ Adicionar pooling de conexões SQLite no backend para evitar locks

### 2. **Erro "Already Logged In" Durante QR Polling** ✅ PARCIALMENTE RESOLVIDO
**Problema:**
- Após conexão bem-sucedida, o sistema continua tentando buscar QR code
- Erro 500: `already logged in`
- Polling não para automaticamente

**Correção Necessária:**
- ✅ Detectar erro "already logged in" e parar o polling
- ✅ Limpar intervalo de polling quando status mudar para conectado
- ✅ Recarregar instâncias após detectar conexão

### 3. **Botão "Conectar com Código" Sumiu** ⚠️ PENDENTE
**Problema:**
- Interface anterior tinha opção de conectar via código de pareamento
- Botão desapareceu após alterações

**Correção Necessária:**
- ✅ Adicionar botão "Código de Pareamento" na interface
- ⚠️ Implementar modal para solicitar código de pareamento
- ⚠️ Implementar endpoint backend para gerar código de pareamento

### 4. **Layout das Instâncias** ⚠️ PENDENTE
**Problema:**
- Solicitado layout em grade de 3 colunas com bordas arredondadas
- Atualmente pode não estar otimizado

**Correção Necessária:**
- ✅ CSS já implementa grade de 3 colunas com bordas arredondadas
- ⚠️ Verificar responsividade em diferentes tamanhos de tela
- ⚠️ Ajustar espaçamentos e paddings se necessário

### 5. **Token Admin Automático** ⚠️ NÃO IMPLEMENTADO
**Problema:**
- Após cadastro/login, usuário ainda precisa lidar com tokens manualmente
- Experiência ruim ter que copiar token

**Correção Necessária:**
- ⚠️ Gerar automaticamente token admin para cada usuário no cadastro
- ⚠️ Armazenar token admin no banco associado ao usuário
- ⚠️ Fazer login automático após cadastro direcionando para dashboard
- ⚠️ Remover necessidade de exibir/copiar token na interface

### 6. **QR Code não Aparece no Frontend** ✅ CORRIGIDO
**Problema:**
- Backend gera QR code corretamente
- Frontend não exibe a imagem do QR code
- Erros de JavaScript no console

**Causa Raiz:**
- Frontend está buscando `qrJson.QRCode` mas o backend retorna dentro de `qrJson.data.QRCode`
- Formato da resposta não está sendo tratado corretamente

**Correção Aplicada:**
- ✅ Atualizar parsing da resposta do QR code no frontend
- ✅ Adicionar fallbacks para diferentes formatos de resposta
- ✅ Validar que QR code é uma imagem válida antes de exibir

### 7. **Erro de Conexão 400 em /session/connect** ⚠️ INVESTIGAR
**Problema:**
- Erro 400 (Bad Request) ao tentar conectar instância
- Algumas instâncias mostram "Conectado" mesmo sem estar

**Correção Necessária:**
- ⚠️ Verificar payload sendo enviado para `/session/connect`
- ⚠️ Adicionar validação de campos obrigatórios
- ⚠️ Melhorar mensagens de erro para debug

### 8. **SQLite Database Lock** 🔴 CRÍTICO
**Problema:**
- `database is locked (5) (SQLITE_BUSY)`
- Erro fatal ao criar sqlstore
- Múltiplos acessos simultâneos ao banco

**Correção Necessária:**
- 🔴 Implementar connection pooling para SQLite
- 🔴 Adicionar timeout e retry em operações de banco
- 🔴 Considerar migração para PostgreSQL para produção
- ⚠️ Adicionar `PRAGMA busy_timeout` no SQLite
- ⚠️ Garantir que conexões são fechadas corretamente

### 9. **Histórico de Mensagens** ⚠️ NÃO IMPLEMENTADO
**Solicitação:**
- Buscar últimas 100 mensagens por conversa ao fazer login
- Atualmente só armazena mensagens após login

**Correção Necessária:**
- ⚠️ Implementar busca de histórico no backend
- ⚠️ Configurar parâmetro de histórico ao conectar
- ⚠️ Armazenar mensagens históricas no banco
- ⚠️ Exibir histórico na interface (se necessário)

### 10. **Envio Manual de Compilado para Webhook** ⚠️ FUNCIONALIDADE NOVA
**Solicitação:**
- Criar endpoint para enviar manualmente compilado de mensagens
- Testar envio para webhook sem esperar agendamento

**Correção Necessária:**
- ⚠️ Criar endpoint `/api/send-daily-now` ou similar
- ⚠️ Reaproveitar lógica do daily sender
- ⚠️ Adicionar botão na interface para trigger manual
- ⚠️ Proteger com autenticação

## Prioridades de Implementação

### Alta Prioridade (URGENTE)
1. 🔴 **SQLite Database Lock** - Sistema pode travar
2. ✅ **Status de Conexão não Atualiza** - Experiência do usuário comprometida
3. ✅ **QR Code não Aparece** - Impossível conectar WhatsApp

### Média Prioridade
4. ⚠️ **Token Admin Automático** - Melhora significativa de UX
5. ⚠️ **Botão Conectar com Código** - Funcionalidade ausente
6. ⚠️ **Erro 400 em /session/connect** - Pode impedir conexões

### Baixa Prioridade
7. ⚠️ **Histórico de Mensagens** - Feature adicional
8. ⚠️ **Envio Manual para Webhook** - Ferramenta de teste
9. ⚠️ **Layout das Instâncias** - Ajustes cosméticos

## Tecnologia Frontend

**Stack Atual:** HTML puro + JavaScript vanilla + Semantic UI (Fomantic UI)
- Não utiliza React ou outro framework
- JavaScript modular em arquivos separados
- CSS customizado + framework Semantic UI

## Próximos Passos

1. ✅ Corrigir parsing de status no frontend
2. ✅ Corrigir exibição de QR code
3. 🔴 Implementar solução para SQLite locks
4. ⚠️ Implementar token admin automático
5. ⚠️ Adicionar botão e modal de código de pareamento
6. ⚠️ Testar e validar todas as correções

## Notas Técnicas

- Backend em Go com WhatsApp Web Multi-Device API (whatsmeow)
- Frontend em HTML/JS/CSS com Semantic UI
- Banco de dados SQLite (considerar PostgreSQL para produção)
- Sistema de autenticação JWT
- API RESTful
