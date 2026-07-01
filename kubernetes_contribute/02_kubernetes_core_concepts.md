# Chapter 2: Kubernetes Core Concepts

## Pod

A **Pod** is the smallest unit Kubernetes runs.

Usually:

- 1 Pod = 1 application container
- Sometimes 1 Pod = app container + helper sidecar container

Example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
    - name: nginx
      image: nginx:1.27
      ports:
        - containerPort: 80
```

You rarely manage Pods directly in production. You usually create a Deployment.

---

## Deployment

A **Deployment** manages replicated Pods and rollouts.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.27
          ports:
            - containerPort: 80
```

If one Pod dies, the Deployment controller creates another.

---

## Service

Pods can be replaced, so their IPs are not stable. A **Service** gives a stable network identity.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
    - port: 80
      targetPort: 80
```

Think:

```text
Deployment = keep app running
Service    = let others reach app
```

---

## ConfigMap and Secret

Use **ConfigMaps** for non-sensitive configuration.

Use **Secrets** for sensitive values, but remember: default Kubernetes Secrets are encoded, not magically secure. Real security depends on cluster configuration.

---

## Namespace

A **Namespace** groups resources.

```bash
kubectl create namespace learning
kubectl get pods -n learning
```

Namespaces are useful for separating environments, teams, and experiments.

---

## Labels and Selectors

Labels are key-value metadata.

Selectors find resources by labels.

```text
Pod label:
  app=nginx

Service selector:
  app=nginx

Result:
  Service sends traffic to Pods with app=nginx
```

Many Kubernetes systems depend on labels.

---

## The Contributor Angle

When reading Kubernetes code, always ask:

- What object is this code working with?
- Which controller owns the behavior?
- Is this desired state or actual state?
- Is this validation, scheduling, reconciliation, networking, storage, or node behavior?

