# Stage 1: Build the SPA
FROM node:20-alpine AS builder

WORKDIR /build

# Copy core module first (must be built before web-client)
COPY nkrypt-xyz-core-nodejs/ ./nkrypt-xyz-core-nodejs/
RUN cd nkrypt-xyz-core-nodejs && npm ci && npm run build

# Copy web-client (needs devDependencies for build)
COPY web-client/package*.json ./web-client/
RUN cd web-client && npm ci

COPY web-client/ ./web-client/

# Build args for configurable server URL (baked in at build time)
ARG VITE_DEFAULT_SERVER_URL=http://localhost:9041
ARG VITE_CLIENT_APPLICATION_NAME=nkrypt-web-client
ENV VITE_DEFAULT_SERVER_URL=$VITE_DEFAULT_SERVER_URL
ENV VITE_CLIENT_APPLICATION_NAME=$VITE_CLIENT_APPLICATION_NAME

RUN cd web-client && npm run build:core && npm run build-spa

# Stage 2: Serve with Nginx
FROM nginx:alpine

# Copy built SPA from builder
COPY --from=builder /build/web-client/dist/spa /usr/share/nginx/html

# SPA with hash router - nginx config for client-side routing
RUN echo 'server { \
    listen 80; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
    location /healthz { \
        return 200 "ok"; \
        add_header Content-Type text/plain; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/healthz || exit 1
