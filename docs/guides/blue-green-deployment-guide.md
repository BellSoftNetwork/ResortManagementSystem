# Blue-Green Deployment Guide: api-legacy to api-core

## Overview

This guide covers the blue-green deployment process for migrating from api-legacy (Kotlin + Spring Boot) to api-core (Go + Gin).

## Prerequisites

- [ ] api-core deployed and running
- [ ] api-legacy running (for rollback)
- [ ] Both backends using the same database
- [ ] Load balancer configured for traffic switching
- [ ] Monitoring/alerting set up

## Database Considerations

### Shared Database

Both api-core and api-legacy can connect to the same MySQL database simultaneously. No data migration is needed.

### Schema Compatibility

The following database aspects are verified compatible:

| Aspect | Details |
|--------|---------|
| Soft Delete | `deleted_at = '1970-01-01 00:00:00'` for active records |
| Timestamps | `YYYY-MM-DD HH:MM:SS` format |
| Password Storage | `{bcrypt}$2a$10$...` format |
| Foreign Keys | All constraints preserved |

### History/Audit Data

| Scenario | Behavior |
|----------|----------|
| Old history data (Envers) | Preserved, read-only via api-core |
| New changes via api-core | Logged to `audit_logs` table |
| New changes via api-legacy | Logged to `*_history` tables via Envers |

**Note**: During the transition period, audit data may be split between the two systems. This is expected and acceptable.

## Pre-Deployment Checklist

### 1. Optional Data Fixes

Fix invalid room status (if desired):

```sql
-- Fix room 208호 with invalid status value
UPDATE room SET status = 1 WHERE id = 44 AND status = 2;
```

### 2. Environment Variables

Ensure api-core has all required environment variables:

```bash
# Database
DATABASE_MYSQL_HOST=<production-mysql-host>
DATABASE_MYSQL_PORT=3306
DATABASE_MYSQL_USER=<user>
DATABASE_MYSQL_PASSWORD=<password>
DATABASE_MYSQL_DATABASE=<database-name>

# Redis
REDIS_HOST=<production-redis-host>
REDIS_PORT=6379

# JWT (must match api-legacy for token compatibility)
JWT_SECRET=<same-secret-as-api-legacy>
JWT_ACCESS_TOKEN_EXPIRY=900
JWT_REFRESH_TOKEN_EXPIRY=604800

# API
API_PORT=8080
API_PROFILE=production
```

### 3. JWT Token Compatibility

Both backends must use the same JWT secret for seamless user sessions during transition.

## Deployment Steps

### Step 1: Deploy api-core (Green)

1. Deploy api-core to production
2. Verify health endpoint: `GET /api/v1/env`
3. Run smoke tests against api-core directly

### Step 2: Switch Traffic (10%)

1. Configure load balancer to route 10% traffic to api-core
2. Monitor error rates and response times
3. Check application logs for errors

### Step 3: Gradual Rollout

| Phase | api-core | api-legacy | Duration |
|-------|----------|------------|----------|
| 1 | 10% | 90% | 1 hour |
| 2 | 25% | 75% | 2 hours |
| 3 | 50% | 50% | 4 hours |
| 4 | 75% | 25% | 4 hours |
| 5 | 100% | 0% | - |

### Step 4: Full Cutover

1. Route 100% traffic to api-core
2. Keep api-legacy running for 24-48 hours
3. Monitor for any issues

### Step 5: Decommission api-legacy

After stable operation:

1. Stop api-legacy instances
2. Archive api-legacy deployment
3. Update documentation

## Rollback Procedure

If issues are detected:

### Immediate Rollback

1. Switch load balancer to route 100% traffic to api-legacy
2. Investigate api-core issues
3. No data migration needed (shared database)

### Monitoring Triggers for Rollback

- Error rate > 1%
- Response time > 2x baseline
- Critical user complaints
- Data integrity issues

## Post-Deployment Verification

### API Endpoint Checks

```bash
# Health check
curl https://api.example.com/api/v1/env

# Config check
curl https://api.example.com/api/v1/config

# Authenticated endpoint (with valid token)
curl https://api.example.com/api/v1/my \
  -H "Authorization: Bearer <token>"
```

### Frontend Verification

1. Login with existing credentials
2. View reservations list
3. View rooms list
4. Create/edit/delete a test record
5. Verify statistics page loads

## Known Differences

### Error Response Format

| Aspect | api-legacy | api-core |
|--------|------------|----------|
| Null fields | Included (`"errors": null`) | Omitted |
| Message language | Korean | English |

Frontend handles both formats.

### History API

| Aspect | api-legacy | api-core |
|--------|------------|----------|
| Data source | Hibernate Envers tables | audit_logs table |
| Old history | Returns data | Returns empty (expected) |
| New history | Returns data | Returns data |

## Contact

For issues during deployment:

- Engineering team lead
- On-call engineer
- Platform team (for infrastructure issues)
