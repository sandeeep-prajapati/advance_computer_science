# Chapter 16: CI/CD with Kubernetes

## CI/CD Pipeline Overview

```
Developer pushes code
    ↓
CI (GitHub Actions / GitLab CI)
    ├── Run tests
    ├── Build Docker image
    ├── Push to registry
    └── Deploy to Kubernetes
         ├── kubectl apply (imperative)
         └── ArgoCD sync (GitOps)
```

---

## GitHub Actions: Build & Deploy

### Complete Workflow

**File:** `.github/workflows/deploy.yml`

```yaml
name: Build and Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: test
          MYSQL_DATABASE: testdb
        ports:
          - 3306:3306
        options: --health-cmd="mysqladmin ping" --health-interval=10s

    steps:
      - uses: actions/checkout@v4

      - name: Run tests
        run: |
          composer install
          vendor/bin/phpunit

  build:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    outputs:
      image-tag: ${{ steps.meta.outputs.tags }}

    steps:
      - uses: actions/checkout@v4

      - name: Log in to registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=sha,prefix=sha-
            type=ref,event=branch
            type=semver,pattern={{version}}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: staging

    steps:
      - uses: actions/checkout@v4

      - name: Configure kubectl
        uses: azure/k8s-set-context@v3
        with:
          method: kubeconfig
          kubeconfig: ${{ secrets.KUBE_CONFIG_STAGING }}

      - name: Deploy to staging
        run: |
          IMAGE_TAG="sha-${{ github.sha }}"
          kubectl set image deployment/my-app \
            my-app=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:$IMAGE_TAG \
            -n staging
          kubectl rollout status deployment/my-app -n staging --timeout=5m

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production   # Requires manual approval

    steps:
      - uses: actions/checkout@v4

      - name: Configure kubectl
        uses: azure/k8s-set-context@v3
        with:
          method: kubeconfig
          kubeconfig: ${{ secrets.KUBE_CONFIG_PRODUCTION }}

      - name: Deploy to production
        run: |
          IMAGE_TAG="sha-${{ github.sha }}"
          kubectl set image deployment/my-app \
            my-app=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:$IMAGE_TAG \
            -n production
          kubectl rollout status deployment/my-app -n production --timeout=10m
```

---

## ArgoCD — GitOps Continuous Delivery

**ArgoCD** watches a Git repository and automatically syncs the cluster to match what's in Git.

```
Git repo (source of truth)
    ↓ ArgoCD watches for changes
    ↓ Detects drift between Git state and cluster state
    ↓ Auto-syncs cluster to match Git
Kubernetes cluster
```

### Install ArgoCD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f \
  https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for pods
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=argocd-server -n argocd

# Get admin password
kubectl get secret argocd-initial-admin-secret -n argocd \
  -o jsonpath='{.data.password}' | base64 -d

# Access UI
kubectl port-forward svc/argocd-server 8080:443 -n argocd
# Open: https://localhost:8080 (admin / <password above>)
```

### Create an ArgoCD Application

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/my-org/my-app-config
    targetRevision: HEAD
    path: k8s/production         # Path to K8s manifests in the repo

  destination:
    server: https://kubernetes.default.svc
    namespace: production

  syncPolicy:
    automated:
      prune: true                # Delete resources removed from Git
      selfHeal: true             # Revert manual changes to cluster
    syncOptions:
      - CreateNamespace=true
```

### GitOps Workflow

```
Developer pushes YAML change to Git
    ↓
ArgoCD detects difference
    ↓
ArgoCD applies change to cluster
    ↓
ArgoCD shows sync status: Synced ✓
```

### ArgoCD CLI

```bash
# Install ArgoCD CLI
curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install argocd /usr/local/bin/argocd

# Login
argocd login localhost:8080

# List apps
argocd app list

# Sync an app manually
argocd app sync my-app

# Check sync status
argocd app get my-app

# Rollback to previous version
argocd app rollback my-app
```

---

## Blue-Green Deployment

Two identical environments: only one receives traffic at a time.

```yaml
# Blue deployment (current live)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
      version: blue
  template:
    metadata:
      labels:
        app: my-app
        version: blue
    spec:
      containers:
        - name: my-app
          image: my-app:1.0

---
# Green deployment (new version)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
      version: green
  template:
    metadata:
      labels:
        app: my-app
        version: green
    spec:
      containers:
        - name: my-app
          image: my-app:2.0

---
# Service: switch traffic from blue to green
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
spec:
  selector:
    app: my-app
    version: blue    # Change to 'green' to switch traffic
  ports:
    - port: 80
```

```bash
# Switch to green
kubectl patch service my-app-service \
  -p '{"spec":{"selector":{"version":"green"}}}'

# Rollback: switch back to blue
kubectl patch service my-app-service \
  -p '{"spec":{"selector":{"version":"blue"}}}'
```

---

## Canary Deployment

Send a small percentage of traffic to the new version:

```yaml
# Current: 3 replicas of v1 (75% traffic)
# Canary:  1 replica of v2 (25% traffic)

apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-stable
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:1.0
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-canary
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:2.0   # New version
---
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
spec:
  selector:
    app: my-app     # Routes to BOTH deployments (both share this label)
  ports:
    - port: 80
```

---

## Docker Image Best Practices for CI/CD

```dockerfile
# Multi-stage build — keep final image small
FROM composer:2.7 AS composer
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --no-scripts --prefer-dist

FROM php:8.3-fpm-alpine AS production
WORKDIR /var/www/html

# Install runtime dependencies only
RUN apk add --no-cache nginx

COPY --from=composer /app/vendor ./vendor
COPY . .

# Non-root user
RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 9000
CMD ["php-fpm"]
```

---

## CI/CD Best Practices

- **Build once, deploy many:** Same image to staging then production (just change env)
- **Tag images with Git SHA:** `my-app:sha-abc1234` — immutable, traceable
- **Never use `latest` in production:** Immutable tags only
- **Rollback = re-deploy old tag:** Not a special process — just set the old image tag
- **Store Kubernetes manifests in Git:** Enable ArgoCD GitOps
- **Separate CI (test/build) from CD (deploy):** Different workflows, different secrets
- **Use environments + approvals:** Staging auto-deploys, production requires approval
