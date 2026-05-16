# Chapter 18: Kubernetes Setup & Debugging

## Overview

Running Pimcore on Kubernetes (K8s) gives you:
- **Horizontal scaling** — multiple PHP-FPM pods
- **High availability** — automatic restarts, rolling updates
- **Infrastructure as code** — all config in YAML files
- **Environment isolation** — dev, staging, prod on same cluster

Pimcore on Kubernetes requires: a PHP-FPM app pod, a MySQL pod (or external RDS), a Redis pod, a shared persistent volume for assets, and an Nginx ingress controller.

---

## Architecture

```
Internet
    │
    ▼
┌─────────────────────────────┐
│   Ingress (Nginx Controller) │
└──────────────┬──────────────┘
               │
    ┌──────────▼──────────┐
    │   Nginx Service     │   (ClusterIP)
    └──────────┬──────────┘
               │
    ┌──────────▼──────────┐
    │   PHP-FPM Pods      │   (Deployment, 2+ replicas)
    │   (pimcore app)     │
    └──────┬──────┬───────┘
           │      │
    ┌──────▼──┐ ┌─▼──────┐  ┌──────────────┐
    │  MySQL  │ │ Redis  │  │ Shared PVC   │
    │  Pod    │ │ Pod    │  │ (assets,var) │
    └─────────┘ └────────┘  └──────────────┘
```

---

## Step 1: Docker Image for Pimcore

Create a production Docker image:

**`Dockerfile`:**
```dockerfile
FROM pimcore/pimcore:php8.3-latest

WORKDIR /var/www/html

# Copy project files
COPY . .

# Install Composer dependencies (no dev)
RUN composer install --no-dev --optimize-autoloader --prefer-dist

# Set permissions
RUN chown -R www-data:www-data var/ public/var/ && \
    chmod -R 775 var/ public/var/

EXPOSE 9000

CMD ["php-fpm"]
```

**Nginx sidecar `nginx.conf`:**
```nginx
server {
    listen 80;
    server_name _;
    root /var/www/html/public;
    index index.php;

    location / {
        try_files $uri $uri/ /index.php$is_args$args;
    }

    location ~ \.php$ {
        fastcgi_pass 127.0.0.1:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
        fastcgi_read_timeout 300;
    }

    location ~* \.(jpg|jpeg|png|gif|ico|css|js|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

---

## Step 2: Namespace & ConfigMap

```bash
kubectl create namespace pimcore
```

**`k8s/configmap.yaml`:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: pimcore-config
  namespace: pimcore
data:
  APP_ENV: "prod"
  APP_DEBUG: "0"
  PIMCORE_ENVIRONMENT: "prod"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  MYSQL_HOST: "mysql-service"
  MYSQL_DATABASE: "pimcore"
  MYSQL_USER: "pimcore"
```

---

## Step 3: Secrets

**`k8s/secret.yaml`:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: pimcore-secrets
  namespace: pimcore
type: Opaque
stringData:
  MYSQL_PASSWORD: "your_secure_password"
  DATABASE_URL: "mysql://pimcore:your_secure_password@mysql-service:3306/pimcore"
  APP_SECRET: "your_symfony_app_secret_32chars"
  PIMCORE_IMAGE_PROXY_SECRET: "your_proxy_secret"
```

```bash
kubectl apply -f k8s/secret.yaml -n pimcore
```

---

## Step 4: Persistent Volume for Assets

Pimcore assets (`var/assets`) must be **shared across all PHP pods** — use a ReadWriteMany PVC.

**`k8s/pvc.yaml`:**
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pimcore-assets-pvc
  namespace: pimcore
spec:
  accessModes:
    - ReadWriteMany        # Required for multiple pods to share
  storageClassName: nfs    # Use NFS, EFS (AWS), or Filestore (GCP)
  resources:
    requests:
      storage: 50Gi
```

> **Important:** Use a `ReadWriteMany` storage class. Options:
> - On-premise: NFS server with `nfs-subdir-external-provisioner`
> - AWS: EFS with `aws-efs-csi-driver`
> - GCP: Filestore
> - Azure: Azure Files

---

## Step 5: MySQL Deployment

