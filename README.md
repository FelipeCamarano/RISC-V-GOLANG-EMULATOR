# RISC-V Golang Emulator

[![Ask DeepWiki](https://devin.ai/assets/askdeepwiki.png)](https://deepwiki.com/felipecamarano/risc-v-golang-emulator)

Este repositório contém um emulador RISC-V (RV32I) escrito em Go. Ele simula de forma abrangente o hardware, incluindo CPU, barramento do sistema e vários tipos de memória (RAM, ROM, VRAM). O emulador possui uma interface gráfica construída com o toolkit Fyne, oferecendo visualização em tempo real e controle sobre o processo de emulação.

## Funcionalidades

* **Emulação RV32I:** Implementa um subconjunto principal do conjunto de instruções inteiras de 32 bits do RISC-V.
* **Arquitetura Modular:** Organizado em componentes de hardware distintos: CPU, Bus, RAM, ROM e VRAM.
* **GUI Interativa:** Interface baseada em Fyne para uma experiência amigável, com:
  * **Monitor de VRAM:** Display gráfico de 320x200 pixels renderizando o conteúdo da memória de vídeo.
  * **Visualizador de Registradores:** Exibição em tempo real dos 32 registradores de uso geral da CPU.
  * **Visualizador de Memória:** Visualizador em formato hexadecimal de todo o mapa de memória, incluindo RAM, VRAM, ROM de Cartucho e espaço de I/O.
  * **Terminal UART:** Terminal serial simples para saída de caracteres do sistema emulado.
* **Controle de Emulação:** Controles para executar, pausar e avançar instrução a instrução.
* **Velocidade de Clock Ajustável:** Um slider para controlar a velocidade de emulação.
* **Carregamento de Cartuchos:** Suporte ao carregamento e execução de programas externos a partir de arquivos binários (por exemplo, `bios.bin`, `hello_world.bin`).
* **Suíte de Testes Integrada:** Inclui uma janela dedicada para rodar a suíte `riscv-tests` (`rv32ui-p-*`) a fim de verificar a correção das instruções implementadas.

## Primeiros Passos

### Pré-requisitos

* Go 1.25 ou superior.
* Compilador C 64 bits compatível com GCC, necessário para compilar as dependências nativas da GUI (Fyne).
  * **Windows:** testado com `gcc (tdm64-1) 10.3.0`.

### Instalação e Execução

1. Clone o repositório:
  
      git clone https://github.com/felipecamarano/risc-v-golang-emulator.git
      cd risc-v-golang-emulator
  
2. Execute o emulador:
  
      go run main.go
  
  Isso abrirá a janela principal da GUI. O emulador inicia carregando o arquivo `bios.bin`.
  

## Uso

A janela principal é o ponto central para interagir com o emulador.

* **Insert Cartridge:** Clique neste botão para abrir um diálogo de arquivo e carregar um binário RISC-V customizado (por exemplo, `hello_world.bin` ou `paint_screen.bin` incluídos no repositório). O sistema será resetado e começará a execução a partir da BIOS, que então salta para o cartucho.
* **Run/Stop:** Alterna entre execução contínua do programa carregado e pausa. A velocidade é controlada pelo slider "Clock Speed".
* **Step:** Executa uma única instrução da CPU, permitindo depuração detalhada e observação do estado. A interface é atualizada a cada passo.
* **Debug:** Abre a janela da Suíte de Testes, onde você pode rodar testes automatizados para instruções específicas.

### Suíte de Testes

O emulador inclui uma suíte de testes para validar a implementação das instruções da CPU.

1. Clique no botão **Debug** na barra de ferramentas principal para abrir a janela **Suite de Testes RISC-V**.
2. Essa janela lista todos os programas de teste disponíveis no diretório `TESTES HEX RISCV/`.
3. Clique em **Executar Todos** para rodar todos os testes automaticamente. Os resultados (Pass/Fail), número de ciclos e valores finais de registradores são exibidos.
4. Para depurar um teste específico, clique no ícone de lupa (🔍) na linha correspondente. Isso carregará o programa de teste na janela principal do emulador, permitindo que você avance passo a passo na execução.

## Estrutura do Projeto

* `main.go`: Ponto de entrada da aplicação.
* `/cpu`: Implementação da CPU, incluindo lógica de busca, decodificação e execução de instruções.
* `/memory`: Implementa os dispositivos mapeados em memória, como RAM, ROM e VRAM.
* `/bus`: Define o barramento do sistema que conecta a CPU a todas as memórias e dispositivos de I/O.
* `/motherboard`: Configura o mapa de memória do sistema e integra todos os componentes de hardware.
* `/gui`: Interface gráfica baseada em Fyne, com todos os componentes e janelas.
* `/TESTES HEX RISCV`: Coleção de programas de teste `rv32ui-p-*` usados para verificar a correção das instruções.
* `*.bin`: Programas de exemplo pré-compilados (`bios.bin`, `hello_world.bin`, `paint_screen.bin`).

## Licença

Este projeto é licenciado sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.
