# Chapter 15: Security & Best Practices

## Security Layers in Kubernetes

```
1. Cluster Infrastructure  — secure nodes, network, etcd encryption
2. Authentication          — who can connect to the API server
3. Authorization (RBAC)    — what they can do
4. Admission Control       — validate/mutate resources before creation
5. Network Policies        — firewall between pods
6. Pod Security            — what containers can do on the node
7. Image Security          — trusted, scanned images only
8. Secrets Management      — proper handling of sensitive data
```

---

## SecurityContext

**SecurityContext** controls what a container can do on the host.

### Pod-level SecurityContext

```yaml
spec:
  securityContext:
    runAsNonRoot: true        # Container must not run as root
    runAsUser: 1000           # Run as user ID 1000
    runAsGroup: 3000
    fsGroup: 2000             # All volumes owned by this group
    seccompProfile:
      type: RuntimeDefault    # Apply default seccomp filter
```

### Container-level SecurityContext

```yaml
containers:
  - name: app
    securityContext:
      allowPrivilegeEscalation: false   # Can't gain more privileges
      readOnlyRootFilesystem: true      # Filesystem is read-only
      capabilities:
        drop:
          - ALL                         # Drop all Linux capabilities
        add:
          - NET_BIND_SERVICE            # Only add what's needed
```

### Recommended Secure Container Template

```yaml
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    fsGroup: 2000
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: app
      image: my-app:latest
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop: ["ALL"]
      volumeMounts:
        - name: tmp
          mountPath: /tmp          # Writable tmp directory
        - name: cache
          mountPath: /app/cache
  volumes:
    - name: tmp
      emptyDir: {}
    - name: cache
      emptyDir: {}
```

---

## Pod Security Standards (PSS)

Kubernetes enforces security policies via **Pod Security Admission** (replaced PodSecurityPolicy in K8s 1.25+).

Three levels:
| Level | Behavior |
|---|---|
| `privileged` | No restrictions |
| `baseline` | Minimal restrictions, prevents known privilege escalations |
| `restricted` | Heavily restricted, follows current pod hardening best practices |

Apply at the namespace level:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

---

## Network Policies

By default, all pods can communicate with all other pods. Use NetworkPolicy to restrict this.

### Zero-Trust: Deny All, Then Allow Selectively

```yaml
# 1. Deny all ingress and egress
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
  namespace: production
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
---
# 2. Allow frontend to reach backend
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - port: 8080
---
# 3. Allow all pods to reach DNS (required for service discovery)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
  namespace: production
spec:
  podSelector: {}
  policyTypes:
    - Egress
  egress:
    - ports:
        - port: 53
          protocol: UDP
        - port: 53
          protocol: TCP
```

---

## Image Security

### Use Specific Tags (Never `latest` in Production)

```yaml
# Bad
image: nginx:latest

# Good
image: nginx:1.25.3
# Or use digest (immutable)
image: nginx@sha256:abc123...
```

### Image Scanning

```bash
# Trivy: open-source vulnerability scanner
trivy image nginx:1.25.3
trivy image my-app:latest

# In CI/CD pipeline
trivy image --exit-code 1 --severity HIGH,CRITICAL my-app:latest
```

### Limit Image Registries (OPA/Kyverno)

```yaml
# Kyverno policy: only allow images from trusted registry
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: allowed-registries
spec:
  validationFailureAction: Enforce
  rules:
    - name: check-registry
      match:
        resources:
          kinds: [Pod]
      validate:
        message: "Image must come from registry.company.com"
        pattern:
          spec:
            containers:
              - image: "registry.company.com/*"
```

---

## etcd Encryption at Rest

```yaml
# /etc/kubernetes/encryption-config.yaml
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
      - secrets
    providers:
      - aescbc:
          keys:
            - name: key1
              secret: <base64-encoded-32-byte-key>
      - identity: {}
```

Add to kube-apiserver: `--encryption-provider-config=/etc/kubernetes/encryption-config.yaml`

---

## TLS Everywhere

```bash
# All internal K8s communication is TLS by default
# For your apps: use cert-manager (Chapter 9) to manage certs
# For ingress: always enable TLS termination at the ingress layer

# Rotate kubelet certificates
kubeadm certs renew all
```

---

## Admission Controllers

Admission controllers intercept requests before they're saved to etcd.

**OPA Gatekeeper** and **Kyverno** are popular policy engines:

```bash
# Install Kyverno
helm install kyverno kyverno/kyverno -n kyverno --create-namespace
```

```yaml
# Require all pods to have resource limits
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-resource-limits
spec:
  validationFailureAction: Enforce
  rules:
    - name: check-limits
      match:
        resources:
          kinds: [Pod]
      validate:
        message: "Resource limits are required"
        pattern:
          spec:
            containers:
              - resources:
                  limits:
                    memory: "?*"
                    cpu: "?*"
```

---

## Security Audit Commands

```bash
# Check all privileged pods
kubectl get pods -A -o json | jq '.items[] | select(.spec.containers[].securityContext.privileged == true) | .metadata.name'

# Check pods running as root
kubectl get pods -A -o json | jq '.items[] | select(.spec.securityContext.runAsUser == 0 or .spec.securityContext.runAsNonRoot != true) | .metadata.name'

# Who has cluster-admin?
kubectl get clusterrolebindings -o json | jq '.items[] | select(.roleRef.name == "cluster-admin") | .subjects'

# Check exposed secrets
kubectl get secrets -A

# Audit API server access logs
kubectl logs -n kube-system kube-apiserver-<node> | grep "authentication"

# kube-bench: CIS Kubernetes benchmark
kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job.yaml
kubectl logs job/kube-bench
```

---

## Security Checklist

- [ ] No containers running as root (`runAsNonRoot: true`)
- [ ] All containers have `allowPrivilegeEscalation: false`
- [ ] `readOnlyRootFilesystem: true` where possible
- [ ] All capabilities dropped (`drop: [ALL]`)
- [ ] NetworkPolicies restrict east-west traffic
- [ ] RBAC follows least privilege
- [ ] Secrets not in Git (use Sealed Secrets / Vault)
- [ ] etcd encryption at rest enabled
- [ ] Images scanned for CVEs before deployment
- [ ] Pod Security Standards enforced (`restricted` for prod)
- [ ] Ingress has TLS termination
- [ ] Admission controller (Kyverno/Gatekeeper) enforcing policies
- [ ] Regular CIS benchmark audits (kube-bench)
