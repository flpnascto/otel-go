# Sistema em Go com Observabilidade e Telemetria

## :notebook_with_decorative_cover: Sobre o Projeto

Aplicação desenvolvida no curso **Pós Graduação Go Expert - Full Cycle** na linguagem Go.
O sistema recebe uma requisição de CEP, que após ser validado consulta uma API externa para coletar informações de endereço.
Com a informação da cidade referente ao CEP, é realizado uma consulta em outra API externa para coletar informações sobre o clima.
E por fim é retornado a cidade e as temperaturas da localidade em Célsius, Fahrenheit e Kelvin.
O sistema é composto de dois microsserviços
  - **Input**: responsável por receber um input POST com o schema `{ cep: string}` e realizar uam requisição no microsserviço **Climate**.
  - **Climate**: responsável por consultar API externas e retornar o resultado

## :wrench: Requisitos

### :sparkles: Funcionalidades
- Através do CEP retorna a temperatura de uma cidade
- Observabilidade
- Telemetria

### :computer: Tecnologias Aplicadas
* Go
* Net/Http
* Docker
* Open Telemetry (OTEL) + Zipkin
* Jaeguer
* Prometheus com Grafana

## :arrow_forward: Executando a aplicação

Para executar a aplicação siga as instruções abaixo.

### Pré-requisitos

Primeiramente é necessário que possua instalado as seguintes ferramentas: Go, Git, Docker.
Além disto é bom ter um editor para trabalhar com o código como VSCode.

### Instalação

1. Faça uma cópia do repositório (HTTPS ou SSH)
   ```sh
   git clone https://github.com/flpnascto/otel-go.git
   ```
   ```sh
   git clone git@github.com:flpnascto/otel-go.git
   ```
2. Acessar a pasta do repositório local e inicie o docker para instanciar as aplicações
  ```sh
   docker compose up -d
   ```

### Aplicações

- Prometheus: pode ser acessado via http://localhost:9090
- Grafana: pode ser acessado via http://localhost:3000 (user: admin, password: admin)
- JAEGER: pode ser acessado via http://localhost:16686
- O microsserviço **Input** é executado em http://localhost:8081
- O microsserviço **Climate** é executado em http://localhost:8080

### Realizando requisições

Existe na raiz do projeto um arquivo `api.http` com exemplos de requisições. Para utilizá-lo é necessário a extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) instalada no VSCode

