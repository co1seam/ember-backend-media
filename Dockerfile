FROM golang:1.24.1-bookworm AS base
WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

COPY . .

COPY internal/adapters/repository/migrations .

FROM base AS development

RUN apt-get install git bash curl

FROM development AS debug

FROM base AS production

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/cmd ./cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=production /app/bin/main /app/main

RUN adduser -D -u 1001 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 50051

CMD ["/app/main"]