# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git bash

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o polygon-edge \
    -ldflags="-X 'github.com/0xPolygon/polygon-edge/versioning.Version=v1.0.0' \
              -X 'github.com/0xPolygon/polygon-edge/versioning.Branch=main'" \
    main.go

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.18

RUN apk add --no-cache ca-certificates bash

WORKDIR /app

# Copy binary
COPY --from=builder /app/polygon-edge ./polygon-edge

# Copy available genesis files. The Railway testnet genesis must be generated
# from fresh test-only validator keys before the service can start.
COPY genesis*.json ./

RUN mkdir -p /data

# Copy startup script
COPY start-railway.sh ./start-railway.sh
RUN chmod +x ./start-railway.sh ./polygon-edge

# P2P listens on this container port. Railway service networking must be
# configured separately; EXPOSE does not make the port public.
EXPOSE 30303/tcp

CMD ["./start-railway.sh"]
