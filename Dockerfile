# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /src/app/frontend

COPY app/frontend/package.json app/frontend/package-lock.json ./
RUN npm ci

COPY app/frontend ./
ARG PUBLIC_BASE_URL=http://localhost:3003
ENV PUBLIC_BASE_URL=${PUBLIC_BASE_URL}
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.22-alpine AS backend-builder
WORKDIR /src/app

# Install git for go mod download if needed
RUN apk add --no-cache git

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app .
# Copy built frontend assets to the expected location for embedding
COPY --from=frontend-builder /src/app/frontend/dist ./frontend/dist

# Build the binary
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 go build -o /bin/coveritup main.go

# Stage 3: Final Runtime Image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS calls and curl for healthcheck/debugging
RUN apk add --no-cache ca-certificates curl

# Copy binary
COPY --from=backend-builder /bin/coveritup /bin/coveritup

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port
EXPOSE 3003

ENTRYPOINT ["/entrypoint.sh"]
