# Chapter 20: Learning Roadmap — 5 Weeks to Kubernetes Proficiency

## Overview

A structured, hands-on 5-week plan that takes you from zero → running production-grade workloads, with concrete daily tasks and clear success criteria.

---

## Week 1: Foundations & First Deployment

**Goal:** Run a cluster locally, understand pods and services, deploy your first app.

| Day | Task |
|---|---|
| Day 1 | Install kubectl + Minikube. Run `minikube start`. Deploy nginx with `kubectl create deployment`. |
| Day 2 | Write your first Pod YAML. Apply it. Get logs, exec into it, delete it. |
| Day 3 | Write a Deployment YAML (3 replicas). Scale it. Update the image. Watch rolling update. |
| Day 4 | Create a ClusterIP Service. Connect to it from another pod. Try NodePort. |
| Day 5 | Create a ConfigMap and a Secret. Load them as env vars in a Deployment. |
| Day 6–7 | Deploy a full stack: Nginx + MySQL (each with a Service). Connect one to the other. |

**Success Criteria:**
- Cluster running locally
- Can apply, scale, update, and rollback a Deployment
- Understands Pod lifecycle and pod status debugging
- ConfigMap/Secret mounted in a running pod

---

## Week 2: Networking & Storage

**Goal:** Route external traffic, persist data, understand namespace isolation.

| Day | Task |
|---|---|
| Day 1–2 | Install ingress-nginx. Create an Ingress for your app. Test path-based and host-based routing. |
| Day 3 | Install cert-manager. Create a ClusterIssuer. Add TLS to your Ingress. |
| Day 4 | Create a PersistentVolumeClaim. Mount it in a MySQL StatefulSet. Verify data persists after pod restart. |
| Day 5 | Create `dev` and `staging` namespaces. Deploy app in each. Verify cross-namespace DNS. |
| Day 6 | Add ResourceQuota and LimitRange to each namespace. |
| Day 7 | Implement a NetworkPolicy: deny all → allow only frontend→backend. |

**Success Criteria:**
- External traffic reaches app via Ingress with TLS
- Database data survives pod restart
- Namespaces with quotas in place
- NetworkPolicy restricting pod communication

---

## Week 3: RBAC, Helm, and Autoscaling

**Goal:** Secure the cluster, package your app, automate scaling.

| Day | Task |
|---|---|
| Day 1–2 | Create a ServiceAccount. Bind a read-only Role to it. Test permissions with `kubectl auth can-i`. |
| Day 3 | Install Helm. Deploy MySQL with `helm install bitnami/mysql`. Upgrade it with new values. |
| Day 4 | Run `helm create my-app`. Convert your Deployment into a Helm chart with values.yaml. |
| Day 5 | Install metrics-server. Set up HPA for your app (target 70% CPU). |
| Day 6 | Generate load with a busybox loop. Watch HPA scale pods up. Stop load, watch scale down. |
| Day 7 | Set up RBAC: one ServiceAccount per namespace with least-privilege permissions. |

**Success Criteria:**
- RBAC with least-privilege ServiceAccounts
- App packaged as a Helm chart
- HPA scaling in response to load
- `kubectl auth can-i` returns correct results

---

## Week 4: Monitoring, Security & CI/CD

**Goal:** Observe what's happening, harden security, automate deployments.

| Day | Task |
|---|---|
| Day 1–2 | Install kube-prometheus-stack. Open Grafana. Import K8s cluster dashboard. |
| Day 3 | Create a PrometheusRule: alert when pod restarts > 0. Configure Alertmanager to print to a log. |
| Day 4 | Apply PodSecurity labels on namespaces (`restricted` mode). Fix any policy violations. |
| Day 5 | Add SecurityContext to your app: `runAsNonRoot`, `readOnlyRootFilesystem`, drop all capabilities. |
| Day 6 | Write a GitHub Actions workflow: test → build Docker image → push to GHCR → deploy to cluster. |
| Day 7 | Install ArgoCD. Create an Application pointing to your Helm chart repo. Watch GitOps auto-sync. |

**Success Criteria:**
- Grafana dashboards showing cluster health
- Alert fires on pod restart
- Pods passing `restricted` Pod Security Standard
- GitHub Actions CI/CD pipeline deploys to cluster

