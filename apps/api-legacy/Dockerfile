# Multi-stage Dockerfile for api-legacy
# Build arguments for multi-platform support
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Use Amazon Corretto for better multi-platform support
FROM amazoncorretto:17-alpine AS base

# Install base packages
RUN apk add --no-cache bash

# Set working directory
WORKDIR /app

# Builder stage
FROM base AS builder

# Set working directory for builder stage
WORKDIR /app/apps/api-legacy

# Copy gradle files
COPY gradlew .
COPY gradle gradle
COPY build.gradle.kts .
COPY settings.gradle.kts .

# Copy source code
COPY . .

# Build the application
RUN ./gradlew bootJar --no-daemon -PbootJar.archiveFileName=app.jar
RUN ls -lR /app

# Development stage
FROM base AS development

# Install gradle
RUN apk add --no-cache gradle

# Expose ports
EXPOSE 8080 5005

# Default command (can be overridden)
CMD ["./gradlew", "bootRun", "--no-daemon"]

# Production stage
FROM amazoncorretto:17-alpine AS production

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy JAR from builder
COPY --from=builder /app/apps/api-legacy/build/libs/app.jar app.jar

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variable
ENV SPRING_PROFILES_ACTIVE=production
ENV JAVA_OPTS="-Xmx512m -Xms256m"

# Run the application
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]