# Chapter 12: Spring Boot with Kubernetes

## Deploying Spring Boot to Kubernetes

This chapter covers deploying a Spring Boot application to Kubernetes, including best practices for health checks, configuration, secrets, and scaling.

---

## Basic Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-spring-app
  namespace: production
  labels:
    app: my-spring-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-spring-app
  template:
    metadata:
      labels:
        app: my-spring-app
    spec:
      containers:
        - name: my-spring-app
          image: myregistry/my-spring-app:1.0.0
          ports:
            - containerPort: 8080
          env:
            - name: SPRING_PROFILES_ACTIVE
              value: prod
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: db-secret
                  key: url
            - name: DATABASE_USER
              valueFrom:
                secretKeyRef:
                  name: db-secret
                  key: username
            - name: DATABASE_PASS
              valueFrom:
                secretKeyRef:
                  name: db-secret
                  key: password
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /actuator/health/liveness
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /actuator/health/readiness
              port: 8080
            initialDelaySeconds: 20
            periodSeconds: 5
            failureThreshold: 3
```

---

## Liveness vs Readiness Probes

Spring Boot Actuator provides separate health groups for Kubernetes probes:

```yaml
# application.yml
management:
  endpoint:
    health:
      probes:
        enabled: true     # enables /actuator/health/liveness and /actuator/health/readiness
  health:
    livenessstate:
      enabled: true
    readinessstate:
      enabled: true
```

| Probe | Path | Meaning | Failure Action |
|---|---|---|---|
| **Liveness** | `/actuator/health/liveness` | App is alive (not deadlocked/crashed) | Restart the pod |
| **Readiness** | `/actuator/health/readiness` | App is ready to serve traffic | Remove from Service endpoints |
| **Startup** | `/actuator/health` | App has started | Wait longer before checking liveness |

```java
// Custom liveness indicator
@Component
public class DeadlockLivenessIndicator implements LivenessStateReporter {
    // Only marks DOWN if the app is truly broken (needs restart)
}

// Custom readiness indicator
@Component
public class DatabaseReadinessIndicator implements ReadinessStateReporter {
    // Marks DOWN if the app can't serve requests (DB is down)
}
```

---

## ConfigMaps for Application Config

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-app-config
  namespace: production
data:
  application.yaml: |
    server:
      port: 8080
    spring:
      jpa:
        show-sql: false
    logging:
      level:
        root: INFO
        com.example: INFO
```

Mount as a file:
```yaml
containers:
  - name: my-spring-app
    volumeMounts:
      - name: config
        mountPath: /config
volumes:
  - name: config
    configMap:
      name: my-app-config
```

```yaml
# Tell Spring Boot to load from this path
SPRING_CONFIG_LOCATION: /config/application.yaml
```

---

## Secrets

```bash
# Create from CLI
kubectl create secret generic db-secret \
  --from-literal=url="jdbc:postgresql://postgres:5432/mydb" \
  --from-literal=username="myuser" \
  --from-literal=password="supersecret"

# Or from YAML (values must be base64 encoded)
echo -n "supersecret" | base64   # c3VwZXJzZWNyZXQ=
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
data:
  url: amRiYzpwb3N0Z3Jlc3FsOi8vcG9zdGdyZXM6NTQzMi9teWRi
  username: bXl1c2Vy
  password: c3VwZXJzZWNyZXQ=
```

---

## Service

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-spring-app
  namespace: production
spec:
  selector:
    app: my-spring-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP   # internal only; use LoadBalancer for external
```

---

## Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-spring-app
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - api.example.com
      secretName: api-tls
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-spring-app
                port:
                  number: 80
```

---

## Horizontal Pod Autoscaler (HPA)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: my-spring-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-spring-app
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

## Graceful Shutdown in Kubernetes

```yaml
# application.yml
server:
  shutdown: graceful
spring:
  lifecycle:
    timeout-per-shutdown-phase: 30s
```

```yaml
# In Deployment spec
spec:
  template:
    spec:
      terminationGracePeriodSeconds: 60   # must be > timeout-per-shutdown-phase
      containers:
        - name: my-spring-app
          lifecycle:
            preStop:
              exec:
                command: ["sleep", "5"]   # wait for load balancer to remove pod
```

When Kubernetes sends `SIGTERM`:
1. Pod is removed from Service endpoints (stop receiving traffic)
2. `preStop` hook runs (sleep 5s)
3. Spring Boot receives SIGTERM → graceful shutdown starts
4. Active requests finish → app exits

---

## Helm Chart for Spring Boot

```
my-spring-app/
├── Chart.yaml
├── values.yaml
└── templates/
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    ├── configmap.yaml
    ├── secret.yaml
    └── hpa.yaml
```

```yaml
# values.yaml
replicaCount: 2

image:
  repository: myregistry/my-spring-app
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  port: 80

ingress:
  enabled: true
  host: api.example.com

resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"

env:
  SPRING_PROFILES_ACTIVE: prod

secrets:
  DATABASE_URL: ""
  DATABASE_USER: ""
  DATABASE_PASS: ""

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

```bash
# Install
helm install my-spring-app ./my-spring-app -f values-prod.yaml

# Upgrade
helm upgrade my-spring-app ./my-spring-app --set image.tag=1.1.0

# Rollback
helm rollback my-spring-app 1
```

---

## Zero-Downtime Deployments

```yaml
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1          # allow 1 extra pod during update
      maxUnavailable: 0    # never reduce below desired replicas
```

With `readinessProbe` configured, Kubernetes won't route traffic to a new pod until it passes readiness — ensuring zero downtime.

---

## Best Practices for Spring Boot on K8s

| Practice | Why |
|---|---|
| Always set `resources.requests` and `limits` | Prevents noisy-neighbor issues, required for HPA |
| Use `livenessProbe` and `readinessProbe` | Prevents traffic to unhealthy pods |
| Use `server.shutdown=graceful` | Prevents dropped requests during deploys |
| Store secrets in K8s Secrets or external vault | Never bake secrets into images |
| Run as non-root user | Security best practice |
| Use `readOnly: true` filesystem | Limit attack surface |
| Set `terminationGracePeriodSeconds > shutdown timeout` | Prevents forced kills |
