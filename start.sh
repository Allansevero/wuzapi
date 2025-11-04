#!/bin/bash

# Script de Inicialização Rápida do Wuzapi
# Atualizado em: 04/11/2025

set -e

echo "================================================"
echo "🚀 INICIANDO WUZAPI - SISTEMA COMPLETO"
echo "================================================"
echo ""

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar se já está rodando
if pgrep -f "wuzapi" > /dev/null; then
    echo -e "${YELLOW}⚠️  Wuzapi já está rodando!${NC}"
    echo "Matando processo anterior..."
    pkill -f wuzapi
    sleep 2
fi

# Compilar se necessário
if [ ! -f "./wuzapi" ] || [ "./main.go" -nt "./wuzapi" ]; then
    echo -e "${BLUE}📦 Compilando...${NC}"
    go build -o wuzapi .
    echo -e "${GREEN}✅ Compilação completa!${NC}"
    echo ""
fi

# Mostrar informações
echo -e "${GREEN}✅ Sistema Pronto!${NC}"
echo ""
echo "📋 Funcionalidades Ativas:"
echo "  ✅ Sistema de usuários com autenticação"
echo "  ✅ 3 planos (Gratuito, Pro, Analista)"
echo "  ✅ Envio diário automático às 18h"
echo "  ✅ Webhook fixo configurado"
echo "  ✅ Parâmetro enviar_para"
echo ""
echo "🌐 Endpoints Disponíveis:"
echo "  • http://localhost:8080 - Dashboard"
echo "  • http://localhost:8080/auth/register - Cadastro"
echo "  • http://localhost:8080/auth/login - Login"
echo "  • http://localhost:8080/my/plans - Ver planos"
echo "  • http://localhost:8080/my/subscription - Ver assinatura"
echo ""
echo "📚 Documentação:"
echo "  • README_IMPLEMENTACAO.md - Resumo rápido"
echo "  • LEIA_ISTO_PRIMEIRO_FINAL.md - Guia completo"
echo "  • GUIA_TESTE_SISTEMA_COMPLETO.md - Testes"
echo ""
echo -e "${BLUE}🕐 Cron Job: Envio diário às 18:00 (Brasília)${NC}"
echo "🔗 Webhook: https://n8n-webhook.fmy2un.easypanel.host/webhook/..."
echo ""
echo "================================================"
echo "🎯 INICIANDO SERVIDOR..."
echo "================================================"
echo ""

# Executar
./wuzapi

# Se chegou aqui, o servidor foi parado
echo ""
echo "================================================"
echo "⏹️  Servidor Encerrado"
echo "================================================"
