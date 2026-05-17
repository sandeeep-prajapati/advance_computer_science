# Chapter 14: Custom Resources (CRDs)

## Three Ways to Access CRDs

| Approach | Use When | Type Safety |
|---|---|---|
| Dynamic client | Unknown schema, generic tools | None (unstructured) |
| Generated typed client | Known CRD, production operator | Full |
| controller-runtime | Building operators | Full + helpers |

---

## Dynamic Client

The dynamic client works with any resource — standard or custom — as `unstructured.Unstructured`.

```go
import (
    "k8s.io/client-go/dynamic"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

dynClient, err := dynamic.NewForConfig(config)
if err != nil {
    log.Fatal(err)
}
```

### List Custom Resources

```go
// Define the GVR (GroupVersionResource)
certificateGVR := schema.GroupVersionResource{
    Group:    "cert-manager.io",
    Version:  "v1",
    Resource: "certificates",
}

list, err := dynClient.Resource(certificateGVR).Namespace("default").List(ctx, metav1.ListOptions{})
if err != nil {
    log.Fatal(err)
}

for _, item := range list.Items {
    name := item.GetName()
    // Access spec fields
    dnsNames, _, _ := unstructured.NestedStringSlice(item.Object, "spec", "dnsNames")
    ready := "Unknown"
    conditions, _, _ := unstructured.NestedSlice(item.Object, "status", "conditions")
    for _, c := range conditions {
        cond := c.(map[string]interface{})
        if cond["type"] == "Ready" {
            ready = cond["status"].(string)
        }
    }
    fmt.Printf("%-30s ready=%-8s dnsNames=%v\n", name, ready, dnsNames)
}
```

### Get a Custom Resource

```go
cert, err := dynClient.Resource(certificateGVR).Namespace("default").Get(ctx, "myapp-tls", metav1.GetOptions{})

fmt.Println("Name:", cert.GetName())
fmt.Println("Labels:", cert.GetLabels())

secretName, _, _ := unstructured.NestedString(cert.Object, "spec", "secretName")
fmt.Println("SecretName:", secretName)
```

### Create a Custom Resource

```go
cert := &unstructured.Unstructured{
    Object: map[string]interface{}{
        "apiVersion": "cert-manager.io/v1",
        "kind":       "Certificate",
        "metadata": map[string]interface{}{
            "name":      "myapp-tls",
            "namespace": "default",
        },
        "spec": map[string]interface{}{
            "secretName": "myapp-tls-secret",
            "issuerRef": map[string]interface{}{
                "name": "letsencrypt-prod",
                "kind": "ClusterIssuer",
            },
            "dnsNames": []interface{}{
                "myapp.example.com",
                "www.myapp.example.com",
            },
            "duration":    "2160h",
            "renewBefore": "360h",
        },
    },
}

created, err := dynClient.Resource(certificateGVR).Namespace("default").Create(ctx, cert, metav1.CreateOptions{})
fmt.Printf("Certificate created: %s\n", created.GetName())
```

### Update a Custom Resource

```go
import "k8s.io/client-go/util/retry"

err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
    cert, err := dynClient.Resource(certificateGVR).Namespace("default").Get(ctx, "myapp-tls", metav1.GetOptions{})
    if err != nil {
        return err
    }

    // Update a field
    unstructured.SetNestedField(cert.Object, "720h", "spec", "renewBefore")

    _, err = dynClient.Resource(certificateGVR).Namespace("default").Update(ctx, cert, metav1.UpdateOptions{})
    return err
})
```

### Delete a Custom Resource

```go
err := dynClient.Resource(certificateGVR).Namespace("default").Delete(ctx, "myapp-tls", metav1.DeleteOptions{})
```

---

## Listing CRDs

```go
import (
    apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
    apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// go get k8s.io/apiextensions-apiserver@v0.29.0
apiextCS, err := apiextclient.NewForConfig(config)

crds, err := apiextCS.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
for _, crd := range crds.Items {
    fmt.Printf("%-55s %-30s %s\n", crd.Name, crd.Spec.Group, crd.Spec.Scope)
}
```

