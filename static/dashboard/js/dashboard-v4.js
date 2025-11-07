// Dashboard V4 - Main JavaScript
// Integração com o backend WuzAPI

(function() {
    'use strict';

    // Estado da aplicação
    const state = {
        user: null,
        instances: [],
        subscription: null,
        currentFilter: 'all',
        searchQuery: '',
        currentInstanceId: null
    };

    // Utilitários
    const API = {
        getToken: () => localStorage.getItem('token') || localStorage.getItem('auth_token'),
        
        async request(url, options = {}) {
            const token = this.getToken();
            if (!token) {
                window.location.href = '/login/';
                return;
            }

            const headers = {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
                ...options.headers
            };

            try {
                const response = await fetch(url, { ...options, headers });
                
                if (response.status === 401) {
                    localStorage.removeItem('token');
                    localStorage.removeItem('auth_token');
                    localStorage.removeItem('authToken');
                    window.location.href = '/login/';
                    return;
                }

                const result = await response.json();
                
                // WuzAPI retorna { code, data, success } ou { code, error, success }
                if (result.success === false) {
                    throw new Error(result.error || 'Request failed');
                }
                
                if (!response.ok) {
                    throw new Error(result.error || `HTTP error! status: ${response.status}`);
                }
                
                return result.data || result;
            } catch (error) {
                console.error('API request failed:', error);
                throw error;
            }
        },
        
        // Request using instance token instead of JWT
        async requestWithInstanceToken(url, instanceToken, options = {}) {
            const headers = {
                'token': instanceToken,
                'Content-Type': 'application/json',
                ...options.headers
            };

            try {
                const response = await fetch(url, { ...options, headers });
                
                const result = await response.json();
                
                if (result.success === false) {
                    throw new Error(result.error || 'Request failed');
                }
                
                if (!response.ok) {
                    throw new Error(result.error || `HTTP error! status: ${response.status}`);
                }
                
                return result.data || result;
            } catch (error) {
                console.error('API request with instance token failed:', error);
                throw error;
            }
        },

        // User endpoints
        getUser: () => API.request('/my/profile'),
        updateUser: (data) => API.request('/my/profile', { method: 'PUT', body: JSON.stringify(data) }),
        
        // Force profile refresh endpoint
        refreshProfile: () => API.request('/my/profile'),
        
        // Instance endpoints
        getInstances: () => API.request('/my/instances'),
        createInstance: (name) => API.request('/my/instances', { 
            method: 'POST', 
            body: JSON.stringify({ name }) 
        }),
        deleteInstance: (id) => API.request(`/my/instances/${id}`, { method: 'DELETE' }),
        
        // Instance operations (use instance token)
        connectInstance: (instanceToken) => API.requestWithInstanceToken('/session/connect', instanceToken, { 
            method: 'POST',
            body: JSON.stringify({ 
                Subscribe: ['Message'],
                Immediate: true 
            })
        }),
        disconnectInstance: (instanceToken) => API.requestWithInstanceToken('/session/logout', instanceToken, { method: 'POST' }),
        getQRCode: (instanceToken) => API.requestWithInstanceToken('/session/qr', instanceToken),
        
        // Subscription endpoints
        getSubscription: () => API.request('/my/subscription'),
        getPlans: () => API.request('/my/plans')
    };

    // Formatação de dados
    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';
        const date = new Date(dateString);
        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const year = date.getFullYear();
        return `${day}/${month}/${year}`;
    };

    const getFirstName = (fullName) => {
        if (!fullName) return 'Usuário';
        return fullName.split(' ')[0];
    };

    const getInitials = (name) => {
        if (!name) return 'U';
        const parts = name.split(' ');
        if (parts.length >= 2) {
            return (parts[0][0] + parts[1][0]).toUpperCase();
        }
        return name[0].toUpperCase();
    };

    // Renderização de componentes
    const UI = {
        updateUserInfo() {
            if (!state.user) return;

            const firstName = getFirstName(state.user.name);
            const initials = getInitials(state.user.name);

            document.getElementById('userGreeting').textContent = `Olá, ${firstName} 👋`;
            document.getElementById('userDataGreeting').textContent = `${firstName}, seus dados 👋`;
            document.getElementById('userAvatar').src = `https://placehold.co/48x48/E2E8F0/4A5568?text=${initials}`;
            
            // Preencher formulário de dados
            document.getElementById('nome').value = state.user.name || '';
            document.getElementById('email').value = state.user.email || '';
            document.getElementById('whatsapp').value = state.user.whatsapp_number || '';
        },

        updateInstancesProgress() {
            if (!state.subscription) return;

            const planName = state.subscription.plan?.name || '';
            const maxInstances = state.subscription.max_instances || state.subscription.plan?.max_instances || 0;
            const connectedCount = state.subscription.connected_count || 0;
            const remaining = state.subscription.instances_remaining ?? Math.max(0, maxInstances - connectedCount);
            const percentage = maxInstances > 0 ? (connectedCount / maxInstances) * 100 : 0;
            
            const progressText = document.getElementById('remainingInstances');
            const progressBar = document.getElementById('progressBar');

            // Mostrar WhatsApp conectados e restantes
            progressText.textContent = `${connectedCount} de ${maxInstances} WhatsApp conectado${connectedCount !== 1 ? 's' : ''} (${remaining} disponível${remaining !== 1 ? 'eis' : ''})`;
            progressBar.style.width = `${percentage}%`;
            progressBar.classList.add('bg-mz-green');
            progressBar.classList.remove('bg-mz-orange-dark', 'bg-mz-red');

            // Verificar se o plano expirou
            const isExpired = state.subscription.is_expired === true;
            const expiredAlert = document.getElementById('expiredPlanAlert');
            const addInstanceBtn = document.getElementById('addInstanceBtn');

            if (isExpired) {
                // Mostrar alerta
                expiredAlert?.classList.remove('hidden');
                
                // Desabilitar botão de adicionar
                if (addInstanceBtn) {
                    addInstanceBtn.disabled = true;
                    addInstanceBtn.classList.add('opacity-50', 'cursor-not-allowed');
                    addInstanceBtn.title = 'Seu plano expirou. Faça upgrade para continuar.';
                }

                // Link do alerta para ir para página de planos
                document.getElementById('upgradeLinkFromAlert')?.addEventListener('click', (e) => {
                    e.preventDefault();
                    document.getElementById('linkDados')?.click();
                    setTimeout(() => {
                        document.getElementById('plansContainer')?.scrollIntoView({ behavior: 'smooth' });
                    }, 100);
                });
            } else {
                // Esconder alerta
                expiredAlert?.classList.add('hidden');
                
                // Habilitar botão
                if (addInstanceBtn) {
                    addInstanceBtn.disabled = false;
                    addInstanceBtn.classList.remove('opacity-50', 'cursor-not-allowed');
                    addInstanceBtn.title = '';
                }
            }
        },

        createInstanceCard(instance) {
            // Map API fields: connected (boolean) to status string
            const isConnected = instance.connected?.Bool === true || instance.connected === true;
            const isPaused = instance.paused === true;
            
            let statusClass, statusBg, statusText, borderClass;
            if (isConnected) {
                statusClass = 'text-mz-green';
                statusBg = 'bg-mz-green-light';
                statusText = 'Conectado';
                borderClass = 'border-mz-green';
            } else if (isPaused) {
                statusClass = 'text-yellow-600';
                statusBg = 'bg-yellow-50';
                statusText = 'Em pausa';
                borderClass = 'border-gray-200';
            } else {
                statusClass = 'text-gray-600';
                statusBg = 'bg-gray-100';
                statusText = 'Desconectado';
                borderClass = 'border-gray-200';
            }

            const card = document.createElement('div');
            card.className = `bg-white p-5 rounded-xl shadow-sm border ${borderClass} flex flex-col`;
            card.dataset.instanceId = instance.id;
            card.dataset.status = isConnected ? 'connected' : isPaused ? 'paused' : 'disconnected';

            card.innerHTML = `
                <div class="flex justify-between items-center mb-4">
                    <h2 class="font-semibold text-lg text-gray-800">${instance.name || 'Sem nome'}</h2>
                    <span class="${statusBg} ${statusClass} text-xs font-bold px-3 py-1 rounded-full">${statusText}</span>
                </div>
                <div class="text-xs text-gray-400 mb-4">
                    ID: ${instance.token || 'N/A'}
                </div>
                <div class="flex justify-between items-center mb-6">
                    <div>
                        <span class="text-sm text-gray-500">Data da criação</span>
                        <p class="font-semibold text-gray-800">${formatDate(instance.created_at)}</p>
                    </div>
                    <div>
                        <span class="text-sm text-gray-500">Análises concluídas</span>
                        <p class="font-semibold text-gray-800">${instance.analysis_count || 0}</p>
                    </div>
                </div>
                <div class="flex space-x-2 mt-auto">
                    ${isConnected ? `
                        <button class="flex-1 bg-gray-800 text-white font-medium py-2 px-4 rounded-lg hover:bg-gray-900 transition-colors disconnect-btn">
                            Desconectar
                        </button>
                    ` : `
                        <button class="flex-1 bg-mz-green text-white font-medium py-2 px-4 rounded-lg hover:bg-green-700 transition-colors connect-btn">
                            Conectar WhatsApp
                        </button>
                    `}
                    <button class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 transition-colors delete-btn">
                        Excluir
                    </button>
                </div>
            `;

            // Event listeners
            const connectBtn = card.querySelector('.connect-btn');
            const disconnectBtn = card.querySelector('.disconnect-btn');
            const deleteBtn = card.querySelector('.delete-btn');

            if (connectBtn) {
                connectBtn.addEventListener('click', () => Handlers.connectInstance(instance.id));
            }
            if (disconnectBtn) {
                disconnectBtn.addEventListener('click', () => Handlers.disconnectInstance(instance.id));
            }
            if (deleteBtn) {
                deleteBtn.addEventListener('click', () => Handlers.openDeleteModal(instance.id));
            }

            return card;
        },

        renderInstances() {
            const grid = document.getElementById('cardGrid');
            grid.innerHTML = '';

            const filtered = state.instances.filter(instance => {
                // Map connected field (boolean) to status for filtering
                const isConnected = instance.connected?.Bool === true || instance.connected === true;
                const isPaused = instance.paused === true;
                
                // Filtro de status
                let matchesStatus = false;
                if (state.currentFilter === 'all') {
                    matchesStatus = true; // Mostra todas as instâncias
                } else if (state.currentFilter === 'connected') {
                    matchesStatus = isConnected;
                } else if (state.currentFilter === 'disconnected') {
                    matchesStatus = !isConnected && !isPaused;
                } else if (state.currentFilter === 'paused') {
                    matchesStatus = isPaused;
                }

                // Filtro de busca
                const matchesSearch = !state.searchQuery || 
                    instance.name?.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
                    instance.jid?.includes(state.searchQuery);

                return matchesStatus && matchesSearch;
            });

            if (filtered.length === 0) {
                grid.innerHTML = `
                    <div class="col-span-full text-center py-12 text-gray-500">
                        Nenhuma instância encontrada
                    </div>
                `;
                return;
            }

            filtered.forEach(instance => {
                grid.appendChild(UI.createInstanceCard(instance));
            });
        },

        renderPlans(plans) {
            const container = document.getElementById('plansContainer');
            if (!container || !plans) return;

            container.innerHTML = '';

            plans.forEach(plan => {
                const isActive = state.subscription?.plan_id === plan.id;
                const features = JSON.parse(plan.features || '[]');

                const planCard = document.createElement('div');
                planCard.className = `bg-white p-6 rounded-xl shadow-sm border ${isActive ? 'border-2 border-mz-green' : 'border border-gray-200'} flex flex-col space-y-4`;

                const featuresList = features.map(feature => `
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                        </svg>
                        <span>${feature}</span>
                    </li>
                `).join('');

                planCard.innerHTML = `
                    <h3 class="text-xl font-semibold text-gray-800">${plan.name}</h3>
                    <p class="text-3xl font-bold text-gray-800">R$${plan.price}</p>
                    <ul class="space-y-2 text-gray-600">
                        ${featuresList}
                    </ul>
                    ${!isActive ? `
                        <button class="bg-mz-green text-white font-semibold py-3 px-5 rounded-lg shadow-sm hover:bg-green-700 transition-colors mt-auto">
                            Fazer upgrade
                        </button>
                    ` : ''}
                `;

                container.appendChild(planCard);
            });
        }
    };

    // Manipuladores de eventos
    const Handlers = {
        async connectInstance(instanceId) {
            try {
                const instance = state.instances.find(i => i.id === instanceId);
                if (!instance || !instance.token) {
                    throw new Error('Token da instância não encontrado');
                }

                if (instance.connected?.Bool === true || instance.connected === true) {
                    alert('Esta instância já está conectada!');
                    return;
                }

                state.currentInstanceId = instanceId;
                
                // Abre o modal imediatamente com um placeholder
                const qrContainer = document.getElementById('qrCodeContainer');
                qrContainer.innerHTML = '<div class="w-full h-full flex items-center justify-center text-gray-500">Gerando QR Code...</div>';
                Modals.open('qrModal');

                // Inicia a conexão em segundo plano
                API.connectInstance(instance.token).catch(error => {
                    console.error('Falha ao iniciar a conexão:', error);
                    Modals.close('qrModal');
                    alert('Erro ao iniciar a conexão: ' + (error.message || 'Erro desconhecido'));
                });

                // Função para tentar obter o QR Code
                const pollForQRCode = async () => {
                    try {
                        const qrData = await API.getQRCode(instance.token);
                        const qrCode = qrData.QRCode || qrData.qr || qrData;

                        if (qrCode && qrCode !== '') {
                            qrContainer.innerHTML = `<img src="${qrCode}" alt="QR Code" class="w-full h-full object-contain rounded-lg">`;
                            // Inicia o polling de status APENAS se o QR Code for recebido
                            Handlers.pollConnectionStatus(instanceId);
                            return; // Para o polling do QR Code
                        }
                    } catch (error) {
                        // Ignora erros 404 ou similares, pois o QR code pode não estar pronto
                        console.warn('Ainda não há QR Code, tentando novamente...');
                    }
                    
                    // Se o modal foi fechado, para de tentar
                    if (document.getElementById('qrModal').classList.contains('hidden')) {
                        return;
                    }
                    
                    // Tenta novamente após um curto período
                    setTimeout(pollForQRCode, 2000);
                };

                // Inicia a busca pelo QR Code
                pollForQRCode();

            } catch (error) {
                console.error('Erro ao conectar instância:', error);
                alert('Erro ao conectar instância: ' + (error.message || 'Erro desconhecido'));
                Modals.close('qrModal');
            }
        },

        async pollConnectionStatus(instanceId) {
            const interval = setInterval(async () => {
                try {
                    const instances = await API.getInstances();
                    const instance = instances.find(i => i.id === instanceId);
                    
                    // Check if connected using the boolean field
                    if (instance && (instance.connected?.Bool === true || instance.connected === true)) {
                        clearInterval(interval);
                        
                        // Mostrar mensagem de instrução com GIF
                        const qrContainer = document.getElementById('qrCodeContainer');
                        qrContainer.innerHTML = `
                            <div class="flex flex-col items-center justify-center space-y-4">
                                <img src="/images/aba.gif" alt="Feche o WhatsApp" class="w-48 h-48 object-contain">
                                <h3 class="text-lg font-bold text-gray-800">Feche seu WhatsApp para a conexão efetivar</h3>
                                <p class="text-sm text-gray-600 text-center">Ao abrir novamente estará conectado corretamente</p>
                            </div>
                        `;
                        
                        // Fechar modal após 10 segundos
                        setTimeout(async () => {
                            Modals.close('qrModal');
                            await loadInstances();
                            alert('WhatsApp conectado com sucesso!');
                        }, 10000);
                    }
                } catch (error) {
                    console.error('Error polling status:', error);
                    clearInterval(interval);
                }
            }, 3000);

            // Stop polling after 5 minutes
            setTimeout(() => clearInterval(interval), 300000);
        },

        async disconnectInstance(instanceId) {
            if (!confirm('Deseja realmente desconectar esta instância?')) return;

            try {
                // Buscar a instância para pegar o token
                const instance = state.instances.find(i => i.id === instanceId);
                if (!instance || !instance.token) {
                    throw new Error('Instance token not found');
                }
                
                await API.disconnectInstance(instance.token);
                await loadInstances();
                alert('Instância desconectada com sucesso!');
            } catch (error) {
                console.error('Error disconnecting instance:', error);
                alert('Erro ao desconectar instância: ' + (error.message || 'Erro desconhecido'));
            }
        },

        openDeleteModal(instanceId) {
            state.currentInstanceId = instanceId;
            Modals.open('deleteModal');
        },

        async deleteInstance() {
            if (!state.currentInstanceId) return;

            try {
                await API.deleteInstance(state.currentInstanceId);
                Modals.close('deleteModal');
                await loadInstances();
                alert('Instância excluída com sucesso!');
            } catch (error) {
                console.error('Error deleting instance:', error);
                alert('Erro ao excluir instância.');
            }
        },

        async addInstance() {
            const name = document.getElementById('instanceNameInput').value.trim();
            if (!name) {
                alert('Por favor, insira um nome para a instância.');
                return;
            }

            try {
                await API.createInstance(name);
                Modals.close('instanceModal');
                document.getElementById('instanceNameInput').value = '';
                await loadInstances();
                alert('Instância criada com sucesso!');
            } catch (error) {
                console.error('Error creating instance:', error);
                const errorMessage = error.message || 'Erro ao criar instância. Verifique se você ainda tem slots disponíveis.';
                alert(errorMessage);
            }
        },

        setFilter(filter) {
            state.currentFilter = filter;
            
            // Update tabs
            ['tabAll', 'tabConnected', 'tabDisconnected', 'tabPaused'].forEach(id => {
                const btn = document.getElementById(id);
                if (btn) {
                    btn.className = 'text-gray-500 font-medium px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors';
                }
            });

            const activeTab = {
                'all': 'tabAll',
                'connected': 'tabConnected',
                'disconnected': 'tabDisconnected',
                'paused': 'tabPaused'
            }[filter];

            const activeBtn = document.getElementById(activeTab);
            if (activeBtn) {
                activeBtn.className = 'bg-white text-gray-800 font-medium px-4 py-2 rounded-lg shadow-sm border border-gray-200';
            }

            UI.renderInstances();
        },

        search(query) {
            state.searchQuery = query;
            UI.renderInstances();
        }
    };

    // Gerenciamento de modais
    const Modals = {
        open(modalId) {
            document.getElementById(modalId)?.classList.remove('hidden');
        },

        close(modalId) {
            document.getElementById(modalId)?.classList.add('hidden');
        },

        init() {
            // QR Modal
            document.getElementById('closeModalBtn')?.addEventListener('click', () => Modals.close('qrModal'));
            document.getElementById('qrModal')?.addEventListener('click', (e) => {
                if (e.target.id === 'qrModal') Modals.close('qrModal');
            });

            // Delete Modal
            document.getElementById('cancelDeleteBtn')?.addEventListener('click', () => Modals.close('deleteModal'));
            document.getElementById('confirmDeleteBtn')?.addEventListener('click', () => Handlers.deleteInstance());
            document.getElementById('deleteModal')?.addEventListener('click', (e) => {
                if (e.target.id === 'deleteModal') Modals.close('deleteModal');
            });

            // Instance Modal
            document.getElementById('addInstanceBtn')?.addEventListener('click', () => Modals.open('instanceModal'));
            document.getElementById('cancelAddInstanceBtn')?.addEventListener('click', () => Modals.close('instanceModal'));
            document.getElementById('confirmAddInstanceBtn')?.addEventListener('click', () => Handlers.addInstance());
            document.getElementById('instanceModal')?.addEventListener('click', (e) => {
                if (e.target.id === 'instanceModal') Modals.close('instanceModal');
            });

            // Orange Alert
            document.getElementById('closeOrangeAlertBtn')?.addEventListener('click', () => {
                document.getElementById('orangeAlert')?.classList.add('hidden');
            });
        }
    };

    // Navegação entre páginas
    const Navigation = {
        init() {
            const linkContas = document.getElementById('linkContas');
            const linkDados = document.getElementById('linkDados');
            const paginaDashboard = document.getElementById('paginaDashboard');
            const paginaDados = document.getElementById('paginaDados');

            const activeClass = 'sidebar-link flex items-center space-x-3 font-medium px-4 py-3 rounded-lg bg-mz-green-light text-mz-green';
            const inactiveClass = 'sidebar-link flex items-center space-x-3 font-medium px-4 py-3 rounded-lg text-gray-600 hover:bg-gray-100';

            linkContas?.addEventListener('click', (e) => {
                e.preventDefault();
                paginaDashboard?.classList.remove('hidden');
                paginaDados?.classList.add('hidden');
                linkContas.className = activeClass;
                linkDados.className = inactiveClass;
            });

            linkDados?.addEventListener('click', (e) => {
                e.preventDefault();
                paginaDashboard?.classList.add('hidden');
                paginaDados?.classList.remove('hidden');
                linkContas.className = inactiveClass;
                linkDados.className = activeClass;
            });

            // Set initial state
            linkContas.className = activeClass;
            linkDados.className = inactiveClass;
        }
    };

    // WhatsApp number editing
    const WhatsAppEditor = {
        init() {
            const input = document.getElementById('whatsapp');
            const editBtn = document.getElementById('editWhatsappBtn');
            const saveBtn = document.getElementById('saveWhatsappBtn');

            editBtn?.addEventListener('click', () => {
                input.disabled = false;
                input.classList.remove('bg-gray-100');
                input.classList.add('bg-white', 'focus:ring-2', 'focus:ring-mz-green');
                input.focus();
                saveBtn.classList.remove('hidden');
            });

            saveBtn?.addEventListener('click', async () => {
                const whatsapp = input.value.trim();
                if (!whatsapp) {
                    alert('Por favor, insira um número de WhatsApp válido');
                    return;
                }

                try {
                    await API.updateUser({ whatsapp_number: whatsapp });
                    
                    // Disable input and hide save button
                    input.disabled = true;
                    input.classList.add('bg-gray-100');
                    input.classList.remove('bg-white', 'focus:ring-2', 'focus:ring-mz-green');
                    saveBtn.classList.add('hidden');
                    
                    alert('Número de WhatsApp salvo com sucesso!');
                } catch (error) {
                    console.error('Error saving WhatsApp number:', error);
                    alert('Erro ao salvar número de WhatsApp');
                }
            });
        }
    };

    // Carregamento de dados
    async function loadUser() {
        try {
            state.user = await API.getUser();
            UI.updateUserInfo();
        } catch (error) {
            console.error('Error loading user:', error);
        }
    }

    async function loadInstances() {
        try {
            state.instances = await API.getInstances();
            console.log('Instances loaded:', state.instances); // Debug log
            UI.renderInstances();
            UI.updateInstancesProgress();
        } catch (error) {
            console.error('Error loading instances:', error);
        }
    }

    async function loadSubscription() {
        try {
            const response = await API.getSubscription();
            // A resposta vem com { subscription: {...}, instance_count, instances_remaining, etc }
            if (response.subscription) {
                state.subscription = {
                    ...response.subscription,
                    connected_count: response.connected_count,
                    instances_remaining: response.instances_remaining,
                    max_instances: response.max_instances,
                    plan_id: response.plan_id,
                    is_expired: response.is_expired,
                    expires_at: response.subscription.expires_at,
                    plan: response.subscription.plan
                };
            } else {
                state.subscription = response;
            }
            UI.updateInstancesProgress();
        } catch (error) {
            console.error('Error loading subscription:', error);
        }
    }

    async function loadPlans() {
        try {
            const plans = await API.getPlans();
            UI.renderPlans(plans);
        } catch (error) {
            console.error('Error loading plans:', error);
        }
    }

    // Inicialização
    async function init() {
        // Check auth
        const token = API.getToken();
        console.log('Dashboard V4 initialized with token:', token ? 'Token exists' : 'No token');
        
        if (!token) {
            window.location.href = '/login/';
            return;
        }

        // Setup avatar dropdown
        const avatarButton = document.getElementById('avatarButton');
        const avatarDropdown = document.getElementById('avatarDropdown');
        const logoutButton = document.getElementById('logoutButton');

        // Toggle dropdown
        avatarButton?.addEventListener('click', (e) => {
            e.stopPropagation();
            avatarDropdown.classList.toggle('hidden');
        });

        // Close dropdown when clicking outside
        document.addEventListener('click', (e) => {
            if (!avatarButton?.contains(e.target) && !avatarDropdown?.contains(e.target)) {
                avatarDropdown?.classList.add('hidden');
            }
        });

        // Logout handler
        logoutButton?.addEventListener('click', () => {
            localStorage.removeItem('token');
            localStorage.removeItem('auth_token');
            localStorage.removeItem('authToken');
            window.location.href = '/login/';
        });

        // Setup event listeners
        Navigation.init();
        Modals.init();
        WhatsAppEditor.init();

        // Tab filters
        document.getElementById('tabAll')?.addEventListener('click', () => Handlers.setFilter('all'));
        document.getElementById('tabConnected')?.addEventListener('click', () => Handlers.setFilter('connected'));
        document.getElementById('tabDisconnected')?.addEventListener('click', () => Handlers.setFilter('disconnected'));
        document.getElementById('tabPaused')?.addEventListener('click', () => Handlers.setFilter('paused'));

        // Search
        document.getElementById('searchInput')?.addEventListener('input', (e) => Handlers.search(e.target.value));

        // Load data
        await Promise.all([
            loadUser(),
            loadSubscription(),
            loadInstances(),
            loadPlans()
        ]);
    }

    // Start application
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
