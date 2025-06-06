# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

This is a monorepo containing a resort management system with:
- **Backend API**: Kotlin Spring Boot application (`apps/api-legacy`)
- **Frontend Web**: Vue.js/Quasar SPA (`apps/frontend-web`) 
- **Mobile App**: Capacitor-based Android app (built from frontend)

The backend follows domain-driven design with modules for authentication, reservations, rooms, users, payments, and revision tracking.

## Development Commands

### Backend (Kotlin Spring Boot)
```bash
# Run application locally (profile: local)
./gradlew bootRun

# Run tests  
./gradlew test

# Generate test coverage report
./gradlew jacocoTestReport

# Lint check
./gradlew ktlintCheck

# Build Docker image
./gradlew bootBuildImage
```

### Frontend (Vue.js/Quasar)
```bash
cd apps/frontend-web

# Install dependencies
yarn

# Development server (http://localhost:9000)
yarn dev

# Build for production
yarn build

# Lint and format
yarn lint
yarn format
```

### Mobile (Android)
```bash
cd apps/frontend-web

# Add Android platform (first time only)
yarn cap:add-android

# Build and sync to Android
yarn android:build

# Open in Android Studio
yarn cap:open-android

# Build and run (different environments)
yarn android:run        # Local emulator
yarn android:run:dev    # Development API
yarn android:run:prod   # Production API
```

## Environment Setup

### Required Environment Variables (Backend)
- `DATABASE_MYSQL_HOST`: MySQL server IP
- `DATABASE_MYSQL_USER`: MySQL username  
- `DATABASE_MYSQL_PASSWORD`: MySQL password
- `REDIS_HOST`: Redis server IP

### Database Migration
Database schema is managed with Liquibase. Changelog files are in `apps/api-legacy/src/main/resources/db/changelog/`. 

### Testing Framework
Backend uses Kotest for testing. Test classes follow the pattern `*Test.kt` and fixtures are in `*Fixture.kt` files.

## Code Quality Standards

### Backend
- JaCoCo coverage requirements:
  - Overall: 30% minimum
  - Business logic services: 80% line coverage, 90% branch coverage
  - Max 200 lines per class
- Ktlint for code formatting
- MapStruct for entity-DTO mapping

### Frontend  
- ESLint + Prettier for code formatting
- TypeScript for type safety
- Pinia for state management
- Quasar components for UI

## Key Technology Stack

### Backend
- Kotlin + Spring Boot 3
- Spring Security with JWT
- JPA + QueryDSL for data access
- MySQL + Redis
- Liquibase for migrations
- MapStruct for mapping

### Frontend
- Vue 3 + TypeScript
- Quasar Framework (Material Design)
- Pinia (state management)
- Capacitor (mobile app platform)
- ApexCharts for data visualization

## Development Workflow

1. Backend runs on port 8080 (local profile)
2. Frontend dev server runs on port 9000 with proxy to backend
3. Mobile apps use different API endpoints based on environment
4. Git pre-commit hooks enforce ktlint formatting on backend
