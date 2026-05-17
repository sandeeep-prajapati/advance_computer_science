# Chapter 3: Connecting to a Cluster

## rest.Config — The Foundation

All client-go connections start with a `*rest.Config`. Everything else builds on top of it.

```go
import (
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/kubernetes"
)
```

---

## Method 1: kubeconfig File (Local Development)

```go
import (
    "os"
    "path/filepath"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/kubernetes"
)

func newClientFromKubeconfig() (*kubernetes.Clientset, error) {
    home, _ := os.UserHomeDir()
    kubeconfig := filepath.Join(home, ".kube", "config")

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        return nil, err
    }

    return kubernetes.NewForConfig(config)
}
```

### Custom kubeconfig path from flag

```go
import "flag"

var kubeconfig *string

func init() {
    if home, _ := os.UserHomeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "path to kubeconfig")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig")
    }
    flag.Parse()
}

config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
```

---

## Method 2: In-Cluster Config (Pod Inside Kubernetes)

```go
import "k8s.io/client-go/rest"

func newInClusterClient() (*kubernetes.Clientset, error) {
    // Reads token + CA from /var/run/secrets/kubernetes.io/serviceaccount/
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }
    return kubernetes.NewForConfig(config)
}
```

---

## Method 3: Auto-detect (In-cluster or kubeconfig)

The most common production pattern — works both locally and in-cluster:

```go
func buildClient() (*kubernetes.Clientset, error) {
    // Try in-cluster first
    config, err := rest.InClusterConfig()
    if err != nil {
        // Fall back to kubeconfig
        home, _ := os.UserHomeDir()
        kubeconfig := filepath.Join(home, ".kube", "config")
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err != nil {
            return nil, fmt.Errorf("cannot build config: %w", err)
        }
    }
    return kubernetes.NewForConfig(config)
}
```

---

## Method 4: Specific Context from kubeconfig

```go
func clientForContext(kubeconfigPath, context string) (*kubernetes.Clientset, error) {
    loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
    loadingRules.ExplicitPath = kubeconfigPath

    configOverrides := &clientcmd.ConfigOverrides{
        CurrentContext: context,
    }

    config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
        loadingRules, configOverrides,
    ).ClientConfig()
    if err != nil {
        return nil, err
    }

    return kubernetes.NewForConfig(config)
}

// Usage
cs, err := clientForContext("~/.kube/config", "prod-cluster")
```

---

## Method 5: Manual REST Config (Token / Cert)

```go
// Bearer token
config := &rest.Config{
    Host:        "https://k8s.example.com:6443",
    BearerToken: "eyJhbGci...",
    TLSClientConfig: rest.TLSClientConfig{
        CAFile: "/path/to/ca.crt",
    },
}

// Client certificate
config := &rest.Config{
    Host: "https://k8s.example.com:6443",
    TLSClientConfig: rest.TLSClientConfig{
        CAFile:   "/path/to/ca.crt",
        CertFile: "/path/to/client.crt",
        KeyFile:  "/path/to/client.key",
    },
}

clientset, err := kubernetes.NewForConfig(config)
```

---

## Multi-Cluster Management

```go
type MultiClusterManager struct {
    clients map[string]*kubernetes.Clientset
}

func NewMultiClusterManager(contexts map[string]string) (*MultiClusterManager, error) {
    m := &MultiClusterManager{clients: make(map[string]*kubernetes.Clientset)}

    for name, ctx := range contexts {
        cs, err := clientForContext("~/.kube/config", ctx)
        if err != nil {
            return nil, fmt.Errorf("cluster %s: %w", name, err)
        }
        m.clients[name] = cs
    }
    return m, nil
}

func (m *MultiClusterManager) For(cluster string) *kubernetes.Clientset {
    return m.clients[cluster]
}

// Usage
mgr, _ := NewMultiClusterManager(map[string]string{
    "us-east": "eks-us-east-1",
    "eu-west": "gke-europe-west1",
})

pods, _ := mgr.For("us-east").CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
```

---

## Tuning the REST Config

```go
config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)

// Increase QPS and Burst for high-throughput controllers
config.QPS   = 100    // default: 5 requests/sec
config.Burst = 200    // default: 10 burst

// Timeout per request
config.Timeout = 30 * time.Second

clientset, _ := kubernetes.NewForConfig(config)
```

---

## Verifying the Connection

```go
func verifyConnection(clientset *kubernetes.Clientset) error {
    version, err := clientset.Discovery().ServerVersion()
    if err != nil {
        return fmt.Errorf("cannot reach API server: %w", err)
    }
    fmt.Printf("Connected to Kubernetes %s\n", version.GitVersion)
    return nil
}
```

---

## Connection Troubleshooting

| Error | Cause | Fix |
|---|---|---|
| `connection refused` | Wrong host/port | Check `kubectl cluster-info` |
| `x509: certificate signed by unknown authority` | Wrong CA | Verify `CAFile` or `InsecureSkipVerify` (dev only) |
| `Unauthorized (401)` | Bad/expired token | Regenerate token |
| `Forbidden (403)` | No RBAC permission | Grant Role to ServiceAccount |
| `dial tcp: no route to host` | Firewall / VPN | Check network access |

```go
// Dev only — skip TLS verification
config.TLSClientConfig = rest.TLSClientConfig{Insecure: true}
```
