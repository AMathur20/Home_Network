# --- Stage 1: Build the UI ---
FROM node:20-alpine AS ui-builder
WORKDIR /app/ui
COPY ui/package*.json ./
RUN npm install
COPY ui/ ./
RUN npm run build

# --- Stage 2: Build the Go Backend ---
FROM golang:1.25-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o hnm-core ./cmd/hnm-core/main.go

# --- Stage 3: Final Image ---
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy artifacts from builders
COPY --from=backend-builder /app/hnm-core .
COPY --from=ui-builder /app/ui/dist ./ui/dist

# Default config & topology (optional, can be mounted)
RUN mkdir -p config config/topology

# Environment variables
ENV HNM_CONFIG_PATH=/app/config/config.yaml
ENV HNM_DB_PATH=/app/data/hnm.db
ENV HNM_UI_PATH=/app/ui/dist

EXPOSE 8080
ENTRYPOINT ["./hnm-core"]
