# Sistema de Compra de Passagens Aéreas
Descrição do problema disponível em: [TEC502 - Problema 1 - Venda de Passagens.pdf](https://github.com/user-attachments/files/16853459/TEC502.-.Problema.1.-.Venda.de.Passagens.-.Versao.2.pdf)

## Introdução

Este projeto foi desenvolvido para implementar a comunicação entre cliente e servidor no contexto de compra e venda de passagens aéreas no setor de aviação de baixo custo (low-cost carriers - LCC). O sistema torna possível o processo de compra de passagens, utilizando o protocolo TCP/IP e uma API desenvolvida em Go, com suporte para múltiplas conexões simultâneas. A aplicação está contida em containers Docker, que isolam e orquestram a execução dos serviços.

## Metodologia

A metodologia abordada foi a partir da construção de fluxogramas e diagramas de sequência da comunicação e fluxo do programa, juntamente com sessões de discussão e desenvolvimento em grupo sobre o problema. Serão abordados nesta sessão tópicos de alta importância para o desenvolvimento do projeto. 

### Arquitetura do Projeto

A arquitetura é dividida em dois principais componentes: **servidor** e **cliente**, ambos encapsulados em containers Docker, que se comunicam através de uma rede interna.

#### Servidor

O servidor é responsável pelas seguintes funcionalidades:

- Cadastro de usuários
- Validação e registro de compra de passagens
- Consulta de rotas disponíveis
- Persistência de dados em arquivo JSON

O servidor escuta requisições na porta TCP configurada, exposta via Docker.

#### Cliente

O cliente permite aos usuários:

- Consultar rotas disponíveis
- Comprar passagens
- Consultar compras anteriores

O cliente se conecta ao servidor via endereço IP do container (configurado no Docker Compose), enviando solicitações seguindo um protocolo de comunicação específico.

### Paradigma de Comunicação

O sistema utiliza o paradigma de comunicação **síncrono** com o protocolo TCP. O cliente envia uma requisição ao servidor, que a processa e retorna uma resposta. A comunicação é orientada a conexão, garantindo entrega confiável dos pacotes, sem perdas ou duplicações, o que é crucial para transações de compra de passagens.

A aplicação é **stateless**, ou seja, os dados de conexão não são armazenados, e cada requisição é tratada como uma nova sessão.

### Protocolo de Comunicação

A comunicação entre cliente e servidor segue um protocolo baseado em mensagens **JSON**, que encapsulam as requisições e respostas trocadas. O formato JSON foi escolhido por ser leve, legível e amplamente utilizado em sistemas de rede.

As principais operações suportadas são:

- **ROTAS**: Solicitação da lista de rotas disponíveis, com resposta contendo as informações de rotas e vagas.
- **COMPRA**: Solicitação de compra de um trecho de rota. O servidor valida a solicitação, ajusta o número de vagas e persiste a compra.
- **CADASTRO**: Cadastro de um novo usuário no sistema.
- **LERCOMPRAS**: Solicitação da lista de compras do usuário.

### Formatação e Tratamento de Dados

O JSON foi escolhido como formato de dados por sua simplicidade, flexibilidade, leiturabilidade e ampla adoção. Essas características tornam-o ideal para a troca de dados entre diferentes sistemas e linguagens de programação, garantindo a interoperabilidade e a facilidade de desenvolvimento.
Sendo assim, as mensagens estão sendo enviadas neste formato e seguem o seguinte padrão:
<p align="center">
    <img src="img/dadosjson.png" />
</p>

Ademais, para garantir a persistência dos dados, as informações sobre usuários, rotas e compras também foram armazenadas em arquivos JSON, devido a simplicidade na leitura e escrita.


### Tratamento de Conexões Simultâneas

O servidor foi projetado para suportar múltiplas conexões simultâneas, utilizando **goroutines**, que são leves e eficientes para concorrência no Go. Para cada nova conexão, uma goroutine é criada para processar as mensagens do cliente, permitindo que o servidor continue aceitando outras conexões em paralelo.

O uso de goroutines garante que o servidor atenda a múltiplos clientes simultaneamente, sem impactar a experiência dos demais usuários.

### Tratamento de Concorrência

Para gerenciar o acesso concorrente à função que valida as compras e gerencia as vagas, foi implementado o uso de um mutex (mutual exclusion). O mutex permite que apenas uma goroutine acesse a função ao mesmo tempo, bloqueando outras goroutines até que a primeira complete sua execução. Esse controle é crucial, especialmente em cenários onde múltiplos usuários podem tentar comprar o mesmo trecho de uma rota que possui vagas limitadas.

A utilização do mutex foi realizada da seguinte maneira:

**Bloqueio do Mutex:** Antes de executar a lógica de validação, a goroutine chama mu.Lock(), que impede que outras goroutines entrem na função enquanto uma instância dela estiver em execução.

**Operação Crítica:** A lógica de validação e atualização das vagas disponíveis ocorre dentro da seção crítica, garantindo que os dados não sejam corrompidos por acessos simultâneos.

**Desbloqueio do Mutex:** Após a conclusão da operação, mu.Unlock() é chamado para liberar o mutex, permitindo que outras goroutines possam acessar a função.

Com essa abordagem, é possível garantir que a validação de compras ocorra de maneira segura e eficiente, evitando problemas como a venda excessiva de vagas e mantendo a integridade do sistema como um todo.

--- 

