<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard Metrizap</title>
    
    <!-- Importação do Tailwind CSS --><script src="https://cdn.tailwindcss.com"></script>
    
    <!-- Importação da fonte Inter --><link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    
    <!-- Configuração do Tailwind para usar a fonte Inter e cores personalizadas --><script>
      tailwind.config = {
        theme: {
          extend: {
            fontFamily: {
              'sans': ['Inter', 'sans-serif'],
            },
            colors: {
              // Cores personalizadas do design
              'mz-green': '#28a745',
              'mz-green-light': '#e9f7ec',
              'mz-red': '#dc3545',
              'mz-orange-light': '#fff3e0', // Nova cor para o alerta
              'mz-orange-dark': '#fd7e14',   // Nova cor para o texto do alerta
            }
          }
        }
      }
    </script>
</head>
<body class="bg-gray-50 font-sans">

    <div class="flex h-screen w-full">
        
        <!-- ==== Barra Lateral ==== -->
        <aside class="w-64 flex-shrink-0 bg-white border-r border-gray-200 flex flex-col h-screen">
            
            <!-- Logo -->
            <div class="h-20 flex items-center px-6">
                <h1 class="text-3xl font-bold text-gray-800">metrizap</h1>
            </div>

            <!-- Navegação -->
            <nav class="flex-1 px-4 py-4 space-y-6">
                
                <!-- Grupo Principal -->
                <div>
                    <h3 class="px-3 mb-2 text-xs text-gray-500 uppercase tracking-wider font-semibold">PRINCIPAL</h3>
                    <!-- Link Contas Conectadas -->
                    <a href="#" id="linkContas" class="sidebar-link flex items-center space-x-3 font-medium px-4 py-3 rounded-lg">
                        <!-- Ícone SVG Grid -->
                        <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
                        <span>Contas conectadas</span>
                    </a>
                </div>

                <!-- Grupo Perfil -->
                <div>
                    <h3 class="px-3 mb-2 text-xs text-gray-500 uppercase tracking-wider font-semibold">PERFIL</h3>
                    <!-- Link Seus Dados -->
                    <a href="#" id="linkDados" class="sidebar-link flex items-center space-x-3 font-medium px-4 py-3 rounded-lg">
                        <!-- Ícone SVG User -->
                        <svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                        <span>Seus dados</span>
                    </a>
                </div>
            </nav>

            <!-- Rodapé da Barra Lateral -->
            <div class="p-6 border-t border-gray-200">
                <p class="text-sm text-gray-600 mb-2">4 contas conectadas restantes</p>
                <!-- Aqui é uma barra para mostrar quantas contas o usuário ainda pode conectar com base no seu plano -->
                <div class="w-full bg-gray-200 rounded-full h-1.5">
                    <!-- Barra de progresso -->
                    <div class="bg-mz-green h-1.5 rounded-full" style="width: 60%"></div>
                </div>
            </div>
        </aside>

        <!-- ==== Conteúdo Principal ==== -->
        <main class="flex-1 p-8 overflow-y-auto flex flex-col">
            
            <!-- Cabeçalho (Fixo para ambas as páginas) -->
            <header class="flex justify-between items-center mb-8">
                <!-- Após o Olá, é o primeiro nome do usuário que aparece -->
                <h1 class="text-4xl font-semibold text-gray-800">Olá, Allan 👋</h1>
                <!-- Avatar (Placeholder) -->
                <img class="w-12 h-12 rounded-full" src="https://placehold.co/48x48/E2E8F0/4A5568?text=A" alt="Avatar do usuário">
            </header>

            <!-- ==== Página Dashboard (Contas) ==== -->
            <div id="paginaDashboard" class="">
                
                <!-- Controles do Dashboard -->
                <section class="flex justify-between items-center mb-6">
                    <div class="flex items-center space-x-4 flex-grow">
                        <!-- Campo de Busca -->
                        <div class="relative flex-grow max-w-sm">
                            <!-- Ícone SVG Search -->
                            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
                            <input type="text" placeholder="Pesquise por nome ou número" class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-mz-green focus:border-transparent transition-all">
                        </div>
                        
                        <!-- Abas -->
                        <div class="flex space-x-2">
                            <button class="bg-white text-gray-800 font-medium px-4 py-2 rounded-lg shadow-sm border border-gray-200">Conectados</button>
                            <button class="text-gray-500 font-medium px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors">Desconectadas</button>
                            <button class="text-gray-500 font-medium px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors">Em pausa</button>
                        </div>
                    </div>

                    <!-- Botão Adicionar WhatsApp -->
                    <button id="addInstanceBtn" class="ml-4 bg-mz-green text-white font-semibold px-5 py-2.5 rounded-lg shadow-sm hover:bg-green-700 transition-colors">
                        Adicionar WhatsApp
                    </button>
                </section>

                <!-- Grid de Cards --><section id="cardGrid" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">

                    <!-- Card 1: Conectado --><div class="bg-white p-5 rounded-xl shadow-sm border border-mz-green flex flex-col">
                        <div class="flex justify-between items-center mb-4">
                            <h2 class="font-semibold text-lg text-gray-800">Allan</h2>
                            <span class="bg-mz-green-light text-mz-green text-xs font-bold px-3 py-1 rounded-full">Conectado</span>
                        </div>
                        <div class="flex justify-between items-center mb-6">
                            <div>
                                <span class="text-sm text-gray-500">Data da criação</span>
                                <p class="font-semibold text-gray-800">03/03/2025</p>
                            </div>
                            <div>
                                <span class="text-sm text-gray-500">Análises concluídas</span>
                                <p class="font-semibold text-gray-800">110</p>
                            </div>
                        </div>
                        <div class="flex space-x-2 mt-auto">
                            <button class="flex-1 bg-gray-800 text-white font-medium py-2 px-4 rounded-lg hover:bg-gray-900 transition-colors">Desconectar</button>
                            <button class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 transition-colors open-delete-modal-btn">Excluir</button>
                        </div>
                    </div>
    
                    <!-- Card 2: Desconectado --><div class="bg-white p-5 rounded-xl shadow-sm border border-gray-200 flex flex-col">
                        <div class="flex justify-between items-center mb-4">
                            <h2 class="font-semibold text-lg text-gray-800">Emellyn Severo</h2>
                            <span class="bg-gray-100 text-gray-600 text-xs font-bold px-3 py-1 rounded-full">Desconectado</span>
                        </div>
                        <div class="flex justify-between items-center mb-6">
                            <div>
                                <span class="text-sm text-gray-500">Data da criação</span>
                                <p class="font-semibold text-gray-800">10/03/2025</p>
                            </div>
                            <div>
                                <span class="text-sm text-gray-500">Análises concluídas</span>
                                <p class="font-semibold text-gray-800">0</p>
                            </div>
                        </div>
                        <div class="flex space-x-2 mt-auto">
                            <button class="flex-1 bg-mz-green text-white font-medium py-2 px-4 rounded-lg hover:bg-green-700 transition-colors open-modal-btn">Conectar WhatsApp</button>
                            <button class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 transition-colors open-delete-modal-btn">Excluir</button>
                        </div>
                    </div>
    
                    <!-- Card 3: Desconectado --><div class="bg-white p-5 rounded-xl shadow-sm border border-gray-200 flex flex-col">
                        <div class="flex justify-between items-center mb-4">
                            <h2 class="font-semibold text-lg text-gray-800">Emellyn Severo</h2>
                            <span class="bg-gray-100 text-gray-600 text-xs font-bold px-3 py-1 rounded-full">Desconectado</span>
                        </div>
                        <div class="flex justify-between items-center mb-6">
                            <div>
                                <span class="text-sm text-gray-500">Data da criação</span>
                                <p class="font-semibold text-gray-800">10/03/2025</p>
                            </div>
                            <div>
                                <span class="text-sm text-gray-500">Análises concluídas</span>
                                <p class="font-semibold text-gray-800">0</p>
                            </div>
                        </div>
                        <div class="flex space-x-2 mt-auto">
                            <button class="flex-1 bg-mz-green text-white font-medium py-2 px-4 rounded-lg hover:bg-green-700 transition-colors open-modal-btn">Conectar WhatsApp</button>
                            <button class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 transition-colors open-delete-modal-btn">Excluir</button>
                        </div>
                    </div>
    
                </section>

                <!-- Card de Alerta na parte inferior -->
                <div id="orangeAlert" class="bg-mz-orange-light border border-mz-orange-dark text-mz-orange-dark p-4 rounded-lg flex items-start space-x-3 text-sm mt-auto">
                    <!-- Ícone SVG Info -->
                    <svg class="w-5 h-5 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
                    <p class="flex-grow">Todos os dias às 18 horas você receberá no WhatsApp a análise do dia de cada número conectado.</p>
                    <!-- Botão de Fechar (X) -->
                    <button id="closeOrangeAlertBtn" class="ml-auto -mr-1 -mt-1 p-1 rounded-md text-mz-orange-dark hover:bg-mz-orange-light/50 transition-colors">
                        <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                    </button>
                </div>
            </div>

            <!-- ==== Página Seus Dados ==== -->
            <div id="paginaDados" class="hidden">
                <h1 class="text-4xl font-semibold text-gray-800 mb-8">Allan, seus dados 👋</h1>

                <div class="space-y-10 max-w-5xl">
                    
                    <!-- Formulário de Dados -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
                        <!-- Seus dados -->
                        <div>
                            <label for="nome" class="block text-sm font-medium text-gray-700 mb-1">Seus dados</label>
                            <input type="text" id="nome" value="Allan Miranda Severo Rodrigues" class="w-full p-3 border border-gray-200 rounded-lg bg-gray-100 text-gray-500" disabled>
                        </div>
                        <!-- E-mail -->
                        <div>
                            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">E-mail</label>
                            <input type="email" id="email" value="grupoteaser@gmail.com" class="w-full p-3 border border-gray-200 rounded-lg bg-gray-100 text-gray-500" disabled>
                        </div>
                        <!-- Senha atual -->
                        <div>
                            <label for="senha" class="block text-sm font-medium text-gray-700 mb-1">Senha atual</label>
                            <input type="password" id="senha" value="************" class="w-full p-3 border border-gray-200 rounded-lg bg-gray-100 text-gray-500" disabled>
                            <a href="#" class="text-sm text-blue-600 hover:underline mt-1 inline-block">Alterar senha</a>
                        </div>
                        <!-- Quero receber análises no: -->
                        <!-- Ao invés de modal dentro da instância, o número fica registrado aqui e servirá globalmente para este usuário como "enviar_para" -->
                        <div>
                            <label for="whatsapp" class="block text-sm font-medium text-gray-700 mb-1">Quero receber análises no:</label>
                            <div class="relative">
                                <input type="text" id="whatsapp" value="+55 51 98193-6133" class="w-full p-3 border border-gray-300 rounded-lg bg-white text-gray-800 pr-10 focus:outline-none focus:ring-2 focus:ring-mz-green focus:border-transparent transition-all">
                                <!-- Ícone SVG Pencil -->
                                <svg class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path></svg>
                            </div>
                        </div>
                    </div>

                    <!-- Plano Atual -->
                    <div>
                        <h2 class="text-2xl font-semibold text-gray-800 mb-6">Plano atual</h2>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                            
                            <!-- Card Plano Pro -->
                            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 flex flex-col space-y-4">
                                <h3 class="text-xl font-semibold text-gray-800">Análise Pro</h3>
                                <p class="text-3xl font-bold text-gray-800">R$29</p>
                                <ul class="space-y-2 text-gray-600">
                                    <!-- Ícone SVG Check Circle -->
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Análise diária</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>8 contas conectadas</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Agendamento de mensgens</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Gerente de conta</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Suporte prioritário</span></li>
                                </ul>
                            </div>
                            
                            <!-- Card Plano Analista (Ativo) -->
                            <div class="bg-white p-6 rounded-xl shadow-sm border-2 border-mz-green flex flex-col space-y-4">
                                <h3 class="text-xl font-semibold text-gray-800">Análise Analista</h3>
                                <p class="text-3xl font-bold text-gray-800">R$97</p>
                                <ul class="space-y-2 text-gray-600">
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Análise diária</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>20 contas conectadas</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Métricas personalizadas</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Análise mensal</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Gerente de conta</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Agendamento de mensgens</span></li>
                                    <li class="flex items-center space-x-2"><svg class="w-5 h-5 text-mz-green" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg><span>Suporte prioritário</span></li>
                                </ul>
                                <button class="bg-mz-green text-white font-semibold py-3 px-5 rounded-lg shadow-sm hover:bg-green-700 transition-colors mt-auto">
                                    Fazer upgrade
                                </button>
                            </div>
                        </div>
                    </div>

                </div>
            </div>

        </main>
    </div>

    <!-- ==== Modal QR Code ==== -->
    <!-- Fundo do Overlay -->
    <div id="qrModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center p-4 z-50 hidden">
        <!-- Conteúdo do Modal -->
        <div class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-md relative text-center">
            <!-- Botão de Fechar -->
            <button id="closeModalBtn" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition-colors">
                <!-- Ícone SVG X -->
                <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
            </button>
            
            <!-- Instruções -->
            <p class="text-sm text-gray-600 mb-6">
                Abra seu WhatsApp → Dispositivos conectados → Aponte a câmera
            </p>
            
            <!-- Imagem do QR Code (Placeholder) -->
            <!-- AQUI VAI O QRCODE PARA CONECTAR O WHATSAPP -->
            <div class="w-64 h-64 mx-auto my-4">
                <img src="https://placehold.co/256x256/000000/FFFFFF?text=QR+CODE" alt="QR Code para conectar WhatsApp" class="w-full h-full object-cover rounded-lg">
            </div>
            
            <!-- Aviso Verde -->
            <div class="bg-mz-green-light border border-mz-green text-mz-green p-3 rounded-lg flex items-start space-x-2 text-sm text-left mt-6">
                <!-- Ícone SVG Info -->
                <svg class="w-5 h-5 flex-shrink-0 mt-0.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
                <p>Diariamente você receberá análises desse número até desconecta-lo.</p>
            </div>
        </div>
    </div>

    <!-- ==== Modal Excluir ==== -->
    <div id="deleteModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center p-4 z-50 hidden">
        <!-- Conteúdo do Modal -->
        <div class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-sm relative text-center">
            <h3 class="text-lg font-semibold text-gray-800 mb-6">Tem certeza que é isso que deseja?</h3>
            
            <div class="flex space-x-4">
                <button id="cancelDeleteBtn" class="flex-1 bg-gray-200 text-gray-700 font-medium py-2 px-4 rounded-lg hover:bg-gray-300 transition-colors">
                    Cancelar
                </button>
                <button id="confirmDeleteBtn" class="flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 transition-colors">
                    Sim, excluir
                </button>
            </div>
        </div>
    </div>

    <!-- ==== Modal Adicionar Instância ==== -->
    <div id="instanceModal" class="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center p-4 z-50 hidden">
        <!-- Conteúdo do Modal -->
        <div class="bg-white p-8 rounded-2xl shadow-xl w-full max-w-sm relative text-center">
            <h3 class="text-xl font-semibold text-gray-800 mb-6">Nome da instância</h3>
            
            <!-- Campo de Input -->
            <div class="relative mb-6">
                <input id="instanceNameInput" type="text" placeholder="Armelin Nunes" class="w-full pl-4 pr-10 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-mz-green focus:border-transparent transition-all">
                <!-- Ícone SVG Pencil -->
                <svg class="absolute right-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path></svg>
            </div>
            
            <div class="flex space-x-4">
                <button id="cancelAddInstanceBtn" class="flex-1 bg-gray-200 text-gray-700 font-medium py-2 px-4 rounded-lg hover:bg-gray-300 transition-colors">
                    Cancelar
                </button>
                <button id="confirmAddInstanceBtn" class="flex-1 bg-mz-green text-white font-medium py-2 px-4 rounded-lg hover:bg-green-700 transition-colors">
                    Sim, adicionar
                </button>
            </div>
        </div>
    </div>


    <!-- Script de interatividade -->
    <script>
      // Espera o DOM estar completamente carregado antes de executar o script
      window.onload = () => {
        
        // --- Lógica do Modal ---
        const modal = document.getElementById('qrModal');
        const closeModalBtn = document.getElementById('closeModalBtn');
        const openModalBtns = document.querySelectorAll('.open-modal-btn');

        const openModal = () => {
          if (modal) modal.classList.remove('hidden');
        };

        const closeModal = () => {
          if (modal) modal.classList.add('hidden');
        };

        // Adiciona evento para todos os botões "Conectar WhatsApp"
        openModalBtns.forEach(btn => {
          btn.addEventListener('click', openModal);
        });

        // Adiciona evento para o botão de fechar (X)
        if (closeModalBtn) {
          closeModalBtn.addEventListener('click', closeModal);
        }

        // Adiciona evento para fechar ao clicar fora do modal (no overlay)
        if (modal) {
          modal.addEventListener('click', (e) => {
            if (e.target === modal) {
              closeModal();
            }
          });
        }
        // --- Fim da Lógica do Modal QR Code ---

        // --- Lógica do Modal Excluir ---
        const deleteModal = document.getElementById('deleteModal');
        const openDeleteModalBtns = document.querySelectorAll('.open-delete-modal-btn');
        const cancelDeleteBtn = document.getElementById('cancelDeleteBtn');
        const confirmDeleteBtn = document.getElementById('confirmDeleteBtn');
        
        let cardToDelete = null; // Variável para guardar qual card será excluído

        const openDeleteModal = (e) => {
          // Encontra o elemento 'card' (div) pai do botão que foi clicado
          cardToDelete = e.target.closest('.flex.flex-col'); 
          if (deleteModal) deleteModal.classList.remove('hidden');
        };

        const closeDeleteModal = () => {
          if (deleteModal) deleteModal.classList.add('hidden');
          cardToDelete = null; // Limpa a referência ao card
        };

        // Adiciona evento para todos os botões "Excluir"
        openDeleteModalBtns.forEach(btn => {
          btn.addEventListener('click', openDeleteModal);
        });

        // Adiciona evento para o botão "Cancelar"
        if (cancelDeleteBtn) {
          cancelDeleteBtn.addEventListener('click', closeDeleteModal);
        }

        // Adiciona evento para o botão "Sim, excluir"
        if (confirmDeleteBtn) {
          confirmDeleteBtn.addEventListener('click', () => {
            if (cardToDelete) {
              cardToDelete.remove(); // Remove o card do DOM
            }
            closeDeleteModal(); // Fecha o modal e limpa a variável
          });
        }

        // Adiciona evento para fechar (cancelar) ao clicar fora do modal
        if (deleteModal) {
          deleteModal.addEventListener('click', (e) => {
            if (e.target === deleteModal) {
              closeDeleteModal();
            }
          });
        }
        // --- Fim da Lógica do Modal Excluir ---

        // --- Lógica do Modal Adicionar Instância ---
        const addInstanceModal = document.getElementById('instanceModal');
        const addInstanceBtn = document.getElementById('addInstanceBtn'); // O botão principal "Adicionar WhatsApp"
        const cancelAddInstanceBtn = document.getElementById('cancelAddInstanceBtn');
        const confirmAddInstanceBtn = document.getElementById('confirmAddInstanceBtn');
        const instanceNameInput = document.getElementById('instanceNameInput');
        const cardGrid = document.getElementById('cardGrid'); // A section que agrupa os cards

        /**
         * Cria um novo elemento de card (desconectado)
         * @param {string} name - O nome da instância
         * @returns {HTMLElement} - O elemento div do card
         */
        const createNewCard = (name) => {
            // 1. Cria o container principal do card
            const card = document.createElement('div');
            card.className = 'bg-white p-5 rounded-xl shadow-sm border border-gray-200 flex flex-col';

            // 2. Cria o Cabeçalho
            const cardHeader = document.createElement('div');
            cardHeader.className = 'flex justify-between items-center mb-4';
            
            const title = document.createElement('h2');
            title.className = 'font-semibold text-lg text-gray-800';
            title.textContent = name; // Seguro contra XSS
            
            const badge = document.createElement('span');
            badge.className = 'bg-gray-100 text-gray-600 text-xs font-bold px-3 py-1 rounded-full';
            badge.textContent = 'Desconectado';
            
            cardHeader.appendChild(title);
            cardHeader.appendChild(badge);

            // 3. Cria o Corpo
            const cardBody = document.createElement('div');
            cardBody.className = 'flex justify-between items-center mb-6';
            
            const creationDateDiv = document.createElement('div');
            const creationDateLabel = document.createElement('span');
            creationDateLabel.className = 'text-sm text-gray-500';
            creationDateLabel.textContent = 'Data da criação';
            const creationDateValue = document.createElement('p');
            creationDateValue.className = 'font-semibold text-gray-800';
            // Formata a data atual
            const today = new Date();
            const yyyy = today.getFullYear();
            const mm = String(today.getMonth() + 1).padStart(2, '0'); // Meses são 0-indexados
            const dd = String(today.getDate()).padStart(2, '0');
            creationDateValue.textContent = `${dd}/${mm}/${yyyy}`;
            creationDateDiv.appendChild(creationDateLabel);
            creationDateDiv.appendChild(creationDateValue);

            const analysisDiv = document.createElement('div');
            const analysisLabel = document.createElement('span');
            analysisLabel.className = 'text-sm text-gray-500';
            analysisLabel.textContent = 'Análises concluídas';
            const analysisValue = document.createElement('p');
            analysisValue.className = 'font-semibold text-gray-800';
            analysisValue.textContent = '0';
            analysisDiv.appendChild(analysisLabel);
            analysisDiv.appendChild(analysisValue);

            cardBody.appendChild(creationDateDiv);
            cardBody.appendChild(analysisDiv);

            // 4. Cria o Rodapé com botões
            const cardFooter = document.createElement('div');
            cardFooter.className = 'flex space-x-2 mt-auto';

            const connectBtn = document.createElement('button');
            connectBtn.className = 'flex-1 bg-mz-green text-white font-medium py-2 px-4 rounded-lg hover:bg-green-700 transition-colors open-modal-btn';
            connectBtn.textContent = 'Conectar WhatsApp';
            // Anexa o listener do modal de QR Code (função já existente)
            connectBtn.addEventListener('click', openModal); 
                                                    
            const deleteBtn = document.createElement('button');
            deleteBtn.className = 'flex-1 bg-mz-red text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 transition-colors open-delete-modal-btn';
            deleteBtn.textContent = 'Excluir';
            // Anexa o listener do modal de exclusão (função já existente)
            deleteBtn.addEventListener('click', openDeleteModal); 

            cardFooter.appendChild(connectBtn);
            cardFooter.appendChild(deleteBtn);

            // 5. Monta o card
            card.appendChild(cardHeader);
            card.appendChild(cardBody);
            card.appendChild(cardFooter);

            return card;
        };

        const openAddInstanceModal = () => {
            if (instanceNameInput) instanceNameInput.value = ''; // Limpa o input
            if (addInstanceModal) addInstanceModal.classList.remove('hidden');
        };

        const closeAddInstanceModal = () => {
            if (addInstanceModal) addInstanceModal.classList.add('hidden');
            if (instanceNameInput) instanceNameInput.value = ''; // Limpa o input
        };

        // Botão "Adicionar WhatsApp" (principal)
        if (addInstanceBtn) {
            addInstanceBtn.addEventListener('click', openAddInstanceModal);
        }

        // Botão "Cancelar" no modal
        if (cancelAddInstanceBtn) {
            cancelAddInstanceBtn.addEventListener('click', closeAddInstanceModal);
        }

        // Botão "Sim, adicionar" no modal
        if (confirmAddInstanceBtn) {
            confirmAddInstanceBtn.addEventListener('click', () => {
                const instanceName = instanceNameInput.value.trim();
                if (instanceName && cardGrid) { // Só adiciona se tiver um nome
                    const newCard = createNewCard(instanceName);
                    cardGrid.appendChild(newCard);
                    closeAddInstanceModal();
                }
            });
        }
        
        // Fechar ao clicar no overlay
        if (addInstanceModal) {
          addInstanceModal.addEventListener('click', (e) => {
            if (e.target === addInstanceModal) {
              closeAddInstanceModal();
            }
          });
        }
        // --- Fim da Lógica do Modal Adicionar Instância ---


        // --- Lógica do Alerta Laranja ---
        const orangeAlert = document.getElementById('orangeAlert');
        const closeOrangeAlertBtn = document.getElementById('closeOrangeAlertBtn');

        if (closeOrangeAlertBtn && orangeAlert) {
          closeOrangeAlertBtn.addEventListener('click', () => {
            orangeAlert.classList.add('hidden');
          });
        }
        // --- Fim da Lógica do Alerta Laranja ---


        // --- Lógica de Navegação das Páginas ---
        const linkContas = document.getElementById('linkContas');
        const linkDados = document.getElementById('linkDados');
        const paginaDashboard = document.getElementById('paginaDashboard');
        const paginaDados = document.getElementById('paginaDados');

        // Classes para estilo ativo/inativo
        const classeAtiva = 'sidebar-link flex items-center space-x-3 font-medium px-4 py-3 rounded-lg bg-mz-green-light text-mz-green';
        const classeInativa = 'sidebar-link flex items-center space-x-3 font-medium px-4 py-3 rounded-lg text-gray-600 hover:bg-gray-100';

        // Define o estado inicial (Página Dashboard ativa)
        if (linkContas) linkContas.className = classeAtiva; // <-- Correção: Verificações de nulo
        if (linkDados) linkDados.className = classeInativa;
        if (paginaDashboard) paginaDashboard.classList.remove('hidden');
        if (paginaDados) paginaDados.classList.add('hidden');

        // Clique em "Contas conectadas"
        if (linkContas) {
          linkContas.addEventListener('click', (e) => {
            e.preventDefault();
            if (paginaDashboard) paginaDashboard.classList.remove('hidden');
            if (paginaDados) paginaDados.classList.add('hidden');
            if (linkContas) linkContas.className = classeAtiva;
            if (linkDados) linkDados.className = classeInativa;
          });
        }

        // Clique em "Seus dados"
        if (linkDados) {
          linkDados.addEventListener('click', (e) => {
            e.preventDefault();
            if (paginaDashboard) paginaDashboard.classList.add('hidden');
            if (paginaDados) paginaDados.classList.remove('hidden');
            if (linkContas) linkContas.className = classeInativa;
            if (linkDados) linkDados.className = classeAtiva;
          });
        }
        // --- Fim da Lógica de Navegação ---

      };
    </script>
</body>
</html>