---

## Defining a CRD via Go

```go
import apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

crd := &apiextv1.CustomResourceDefinition{
    ObjectMeta: metav1.ObjectMeta{Name: "myapps.mycompany.io"},
    Spec: apiextv1.CustomResourceDefinitionSpec{
        Group: "mycompany.io",
        Scope: apiextv1.NamespaceScoped,
        Names: apiextv1.CustomResourceDefinitionNames{
            Plural:   "myapps",
            Singular: "myapp",
            Kind:     "MyApp",
        },
        Versions: []apiextv1.CustomResourceDefinitionVersion{{
            Name:    "v1",
            Served:  true,
            Storage: true,
            Schema: &apiextv1.CustomResourceValidation{
                OpenAPIV3Schema: &apiextv1.JSONSchemaProps{
                    Type: "object",
                    Properties: map[string]apiextv1.JSONSchemaProps{
                        "spec": {
                            Type: "object",
                            Properties: map[string]apiextv1.JSONSchemaProps{
                                "replicas": {Type: "integer"},
                                "image":    {Type: "string"},
                            },
                        },
                    },
                },
            },
        }},
    },
}

apiextCS.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crd, metav1.CreateOptions{})
```

---

## Watch Custom Resources

```go
watcher, err := dynClient.Resource(certificateGVR).Namespace("default").Watch(ctx, metav1.ListOptions{})
if err != nil {
    log.Fatal(err)
}
defer watcher.Stop()

for event := range watcher.ResultChan() {
    obj  := event.Object.(*unstructured.Unstructured)
    name := obj.GetName()

    ready := "Unknown"
    conditions, _, _ := unstructured.NestedSlice(obj.Object, "status", "conditions")
    for _, c := range conditions {
        if cmap, ok := c.(map[string]interface{}); ok {
            if cmap["type"] == "Ready" {
                ready = fmt.Sprintf("%v", cmap["status"])
            }
        }
    }

    fmt.Printf("[%s] Certificate: %s — Ready: %s\n", event.Type, name, ready)
}
```

---

## ArgoCD Application via Dynamic Client

```go
argoAppGVR := schema.GroupVersionResource{
    Group:    "argoproj.io",
    Version:  "v1alpha1",
    Resource: "applications",
}

app := &unstructured.Unstructured{
    Object: map[string]interface{}{
        "apiVersion": "argoproj.io/v1alpha1",
        "kind":       "Application",
        "metadata": map[string]interface{}{
            "name":      "my-app",
            "namespace": "argocd",
        },
        "spec": map[string]interface{}{
            "project": "default",
            "source": map[string]interface{}{
                "repoURL":        "https://github.com/myorg/myrepo",
                "targetRevision": "HEAD",
                "path":           "k8s/production",
            },
            "destination": map[string]interface{}{
                "server":    "https://kubernetes.default.svc",
                "namespace": "production",
            },
            "syncPolicy": map[string]interface{}{
                "automated": map[string]interface{}{
                    "prune":    true,
                    "selfHeal": true,
                },
            },
        },
    },
}

dynClient.Resource(argoAppGVR).Namespace("argocd").Create(ctx, app, metav1.CreateOptions{})
```

---

## Common GVRs

```go
// cert-manager
schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}
schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "clusterissuers"}

// ArgoCD
schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}

// Prometheus Operator
schema.GroupVersionResource{Group: "monitoring.coreos.com", Version: "v1", Resource: "prometheusrules"}

// Flux
schema.GroupVersionResource{Group: "helm.toolkit.fluxcd.io", Version: "v2beta1", Resource: "helmreleases"}

// Traefik
schema.GroupVersionResource{Group: "traefik.containo.us", Version: "v1alpha1", Resource: "ingressroutes"}

// VPA
schema.GroupVersionResource{Group: "autoscaling.k8s.io", Version: "v1", Resource: "verticalpodautoscalers"}
```
