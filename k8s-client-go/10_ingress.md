# Chapter 10: Ingress

## Ingress Interface

```go
clientset.NetworkingV1().Ingresses(namespace)
```

---

## Listing Ingresses

```go
import netv1 "k8s.io/api/networking/v1"

ingresses, err := clientset.NetworkingV1().Ingresses("default").List(ctx, metav1.ListOptions{})
for _, ing := range ingresses.Items {
    for _, rule := range ing.Spec.Rules {
        host := rule.Host
        if rule.HTTP == nil {
            continue
        }
        for _, path := range rule.HTTP.Paths {
            backend := path.Backend.Service
            fmt.Printf("%-20s %-25s %-10s → %s:%d\n",
                ing.Name, host, path.Path,
                backend.Name, backend.Port.Number)
        }
    }
}
```

---

## Creating a Basic Ingress

```go
import (
    netv1 "k8s.io/api/networking/v1"
)

ingressClass  := "nginx"
pathTypePrefix := netv1.PathTypePrefix

ingress := &netv1.Ingress{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-app-ingress",
        Namespace: "default",
        Annotations: map[string]string{
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    },
    Spec: netv1.IngressSpec{
        IngressClassName: &ingressClass,
        Rules: []netv1.IngressRule{
            {
                Host: "myapp.example.com",
                IngressRuleValue: netv1.IngressRuleValue{
                    HTTP: &netv1.HTTPIngressRuleValue{
                        Paths: []netv1.HTTPIngressPath{
                            {
                                Path:     "/api",
                                PathType: &pathTypePrefix,
                                Backend: netv1.IngressBackend{
                                    Service: &netv1.IngressServiceBackend{
                                        Name: "api-service",
                                        Port: netv1.ServiceBackendPort{Number: 80},
                                    },
                                },
                            },
                            {
                                Path:     "/",
                                PathType: &pathTypePrefix,
                                Backend: netv1.IngressBackend{
                                    Service: &netv1.IngressServiceBackend{
                                        Name: "frontend-service",
                                        Port: netv1.ServiceBackendPort{Number: 80},
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    },
}

created, err := clientset.NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})
fmt.Printf("Ingress created: %s\n", created.Name)
```

---

## Creating an Ingress with TLS

```go
ingress := &netv1.Ingress{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-app-tls-ingress",
        Namespace: "default",
        Annotations: map[string]string{
            "cert-manager.io/cluster-issuer":           "letsencrypt-prod",
            "nginx.ingress.kubernetes.io/ssl-redirect": "true",
        },
    },
    Spec: netv1.IngressSpec{
        IngressClassName: &ingressClass,
        TLS: []netv1.IngressTLS{
            {
                Hosts:      []string{"myapp.example.com", "api.example.com"},
                SecretName: "myapp-tls-secret",
            },
        },
        Rules: []netv1.IngressRule{
            {
                Host: "myapp.example.com",
                IngressRuleValue: netv1.IngressRuleValue{
                    HTTP: &netv1.HTTPIngressRuleValue{
                        Paths: []netv1.HTTPIngressPath{{
                            Path: "/", PathType: &pathTypePrefix,
                            Backend: netv1.IngressBackend{
                                Service: &netv1.IngressServiceBackend{
                                    Name: "frontend", Port: netv1.ServiceBackendPort{Number: 80},
                                },
                            },
                        }},
                    },
                },
            },
            {
                Host: "api.example.com",
                IngressRuleValue: netv1.IngressRuleValue{
                    HTTP: &netv1.HTTPIngressRuleValue{
                        Paths: []netv1.HTTPIngressPath{{
                            Path: "/", PathType: &pathTypePrefix,
                            Backend: netv1.IngressBackend{
                                Service: &netv1.IngressServiceBackend{
                                    Name: "api", Port: netv1.ServiceBackendPort{Number: 80},
                                },
                            },
                        }},
                    },
                },
            },
        },
    },
}

clientset.NetworkingV1().Ingresses("default").Create(ctx, ingress, metav1.CreateOptions{})
```

---

## Getting Ingress Address

```go
ing, err := clientset.NetworkingV1().Ingresses("default").Get(ctx, "my-app-ingress", metav1.GetOptions{})

for _, lb := range ing.Status.LoadBalancer.Ingress {
    addr := lb.IP
    if addr == "" {
        addr = lb.Hostname
    }
    fmt.Printf("Address: %s\n", addr)
}
```

---

## Updating an Ingress (Add a Host)

```go
import "k8s.io/client-go/util/retry"

err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
    ing, err := clientset.NetworkingV1().Ingresses("default").Get(ctx, "my-app-ingress", metav1.GetOptions{})
    if err != nil {
        return err
    }

    ing.Spec.Rules = append(ing.Spec.Rules, netv1.IngressRule{
        Host: "new-service.example.com",
        IngressRuleValue: netv1.IngressRuleValue{
            HTTP: &netv1.HTTPIngressRuleValue{
                Paths: []netv1.HTTPIngressPath{{
                    Path: "/", PathType: &pathTypePrefix,
                    Backend: netv1.IngressBackend{
                        Service: &netv1.IngressServiceBackend{
                            Name: "new-service", Port: netv1.ServiceBackendPort{Number: 80},
                        },
                    },
                }},
            },
        },
    })

    _, err = clientset.NetworkingV1().Ingresses("default").Update(ctx, ing, metav1.UpdateOptions{})
    return err
})
```

---

## Delete an Ingress

```go
err := clientset.NetworkingV1().Ingresses("default").Delete(ctx, "my-app-ingress", metav1.DeleteOptions{})
```

---

## Common Nginx Annotations

```go
annotations := map[string]string{
    // Rewrite
    "nginx.ingress.kubernetes.io/rewrite-target": "/$2",
    // Rate limiting
    "nginx.ingress.kubernetes.io/limit-rps":        "10",
    "nginx.ingress.kubernetes.io/limit-connections": "20",
    // Proxy timeouts
    "nginx.ingress.kubernetes.io/proxy-read-timeout": "300",
    "nginx.ingress.kubernetes.io/proxy-body-size":    "50m",
    // SSL
    "nginx.ingress.kubernetes.io/ssl-redirect":       "true",
    "nginx.ingress.kubernetes.io/force-ssl-redirect": "true",
    // CORS
    "nginx.ingress.kubernetes.io/enable-cors":         "true",
    "nginx.ingress.kubernetes.io/cors-allow-origin":   "https://myapp.com",
}
```

---

## PathType Reference

```go
netv1.PathTypeExact                  // exact path match
netv1.PathTypePrefix                 // prefix match
netv1.PathTypeImplementationSpecific // controller-specific
```
