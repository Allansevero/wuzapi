// User Dashboard v2 - Fixed version
const authToken = localStorage.getItem('auth_token');
let instances = [];
let qrPollingIntervals = {};

if (!authToken) {
    window.location.href = '/user-login.html';
}

document.getElementById('user-email').textContent = localStorage.getItem('user_email');

// Show alert
function showAlert(message, type) {
    const alert = $('#alert');
    alert.removeClass('success error warning info');
    alert.addClass(type);
    alert.html(`<i class="close icon"></i><div class="header">${message}</div>`);
    alert.show();
    $('.message .close').on('click', function() {
        $(this).closest('.message').hide();
    });
    
    // Auto hide after 5 seconds
    setTimeout(() => alert.fadeOut(), 5000);
}

// Load instances
async function loadInstances() {
    try {
        const response = await fetch('/my/instances', {
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        const data = await response.json();

        if (data.success) {
            instances = data.data || [];
            
            console.log('=== LOADING INSTANCES ===');
            console.log('Total instances:', instances.length);
            
            // Get real-time status for each instance
            for (let instance of instances) {
                try {
                    const statusRes = await fetch(`/session/status?token=${instance.token}`);
                    const statusData = await statusRes.json();
                    
                    console.log('----------------------------');
                    console.log('Instance:', instance.name);
                    console.log('Status Response:', statusData);
                    console.log('data.connected:', statusData.data?.connected);
                    console.log('data.loggedIn:', statusData.data?.loggedIn);
                    console.log('data.jid:', statusData.data?.jid);
                    console.log('----------------------------');
                    
                    // Use the actual response fields from data object
                    instance.connected = statusData.data?.connected || false;
                    instance.loggedIn = statusData.data?.loggedIn || false;
                    instance.jid = statusData.data?.jid || '';
                    instance.qrcode = '';
                    
                    console.log('Instance after update:', {
                        name: instance.name,
                        connected: instance.connected,
                        loggedIn: instance.loggedIn,
                        jid: instance.jid
                    });
                } catch (e) {
                    console.log('Error getting status for instance', instance.name, e);
                    instance.connected = false;
                    instance.loggedIn = false;
                    instance.jid = '';
                }
            }
            
            console.log('=== RENDERING INSTANCES ===');
            renderInstances();
        } else {
            showAlert('Erro ao carregar instâncias', 'error');
        }
    } catch (error) {
        console.error('Error loading instances:', error);
    }
}

// Render instances
function renderInstances() {
    const container = $('#instances-container');
    const emptyState = $('#empty-state');
    const welcomeMsg = $('#welcome-message');

    if (!instances || instances.length === 0) {
        container.empty();
        emptyState.show();
        welcomeMsg.hide();
        return;
    }

    emptyState.hide();
    
    // Check if first login
    const isFirstLogin = instances.length === 1 && instances[0].name === 'Instância Padrão' && !instances[0].loggedIn;
    if (isFirstLogin) {
        welcomeMsg.show();
    } else {
        welcomeMsg.hide();
    }

    container.empty();
    
    instances.forEach(instance => {
        const card = createInstanceCard(instance, isFirstLogin);
        container.append(card);
    });
}

// Create instance card
function createInstanceCard(instance, isFirstLogin) {
    const connected = instance.connected || false;
    const loggedIn = instance.loggedIn || false;
    
    const column = $('<div class="column"></div>');
    const card = $(`
        <div class="instance-card ${isFirstLogin ? 'first-login' : ''}">
            <div class="instance-header">
                <div class="instance-name">${instance.name}</div>
                <div class="ui label ${connected && loggedIn ? 'green' : 'grey'}">
                    ${connected && loggedIn ? 'Conectado' : 'Desconectado'}
                </div>
            </div>
            
            <div class="instance-info">
                <div class="info-row">
                    <div class="info-label">ID:</div>
                    <div class="info-value">${instance.id.substring(0, 16)}...</div>
                </div>
                <div class="info-row">
                    <div class="info-label">Status:</div>
                    <div class="info-value">
                        ${loggedIn ? '<span style="color: #21ba45;"><i class="check circle icon"></i>Logado</span>' : '<span style="color: #db2828;"><i class="times circle icon"></i>Não logado</span>'}
                    </div>
                </div>
                ${loggedIn ? `
                <div class="info-row">
                    <div class="info-label">Número:</div>
                    <div class="info-value">${instance.jid ? instance.jid.split('@')[0] : 'N/A'}</div>
                </div>
                ` : ''}
                <div class="info-row">
                    <div class="info-label">Destino:</div>
                    <div class="info-value">${instance.destination_number || 'Não configurado'}</div>
                </div>
            </div>
            
            <div id="qr-section-${instance.id}" class="qr-section" style="display: none;">
                <div class="qr-placeholder">
                    <i class="huge qrcode icon"></i>
                    <p>Aguardando QR Code...</p>
                </div>
            </div>
            
            <div class="instance-actions">
                <div class="btn-group">
                    ${!loggedIn ? `
                        <button class="ui positive button btn-full" onclick="connectInstance('${instance.token}', '${instance.id}')">
                            <i class="plug icon"></i> Conectar WhatsApp
                        </button>
                        <button class="ui blue button btn-full" onclick="showPairingModal('${instance.token}')">
                            <i class="mobile icon"></i> Código de Pareamento
                        </button>
                    ` : ''}
                    
                    ${loggedIn ? `
                        <button class="ui red button" onclick="disconnectInstance('${instance.token}', '${instance.id}')">
                            <i class="sign-out icon"></i> Desconectar
                        </button>
                    ` : ''}
                    
                    <button class="ui button" onclick="openDestinationModal('${instance.token}', '${instance.destination_number || ''}')">
                        <i class="phone icon"></i> Config. Destino
                    </button>
                </div>
                
                ${instances.length > 1 || (isFirstLogin && loggedIn) ? `
                    <button class="ui negative button btn-full" onclick="deleteInstance('${instance.id}')">
                        <i class="trash icon"></i> Deletar Instância
                    </button>
                ` : ''}
            </div>
        </div>
    `);
    
    column.append(card);
    return column;
}

// Connect instance
async function connectInstance(token, instanceId) {
    try {
        // First check if already connected
        const statusResponse = await fetch(`/session/status?token=${token}`);
        const statusData = await statusResponse.json();
        
        if (statusData.loggedIn && statusData.jid) {
            showAlert('WhatsApp já está conectado!', 'info');
            await loadInstances();
            return;
        }
        
        const response = await fetch('/session/connect', {
            method: 'POST',
            headers: {
                'token': token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                Subscribe: ['Message'],
                Immediate: true
            })
        });

        const data = await response.json();
        
        if (data.success || data.code === 200) {
            showAlert('Conectando... QR Code aparecerá em breve.', 'info');
            $(`#qr-section-${instanceId}`).show();
            startQRPolling(token, instanceId);
        } else {
            showAlert(data.error || 'Erro ao conectar', 'error');
        }
    } catch (error) {
        console.error('Connect error:', error);
        showAlert('Erro de conexão', 'error');
    }
}

// Start QR polling
function startQRPolling(token, instanceId) {
    console.log('Starting QR polling for instance:', instanceId);
    
    // Clear any existing interval
    if (qrPollingIntervals[instanceId]) {
        clearInterval(qrPollingIntervals[instanceId]);
    }
    
    const pollQR = async () => {
        try {
            // Check status first
            const statusResponse = await fetch(`/session/status?token=${token}`);
            const statusData = await statusResponse.json();
            
            console.log('=== QR POLL STATUS CHECK ===');
            console.log('Instance ID:', instanceId);
            console.log('Status Data:', statusData);
            console.log('data.connected:', statusData.data?.connected);
            console.log('data.loggedIn:', statusData.data?.loggedIn);
            console.log('data.jid:', statusData.data?.jid);
            console.log('===========================');
            
            // Check if connected and logged in
            if (statusData.data?.loggedIn === true && statusData.data?.jid) {
                console.log('✓ WhatsApp connected successfully! JID:', statusData.data.jid);
                clearInterval(qrPollingIntervals[instanceId]);
                delete qrPollingIntervals[instanceId];
                
                // Hide QR section
                $(`#qr-section-${instanceId}`).hide();
                
                showAlert('WhatsApp conectado com sucesso!', 'success');
                
                // Reload instances after a short delay to update UI
                setTimeout(() => {
                    console.log('Reloading instances after successful connection...');
                    loadInstances();
                }, 1500);
                return;
            }
            
            // Get QR code only if not logged in
            const qrResponse = await fetch(`/session/qr?token=${token}`);
            console.log('QR Response status:', qrResponse.status);
            
            if (qrResponse.ok) {
                const qrJson = await qrResponse.json();
                console.log('QR JSON received:', qrJson);
                
                // Check if already logged in (new format)
                if (qrJson.data?.message === 'already logged in' || qrJson.message === 'already logged in') {
                    console.log('✓ Already logged in, stopping QR polling');
                    clearInterval(qrPollingIntervals[instanceId]);
                    delete qrPollingIntervals[instanceId];
                    $(`#qr-section-${instanceId}`).hide();
                    showAlert('WhatsApp já está conectado!', 'success');
                    setTimeout(() => loadInstances(), 1000);
                    return;
                }
                
                // The QR code can come in different formats
                const qrData = qrJson.data?.QRCode || qrJson.QRCode || qrJson.data || '';
                
                console.log('QR Data exists:', !!qrData);
                console.log('QR Data length:', qrData ? qrData.length : 0);
                
                if (qrData && qrData.length > 0 && qrData.startsWith('data:image')) {
                    console.log('✓ Valid QR code found, displaying...');
                    $(`#qr-section-${instanceId}`).html(`
                        <img src="${qrData}" alt="QR Code" style="max-width: 220px; max-height: 220px; border-radius: 8px;">
                        <p style="margin-top: 10px; font-size: 12px; color: #666;">
                            <i class="whatsapp icon"></i> Abra WhatsApp > Aparelhos conectados > Conectar um aparelho
                        </p>
                    `);
                } else {
                    console.log('✗ No QR code in response yet, waiting...');
                }
            } else {
                const errorJson = await qrResponse.json();
                console.log('QR Error:', errorJson);
                
                // If already logged in (old error format), stop polling and reload instances
                if (errorJson.error === 'already logged in') {
                    console.log('✓ Already logged in (error), stopping QR polling');
                    clearInterval(qrPollingIntervals[instanceId]);
                    delete qrPollingIntervals[instanceId];
                    $(`#qr-section-${instanceId}`).hide();
                    showAlert('WhatsApp já está conectado!', 'success');
                    setTimeout(() => loadInstances(), 1000);
                    return;
                }
                
                // If error is "no session", stop polling
                if (errorJson.error === 'no session') {
                    clearInterval(qrPollingIntervals[instanceId]);
                    delete qrPollingIntervals[instanceId];
                    showAlert('Sessão não encontrada. Tente reconectar.', 'error');
                    return;
                }
            }
        } catch (error) {
            console.error('QR polling error:', error);
        }
    };
    
    // Poll immediately and then every 2 seconds
    pollQR();
    qrPollingIntervals[instanceId] = setInterval(pollQR, 2000);
}

// Disconnect instance
async function disconnectInstance(token, instanceId) {
    if (!confirm('Tem certeza que deseja desconectar?')) {
        return;
    }

    try {
        const response = await fetch('/session/logout', {
            method: 'POST',
            headers: {
                'token': token
            }
        });

        const data = await response.json();
        
        if (data.success || data.code === 200) {
            showAlert('Desconectado com sucesso', 'success');
            
            // Clear QR polling if exists
            if (qrPollingIntervals[instanceId]) {
                clearInterval(qrPollingIntervals[instanceId]);
                delete qrPollingIntervals[instanceId];
            }
            
            setTimeout(() => loadInstances(), 1000);
        } else {
            showAlert(data.error || 'Erro ao desconectar', 'error');
        }
    } catch (error) {
        showAlert('Erro de conexão', 'error');
    }
}

// Show pairing modal
function showPairingModal(token) {
    $('#pairing-instance-token').val(token);
    $('#pairing-modal').modal('show');
}

// Request pairing code
async function requestPairingCode() {
    const token = $('#pairing-instance-token').val();
    const phone = $('#pairing-phone').val().replace(/[^0-9]/g, '');
    
    if (!phone) {
        showAlert('Digite o número de telefone', 'warning');
        return;
    }

    try {
        const response = await fetch('/session/pairphone', {
            method: 'POST',
            headers: {
                'token': token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ phone: phone })
        });

        const data = await response.json();
        
        if (data.success || data.code) {
            $('#pairing-modal').modal('hide');
            showAlert(`Código de pareamento: ${data.code || 'Verifique seu WhatsApp'}`, 'success');
            setTimeout(() => loadInstances(), 3000);
        } else {
            showAlert(data.error || 'Erro ao solicitar código', 'error');
        }
    } catch (error) {
        showAlert('Erro de conexão', 'error');
    }
}

