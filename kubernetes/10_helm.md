# Chapter 10: Helm — Package Manager for Kubernetes

## What is Helm?

**Helm** is the package manager for Kubernetes. Instead of maintaining many separate YAML files for each environment, Helm bundles them into a **Chart** — a reusable, versioned, configurable package.

> **Analogy:** Helm is to Kubernetes what `apt` or `brew` is to Linux/macOS.

---

## Why Use Helm?

| Without Helm | With Helm |
|---|---|
| 10+ YAML files per app | One chart |
| Manual sed/envsubst for different environments | `values.yaml` per environment |
| No versioning | Chart versioning + release history |
| Manual rollback = error-prone | `helm rollback` |
| Copy-paste config per cluster | `helm install -f values.prod.yaml` |

---

## Install Helm

```bash
# Linux/macOS (recommended)
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# macOS via Homebrew
brew install helm

# Verify
helm version
```

---

## Key Concepts

| Concept | Description |
|---|---|
| **Chart** | Package of Kubernetes YAML templates + metadata |
| **Release** | A deployed instance of a chart |
| **Repository** | Storage for charts (like Docker Hub for images) |
| **Values** | Config variables that customize a chart |

---

## Adding Chart Repositories

```bash
# Add popular repos
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add stable https://charts.helm.sh/stable
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

# Update repos
helm repo update

# Search for a chart
helm search repo nginx
helm search repo mysql
helm search hub wordpress   # Search Artifact Hub
```

---

## Installing Charts

```bash
# Install with default values
helm install my-nginx ingress-nginx/ingress-nginx

# Install with custom values
helm install my-mysql bitnami/mysql \
  --set auth.rootPassword=secret123 \
  --set primary.persistence.size=20Gi

# Install from a values file
helm install my-mysql bitnami/mysql -f values.yaml

# Install in a specific namespace
helm install my-mysql bitnami/mysql -n database --create-namespace

# Install and wait for readiness
helm install my-mysql bitnami/mysql --wait --timeout=300s
```

---

## Viewing Releases

```bash
# List all releases
helm list
helm list -A   # all namespaces

# Get details of a release
helm status my-mysql

# Get the values used in a release
helm get values my-mysql

# Get all rendered manifests
helm get manifest my-mysql
```

---

## Upgrading and Rolling Back

```bash
# Upgrade a release (new image, new values)
helm upgrade my-mysql bitnami/mysql \
  --set auth.rootPassword=newpassword

# Upgrade from values file
helm upgrade my-mysql bitnami/mysql -f values.prod.yaml

# Upgrade or install if not exists
helm upgrade --install my-mysql bitnami/mysql -f values.yaml

# View history
helm history my-mysql

# Rollback to previous version
helm rollback my-mysql

# Rollback to specific revision
helm rollback my-mysql 2
```

---

## Uninstalling

```bash
helm uninstall my-mysql
helm uninstall my-mysql -n database
```

---

## Creating Your Own Chart

```bash
# Generate scaffold
helm create my-app

# Structure:
# my-app/
# ├── Chart.yaml          # Chart metadata
# ├── values.yaml         # Default values
# ├── templates/
# │   ├── deployment.yaml
# │   ├── service.yaml
# │   ├── ingress.yaml
# │   ├── _helpers.tpl    # Template helper functions
# │   ├── hpa.yaml
# │   └── serviceaccount.yaml
# └── charts/             # Sub-charts (dependencies)
```

### Chart.yaml

```yaml
apiVersion: v2
name: my-app
description: A Helm chart for my application
type: application
version: 0.1.0        # Chart version
appVersion: "1.0.0"   # App version (image tag)
```

### values.yaml

```yaml
replicaCount: 3

image:
  repository: my-registry/my-app
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: nginx
  host: myapp.example.com
  tls: true

resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"

env:
  APP_ENV: production
  LOG_LEVEL: info

secrets:
  DB_PASSWORD: ""
  API_KEY: ""
```

### templates/deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-app.fullname" . }}
  labels:
    {{- include "my-app.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "my-app.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "my-app.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: {{ .Values.service.port }}
          env:
            {{- range $key, $val := .Values.env }}
            - name: {{ $key }}
              value: {{ $val | quote }}
            {{- end }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
```

---

## Multiple Environments with Values Files

```bash
# values.yaml       — defaults
# values.dev.yaml   — dev overrides
# values.prod.yaml  — production overrides
```

**values.prod.yaml:**
```yaml
replicaCount: 5
image:
  tag: "2.1.0"
resources:
  requests:
    cpu: "500m"
    memory: "512Mi"
  limits:
    cpu: "2"
    memory: "2Gi"
ingress:
  host: myapp.com
```

```bash
# Deploy to dev
helm upgrade --install my-app ./my-app -f values.dev.yaml -n dev

# Deploy to production
helm upgrade --install my-app ./my-app -f values.prod.yaml -n production
```

---

## Lint and Debug

```bash
# Lint the chart for errors
helm lint ./my-app

# Render templates without deploying (dry run)
helm template my-app ./my-app -f values.prod.yaml

# Dry run against a real cluster
helm install my-app ./my-app --dry-run --debug

# Show computed values
helm install my-app ./my-app --debug 2>&1 | head -50
```

---

## Helm Cheatsheet

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm search repo nginx
helm install my-release bitnami/nginx
helm list
helm upgrade my-release bitnami/nginx --set replicaCount=3
helm rollback my-release
helm history my-release
helm uninstall my-release
helm create my-chart
helm lint ./my-chart
helm template ./my-chart
helm package ./my-chart
```
