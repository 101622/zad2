# syntax=docker/dockerfile:1

# Użycie języka golang jako buildera
FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./

# Wyłączamy CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o weatherapp main.go

FROM scratch

LABEL org.opencontainers.image.authors="Wojciech Makowka"

# Kopia certyfikatów buildera
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/weatherapp /weatherapp

# Port 8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s CMD ["/weatherapp", "-health"]

ENTRYPOINT ["/weatherapp"]