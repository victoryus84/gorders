# 1) Build stage
# Folosim varianta alpine și pentru builder pentru consistență și viteză
FROM golang:1.26-alpine AS builder

# Instalăm git doar dacă este necesar pentru module private (folosind apk)
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

# OPTIMIZARE: Folosim un mount cache pentru a nu descărca modulele la fiecare build
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# OPTIMIZARE: Folosim mount cache și pentru build-ul propriu-zis
# Go 1.26 beneficiază enorm de cache-ul de compilare refolosit
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server ./cmd/server

# 2) Runtime stage
# Am updatat alpine la o versiune de început de 2026
FROM alpine:3.23 

# Pasul 1: Instalăm pachetul pentru fus orar
RUN apk add --no-cache tzdata

# Pasul 2: Setăm variabila de mediu pentru Moldova
ENV TZ=Europe/Chisinau

WORKDIR /app

# Luăm binarul din builder
COPY --from=builder /app/server /app/server

# Securitate: Creăm utilizatorul fără privilegii
RUN adduser -D appuser
USER appuser

ENV PORT=8080
EXPOSE ${PORT}

HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-check-certificate --quiet --tries=1 --spider http://localhost:${PORT}/health || exit 1

ENTRYPOINT ["/app/server"]