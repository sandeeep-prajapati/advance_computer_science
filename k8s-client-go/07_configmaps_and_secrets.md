# Chapter 7: ConfigMaps & Secrets

## ConfigMap Interface

```go
clientset.CoreV1().ConfigMaps(namespace)
clientset.CoreV1().Secrets(namespace)
```

---

## Working with ConfigMaps

### List ConfigMaps

```go
cms, err := clientset.CoreV1().ConfigMaps("default").List(ctx, metav1.ListOptions{})
for _, cm := range cms.Items {
    keys := make([]string, 0, len(cm.Data))
    for k := range cm.Data {
        keys = append(keys, k)
    }
    fmt.Printf("%-30s keys: %s\n", cm.Name, strings.Join(keys, ", "))
}
```

### Get a ConfigMap

```go
cm, err := clientset.CoreV1().ConfigMaps("default").Get(ctx, "app-config", metav1.GetOptions{})
if err != nil {
    if errors.IsNotFound(err) {
        fmt.Println("ConfigMap not found")
        return
    }
    log.Fatal(err)
}

fmt.Println("DB_HOST:", cm.Data["DB_HOST"])
fmt.Println("APP_ENV:", cm.Data["APP_ENV"])
```

### Create a ConfigMap

```go
import corev1 "k8s.io/api/core/v1"

cm := &corev1.ConfigMap{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "app-config",
        Namespace: "default",
        Labels:    map[string]string{"app": "my-app"},
    },
    Data: map[string]string{
        "APP_ENV":   "production",
        "DB_HOST":   "postgres.default.svc.cluster.local",
        "DB_PORT":   "5432",
        "CACHE_TTL": "3600",
        "nginx.conf": `server {
    listen 80;
    server_name _;
    location / { proxy_pass http://app:8080; }
}`,
    },
}

created, err := clientset.CoreV1().ConfigMaps("default").Create(ctx, cm, metav1.CreateOptions{})
fmt.Printf("ConfigMap created: %s\n", created.Name)
```

### Update a ConfigMap

```go
import "k8s.io/client-go/util/retry"

err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
    cm, err := clientset.CoreV1().ConfigMaps("default").Get(ctx, "app-config", metav1.GetOptions{})
    if err != nil {
        return err
    }
    cm.Data["CACHE_TTL"] = "7200"
    cm.Data["NEW_FLAG"]  = "true"
    _, err = clientset.CoreV1().ConfigMaps("default").Update(ctx, cm, metav1.UpdateOptions{})
    return err
})
```

### Apply (Create or Update)

```go
func applyConfigMap(ctx context.Context, clientset *kubernetes.Clientset, namespace string, data map[string]string) error {
    name := "app-config"
    cm := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
        Data:       data,
    }

    _, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
    if errors.IsNotFound(err) {
        _, err = clientset.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
        return err
    }

    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        existing, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }
        existing.Data = data
        _, err = clientset.CoreV1().ConfigMaps(namespace).Update(ctx, existing, metav1.UpdateOptions{})
        return err
    })
}
```

### Delete a ConfigMap

```go
err := clientset.CoreV1().ConfigMaps("default").Delete(ctx, "app-config", metav1.DeleteOptions{})
```

---

## Working with Secrets

> **Security Rule:** Never log secret values. Handle them only in memory when needed.

### List Secrets

```go
secrets, err := clientset.CoreV1().Secrets("default").List(ctx, metav1.ListOptions{})
for _, s := range secrets.Items {
    keys := make([]string, 0, len(s.Data))
    for k := range s.Data {
        keys = append(keys, k)
    }
    fmt.Printf("%-30s type=%-40s keys=%s\n", s.Name, s.Type, strings.Join(keys, ","))
}
```

### Get and Decode a Secret

```go
secret, err := clientset.CoreV1().Secrets("default").Get(ctx, "db-credentials", metav1.GetOptions{})
if err != nil {
    log.Fatal(err)
}

// Data is already []byte — no need to base64 decode (client-go does it)
username := string(secret.Data["username"])
password := string(secret.Data["password"])

dsn := fmt.Sprintf("postgres://%s:%s@postgres:5432/mydb", username, password)
```

### Create an Opaque Secret

```go
secret := &corev1.Secret{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "db-credentials",
        Namespace: "default",
        Labels:    map[string]string{"app": "my-app"},
    },
    Type: corev1.SecretTypeOpaque,
    // StringData is auto base64-encoded by the API
    StringData: map[string]string{
        "username": "admin",
        "password": "super-secret-password",
        "db_url":   "postgres://admin:pass@postgres:5432/mydb",
    },
}

clientset.CoreV1().Secrets("default").Create(ctx, secret, metav1.CreateOptions{})
```

### Create a Docker Registry Secret

```go
import "encoding/json"

dockerConfig := map[string]interface{}{
    "auths": map[string]interface{}{
        "registry.example.com": map[string]string{
            "username": "myuser",
            "password": "mypass",
            "email":    "ci@example.com",
            "auth":     base64.StdEncoding.EncodeToString([]byte("myuser:mypass")),
        },
    },
}

dockerConfigJSON, _ := json.Marshal(dockerConfig)

secret := &corev1.Secret{
    ObjectMeta: metav1.ObjectMeta{Name: "registry-credentials", Namespace: "default"},
    Type:       corev1.SecretTypeDockerConfigJson,
    Data: map[string][]byte{
        corev1.DockerConfigJsonKey: dockerConfigJSON,
    },
}

clientset.CoreV1().Secrets("default").Create(ctx, secret, metav1.CreateOptions{})
```

### Create a TLS Secret

```go
secret := &corev1.Secret{
    ObjectMeta: metav1.ObjectMeta{Name: "my-tls-secret", Namespace: "default"},
    Type:       corev1.SecretTypeTLS,
    Data: map[string][]byte{
        corev1.TLSCertKey:       certPEM,   // []byte of the cert
        corev1.TLSPrivateKeyKey: keyPEM,    // []byte of the key
    },
}

clientset.CoreV1().Secrets("default").Create(ctx, secret, metav1.CreateOptions{})
```

### Update a Secret (Rotate Password)

```go
func rotateSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name, key, newValue string) error {
    return retry.RetryOnConflict(retry.DefaultRetry, func() error {
        secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
        if err != nil {
            return err
        }

        // Keep previous value for rollback
        if val, ok := secret.Data[key]; ok {
            secret.Data[key+"_previous"] = val
        }
        secret.Data[key] = []byte(newValue)

        _, err = clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
        return err
    })
}
```

### Delete a Secret

```go
err := clientset.CoreV1().Secrets("default").Delete(ctx, "db-credentials", metav1.DeleteOptions{})
```
