# Chapter 2: Installation & Setup

## Prerequisites

- Go 1.21+
- A running Kubernetes cluster (Minikube, kind, k3s, or cloud)
- kubectl configured (for getting credentials)

---

## Initialize a Go Module

```bash
mkdir my-k8s-tool && cd my-k8s-tool
go mod init github.com/yourname/my-k8s-tool
```

---

## Installing client-go

```bash
# Install client-go (match version to your cluster's Kubernetes version)
go get k8s.io/client-go@v0.29.0

# Core API types (Pod, Deployment, etc.)
go get k8s.io/api@v0.29.0

# Meta types (ObjectMeta, ListOptions, etc.)
go get k8s.io/apimachinery@v0.29.0

# Tidy up
go mod tidy
```

### Version Compatibility

| Kubernetes Cluster | client-go version |
|---|---|
| 1.29.x | v0.29.x |
| 1.28.x | v0.28.x |
| 1.27.x | v0.27.x |
| 1.26.x | v0.26.x |

```bash
# Check your cluster version
kubectl version --short
```

---

## Minimal go.mod

```go
module github.com/yourname/my-k8s-tool

go 1.21

require (
    k8s.io/api            v0.29.0
    k8s.io/apimachinery   v0.29.0
    k8s.io/client-go      v0.29.0
)
```

---

## Project Layout

```
my-k8s-tool/
├── go.mod
├── go.sum
├── main.go
├── pkg/
│   ├── client/
│   │   └── client.go       ← client initialization
│   ├── controller/
│   │   └── controller.go   ← reconciliation logic
│   └── util/
│       └── wait.go         ← helper functions
└── config/
    └── rbac.yaml           ← ServiceAccount + Role for in-cluster
```

---

## Minimal main.go

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Load kubeconfig from default location
    kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %v", err)
    }

    // Create the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating clientset: %v", err)
    }

    // List pods in default namespace
    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error listing pods: %v", err)
    }

    fmt.Printf("Found %d pods in 'default' namespace:\n", len(pods.Items))
    for _, pod := range pods.Items {
        fmt.Printf("  - %s (%s)\n", pod.Name, pod.Status.Phase)
    }
}
```

```bash
go run main.go
```

---

## Key Import Paths

```go
import (
    // Client
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/dynamic"

    // Informers
    "k8s.io/client-go/informers"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"

    // Retry
    "k8s.io/client-go/util/retry"

    // API types
    corev1  "k8s.io/api/core/v1"
    appsv1  "k8s.io/api/apps/v1"
    batchv1 "k8s.io/api/batch/v1"
    netv1   "k8s.io/api/networking/v1"
    rbacv1  "k8s.io/api/rbac/v1"

    // Meta types
    metav1  "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/api/resource"
    "k8s.io/apimachinery/pkg/util/intstr"
    "k8s.io/apimachinery/pkg/runtime/schema"
)
```

---

## Fake Client for Testing

```go
import (
    "k8s.io/client-go/kubernetes/fake"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMyController(t *testing.T) {
    // Create a fake clientset pre-seeded with a pod
    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-pod",
            Namespace: "default",
        },
    }
    fakeClient := fake.NewSimpleClientset(pod)

    // Use fakeClient like a real clientset
    pods, err := fakeClient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    assert.NoError(t, err)
    assert.Equal(t, 1, len(pods.Items))
}
```

---

## Environment Variables

```go
// Useful env vars for configuring client behavior
os.Getenv("KUBECONFIG")             // custom kubeconfig path
os.Getenv("KUBERNETES_SERVICE_HOST") // set when running in-cluster
os.Getenv("KUBERNETES_SERVICE_PORT") // set when running in-cluster
```

---

## Build & Run

```bash
# Run locally
go run main.go

# Build binary
go build -o k8s-tool ./...

# Cross-compile for Linux (for Docker)
GOOS=linux GOARCH=amd64 go build -o k8s-tool-linux ./...

# Docker image
docker build -t myrepo/k8s-tool:latest .
docker push myrepo/k8s-tool:latest
```

### Minimal Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /k8s-tool ./...

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /k8s-tool /k8s-tool
ENTRYPOINT ["/k8s-tool"]
```
