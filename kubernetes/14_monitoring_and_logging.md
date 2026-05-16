# Chapter 14: Monitoring & Logging

## Overview

A production Kubernetes cluster needs:
- **Metrics** — CPU, memory, request rates, error rates
- **Logs** — container stdout/stderr, application logs
- **Alerts** — notify when something goes wrong
- **Dashboards** — visualize everything at a glance

**Standard Stack:**
```
Prometheus (collect metrics)
    + Grafana (visualize)
    + Alertmanager (alerting)
    + Loki (log aggregation)
    + Promtail/Fluentd (log shipping)
```

---

## kubectl Logs

### Basic Log Commands

```bash
# Logs from a pod (latest container)
kubectl logs my-pod

# Follow logs in real time
kubectl logs -f my-pod

# Last 100 lines
kubectl logs my-pod --tail=100

# Logs since a time
kubectl logs my-pod --since=1h
kubectl logs my-pod --since-time=2025-01-01T10:00:00Z

# Logs from a specific container in a multi-container pod
kubectl logs my-pod -c nginx
kubectl logs my-pod -c php-fpm

# Previous container (after crash)
kubectl logs my-pod --previous

# Logs from all pods matching a label
kubectl logs -l app=my-app --all-containers=true

# Logs from a deployment
kubectl logs deployment/my-app
kubectl logs deployment/my-app -c my-app -f
```

---

## metrics-server

**metrics-server** enables `kubectl top` commands for basic CPU/memory metrics.

```bash
# Install
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# On Minikube
minikube addons enable metrics-server

# Use
kubectl top nodes
kubectl top pods -A
kubectl top pods -n production --sort-by=cpu
```

---

## Prometheus + Grafana Stack

### Install with Helm (Recommended)

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install kube-prometheus-stack (includes Prometheus, Grafana, Alertmanager)
helm install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace \
  --set grafana.adminPassword=admin123 \
  --set prometheus.prometheusSpec.retention=30d

# Check pods
kubectl get pods -n monitoring
```

### Access Grafana

```bash
kubectl port-forward -n monitoring svc/monitoring-grafana 3000:80
# Open: http://localhost:3000
# Login: admin / admin123
```

### Access Prometheus

```bash
kubectl port-forward -n monitoring svc/monitoring-kube-prometheus-prometheus 9090:9090
# Open: http://localhost:9090
```

---

## Prometheus Basics

### Key Concepts

| Concept | Description |
|---|---|
| **Scrape** | Prometheus pulls metrics from targets at intervals |
| **Metric** | Named time-series data point |
| **Label** | Key-value tags on metrics |
| **PromQL** | Query language for metrics |
| **Exporter** | App that exposes metrics at `/metrics` |

### Common PromQL Queries

```promql
# CPU usage by pod
sum(rate(container_cpu_usage_seconds_total[5m])) by (pod)

# Memory usage by namespace
sum(container_memory_usage_bytes) by (namespace)

# HTTP error rate
rate(http_requests_total{status=~"5.."}[5m])

# Pod restart count
kube_pod_container_status_restarts_total

# Node disk pressure
node_filesystem_avail_bytes / node_filesystem_size_bytes < 0.10

# Request latency p99
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### ServiceMonitor (Auto-scrape your apps)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: my-app-monitor
  namespace: production
spec:
  selector:
    matchLabels:
      app: my-app
  endpoints:
    - port: http
      path: /metrics
      interval: 30s
```

Your app must expose a `/metrics` endpoint in Prometheus format (use a Prometheus client library).

---

## Alerting with Alertmanager

### PrometheusRule Example

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: my-app-alerts
  namespace: production
spec:
  groups:
    - name: my-app
      rules:
        - alert: HighErrorRate
          expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
          for: 2m
          labels:
            severity: critical
          annotations:
            summary: "High error rate on {{ $labels.pod }}"
            description: "Error rate is {{ $value | humanizePercentage }}"

        - alert: PodCrashLooping
          expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Pod {{ $labels.pod }} is crash looping"

        - alert: HighMemoryUsage
          expr: |
            container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.9
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Container using >90% memory limit"
```

---

## Loki — Log Aggregation

**Loki** is Grafana's log aggregation system — lightweight, stores logs by labels, not full-text indexing.

### Install Loki Stack

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki-stack \
  -n monitoring \
  --set grafana.enabled=false \
  --set promtail.enabled=true
```

**Promtail** auto-discovers pods and ships their logs to Loki.

### Query Logs in Grafana (LogQL)

```logql
# All logs from my-app in production
{app="my-app", namespace="production"}

# Errors only
{app="my-app"} |= "ERROR"

# Count errors per minute
rate({app="my-app"} |= "ERROR" [1m])

# Parse JSON logs
{app="my-app"} | json | level="error"
```

---

## Pre-built Grafana Dashboards

Import these from grafana.com/dashboards:

| Dashboard | ID | What It Shows |
|---|---|---|
| Kubernetes cluster overview | 315 | Node CPU/memory, pod counts |
| Node Exporter | 1860 | Per-node CPU, memory, disk, network |
| Pod monitoring | 6417 | Pod CPU/memory/restarts |
| Nginx Ingress | 9614 | Request rates, errors, latency |
| MySQL | 7362 | Queries, connections, slow queries |
| Redis | 11835 | Operations, memory, hit rate |

```bash
# In Grafana UI:
# Dashboards → Import → Enter ID → Load → Select Prometheus data source → Import
```

---

## Kubernetes Events

Events are lightweight logs of what happened to K8s objects:

```bash
# All events in a namespace
kubectl get events -n production

# Sort by time
kubectl get events -n production --sort-by='.lastTimestamp'

# Events for a specific pod
kubectl describe pod my-pod -n production  # Shows Events at bottom

# Watch events in real time
kubectl get events -w -n production
```

---

## Monitoring Checklist for Production

- [ ] metrics-server installed (`kubectl top` works)
- [ ] Prometheus scraping cluster metrics
- [ ] Grafana dashboards for cluster + applications
- [ ] Alertmanager configured (Slack/PagerDuty/email)
- [ ] Loki/EFK for log aggregation
- [ ] Pod restart alerts set up
- [ ] Node disk/memory pressure alerts
- [ ] Application-level metrics exposed at `/metrics`
- [ ] Log retention policy defined
- [ ] On-call runbook linked in alert annotations
