# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

This is a monorepo containing a resort management system with:
- **Backend API**: Kotlin Spring Boot application (`apps/api-legacy`)
- **Frontend Web**: Vue.js/Quasar SPA (`apps/frontend-web`)
- **Mobile App**: Capacitor-based Android app (built from frontend)
- **Performance Tests**: k6-based load testing (`tests/performance`)

### Backend Domain Structure
The backend follows domain-driven design (DDD) with a layered architecture. Each domain module (`authentication`, `reservation`, `room`, `user`, `payment`, `revision`) contains:
- `controller/`: API endpoints (interface + implementation pattern)
- `service/`: Business logic
- `repository/`: Data access with JPA + QueryDSL
- `entity/`: Domain models
- `dto/`: Data transfer objects
- `mapper/`: MapStruct mappers for entity-DTO conversion
- `exception/`: Domain-specific exceptions
- `converter/`: JPA type converters

### Frontend Structure
- `src/api/`: Backend API client modules
- `src/components/`: Reusable Vue components
- `src/pages/`: Route-based page components
- `src/stores/`: Pinia state management
- `src/schema/`: TypeScript data models
- `src/router/`: Vue Router configuration

## Development Commands

### Backend (Kotlin Spring Boot)
```bash
# Run application locally (profile: local, debug port: 5005)
./gradlew bootRun

# Run all tests
./gradlew test

# Run a single test class
./gradlew test --tests "net.bellsoft.rms.reservation.service.ReservationServiceTest"

# Run tests matching a pattern
./gradlew test --tests "*ReservationService*"

# Generate test coverage report (outputs to build/jacoco/jacoco.xml)
./gradlew jacocoTestReport

# Lint check
./gradlew ktlintCheck

# Auto-format code
./gradlew ktlintFormat

# Build Docker image
./gradlew bootBuildImage
```

### Frontend (Vue.js/Quasar)
```bash
cd apps/frontend-web

# Install dependencies
yarn

# Development server (http://localhost:9000, proxies API to :8080)
yarn dev

# Build for production
yarn build

# Lint check
yarn lint

# Auto-fix lint issues and format
yarn lintfix
yarn format
```

### Performance Tests (k6)
```bash
cd tests/performance

# Run reservation browsing scenario
yarn test:reservation-browsing

# Run load test
yarn test:reservation-load

# Run against local backend
yarn test:local

# Run with debug logging
yarn test:debug
```

### Mobile (Android)
```bash
cd apps/frontend-web

# First-time setup
yarn cap:add-android

# Build and sync
yarn android:build

# Open in Android Studio
yarn cap:open-android

# Build and run (different environments)
yarn android:run        # Local emulator (uses 10.0.2.2:8080)
yarn android:run:dev    # Development API
yarn android:run:prod   # Production API
```

## Environment Setup

### Required Environment Variables (Backend)
- `DATABASE_MYSQL_HOST`: MySQL server IP
- `DATABASE_MYSQL_USER`: MySQL username
- `DATABASE_MYSQL_PASSWORD`: MySQL password
- `REDIS_HOST`: Redis server IP

Optional:
- `JWT_SECRET`: JWT signing key
- `JWT_ACCESS_TOKEN_VALIDITY`: Access token validity in hours (default: 1)
- `JWT_REFRESH_TOKEN_VALIDITY`: Refresh token validity in hours (default: 720)

### Initial Setup
```bash
# Apply ktlint settings to IntelliJ
./gradlew ktlintApplyToIdea

# Add pre-commit hook for ktlint
./gradlew addKtlintCheckGitPreCommitHook
```

### Database Migration
Schema managed with Liquibase. Changelogs in `apps/api-legacy/src/main/resources/db/changelog/`.

### Testing Framework
- **Kotest** for BDD-style tests
- **kotlinfixture** for test data generation
- Each domain has a dedicated `*Fixture.kt` class in `src/test/kotlin/.../fixture/`
- Test classes: `*Test.kt`

## Code Quality Standards

### Backend
- JaCoCo coverage requirements:
  - Overall: 30% minimum
  - Business logic services (`*.service.*`): 80% line, 90% branch coverage
  - Max 200 lines per class (excluding generated code)
- Ktlint for formatting (pre-commit hook enforced)
- MapStruct for entity-DTO mapping (check generated code after creating new mappers)

### Frontend
- ESLint + Prettier for formatting
- TypeScript strict mode
- Vue 3 Composition API

## API Documentation
- Schema: `/docs/schema`
- Swagger UI: `/docs/swagger-ui`

## Development Workflow

1. Backend runs on port 8080 (local profile) with remote debug on 5005
2. Frontend dev server runs on port 9000 with proxy to backend
3. Mobile apps use different API endpoints based on environment
4. Git pre-commit hooks enforce ktlint formatting
