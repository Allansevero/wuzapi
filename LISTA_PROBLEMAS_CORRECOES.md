# Lista de Problemas e Correções Necessárias

## Problemas Identificados

### 1. QR Code não aparece no frontend
**Descrição**: O QR code está sendo gerado no backend mas não aparece no frontend.
**Causa**: O frontend está buscando o QR code em `qrJson.data.QRCode` mas o backend pode estar retornando em formato diferente.
**Status dos logs**: 
- Backend gera QR code com sucesso
- Frontend recebe resposta 200 mas reporta "No QR code in response"
- Erro 500 intermitente: "database is locked", "not connected", "no session"

### 2. Status de conexão não atualiza no frontend
**Descrição**: Mesmo após conectar o WhatsApp escaneando o QR code, o dashboard continua mostrando status "Desconectado".
**Causa**: O polling de status não está interpretando corretamente a resposta do endpoint `/session/status`.
**Logs relevantes**:
```
Connected: undefined
LoggedIn: undefined
JID: undefined
```

### 3. Erro "already logged in" após conexão
**Descrição**: Após conectar com sucesso, o sistema continua tentando buscar QR code e retorna erro 500 "already logged in".
**Causa**: O polling de QR code não para adequadamente após conexão bem-sucedida.

### 4. Erros intermitentes de banco de dados
**Descrição**: Erro "database is locked (5) (SQLITE_BUSY)" durante operações.
**Causa**: Múltiplas requisições simultâneas ao SQLite sem tratamento adequado de concorrência.

## Correções Necessárias

### Frontend (user-dashboard-v2.js)

1. **Corrigir interpretação da resposta do QR code**
   - Verificar estrutura exata da resposta do backend
   - Ajustar caminho de acesso ao QR code (pode ser `data.qr` ao invés de `data.QRCode`)

2. **Corrigir interpretação do status de conexão**
   - Verificar campos corretos retornados pelo `/session/status`
   - Ajustar verificação de `statusData.connected`, `statusData.loggedIn`, `statusData.jid`

3. **Melhorar lógica de parada do polling**
   - Parar polling de QR code imediatamente após detectar `loggedIn: true`
   - Adicionar tratamento para erro "already logged in"

4. **Adicionar tratamento de erros**
   - Tratar erro 500 adequadamente
   - Mostrar mensagens de erro mais descritivas ao usuário

### Backend (handlers.go)

1. **Padronizar resposta do endpoint QR**
   - Garantir que o QR code sempre retorne no mesmo campo
   - Retornar estrutura consistente: `{code: 200, data: {qr: "data:image/png;base64,..."}, success: true}`

2. **Melhorar endpoint de status**
   - Retornar campos claros: `connected`, `loggedIn`, `jid`
   - Incluir informação se tem QR code disponível

3. **Tratamento de concorrência SQLite**
   - Implementar retry logic para "database is locked"
   - Considerar usar WAL mode no SQLite
   - Adicionar timeouts adequados

4. **Evitar geração de QR quando já conectado**
   - Verificar status antes de gerar QR code
   - Retornar erro apropriado com código 200 (não 500) quando já logado

## Alterações da Lista Original

### ✅ Já Implementadas

1. **Autenticação por usuário** - Cada usuário tem email e senha
2. **Isolamento de instâncias** - Usuário vê apenas suas instâncias
3. **Token automático** - Gerado no cadastro/login
4. **Webhook padrão** - Configurado para envio diário às 18h

### 🔨 Pendentes/Correções

1. **Funcionalidade de conexão do WhatsApp** - QUEBRADA, precisa correção
2. **Atualização de status em tempo real** - NÃO FUNCIONA
3. **Interface de instâncias** - Melhorar layout (3 colunas, bordas arredondadas)
4. **Botão para número de recebimento** - Implementar popup

## Prioridade de Correção

### Prioridade ALTA (Bloqueia uso)
1. ✅ Corrigir exibição do QR code
2. ✅ Corrigir atualização de status após conexão
3. ✅ Parar polling adequadamente após conexão

### Prioridade MÉDIA (Melhoria de UX)
4. Melhorar tratamento de erros
5. Otimizar concorrência do banco
6. Melhorar layout das instâncias (grid 3 colunas)

### Prioridade BAIXA (Features adicionais)
7. Implementar popup para número de recebimento
8. Melhorar feedback visual durante conexão

## Próximos Passos

1. Inspecionar resposta exata do endpoint `/session/qr` no backend
2. Inspecionar resposta exata do endpoint `/session/status` no backend
3. Ajustar frontend para interpretar respostas corretamente
4. Testar fluxo completo de conexão
5. Implementar melhorias de layout
