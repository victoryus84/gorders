# 1) Build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# ATENȚIE: Dacă pornești build-ul din rădăcina proiectului (unde e compose-prod.yml)
# trebuie să specifici folderul "server/"
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copiem restul fișierelor din server
COPY . .

# Compilăm binarul. 
# Recomand să-i punem numele "server" ca să bată cu ce ai în docker-compose
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server/cmd/api ./cmd/api

# 2) Runtime stage
FROM alpine:3.23 

RUN apk add --no-cache tzdata
ENV TZ=Europe/Chisinau

WORKDIR /app

# REPARAȚIA: Numele de la sursă trebuie să fie EXACT cel de la -o (mai sus)
# Ambele trebuie să fie /app/server
COPY --from=builder /app/server/cmd/api /app/server

RUN adduser -D appuser
USER appuser

ENV PORT=8080
EXPOSE ${PORT}

HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-check-certificate --quiet --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Pornim binarul numit server
ENTRYPOINT ["/app/server"]