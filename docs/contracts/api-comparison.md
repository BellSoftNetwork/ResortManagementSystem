# API Endpoint Comparison: api-legacy vs api-core

## Summary
This document compares API endpoints between api-legacy (Spring Boot) and api-core (Golang + Gin) to track migration progress.

## ✅ Fully Implemented APIs

### Authentication APIs
- ✅ POST `/api/v1/auth/register`
- ✅ POST `/api/v1/auth/login`
- ✅ POST `/api/v1/auth/refresh`

### User/My Account APIs
- ✅ GET `/api/v1/my`
- ✅ POST `/api/v1/my`
- ✅ PATCH `/api/v1/my`

### Admin Account Management APIs
- ✅ GET `/api/v1/admin/accounts`
- ✅ POST `/api/v1/admin/accounts`
- ✅ PATCH `/api/v1/admin/accounts/{id}`

### Room APIs
- ✅ GET `/api/v1/rooms`
- ✅ GET `/api/v1/rooms/{id}`
- ✅ POST `/api/v1/rooms`
- ✅ PATCH `/api/v1/rooms/{id}`
- ✅ DELETE `/api/v1/rooms/{id}`
- ✅ GET `/api/v1/rooms/{id}/histories`

### Room Group APIs
- ✅ GET `/api/v1/room-groups`
- ✅ GET `/api/v1/room-groups/{id}`
- ✅ POST `/api/v1/room-groups`
- ✅ PATCH `/api/v1/room-groups/{id}`
- ✅ DELETE `/api/v1/room-groups/{id}`

### Reservation APIs
- ✅ GET `/api/v1/reservations`
- ✅ GET `/api/v1/reservations/{id}`
- ✅ POST `/api/v1/reservations`
- ✅ PATCH `/api/v1/reservations/{id}`
- ✅ DELETE `/api/v1/reservations/{id}`
- ✅ GET `/api/v1/reservations/{id}/histories`

### Reservation Statistics APIs
- ✅ GET `/api/v1/reservation-statistics`

### Payment Method APIs
- ✅ GET `/api/v1/payment-methods`
- ✅ GET `/api/v1/payment-methods/{id}`
- ✅ POST `/api/v1/payment-methods`
- ✅ PATCH `/api/v1/payment-methods/{id}`
- ✅ DELETE `/api/v1/payment-methods/{id}`

### Main/System APIs
- ✅ GET `/api/v1/env`
- ✅ GET `/api/v1/config`

## ❌ Missing APIs in api-core

### Documentation APIs
- ❌ GET `/docs/schema` - OpenAPI schema (JSON/YAML)
- ❌ GET `/docs/swagger-ui` - Swagger UI interface

### Health Check APIs (Spring Boot Actuator)
- ❌ GET `/actuator/health` - Basic health check
- ❌ GET `/actuator/health/liveness` - Kubernetes liveness probe
- ❌ GET `/actuator/health/readiness` - Kubernetes readiness probe

## 📋 Implementation Priority

### High Priority (Required for Production)
1. **Health Check Endpoints** - Critical for Kubernetes deployment
   - `/actuator/health`
   - `/actuator/health/liveness`
   - `/actuator/health/readiness`

### Medium Priority (Nice to Have)
1. **API Documentation**
   - OpenAPI/Swagger schema endpoint
   - Swagger UI interface

## 🔍 Response Format Verification Needed

While the endpoints are implemented, we need to verify that response formats match exactly:

1. **Pagination Response Format**
   - Ensure page metadata structure matches
   - Verify sort parameter handling

2. **Error Response Format**
   - Validate error response structure
   - Check field validation error format

3. **Empty Response Handling**
   - Verify empty arrays vs null
   - Check for missing vs empty JSON fields

4. **Date/Time Format**
   - Ensure consistent timezone handling
   - Validate timestamp format

## 🧪 Testing Strategy

1. Create comprehensive integration tests comparing responses between api-legacy and api-core
2. Use actual HTTP requests to both APIs with identical inputs
3. Compare response bodies using JSON diff
4. Test edge cases:
   - Empty results
   - Invalid inputs
   - Authentication failures
   - Permission denied scenarios