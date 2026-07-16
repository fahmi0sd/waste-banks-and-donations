FROM golang:1.25.5 AS builder
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/app/echo-server
RUN CGO_ENABLED=0 go build -o banksampah-api main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl gnupg && \
    install -d /usr/share/postgresql-common/pgdg && \
    curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc && \
    echo "deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] https://apt.postgresql.org/pub/repos/apt bookworm-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends postgresql-client-17 && \
    apt-get purge -y curl gnupg && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/app/echo-server/banksampah-api /banksampah-api

ENTRYPOINT ["/banksampah-api"]
