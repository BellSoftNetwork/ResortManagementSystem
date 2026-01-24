# Multi-stage Dockerfile for frontend-web
# Build arguments for multi-platform support
ARG TARGETPLATFORM
ARG BUILDPLATFORM

FROM node:20-alpine AS base

# Install base packages
RUN apk add --no-cache git bash

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json yarn.lock ./

# Install dependencies
RUN --mount=type=cache,target=/root/.yarn \
    yarn install --frozen-lockfile

# Development stage
FROM base AS development

# Copy source code
COPY . .

# Expose port
EXPOSE 9000

# Run development server
CMD ["yarn", "dev"]

# Builder stage
FROM base AS builder

# Copy source code
COPY . .

# Build the application
RUN --mount=type=cache,target=/root/.yarn \
    --mount=type=cache,target=/app/node_modules/.cache \
    yarn build

# Production stage
FROM nginx:alpine AS production

# Copy built files from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Run nginx
CMD ["nginx", "-g", "daemon off;"]
