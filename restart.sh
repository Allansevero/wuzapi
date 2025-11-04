#!/bin/bash

echo "🔄 Parando processo wuzapi..."
pkill -f ./wuzapi

echo "🔨 Compilando..."
go build -o wuzapi

if [ $? -eq 0 ]; then
    echo "✅ Compilação bem-sucedida!"
    echo "🚀 Iniciando wuzapi..."
    ./wuzapi &
    echo "✅ Wuzapi iniciado em background"
    echo "📝 Para ver logs: tail -f wuzapi.log"
else
    echo "❌ Erro na compilação"
    exit 1
fi
