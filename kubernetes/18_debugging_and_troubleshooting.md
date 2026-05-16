# Chapter 18: Debugging & Troubleshooting

## The Debugging Mindset

When something breaks in Kubernetes, work top-down:

```
1. What is the symptom?   (app down, slow, 502, pod crash)
2. Check the pod status   (kubectl get pods)
3. Describe the pod       (kubectl describe pod)
4. Check logs             (kubectl logs)
5. Check events           (kubectl get events)
6. Check service/ingress  (kubectl describe svc/ingress)
7. Check resource limits  (kubectl top pods/nodes)
8. Exec into the pod      (kubectl exec -it ... -- sh)
```

---

## Common Pod States

| State | Cause | Debug Step |
|---|---|---|
| `Pending` | No node available, PVC not bound | `kubectl describe pod` → Events |
| `CrashLoopBackOff` | Container keeps crashing | `kubectl logs pod --previous` |
| `OOMKilled` | Container exceeded memory limit | Increase limits or fix memory leak |
| `ImagePullBackOff` | Can't pull the image | Check image name, registry auth |
| `ErrImagePull` | First image pull failure | Same as above |
| `CreateContainerConfigError` | Bad config (missing ConfigMap/Secret) | `kubectl describe pod` → Events |
| `Error` | Container exited with non-zero | `kubectl logs` to see error |
| `Terminating` (stuck) | Finalizers blocking deletion | Remove finalizers |

---

## CrashLoopBackOff

```bash
# Step 1: Check current logs
kubectl logs my-pod

# Step 2: Check logs from the PREVIOUS (crashed) container
kubectl logs my-pod --previous

# Step 3: Describe to see last state + exit code
kubectl describe pod my-pod
# Look for:
# Last State: Terminated
#   Reason: Error / OOMKilled
#   Exit Code: 1 (or 137 for OOM)

# Step 4: Test your container command locally
docker run --rm -e DATABASE_URL=... my-app:latest
```

---

## ImagePullBackOff / ErrImagePull

```bash
kubectl describe pod my-pod
# Look for Events section:
# Failed to pull image "registry/my-app:latest": ...

# Common causes and fixes:
# 1. Wrong image name/tag
kubectl set image deployment/my-app my-app=registry/my-app:v1.2.3

# 2. Private registry — need imagePullSecrets
kubectl create secret docker-registry regcred \
  --docker-server=registry.example.com \
  --docker-username=myuser \
  --docker-password=mypassword

# Add to deployment:
spec:
  imagePullSecrets:
    - name: regcred

# 3. For kind: load image manually
kind load docker-image my-app:latest --name my-cluster

# 4. For minikube: load image manually
minikube image load my-app:latest
```

---

## OOMKilled

Container exceeded its memory limit — kernel kills it with exit code 137.

```bash
kubectl describe pod my-pod
# Look for:
# Last State: Terminated
#   Reason: OOMKilled
#   Exit Code: 137

# Fix: increase memory limits
kubectl set resources deployment/my-app \
  --limits=memory=1Gi

# Or edit the deployment
kubectl edit deployment my-app
# Change limits.memory to a higher value

# Find actual usage to set appropriate limits
kubectl top pod my-pod
```

---

## Pending Pods

```bash
kubectl describe pod my-pod
# Look at Events section

# Common causes:
# 1. Insufficient resources
kubectl describe nodes | grep -A5 "Allocated resources"

# 2. PVC not bound
kubectl get pvc
kubectl describe pvc my-pvc

# 3. Node selector doesn't match
kubectl get nodes --show-labels
kubectl describe pod my-pod | grep -A5 "Node-Selectors"

# 4. Taints blocking scheduling
kubectl describe nodes | grep Taints

# 5. No nodes in the cluster
kubectl get nodes
```

---

## Service Not Routing Traffic

```bash
# Step 1: Check if pods are Running
kubectl get pods -l app=my-app

# Step 2: Check service
kubectl get svc my-service
kubectl describe svc my-service

# Step 3: Check endpoints (are pods being selected?)
kubectl get endpoints my-service
# If ENDPOINTS is <none>: label selector doesn't match any pods

# Step 4: Verify labels match
kubectl get pods --show-labels
# Compare with service selector

# Step 5: Test service DNS from inside cluster
kubectl run test --image=busybox --rm -it --restart=Never -- \
  nslookup my-service
kubectl run test --image=busybox --rm -it --restart=Never -- \
  wget -qO- http://my-service

# Step 6: Test directly on pod port
kubectl port-forward pod/my-pod 8080:8080
curl http://localhost:8080
```

