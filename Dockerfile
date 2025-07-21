FROM golang:1.24.4-alpine AS builder

RUN apk --no-cache add make gcc musl-dev

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/marketplace cmd/marketplace/main.go


FROM alpine

COPY --from=builder /app/bin/marketplace /
COPY --from=builder configs/config.yaml configs/config.yaml
COPY --from=builder internal/db/migrations internal/db/migrations

EXPOSE 8080

CMD ["/marketplace"]