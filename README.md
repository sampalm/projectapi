## Iniciando o app

Renomei o .env-example para .env e coloque suas credenciais

```
$ mv .env-example .env
```
As migrations sao geradas automaticamente.

## Compilando o app

Voce pode utilizar o arquivo [Makefile](#) para gerar um ambiente de desenvolvimento com o seguinte comando:

```
$ make dev
```

Sera gerado o app e um banco de dados postgres. Caso preferir configurar o seu proprio ambiente sera necessario adicionar as devidas credenciais no .env com a sua configuraçao.

## Testando

Antes de executar os testes e necessario criar um arquivo test.env com os parametros de um ambiente de teste, utilize o mesmo exemplo do .env-example.

Para realizar um teste em todas as unidades execute o comando:

```
$ go test -v ./tests/*
```

## Importando a Collection do Postman (API's)

Download [Postman](https://www.getpostman.com/) -> Import -> Import From File

Selecione o arquivo `v1.postman_collection.json`

Includes the following:

- Auth
  - Get Token
- Project
  - All
  - One
  - Update
  - Destroy
  - Create
- Api
  - All
  - One
  - Update
  - Destroy
  - Create

> Para realizar uma requisiçao e necessario primeiro gerar um token, e utilizar como Authorization, a requisicao para o endpoint Get Token ira te retornar o seguinte response:

```
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQxMTM1NTUsImlhdCI6MTY0NDEwNjM1NSwiaXNzIjoicHJvamVjdC1hcGkifQ.vhj8ATGX86IFgXlqebkaYDeYpTE1ngcEulnmH3dGvC0"
```

Apos gerar o token voce pode adicionar o token na collection para facilitar o uso nas chamadas da collection.

    Authorization -> Bearer Token with value of {{token}}

E extremamente util para testar as rotas da API sem muito trabalho.