// Show create modal
function showCreateModal() {
    $('#instance-name').val('');
    $('#destination-number').val('');
    $('#create-modal').modal('show');
}

// Create instance
async function createInstance() {
    const name = $('#instance-name').val();
    const destination_number = $('#destination-number').val();

    if (!name) {
        showAlert('Digite o nome da instância', 'warning');
        return;
    }

    try {
        const response = await fetch('/my/instances', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${authToken}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name, destination_number })
        });

        const data = await response.json();

        if (data.success) {
            $('#create-modal').modal('hide');
            showAlert('Instância criada com sucesso!', 'success');
            setTimeout(() => loadInstances(), 1000);
        } else {
            showAlert(data.error || 'Erro ao criar instância', 'error');
        }
    } catch (error) {
        showAlert('Erro de conexão', 'error');
    }
}

// Open destination modal
function openDestinationModal(token, currentNumber) {
    $('#modal-instance-token').val(token);
    $('#modal-destination-number').val(currentNumber);
    $('#destination-modal').modal('show');
}

// Save destination number
async function saveDestinationNumber() {
    const token = $('#modal-instance-token').val();
    const number = $('#modal-destination-number').val();

    try {
        const response = await fetch('/session/destination-number', {
            method: 'POST',
            headers: {
                'token': token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ number })
        });

        const data = await response.json();

        if (data.success) {
            $('#destination-modal').modal('hide');
            showAlert('Número configurado!', 'success');
            setTimeout(() => loadInstances(), 1000);
        } else {
            showAlert(data.error || 'Erro ao configurar', 'error');
        }
    } catch (error) {
        showAlert('Erro de conexão', 'error');
    }
}

