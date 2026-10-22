# =========================================================================
# Sizing Backend Dockerfile (OpenShift & Kubernetes Enterprise Ready)
# Multi-stage build for Enterprise Architecture Architecture & Sizing Hub
# Compatible with OpenShift 4.x Restricted-v2 SCC (Arbitrary UID / GID 0)
# =========================================================================

# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates tzdata

ENV GOTOOLCHAIN=auto

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy backend source
COPY . .

# Compile static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o sizing-backend ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.21 AS runner

LABEL maintainer="Enterprise Architecture & Sizing Team Enterprise Architecture" \
      app="sizing-backend" \
      description="Enterprise Architecture Sizing & Capacity Engine"

WORKDIR /app

# Install security certificates, timezone data, and curl for health probes
RUN apk add --no-cache ca-certificates tzdata curl && \
    mkdir -p /app/templates && \
    chown -R 1001:0 /app && \
    chmod -R g=u /app

# Copy binary from builder
COPY --from=builder --chown=1001:0 /build/sizing-backend /app/sizing-backend

# Copy sizing & capacity templates
COPY --chown=1001:0 templates /app/templates

ENV PORT=8089 \
    TEMPLATES_DIR=/app/templates \
    TZ=Asia/Jakarta

EXPOSE 8089

# OpenShift Restricted SCC uses arbitrary UID belonging to group 0 (root)
USER 1001

ENTRYPOINT ["/app/sizing-backend"]
