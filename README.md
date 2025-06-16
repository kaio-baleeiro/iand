# iand

API para automação da configuração de ambientes de desenvolvimento Android.

## Descrição
O projeto `iand` é uma API escrita em Go, seguindo Clean Architecture, que automatiza tarefas de configuração de ambientes de desenvolvimento, manipulação de sistema de arquivos e execução de comandos no sistema operacional. Não utiliza banco de dados, focando em operações diretas no SO.

## Estrutura do Projeto
- `cmd/` – Ponto de entrada da aplicação.
- `internal/handlers/` – Handlers HTTP (camada de entrega).
- `internal/service/` – Lógica de negócio.
- `internal/repository/` – Abstração de acesso a dados (não utilizado neste projeto).
- `internal/domain/` – Entidades e regras de domínio.
- `logger/` – Logger estruturado.

## Tecnologias
- Go (Golang)
- net/http (API)
- Clean Architecture

## Como rodar
```sh
go run iand-api/main.go
```

## Contribuição
Pull requests são bem-vindos. Para mudanças maiores, abra uma issue primeiro para discutir o que você gostaria de mudar.

## Licença
Veja o arquivo LICENSE para mais detalhes.
