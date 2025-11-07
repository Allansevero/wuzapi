<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Planos - Metrizap</title>
    
    <!-- Importação do Tailwind CSS --><script src="https://cdn.tailwindcss.com"></script>
    
    <!-- Importação das fontes --><link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Krona+One&display=swap" rel="stylesheet">
    
    <!-- Configuração do Tailwind --><script>
      tailwind.config = {
        theme: {
          extend: {
            fontFamily: {
              'sans': ['Inter', 'sans-serif'],
              'krona': ['"Krona One"', 'sans-serif'],
            },
            colors: {
              'mz-green': '#28a745',
              'mz-green-light': '#e9f7ec',
              'mz-orange-light': '#fff3eb', /* Cor para o plano Essencial */
              'mz-blue-light': '#e0f2fe', /* Cor para o plano PRO */
            }
          }
        }
      }
    </script>
</head>
<body class="bg-gray-50 font-sans min-h-screen">

    <!-- Cabeçalho Superior --><header class="flex items-center justify-between p-6 bg-white border-b border-gray-200">
        <a href="#" class="text-gray-600 hover:text-gray-800 transition-colors">
            <!-- Ícone SVG: Seta para a Esquerda --><svg class="w-7 h-7" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="19" y1="12" x2="5" y2="12"></line>
                <polyline points="12 19 5 12 12 5"></polyline>
            </svg>
        </a>
        <img class="w-10 h-10 rounded-full" src="https://placehold.co/40x40/E2E8F0/4A5568?text=A" alt="Avatar do usuário">
    </header>

    <!-- Conteúdo Principal --><main class="container mx-auto p-8 lg:p-12">
        
        <!-- Títulos da Seção --><div class="max-w-4xl mx-auto text-center mb-12">
            <h1 class="text-4xl lg:text-5xl font-bold text-gray-800 mb-4">
                Analise, venda mais e <br> cresce ainda mais
            </h1>
            <p class="text-lg text-gray-600">
                Confira os planos e escolha o que atende suas necessidades hoje.
            </p>
        </div>

        <!-- Cards de Planos --><div class="grid grid-cols-1 lg:grid-cols-2 gap-8 max-w-5xl mx-auto">

            <!-- Card: ESSENCIAL (Gratuito) --><div class="bg-white p-8 rounded-xl border border-gray-200 shadow-sm flex flex-col relative overflow-hidden">
                <!-- Gradiente de cor de fundo (TOP) --><div class="absolute top-0 left-0 w-full h-24 bg-gradient-to-b from-mz-orange-light to-transparent -z-10 rounded-t-xl"></div>
                
                <h3 class="text-sm font-semibold text-gray-500 uppercase mb-1">ESSENCIAL</h3>
                <p class="text-xs text-gray-400 mb-4">Para quem tem poucos WhatsApp</p>
                <p class="text-5xl font-bold text-mz-green mb-8">Gratuito</p>

                <h4 class="text-md font-semibold text-gray-700 mb-4">Funcionalidades essenciais:</h4>
                <ul class="space-y-3 text-sm text-gray-700 flex-grow">
                    <li class="flex items-center space-x-2">
                        <!-- Ícone SVG: Check --><svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Até 2 WhatsApp conectados</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Até 2 WhatsApp conectados</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span class="font-semibold">Análise completa do mês</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span class="font-semibold">Análise diária das conversas no WhatsApp</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Radar Perfomatico de vendas</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>1 número para receber análises</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>1 usuário no sistema</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>5 mensagens agendadas por número mês</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Suporte por e-mail</span>
                    </li>
                </ul>
                
                <button class="mt-8 w-full bg-gray-200 text-gray-500 font-medium py-4 px-6 rounded-lg cursor-not-allowed">
                    Seu plano atual
                </button>
            </div>

            <!-- Card: PRO (Pago) --><div class="bg-white p-8 rounded-xl border border-mz-green shadow-lg flex flex-col relative overflow-hidden">
                <!-- Gradiente de cor de fundo (TOP) --><div class="absolute top-0 left-0 w-full h-24 bg-gradient-to-b from-mz-blue-light to-transparent -z-10 rounded-t-xl"></div>

                <h3 class="text-sm font-semibold text-gray-500 uppercase mb-1">PRO</h3>
                <p class="text-xs text-gray-400 mb-4">Para quem está em crescimento</p>
                <p class="text-gray-800 mb-6">
                    <span id="proPriceDisplay" class="text-5xl font-bold text-mz-green">R$47</span>
                    <span class="text-md text-gray-500">/cobrado ao mês</span>
                </p>
                
                <div class="relative mb-6">
                    <select id="whatsappSelect" class="w-full appearance-none bg-white border border-gray-300 rounded-lg px-4 py-3 text-sm font-medium text-gray-700 focus:outline-none focus:ring-2 focus:ring-mz-green pr-10">
                        <option value="47">Até 8 WhatsApp conectado mês</option>
                        <option value="97">Até 20 WhatsApps conectados mês</option>
                    </select>
                    <!-- Ícone SVG: Chevron Down --><div class="absolute inset-y-0 right-0 flex items-center px-3 pointer-events-none">
                        <svg class="w-5 h-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
                    </div>
                </div>

                <ul class="space-y-3 text-sm text-gray-700 flex-grow">
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span class="font-semibold">Análise completa do mês</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span class="font-semibold">Análise diária das conversas no WhatsApp</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Radar Perfomatico de vendas</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>3 número para receber análises</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Usuário ilimitado no sistema</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>30 mensagens agendadas por número mês</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Suporte por e-mail</span>
                    </li>
                    <li class="flex items-center space-x-2">
                        <svg class="w-5 h-5 text-mz-green flex-shrink-0" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"></path></svg>
                        <span>Gerente de conta direto no WhatsApp</span>
                    </li>
                </ul>
                
                <button class="mt-8 w-full bg-mz-green text-white font-semibold py-4 px-6 rounded-lg shadow-sm hover:bg-green-700 transition-colors">
                    Aumentar meus limites agora
                </button>
            </div>

        </div>

    </main>

    <!-- Script para lógica da página --><script>
        window.onload = () => {
            const whatsappSelect = document.getElementById('whatsappSelect');
            const proPriceDisplay = document.getElementById('proPriceDisplay');

            if (whatsappSelect && proPriceDisplay) {
                whatsappSelect.addEventListener('change', () => {
                    const newPrice = whatsappSelect.value;
                    proPriceDisplay.textContent = `R$${newPrice}`;
                });
            }
        };
    </script>
</body>
</html>