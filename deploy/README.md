# Kubernetes Deployment Configuration

## Directory Structure

```
deploy/
├── templates/          # Kubernetes resource templates
│   ├── api-legacy/     # Spring Boot API resources
│   │   └── deployment.yaml
│   ├── api-core/       # Go API resources
│   │   ├── deployment.yaml
│   │   └── configmap.yaml
│   ├── frontend/       # Vue.js frontend resources
│   │   └── deployment.yaml
│   └── ingress/        # Ingress configuration
│       └── ingress.yaml
└── environments/       # Environment-specific values
    ├── development.env
    ├── staging.env
    └── production.env
```

## API Routing

The Ingress configuration routes traffic as follows:

- `/api/*` → `api-core` (Go API) - Default
- `/*` → `frontend` (Vue.js SPA)

To use api-legacy instead, replace `ingress.yaml` with `ingress-legacy.yaml` in the deployment.

## API Backend Selection

Both api-core and api-legacy use the same `/api/v1` endpoints, allowing seamless switching:

### In Local Development:
```bash
# Use api-core (default)
yarn dev
# or
API_BACKEND=core docker-compose up -d frontend

# Use api-legacy
yarn dev:legacy
# or
API_BACKEND=legacy docker-compose up -d frontend
```

### In Kubernetes:
- Default: Uses api-core
- To switch to api-legacy: Deploy with `ingress-legacy.yaml` instead of `ingress.yaml`

## Environment Variables

### api-core specific
- `PROFILE`: Runtime profile (development/staging/production)

### Common variables
- `APPLICATION_REPLICAS`: Number of pod replicas
- `DEPLOY_ENVIRONMENT`: Deployment environment name

## Deployment Process

1. GitLab CI reads environment-specific values from `environments/*.env`
2. Templates are processed with `envsubst` to replace variables
3. Processed YAML files are applied to the Kubernetes cluster

## Local Development

To test the API locally:

```bash
# Start all services
docker-compose up -d

# Test api-legacy directly
curl http://localhost:8081/api/v1/config

# Test api-core directly
curl http://localhost:8080/api/v1/config

# Test through frontend proxy (defaults to api-core)
curl http://localhost:9000/api/v1/config

# Switch frontend to use api-legacy
API_BACKEND=legacy docker-compose up -d frontend
curl http://localhost:9000/api/v1/config  # Now proxied to api-legacy
```