**`k8s/mysql.yaml`:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: pimcore
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          env:
            - name: MYSQL_DATABASE
              valueFrom:
                configMapKeyRef:
                  name: pimcore-config
                  key: MYSQL_DATABASE
            - name: MYSQL_USER
              valueFrom:
                configMapKeyRef:
                  name: pimcore-config
                  key: MYSQL_USER
            - name: MYSQL_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: pimcore-secrets
                  key: MYSQL_PASSWORD
            - name: MYSQL_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: pimcore-secrets
                  key: MYSQL_PASSWORD
          ports:
            - containerPort: 3306
          volumeMounts:
            - name: mysql-data
              mountPath: /var/lib/mysql
      volumes:
        - name: mysql-data
          persistentVolumeClaim:
            claimName: mysql-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
  namespace: pimcore
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
      targetPort: 3306
  clusterIP: None   # Headless service for stable DNS
```

---

## Step 6: Redis Deployment

**`k8s/redis.yaml`:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: pimcore
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          ports:
            - containerPort: 6379
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: redis-service
  namespace: pimcore
spec:
  selector:
    app: redis
  ports:
    - port: 6379
      targetPort: 6379
```

---

## Step 7: Pimcore App Deployment

**`k8s/pimcore-app.yaml`:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pimcore-app
  namespace: pimcore
spec:
  replicas: 2          # Scale horizontally
  selector:
    matchLabels:
      app: pimcore-app
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0   # Zero-downtime rolling updates
  template:
    metadata:
      labels:
        app: pimcore-app
    spec:
      initContainers:
        # Run migrations and class rebuild before the app starts
        - name: pimcore-init
          image: your-registry/pimcore:latest
          command: ["/bin/sh", "-c"]
          args:
            - |
              bin/console pimcore:deployment:classes-rebuild -c -d
              bin/console doctrine:migrations:migrate --no-interaction --env=prod
              bin/console cache:clear --env=prod
              bin/console cache:warmup --env=prod
          envFrom:
            - configMapRef:
                name: pimcore-config
            - secretRef:
                name: pimcore-secrets
          volumeMounts:
            - name: pimcore-assets
              mountPath: /var/www/html/var/assets
            - name: pimcore-cache
              mountPath: /var/www/html/var/cache

      containers:
        # Nginx container (serves static files, proxies PHP)
        - name: nginx
          image: nginx:1.25-alpine
          ports:
            - containerPort: 80
          volumeMounts:
            - name: nginx-config
              mountPath: /etc/nginx/conf.d
            - name: pimcore-assets
              mountPath: /var/www/html/var/assets
              readOnly: true

        # PHP-FPM container (runs Pimcore)
        - name: php-fpm
          image: your-registry/pimcore:latest
          envFrom:
            - configMapRef:
                name: pimcore-config
            - secretRef:
                name: pimcore-secrets
          ports:
            - containerPort: 9000
          volumeMounts:
            - name: pimcore-assets
              mountPath: /var/www/html/var/assets
            - name: pimcore-cache
              mountPath: /var/www/html/var/cache
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
          livenessProbe:
            exec:
              command: ["php-fpm-healthcheck"]
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command: ["php-fpm-healthcheck", "--accepted-conn=1"]
            initialDelaySeconds: 5
            periodSeconds: 5

      volumes:
        - name: pimcore-assets
          persistentVolumeClaim:
            claimName: pimcore-assets-pvc
        - name: pimcore-cache
          emptyDir: {}      # Cache is ephemeral per pod (Redis handles actual cache)
        - name: nginx-config
          configMap:
            name: nginx-config
---
apiVersion: v1
kind: Service
metadata:
  name: pimcore-service
  namespace: pimcore
spec:
  selector:
    app: pimcore-app
  ports:
    - port: 80
      targetPort: 80
```

---

## Step 8: Ingress (External Access)

**`k8s/ingress.yaml`:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pimcore-ingress
  namespace: pimcore
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/proxy-body-size: "100m"     # allow large file uploads
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
spec:
  tls:
    - hosts:
        - pimcore.yourdomain.com
      secretName: pimcore-tls
  rules:
    - host: pimcore.yourdomain.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: pimcore-service
                port:
                  number: 80
```

---

## Step 9: Deploy Everything