---

## Week 5: Production Readiness & Cloud

**Goal:** Production hardening, cloud deployment, advanced patterns.

| Day | Task |
|---|---|
| Day 1 | Create a cloud cluster (GKE free tier or EKS). Connect kubectl to it. Migrate your app. |
| Day 2 | Set up cloud-native storage (GKE persistent disk or EBS). Test persistence on cloud. |
| Day 3 | Implement blue-green deployment for your app. Test traffic switching. |
| Day 4 | Implement a CronJob for nightly database backup. Trigger it manually and verify backup file. |
| Day 5 | Install Loki + Grafana. Ship container logs to Loki. Query logs in Grafana. |
| Day 6 | Run kube-bench for CIS compliance checks. Fix at least 5 findings. |
| Day 7 | Full DR drill: delete your deployment, restore from backup, verify app recovers cleanly. |

**Success Criteria:**
- App running on real cloud cluster
- Blue-green deployment working
- Automated database backups running via CronJob
- Logs queryable in Grafana/Loki
- kube-bench score improved

---

## Summary Timeline

| Week | Focus | Deliverable |
|---|---|---|
| Week 1 | Pods, Deployments, Services | App running locally |
| Week 2 | Networking, Storage, Namespaces | Full stack with TLS + persistent DB |
| Week 3 | RBAC, Helm, Autoscaling | App as Helm chart + HPA working |
| Week 4 | Monitoring, Security, CI/CD | GitOps pipeline + Grafana dashboards |
| Week 5 | Cloud + Production hardening | Production-grade cloud deployment |

---

## Quick Reference: Most-Used Commands

```bash
# ---- Cluster ----
kubectl cluster-info
kubectl get nodes -o wide
kubectl top nodes

# ---- Pods ----
kubectl get pods -A
kubectl describe pod <pod>
kubectl logs <pod> -f
kubectl logs <pod> --previous
kubectl exec -it <pod> -- sh
kubectl top pod <pod>

# ---- Deployments ----
kubectl get deployments
kubectl scale deployment <name> --replicas=5
kubectl set image deployment/<name> <container>=<image>:<tag>
kubectl rollout status deployment/<name>
kubectl rollout history deployment/<name>
kubectl rollout undo deployment/<name>

# ---- Services ----
kubectl get svc,ep
kubectl port-forward svc/<name> 8080:80

# ---- Debugging ----
kubectl get events --sort-by='.lastTimestamp'
kubectl describe <resource> <name>
kubectl get all -n <namespace>
kubectl run debug --image=busybox --rm -it --restart=Never -- sh

# ---- Helm ----
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install <release> <chart> -f values.yaml
helm upgrade <release> <chart> -f values.yaml
helm rollback <release>
helm list

# ---- Config/Secrets ----
kubectl create configmap <name> --from-literal=key=value
kubectl create secret generic <name> --from-literal=key=value
kubectl get secret <name> -o jsonpath='{.data.key}' | base64 -d

# ---- RBAC ----
kubectl auth can-i <verb> <resource>
kubectl auth can-i get pods --as=user@example.com

# ---- Storage ----
kubectl get pv,pvc
kubectl describe pvc <name>
```

---

## Recommended Resources

| Resource | URL |
|---|---|
| Official K8s Docs | `kubernetes.io/docs` |
| Kubernetes Playground | `labs.play-with-k8s.com` |
| Helm Charts | `artifacthub.io` |
| CNCF Landscape | `landscape.cncf.io` |
| kube.academy (free) | `kube.academy` |
| Killer.sh CKA Practice | `killer.sh` |
| The Kubernetes Book | `thekubernetesbook.com` |

---

## Certifications Path

| Certification | Level | Focus |
|---|---|---|
| **KCNA** | Beginner | Cloud Native fundamentals |
| **CKA** | Intermediate | Cluster administration |
| **CKAD** | Intermediate | Application development |
| **CKS** | Advanced | Security specialist |

**Study tips for CKA/CKAD:**
- All exams are hands-on (open book, real cluster)
- Master `kubectl` — speed matters
- Practice on `killer.sh` (exam simulator)
- Read `kubernetes.io/docs` — that's all you get during the exam
