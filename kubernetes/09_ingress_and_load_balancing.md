# Chapter 9: Ingress & Load Balancing

## What is Ingress?

A **Service** of type LoadBalancer gives you one IP per service — expensive on cloud providers. **Ingress** solves this by providing a single entry point that routes HTTP/HTTPS traffic to multiple services based on hostname or URL path.

```
Internet
    ↓
Ingress (single LoadBalancer IP)
    ├── /api    → api-service
    ├── /admin  → admin-service
    └── /       → frontend-service
```

---

## Ingress Controller

An **Ingress** resource alone does nothing — you need an **Ingress Controller** to implement it.

| Controller | When to Use |
|---|---|
| **Nginx Ingress Controller** | Most common, feature-rich |
| **Traefik** | Auto-discovers services, good for Kubernetes-native |
| **AWS ALB Ingress** | AWS-native, integrates with ALB |
| **GCE Ingress** | GKE-native, integrates with Google Cloud LB |
| **HAProxy** | High-performance edge |

---

## Install Nginx Ingress Controller

### On Minikube

```bash
minikube addons enable ingress
kubectl get pods -n ingress-nginx
```

### On Any Cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/cloud/deploy.yaml

# Wait for controller to be ready
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

---

## Basic Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
  namespace: default
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## Path-Based Routing

Route different paths to different services on the same domain:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    nginx.ingress.kubernetes.io/use-regex: "true"
spec:
  ingressClassName: nginx
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080

          - path: /admin
            pathType: Prefix
            backend:
              service:
                name: admin-service
                port:
                  number: 3000

          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## Host-Based Routing

Route different hostnames to different services:

```yaml
spec:
  rules:
    - host: api.myapp.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080

    - host: admin.myapp.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: admin-service
                port:
                  number: 3000

    - host: myapp.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

---

## TLS / HTTPS with cert-manager

### Install cert-manager

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.5/cert-manager.yaml

kubectl get pods -n cert-manager
```

### Create a Let's Encrypt ClusterIssuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - http01:
          ingress:
            class: nginx
```

### Ingress with TLS

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: secure-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - myapp.example.com
        - api.myapp.example.com
      secretName: myapp-tls        # cert-manager creates this Secret
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
    - host: api.myapp.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 8080
```

---

## Common Nginx Ingress Annotations

```yaml
annotations:
  # Redirect HTTP to HTTPS
  nginx.ingress.kubernetes.io/ssl-redirect: "true"

  # Increase upload size limit
  nginx.ingress.kubernetes.io/proxy-body-size: "100m"

  # Increase timeouts
  nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
  nginx.ingress.kubernetes.io/proxy-send-timeout: "300"

  # Rate limiting
  nginx.ingress.kubernetes.io/limit-rps: "10"

  # CORS
  nginx.ingress.kubernetes.io/enable-cors: "true"
  nginx.ingress.kubernetes.io/cors-allow-origin: "https://myapp.com"

  # Basic auth
  nginx.ingress.kubernetes.io/auth-type: basic
  nginx.ingress.kubernetes.io/auth-secret: basic-auth
  nginx.ingress.kubernetes.io/auth-realm: "Authentication Required"

  # Rewrite URL
  nginx.ingress.kubernetes.io/rewrite-target: /$2
```

---

## Ingress on Minikube (Local Testing)

```bash
# Enable ingress
minikube addons enable ingress

# Get Minikube IP
minikube ip   # e.g., 192.168.49.2

# Add to /etc/hosts for local DNS
echo "192.168.49.2 myapp.local" | sudo tee -a /etc/hosts

# Apply ingress pointing to myapp.local
# Now curl http://myapp.local works
```

---

## Ingress Commands

```bash
# List ingresses
kubectl get ingress
kubectl get ingress -A

# Describe (see rules and backend status)
kubectl describe ingress my-ingress

# Check ingress controller logs
kubectl logs -n ingress-nginx \
  deployment/ingress-nginx-controller

# Check cert-manager certificate status
kubectl get certificate
kubectl describe certificate myapp-tls
```

---

## Full Working Example

```yaml
# All in one: Deployment + Service + Ingress
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-app
  template:
    metadata:
      labels:
        app: hello-app
    spec:
      containers:
        - name: hello-app
          image: hashicorp/http-echo
          args:
            - "-text=Hello from Kubernetes!"
          ports:
            - containerPort: 5678
---
apiVersion: v1
kind: Service
metadata:
  name: hello-service
spec:
  selector:
    app: hello-app
  ports:
    - port: 80
      targetPort: 5678
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: hello.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hello-service
                port:
                  number: 80
```
