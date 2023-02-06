# Solução para a prova do processo seletivo do Studio Sol

Propositalmente deixei o meu .git para caso queira ver histórico de alterações.

## Overview

A aplicação consiste em um servidor HTTP em Go usando pacotes nativos, como `net/http` e `encoding/json`.
Optei por usar pacotes nativos por eles terem suporte oficial e serem mais garantidos de receberem manutenção, além de atenderem bem às necessidades.

O servidor tem um endpoint com o caminho `/verify` que recebe um `POST` contendo a senha e regras de validação da mesma.
Os dados do request e response body são transmitidos utilizando streams para reduzir uso de memória, o que é crucial para serviços com alto número de requisições por segundo.

Poderia também ser implementada uma camada de cache server side ulizando um `map[string]string` e uma `Mutex`, ou utilizando um sistema externo de cache como Redis ou memcached e um cache no client side usando ETag e cache-control, mas fugiria um pouco do escopo inicial da aplicação.

## Rodando a aplicação não conteinerizada
```sh
go run main.go
```
Rodará em localhost:8080

Alternativamente:
```sh
go build -o app main.go
./app
```

### Teste automatizado
```sh
go test -v ./...
```

### Teste manual
```sh
curl -i -X POST 'http://localhost:8080/verify' -d '{"password": "12345", "rules": [{"rule": "minSize", "value": 4}]}'
```
Deve retornar `{"verify":true,"noMatch":[]}`

## Rodando a aplicação conteinerizada
```sh
docker build --tag validator:1.0.0
docker run --rm -p 8080:8080 validator:1.0.0
```

### Teste automatizado
```sh
docker run --rm -it validator:1.0.0 go test -v ./...
```

### Teste manual
```sh
curl -i -X POST 'http://localhost:8080/verify' -d '{"password": "12345", "rules": [{"rule": "minSize", "value": 4}]}'
```
Deve retornar `{"verify":true,"noMatch":[]}`