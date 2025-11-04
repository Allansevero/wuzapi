// Função para teste manual de envio diário de mensagens
// Adicione esta função ao seu arquivo user-dashboard-v2.js

/**
 * Envia manualmente as mensagens compiladas do dia para o webhook
 * @param {string} instanceId - ID da instância (opcional, usa a atual se não fornecido)
 * @param {string} date - Data no formato YYYY-MM-DD (opcional, usa hoje se não fornecido)
 * @returns {Promise<object>} - Resultado do envio
 */
async function sendDailyTestManual(instanceId = null, date = null) {
    const token = getTokenFromURL();
    
    if (!token) {
        alert('Token não encontrado. Faça login novamente.');
        return;
    }
    
    let url = `/session/send-daily-test?token=${token}`;
    
    if (instanceId) {
        url += `&instance_id=${instanceId}`;
    }
    
    if (date) {
        url += `&date=${date}`;
    }
    
    try {
        console.log('Enviando teste de mensagens diárias...');
        console.log('URL:', url);
        
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        const data = await response.json();
        
        if (data.success) {
            console.log('✅ Mensagens enviadas com sucesso!');
            console.log('Instance ID:', data.instance_id);
            console.log('Data:', data.date);
            alert(`Mensagens do dia ${data.date} enviadas com sucesso para o webhook!`);
        } else {
            console.error('❌ Erro ao enviar mensagens:', data.error);
            alert(`Erro: ${data.error}`);
        }
        
        return data;
    } catch (error) {
        console.error('❌ Erro na requisição:', error);
        alert(`Erro ao enviar: ${error.message}`);
        throw error;
    }
}

/**
 * Testa envio para todas as instâncias do usuário
 */
async function sendDailyTestAllInstances() {
    const instances = document.querySelectorAll('.instance-card');
    
    if (instances.length === 0) {
        alert('Nenhuma instância encontrada.');
        return;
    }
    
    const confirmMsg = `Deseja enviar teste de mensagens diárias para todas as ${instances.length} instâncias?`;
    
    if (!confirm(confirmMsg)) {
        return;
    }
    
    let successCount = 0;
    let errorCount = 0;
    
    for (const instanceCard of instances) {
        const instanceId = instanceCard.dataset.instanceId;
        const instanceName = instanceCard.querySelector('.instance-name')?.textContent || instanceId;
        
        try {
            console.log(`Enviando para instância: ${instanceName} (${instanceId})`);
            const result = await sendDailyTestManual(instanceId);
            
            if (result.success) {
                successCount++;
            } else {
                errorCount++;
            }
            
            // Aguarda 1 segundo entre envios para não sobrecarregar
            await new Promise(resolve => setTimeout(resolve, 1000));
            
        } catch (error) {
            console.error(`Erro na instância ${instanceName}:`, error);
            errorCount++;
        }
    }
    
    alert(`Envio concluído!\n✅ Sucesso: ${successCount}\n❌ Erros: ${errorCount}`);
}

/**
 * Adiciona botão de teste ao dashboard
 */
function addDailyTestButton() {
    // Verifica se já existe
    if (document.getElementById('btn-daily-test')) {
        return;
    }
    
    // Cria botão
    const button = document.createElement('button');
    button.id = 'btn-daily-test';
    button.className = 'btn btn-primary';
    button.innerHTML = '📤 Testar Envio Diário';
    button.style.cssText = `
        position: fixed;
        bottom: 20px;
        right: 20px;
        padding: 12px 24px;
        background: #007bff;
        color: white;
        border: none;
        border-radius: 8px;
        cursor: pointer;
        font-size: 14px;
        font-weight: bold;
        box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        z-index: 1000;
        transition: all 0.3s ease;
    `;
    
    button.onmouseover = function() {
        this.style.background = '#0056b3';
        this.style.transform = 'scale(1.05)';
    };
    
    button.onmouseout = function() {
        this.style.background = '#007bff';
        this.style.transform = 'scale(1)';
    };
    
    button.onclick = function() {
        const menu = document.createElement('div');
        menu.style.cssText = `
            position: fixed;
            bottom: 80px;
            right: 20px;
            background: white;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            padding: 16px;
            min-width: 250px;
            z-index: 1001;
        `;
        
        menu.innerHTML = `
            <h4 style="margin: 0 0 12px 0; font-size: 16px;">Testar Envio Diário</h4>
            <button class="test-option" data-action="current">📱 Instância Atual (Hoje)</button>
            <button class="test-option" data-action="all">📱📱 Todas as Instâncias</button>
            <button class="test-option" data-action="custom">📅 Data Personalizada</button>
            <button class="test-option" data-action="cancel" style="background: #dc3545;">❌ Cancelar</button>
        `;
        
        const style = document.createElement('style');
        style.textContent = `
            .test-option {
                display: block;
                width: 100%;
                padding: 10px;
                margin: 6px 0;
                background: #f8f9fa;
                border: 1px solid #dee2e6;
                border-radius: 4px;
                cursor: pointer;
                text-align: left;
                font-size: 14px;
                transition: all 0.2s;
            }
            .test-option:hover {
                background: #e9ecef;
                transform: translateX(4px);
            }
        `;
        document.head.appendChild(style);
        
        menu.querySelectorAll('.test-option').forEach(btn => {
            btn.onclick = async function() {
                const action = this.dataset.action;
                document.body.removeChild(menu);
                
                if (action === 'current') {
                    await sendDailyTestManual();
                } else if (action === 'all') {
                    await sendDailyTestAllInstances();
                } else if (action === 'custom') {
                    const date = prompt('Digite a data (YYYY-MM-DD):', new Date().toISOString().split('T')[0]);
                    if (date) {
                        await sendDailyTestManual(null, date);
                    }
                }
            };
        });
        
        document.body.appendChild(menu);
        
        // Fecha ao clicar fora
        setTimeout(() => {
            document.addEventListener('click', function closeMenu(e) {
                if (!menu.contains(e.target) && e.target !== button) {
                    if (document.body.contains(menu)) {
                        document.body.removeChild(menu);
                    }
                    document.removeEventListener('click', closeMenu);
                }
            });
        }, 100);
    };
    
    document.body.appendChild(button);
}

// Inicializa quando o documento estiver pronto
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', addDailyTestButton);
} else {
    addDailyTestButton();
}

// Expor funções globalmente para uso no console
window.sendDailyTestManual = sendDailyTestManual;
window.sendDailyTestAllInstances = sendDailyTestAllInstances;

console.log('✅ Funções de teste de envio diário carregadas!');
console.log('Use no console:');
console.log('  - sendDailyTestManual() - Testa instância atual');
console.log('  - sendDailyTestManual("instance-id") - Testa instância específica');
console.log('  - sendDailyTestManual(null, "2025-11-03") - Testa com data específica');
console.log('  - sendDailyTestAllInstances() - Testa todas as instâncias');