```bash
# Apply all manifests
kubectl apply -f k8s/ -n pimcore

# Watch rollout
kubectl rollout status deployment/pimcore-app -n pimcore

# Check all pods
kubectl get pods -n pimcore
```

---

## Debugging in Kubernetes

### Get a Shell Inside a Pod

```bash
# Get pod names
kubectl get pods -n pimcore

# Shell into PHP-FPM container
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- bash

# Run Pimcore console inside pod
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- bin/console cache:clear
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- bin/console pimcore:system:requirements:check
```

### View Logs

```bash
# All logs from pimcore-app pods
kubectl logs -l app=pimcore-app -n pimcore --all-containers

# Follow logs from specific pod
kubectl logs -f <pod-name> -n pimcore -c php-fpm

# Previous pod logs (after a crash)
kubectl logs <pod-name> -n pimcore -c php-fpm --previous

# View Pimcore's own log files inside pod
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- tail -f var/log/prod.log
```

### Debug Pod Startup Issues

```bash
# Describe pod (shows events, failed probes, image pull errors)
kubectl describe pod <pod-name> -n pimcore

# Check if init container ran successfully
kubectl logs <pod-name> -n pimcore -c pimcore-init

# Check all events in namespace
kubectl get events -n pimcore --sort-by='.lastTimestamp'
```

### Debug Database Connection

```bash
# Exec into app pod and test MySQL
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- \
  bin/console doctrine:query:sql "SELECT 1"

# Check if MySQL service is reachable
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- \
  nc -zv mysql-service 3306
```

### Debug Redis Connection

```bash
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- \
  php -r "
    \$r = new Redis();
    \$r->connect('redis-service', 6379);
    echo \$r->ping() . PHP_EOL;
  "
```

### Check Persistent Volume

```bash
# Check PVC status
kubectl get pvc -n pimcore

# Check if assets directory is mounted and writable
kubectl exec -it <pod-name> -n pimcore -c php-fpm -- \
  ls -la var/assets/
```

---

## Horizontal Pod Autoscaler (HPA)

Scale PHP pods automatically based on CPU:

**`k8s/hpa.yaml`:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: pimcore-hpa
  namespace: pimcore
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: pimcore-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

---

## Helm Chart (Optional — Recommended for Production)

For repeatable deployments across environments, use a **Helm chart**:

```bash
# Create a Helm chart
helm create pimcore-chart

# Values file for staging
helm install pimcore ./pimcore-chart -f values.staging.yaml -n pimcore

# Values file for production
helm install pimcore ./pimcore-chart -f values.prod.yaml -n pimcore

# Upgrade a release
helm upgrade pimcore ./pimcore-chart -f values.prod.yaml -n pimcore
```

---

## Common Kubernetes Issues for Pimcore

| Issue | Cause | Fix |
|---|---|---|
| Pods in `CrashLoopBackOff` | Migration failed / DB not ready | Check init container logs |
| Assets not shared between pods | PVC not ReadWriteMany | Use NFS/EFS storage class |
| 502 Bad Gateway | PHP-FPM not ready | Check readinessProbe / increase `initialDelaySeconds` |
| Cache inconsistency across pods | Filesystem cache used | Switch to Redis cache adapter |
| Large file upload fails | Nginx body size limit | Set `nginx.ingress.kubernetes.io/proxy-body-size: "100m"` |
| OOM killed pods | PHP memory limit too low | Increase pod memory limits |
| Slow cold start | No OPcache preloading | Add OPcache preload config |

---

## Kubernetes Cheatsheet

```bash
# Get all resources in namespace
kubectl get all -n pimcore

# Watch pods
kubectl get pods -n pimcore -w

# Scale deployment
kubectl scale deployment pimcore-app --replicas=4 -n pimcore

# Rolling restart (pick up new image or config)
kubectl rollout restart deployment/pimcore-app -n pimcore

# View rollout history
kubectl rollout history deployment/pimcore-app -n pimcore

# Rollback to previous version
kubectl rollout undo deployment/pimcore-app -n pimcore

# Port-forward to test locally
kubectl port-forward svc/pimcore-service 8080:80 -n pimcore
# Then open: http://localhost:8080
```
