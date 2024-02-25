# Estágio de construção
FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .
RUN go mod download && go mod verify

# Compilação do código
RUN go build -o /bin/app

# Estágio de produção
FROM scratch

COPY --from=builder /bin/app /bin/app

# Comando padrão para execução do aplicativo
CMD ["/bin/app"]
