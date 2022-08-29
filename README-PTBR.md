[![Go Reference](https://pkg.go.dev/badge/github.com/abu-lang/goabu.svg)](https://pkg.go.dev/github.com/bancodobrasil/featws-ruller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/bancodobrasil/featws-ruller/blob/develop/LICENSE)

# FeatWS API [![About_en](https://github.com/yammadev/flag-icons/blob/master/png/US.png?raw=true)](https://github.com/bancodobrasil/featws-api/blob/develop/README.md)

API Rest do projeto FeatWS, que é responsavél por interligar a UI com o banco de dados.

## Dependências

- Docker compose
- Go 1.18

## Go Modules

``` bash
go mod tidy
```

## Desenvolvimento Local

### Env

- Copiar arquivo `.env.sample` e renomear para `.env`.
- Adicione seu Token GitLab em `FEATWS_API_GITLAB_TOKEN`.

1.Entrar no GitLab: <https://gitlab.com/>.
2. Entre no `Edit Profile`.
3. Clique em `Access Token`.
4. Coloque escolha um token name, selecione os scopes e data de expiração. Em seguida, clique no botão abaixo.
5. Em FEATWS_API_GITLAB_TOKEN coloque o token gerado em Access Token.
6. Crie um grupo, aonde ficarão as regras. Em `FEATWS_API_GITLAB_NAMESPACE` digite o nome do grupo criado.

``` bash
make up-services
```