---

## Ingress Not Working

```bash
# Check ingress resource
kubectl describe ingress my-ingress

# Check if ingress controller is running
kubectl get pods -n ingress-nginx

# Check ingress controller logs
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller

# Verify backend service exists
kubectl get svc my-backend-service

# Check if backend pods are healthy
kubectl get endpoints my-backend-service

# Test path routing
curl -H "Host: myapp.example.com" http://<ingress-ip>/api
```

---

## Persistent Volume Issues

```bash
# PVC stuck in Pending
kubectl describe pvc my-pvc
# Events: no matching StorageClass, no PV available

# List storage classes
kubectl get sc

# Check if PV is Available
kubectl get pv

# Force delete stuck PVC (remove finalizers)
kubectl patch pvc my-pvc -p '{"metadata":{"finalizers":null}}'

# Check mounted volume inside pod
kubectl exec -it my-pod -- ls -la /data
kubectl exec -it my-pod -- df -h
```

---

## Debug Container (kubectl debug)

```bash
# Add a debug container to a running pod (K8s 1.23+)
kubectl debug -it my-pod --image=busybox --target=my-container

# Create a copy of a pod with debug tools
kubectl debug my-pod -it --copy-to=debug-pod --image=ubuntu:22.04

# Debug a node (runs privileged pod on the node)
kubectl debug node/my-node -it --image=ubuntu:22.04
```

---

## Common Application Debugging

### Test DNS Resolution

```bash
kubectl run dns-test --image=busybox --rm -it --restart=Never -- \
  nslookup kubernetes.default.svc.cluster.local
```

### Test HTTP Connectivity

```bash
kubectl run curl-test --image=curlimages/curl --rm -it --restart=Never -- \
  curl -v http://my-service:80/health
```

### Check Environment Variables

```bash
kubectl exec my-pod -- env | sort
```

### Inspect Mounted Files

```bash
kubectl exec my-pod -- ls -la /etc/config/
kubectl exec my-pod -- cat /etc/config/app.conf
```

---

## Node Debugging

```bash
# List nodes
kubectl get nodes

# Node details
kubectl describe node my-node

# Node resource usage
kubectl top node my-node

# Check node conditions
kubectl get nodes -o json | jq '.items[].status.conditions'

# SSH into a node (if you have access)
ssh user@<node-ip>

# Or use kubectl debug node
kubectl debug node/my-node -it --image=ubuntu:22.04
```

---

## Stuck Terminating Pods/Namespaces

```bash
# Pod stuck Terminating
kubectl delete pod my-pod --grace-period=0 --force

# Namespace stuck Terminating — remove finalizers
kubectl get namespace my-namespace -o json \
  | jq '.spec.finalizers = []' \
  | kubectl replace --raw "/api/v1/namespaces/my-namespace/finalize" -f -
```

---

## Debugging Cheatsheet

```bash
# Pod basics
kubectl get pods -A                          # All pods, all namespaces
kubectl get pods -o wide                     # Show node and IP
kubectl describe pod <pod>                   # Full details + events
kubectl logs <pod>                           # Logs
kubectl logs <pod> --previous                # Previous container logs
kubectl logs <pod> -c <container>            # Specific container
kubectl exec -it <pod> -- sh                 # Shell into container
kubectl top pod <pod>                        # CPU/memory usage

# Service/network
kubectl get svc,ep -n <namespace>            # Services + endpoints
kubectl describe svc <service>               # Service details
kubectl port-forward svc/<svc> 8080:80       # Forward to localhost

# Events
kubectl get events -n <namespace> --sort-by='.lastTimestamp'

# Nodes
kubectl get nodes                            # Node status
kubectl describe node <node>                 # Node details + events
kubectl top nodes                            # Node resource usage
kubectl drain <node> --ignore-daemonsets     # Drain node for maintenance

# YAML output
kubectl get pod <pod> -o yaml                # Full YAML
kubectl get all -n <namespace> -o yaml       # All resources as YAML
```
