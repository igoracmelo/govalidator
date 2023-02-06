# Solução para a prova do processo seletivo do Studio Sol

## Rodando testes
```sh
go test ./...
```

## Rodando a aplicação não containerizada
```sh
go run main.go
```
Rodará em localhost:8080


Teste manual:
```sh
curl -i -X POST 'http://localhost:8080/verify' -d '{"password": "12345", "rules": [{"rule": "minSize", "value": 4}]}'
```
Deve retornar `{"verify":true,"noMatch":[]}`
