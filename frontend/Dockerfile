# =========================================================================
# Sizing Frontend Dockerfile (Nuxt 3 - OpenShift & Kubernetes Ready)
# Multi-stage build for Bank Mandiri Architecture & Sizing Hub
# Compatible with OpenShift 4.x Restricted-v2 SCC (Arbitrary UID / GID 0)
# =========================================================================

# Stage 1: Build Nuxt 3 application
FROM node:22-alpine AS builder

WORKDIR /app

# Enable corepack for pnpm package manager (matching lockfile 9.0)
RUN corepack enable && corepack prepare pnpm@9 --activate

# Cache dependencies
COPY package.json pnpm-lock.yaml* pnpm-workspace.yaml* ./
RUN pnpm install --frozen-lockfile

# Copy frontend source code
COPY . .

# Build Nuxt production bundle
RUN pnpm build

# Stage 2: Minimal runtime image
FROM node:22-alpine AS runner

LABEL maintainer="Enterprise Architecture & Sizing Team Bank Mandiri" \
      app="sizing-frontend" \
      description="Enterprise Architecture Sizing & Capacity Hub Frontend"

WORKDIR /app

# OpenShift Restricted SCC security requirements:
# Directories owned by group 0 (root) with group read/write/execute permissions
RUN apk add --no-cache curl tzdata && \
    mkdir -p /app/.output && \
    chown -R 1001:0 /app && \
    chmod -R g=u /app

# Copy production bundle from builder
COPY --from=builder --chown=1001:0 /app/.output /app/.output

ENV NODE_ENV=production \
    NITRO_HOST=0.0.0.0 \
    NITRO_PORT=3000 \
    PORT=3000 \
    TZ=Asia/Jakarta

EXPOSE 3000

# OpenShift non-root user
USER 1001

CMD ["node", ".output/server/index.mjs"]
