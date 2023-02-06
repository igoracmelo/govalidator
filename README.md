# Solução para a prova do processo seletivo do Studio Sol

Propositalmente deixei o meu .git para caso queira ver histórico de alterações.

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
go test ./...
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
docker run --rm -it validator:1.0.0 go test ./...
```

### Teste manual
```sh
curl -i -X POST 'http://localhost:8080/verify' -d '{"password": "12345", "rules": [{"rule": "minSize", "value": 4}]}'
```
Deve retornar `{"verify":true,"noMatch":[]}`