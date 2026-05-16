# Chapter 6: ConfigMaps & Secrets

## Why Externalize Configuration?

Hard-coding config values in container images is an anti-pattern:
- You'd need different images per environment
- Secrets in images = security risk
- Config changes would require rebuilding images

Kubernetes solves this with **ConfigMaps** (non-sensitive config) and **Secrets** (sensitive data).

---

## ConfigMaps

### Create a ConfigMap

**From YAML:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: default
data:
  APP_ENV: "production"
  APP_DEBUG: "false"
  LOG_LEVEL: "info"
  DB_HOST: "mysql-service"
  DB_PORT: "3306"
  db.conf: |
    [client]
    host=mysql-service
    port=3306
    database=myapp
```

```bash
kubectl apply -f configmap.yaml
```

**From literal values:**
```bash
kubectl create configmap app-config \
  --from-literal=APP_ENV=production \
  --from-literal=LOG_LEVEL=info
```

**From a file:**
```bash
kubectl create configmap nginx-conf --from-file=nginx.conf
kubectl create configmap app-props --from-file=config/  # all files in dir
```

### Use ConfigMap as Environment Variables

```yaml
spec:
  containers:
    - name: app
      image: my-app:latest

      # Load all keys as env vars
      envFrom:
        - configMapRef:
            name: app-config

      # Load specific key
      env:
        - name: APP_ENV
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: APP_ENV
```

### Mount ConfigMap as Files

```yaml
spec:
  containers:
    - name: app
      image: my-app:latest
      volumeMounts:
        - name: config-volume
          mountPath: /etc/app/config   # Files appear here

  volumes:
    - name: config-volume
      configMap:
        name: app-config
        items:
          - key: db.conf
            path: database.conf        # Mounted as /etc/app/config/database.conf
```

---

## Secrets

Secrets store sensitive data like passwords, tokens, and keys. They are base64-encoded (not encrypted by default — use encryption at rest for production).

### Create a Secret

**From YAML (base64 values):**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: default
type: Opaque
data:
  DB_PASSWORD: cGFzc3dvcmQxMjM=    # base64 of "password123"
  API_KEY: c2VjcmV0a2V5MTIz          # base64 of "secretkey123"
```

> Encode: `echo -n "password123" | base64`
> Decode: `echo "cGFzc3dvcmQxMjM=" | base64 -d`

**Using `stringData` (plain text — K8s encodes it):**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  DB_PASSWORD: "password123"
  API_KEY: "secretkey123"
  DATABASE_URL: "mysql://user:password123@mysql-service:3306/myapp"
```

**From command line:**
```bash
kubectl create secret generic app-secrets \
  --from-literal=DB_PASSWORD=password123 \
  --from-literal=API_KEY=secretkey123

# From a file (e.g., TLS cert)
kubectl create secret tls my-tls \
  --cert=tls.crt \
  --key=tls.key

# From a Docker registry credential
kubectl create secret docker-registry regcred \
  --docker-server=registry.example.com \
  --docker-username=myuser \
  --docker-password=mypassword
```

### Use Secrets as Environment Variables

```yaml
spec:
  containers:
    - name: app
      image: my-app:latest
      envFrom:
        - secretRef:
            name: app-secrets
      env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
```

### Mount Secret as Files

```yaml
spec:
  containers:
    - name: app
      volumeMounts:
        - name: secret-volume
          mountPath: /etc/secrets
          readOnly: true

  volumes:
    - name: secret-volume
      secret:
        secretName: app-secrets
```

### Pull Private Images with Image Pull Secrets

```yaml
spec:
  imagePullSecrets:
    - name: regcred
  containers:
    - name: app
      image: registry.example.com/my-private-app:latest
```

---

## Combined ConfigMap + Secret Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: web-config
data:
  APP_ENV: "production"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
---
apiVersion: v1
kind: Secret
metadata:
  name: web-secrets
type: Opaque
stringData:
  DB_PASSWORD: "super-secret"
  JWT_SECRET: "my-jwt-signing-key"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
        - name: web-app
          image: my-app:latest
          envFrom:
            - configMapRef:
                name: web-config
            - secretRef:
                name: web-secrets
```

---

## Managing ConfigMaps and Secrets

```bash
# View ConfigMap
kubectl get configmaps
kubectl describe configmap app-config
kubectl get configmap app-config -o yaml

# Edit ConfigMap
kubectl edit configmap app-config

# View Secret (values are base64 encoded)
kubectl get secrets
kubectl describe secret app-secrets
kubectl get secret app-secrets -o yaml

# Decode a secret value
kubectl get secret app-secrets -o jsonpath='{.data.DB_PASSWORD}' | base64 -d

# Delete
kubectl delete configmap app-config
kubectl delete secret app-secrets
```

---

## Production Secret Management

> Default Kubernetes Secrets are only base64-encoded, not encrypted. For production, use:

| Solution | Description |
|---|---|
| **Kubernetes encryption at rest** | Encrypt etcd with `EncryptionConfiguration` |
| **HashiCorp Vault** | External secret store, injects secrets at pod startup |
| **AWS Secrets Manager** | With `external-secrets` operator |
| **Sealed Secrets** | Encrypt secrets in Git safely (by Bitnami) |
| **SOPS** | File-based secret encryption for GitOps |

### Sealed Secrets (Quick Example)

```bash
# Install kubeseal CLI
kubeseal --fetch-cert > pub-sealed-secrets.pem

# Seal a secret (safe to commit to Git)
kubectl create secret generic my-secret --from-literal=password=abc123 --dry-run=client -o yaml | \
  kubeseal --cert pub-sealed-secrets.pem -o yaml > my-sealed-secret.yaml

# Apply (controller decrypts it in cluster)
kubectl apply -f my-sealed-secret.yaml
```

---

## Best Practices

- Never hardcode passwords in Dockerfiles or application code
- Never commit raw Secrets YAML to Git (use Sealed Secrets or Vault)
- Use `readOnly: true` when mounting secrets as volumes
- Rotate secrets regularly and update the Secret object
- Use RBAC to restrict who can `get`/`list` secrets in each namespace
- Consider using `external-secrets` operator to sync from AWS/GCP/Azure secret managers
