# Chapter 6: Services & Networking

## Service Interface

```go
clientset.CoreV1().Services(namespace)
```

---

## Listing Services

```go
services, err := clientset.CoreV1().Services("default").List(ctx, metav1.ListOptions{})
if err != nil {
    log.Fatal(err)
}

for _, svc := range services.Items {
    ports := make([]string, len(svc.Spec.Ports))
    for i, p := range svc.Spec.Ports {
        ports[i] = fmt.Sprintf("%d/%s", p.Port, p.Protocol)
    }
    fmt.Printf("%-35s %-15s %-16s %s\n",
        svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, strings.Join(ports, ","))
}
```

---

## Creating a ClusterIP Service

```go
import (
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

svc := &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-app-svc",
        Namespace: "default",
        Labels:    map[string]string{"app": "my-app"},
    },
    Spec: corev1.ServiceSpec{
        Type:     corev1.ServiceTypeClusterIP,
        Selector: map[string]string{"app": "my-app"},
        Ports: []corev1.ServicePort{
            {
                Name:       "http",
                Port:       80,
                TargetPort: intstr.FromInt(8080),
                Protocol:   corev1.ProtocolTCP,
            },
        },
    },
}

created, err := clientset.CoreV1().Services("default").Create(ctx, svc, metav1.CreateOptions{})
fmt.Printf("Service created: %s (%s)\n", created.Name, created.Spec.ClusterIP)
```

---

## Creating a NodePort Service

```go
nodePort := int32(30080)

svc := &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{Name: "my-app-nodeport", Namespace: "default"},
    Spec: corev1.ServiceSpec{
        Type:     corev1.ServiceTypeNodePort,
        Selector: map[string]string{"app": "my-app"},
        Ports: []corev1.ServicePort{
            {
                Port:       80,
                TargetPort: intstr.FromInt(8080),
                NodePort:   nodePort,
            },
        },
    },
}

clientset.CoreV1().Services("default").Create(ctx, svc, metav1.CreateOptions{})
```

---

## Creating a LoadBalancer Service

```go
svc := &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-app-lb",
        Namespace: "default",
        Annotations: map[string]string{
            "service.beta.kubernetes.io/aws-load-balancer-type": "nlb",
        },
    },
    Spec: corev1.ServiceSpec{
        Type:     corev1.ServiceTypeLoadBalancer,
        Selector: map[string]string{"app": "my-app"},
        Ports: []corev1.ServicePort{
            {Port: 80,  TargetPort: intstr.FromInt(8080)},
            {Port: 443, TargetPort: intstr.FromInt(8443)},
        },
    },
}

clientset.CoreV1().Services("default").Create(ctx, svc, metav1.CreateOptions{})
```

### Wait for LoadBalancer IP

```go
func waitForLoadBalancerIP(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (string, error) {
    for {
        svc, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return "", err
        }

        if len(svc.Status.LoadBalancer.Ingress) > 0 {
            ingress := svc.Status.LoadBalancer.Ingress[0]
            if ingress.IP != "" {
                return ingress.IP, nil
            }
            if ingress.Hostname != "" {
                return ingress.Hostname, nil
            }
        }

        select {
        case <-ctx.Done():
            return "", ctx.Err()
        case <-time.After(5 * time.Second):
            fmt.Println("Waiting for external IP...")
        }
    }
}
```

---

## Creating a Headless Service

```go
svc := &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{Name: "mysql-headless"},
    Spec: corev1.ServiceSpec{
        ClusterIP: "None",   // headless
        Selector:  map[string]string{"app": "mysql"},
        Ports: []corev1.ServicePort{
            {Port: 3306, Name: "mysql"},
        },
    },
}
```

---

## Creating an ExternalName Service

```go
svc := &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{Name: "external-db"},
    Spec: corev1.ServiceSpec{
        Type:         corev1.ServiceTypeExternalName,
        ExternalName: "mydb.rds.amazonaws.com",
    },
}
```

---

## Getting Service Details

```go
svc, err := clientset.CoreV1().Services("default").Get(ctx, "my-app-svc", metav1.GetOptions{})

fmt.Printf("ClusterIP: %s\n", svc.Spec.ClusterIP)
fmt.Printf("Type:      %s\n", svc.Spec.Type)

for _, p := range svc.Spec.Ports {
    fmt.Printf("Port: %d → %s\n", p.Port, p.TargetPort.String())
}

// External IPs for LoadBalancer
for _, ing := range svc.Status.LoadBalancer.Ingress {
    fmt.Printf("External: %s%s\n", ing.IP, ing.Hostname)
}
```

---

## Working with Endpoints

```go
ep, err := clientset.CoreV1().Endpoints("default").Get(ctx, "my-app-svc", metav1.GetOptions{})

for _, subset := range ep.Subsets {
    port := subset.Ports[0].Port
    for _, addr := range subset.Addresses {
        fmt.Printf("Ready:    %s:%d (pod=%s)\n", addr.IP, port, addr.TargetRef.Name)
    }
    for _, addr := range subset.NotReadyAddresses {
        fmt.Printf("NotReady: %s:%d\n", addr.IP, port)
    }
}
```

---

## Update a Service

```go
import "k8s.io/client-go/util/retry"

err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
    svc, err := clientset.CoreV1().Services("default").Get(ctx, "my-app-svc", metav1.GetOptions{})
    if err != nil {
        return err
    }
    // Add annotation
    svc.Annotations["custom.io/updated"] = time.Now().Format(time.RFC3339)
    _, err = clientset.CoreV1().Services("default").Update(ctx, svc, metav1.UpdateOptions{})
    return err
})
```

---

## Delete a Service

```go
err := clientset.CoreV1().Services("default").Delete(ctx, "my-app-svc", metav1.DeleteOptions{})
```

---

## NetworkPolicy (Read)

```go
netPolicies, _ := clientset.NetworkingV1().NetworkPolicies("default").List(ctx, metav1.ListOptions{})

for _, np := range netPolicies.Items {
    fmt.Printf("Policy: %s → podSelector: %v\n",
        np.Name, np.Spec.PodSelector.MatchLabels)
}
```
