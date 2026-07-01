# Chapter 4: Hands-on Kubernetes Basics

## Important kubectl Commands

```bash
kubectl get pods
kubectl get deployments
kubectl get services
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl exec -it <pod-name> -- sh
kubectl delete pod <pod-name>
```

When you delete a Pod owned by a Deployment, Kubernetes creates a replacement. That is reconciliation in action.

---

## Create from YAML

Create `nginx.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 2
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
```

Apply:

```bash
kubectl apply -f nginx.yaml
kubectl get pods
```

Delete:

```bash
kubectl delete -f nginx.yaml
```

---

## Imperative vs Declarative

| Style | Example | Best For |
|---|---|---|
| Imperative | `kubectl create deployment ...` | Quick learning |
| Declarative | `kubectl apply -f app.yaml` | Real work |

Kubernetes is built around declarative APIs. You describe what you want; controllers make it happen.

---

## Debugging Loop

Use this loop when something fails:

```bash
kubectl get pods
kubectl describe pod <pod>
kubectl logs <pod>
kubectl get events --sort-by=.lastTimestamp
```

Common beginner issues:

- image name typo
- container crashes immediately
- port mismatch
- labels do not match Service selector
- resource is created in a different namespace

---

## Skill Goal

Before contributing to Kubernetes code, you should be comfortable with:

- creating a Deployment
- exposing it with a Service
- reading Pod logs
- describing failed Pods
- understanding labels/selectors
- deleting and recreating resources