// Delete instance
async function deleteInstance(instanceId) {
    if (!confirm('Tem certeza que deseja deletar esta instância? Esta ação não pode ser desfeita.')) {
        return;
    }

    try {
        const response = await fetch(`/my/instances/${instanceId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${authToken}`
            }
        });

        const data = await response.json();

        if (data.success) {
            showAlert('Instância deletada!', 'success');
            
            // Clear polling if exists
            if (qrPollingIntervals[instanceId]) {
                clearInterval(qrPollingIntervals[instanceId]);
                delete qrPollingIntervals[instanceId];
            }
            
            setTimeout(() => loadInstances(), 1000);
        } else {
            showAlert(data.error || 'Erro ao deletar', 'error');
        }
    } catch (error) {
        showAlert('Erro de conexão', 'error');
    }
}

// Logout
function logout() {
    // Clear all intervals
    Object.values(qrPollingIntervals).forEach(interval => clearInterval(interval));
    qrPollingIntervals = {};
    
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_email');
    window.location.href = '/user-login.html';
}

// Initialize
loadInstances();

// Refresh every 15 seconds
setInterval(() => {
    // Only reload if not actively polling QR
    if (Object.keys(qrPollingIntervals).length === 0) {
        loadInstances();
    }
}, 15000);

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    Object.values(qrPollingIntervals).forEach(interval => clearInterval(interval));
});
