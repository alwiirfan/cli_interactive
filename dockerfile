FROM golang:1.20.4-alpine AS builder

LABEL author=Mohammad_Alwi_irfani
LABEL email=alwi.irfani1927@gmail.com 
LABEL github=https://github.com/alwi09/cli_interactive

RUN apk update && apk add --no-cache git

WORKDIR /home

COPY go.mod go.sum ./

RUN go mod tidy

RUN go mod download

COPY . .

RUN go build -o cli ./cmd/main.go

FROM alpine:3.15

RUN apk update && apk add --no-cache git

WORKDIR /home

COPY --from=builder /home/cli .
COPY --from=builder /home/internal/database/migrations/ ./internal/database/migrations/

CMD ["./cli", "worker"]