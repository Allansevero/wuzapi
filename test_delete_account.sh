#!/bin/bash

echo "Testando o endpoint DELETE para exclusão de conta..."

# Teste sem token (deve retornar erro de autenticação)
echo "1. Testando sem token (espera-se erro de autenticação)..."
response=$(curl -s -X DELETE http://localhost:8080/my/profile -w "\n%{http_code}" -o /tmp/response.txt)
status_code=$(echo "$response" | tail -n 1)
response_body=$(head -n -1 <<< "$response")

echo "Status Code: $status_code"
echo "Response Body: $response_body"
echo ""

# Se o servidor estiver rodando no localhost:8080, o teste abaixo verificará se o endpoint existe
if [ "$status_code" = "401" ]; then
    echo "✓ Verificação de autenticação funcionando corretamente"
else
    echo "✗ Erro na verificação de autenticação"
fi

echo "Teste do endpoint concluído."