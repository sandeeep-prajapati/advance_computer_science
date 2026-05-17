# Chapter 9: Namespaces & RBAC

## RBAC Interfaces

```go
clientset.CoreV1().Namespaces()
clientset.CoreV1().ServiceAccounts(namespace)
clientset.RbacV1().Roles(namespace)
clientset.RbacV1().ClusterRoles()
clientset.RbacV1().RoleBindings(namespace)
clientset.RbacV1().ClusterRoleBindings()
```

---

## Working with Namespaces

### List Namespaces

```go
nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
for _, ns := range nsList.Items {
    fmt.Printf("[%s] %s\n", ns.Status.Phase, ns.Name)
}
```

### Create a Namespace

```go
ns := &corev1.Namespace{
    ObjectMeta: metav1.ObjectMeta{
        Name: "team-alpha",
        Labels: map[string]string{
            "team":        "alpha",
            "environment": "development",
        },
        Annotations: map[string]string{
            "owner": "alice@example.com",
        },
    },
}

created, err := clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
fmt.Printf("Namespace created: %s\n", created.Name)
```

### Delete a Namespace

```go
// Deletes ALL resources inside it
err := clientset.CoreV1().Namespaces().Delete(ctx, "team-alpha", metav1.DeleteOptions{})
```

### Wait for Namespace Termination

```go
func waitForNamespaceGone(ctx context.Context, clientset *kubernetes.Clientset, name string) error {
    for {
        _, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
        if errors.IsNotFound(err) {
            fmt.Printf("Namespace '%s' fully deleted\n", name)
            return nil
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(3 * time.Second):
            fmt.Printf("Namespace '%s' still terminating...\n", name)
        }
    }
}
```

---

## Working with ServiceAccounts

### Create a ServiceAccount

```go
sa := &corev1.ServiceAccount{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "ci-pipeline",
        Namespace: "default",
        Labels:    map[string]string{"purpose": "ci"},
        Annotations: map[string]string{
            "description": "ServiceAccount for GitHub Actions CI pipeline",
        },
    },
}

clientset.CoreV1().ServiceAccounts("default").Create(ctx, sa, metav1.CreateOptions{})
```

### List ServiceAccounts

```go
sas, _ := clientset.CoreV1().ServiceAccounts("default").List(ctx, metav1.ListOptions{})
for _, sa := range sas.Items {
    fmt.Println(sa.Name)
}
```

---

## Working with Roles

### Create a Role

```go
import rbacv1 "k8s.io/api/rbac/v1"

role := &rbacv1.Role{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "pod-reader",
        Namespace: "default",
    },
    Rules: []rbacv1.PolicyRule{
        {
            APIGroups: []string{""},
            Resources: []string{"pods", "pods/log", "pods/status"},
            Verbs:     []string{"get", "list", "watch"},
        },
    },
}

clientset.RbacV1().Roles("default").Create(ctx, role, metav1.CreateOptions{})
```

### Deployment Manager Role

```go
role := &rbacv1.Role{
    ObjectMeta: metav1.ObjectMeta{Name: "deployment-manager", Namespace: "production"},
    Rules: []rbacv1.PolicyRule{
        {
            APIGroups: []string{"apps"},
            Resources: []string{"deployments", "replicasets"},
            Verbs:     []string{"get", "list", "watch", "create", "update", "patch"},
        },
        {
            APIGroups: []string{""},
            Resources: []string{"pods", "services", "configmaps"},
            Verbs:     []string{"get", "list", "watch"},
        },
    },
}

clientset.RbacV1().Roles("production").Create(ctx, role, metav1.CreateOptions{})
```

---

## Working with ClusterRoles

```go
clusterRole := &rbacv1.ClusterRole{
    ObjectMeta: metav1.ObjectMeta{Name: "cluster-reader"},
    Rules: []rbacv1.PolicyRule{
        {
            APIGroups: []string{"", "apps", "batch"},
            Resources: []string{"*"},
            Verbs:     []string{"get", "list", "watch"},
        },
        {
            APIGroups: []string{""},
            Resources: []string{"nodes", "namespaces", "persistentvolumes"},
            Verbs:     []string{"get", "list"},
        },
    },
}

clientset.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
```

---

## Working with RoleBindings

```go
roleBinding := &rbacv1.RoleBinding{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "ci-pipeline-binding",
        Namespace: "default",
    },
    Subjects: []rbacv1.Subject{
        {
            Kind:      "ServiceAccount",
            Name:      "ci-pipeline",
            Namespace: "default",
        },
    },
    RoleRef: rbacv1.RoleRef{
        APIGroup: "rbac.authorization.k8s.io",
        Kind:     "Role",
        Name:     "deployment-manager",
    },
}

clientset.RbacV1().RoleBindings("default").Create(ctx, roleBinding, metav1.CreateOptions{})
```

### Bind to a User

```go
roleBinding := &rbacv1.RoleBinding{
    ObjectMeta: metav1.ObjectMeta{Name: "alice-binding", Namespace: "team-alpha"},
    Subjects: []rbacv1.Subject{
        {
            Kind:     "User",
            Name:     "alice@example.com",
            APIGroup: "rbac.authorization.k8s.io",
        },
    },
    RoleRef: rbacv1.RoleRef{
        APIGroup: "rbac.authorization.k8s.io",
        Kind:     "Role",
        Name:     "pod-reader",
    },
}
```

---

## Working with ClusterRoleBindings

```go
crb := &rbacv1.ClusterRoleBinding{
    ObjectMeta: metav1.ObjectMeta{Name: "dashboard-cluster-reader"},
    Subjects: []rbacv1.Subject{
        {
            Kind:      "ServiceAccount",
            Name:      "dashboard",
            Namespace: "kube-system",
        },
    },
    RoleRef: rbacv1.RoleRef{
        APIGroup: "rbac.authorization.k8s.io",
        Kind:     "ClusterRole",
        Name:     "cluster-reader",
    },
}

clientset.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{})
```

---

## Multi-Tenant Namespace Provisioning

```go
func provisionTenantNamespace(ctx context.Context, clientset *kubernetes.Clientset, team, ownerEmail string) error {
    ns := team + "-workspace"

    // 1. Create namespace
    _, err := clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name:        ns,
            Labels:      map[string]string{"team": team, "managed-by": "provisioner"},
            Annotations: map[string]string{"owner": ownerEmail},
        },
    }, metav1.CreateOptions{})
    if err != nil && !errors.IsAlreadyExists(err) {
        return err
    }

    // 2. ServiceAccount
    clientset.CoreV1().ServiceAccounts(ns).Create(ctx, &corev1.ServiceAccount{
        ObjectMeta: metav1.ObjectMeta{Name: team + "-sa", Namespace: ns},
    }, metav1.CreateOptions{})

    // 3. Role
    clientset.RbacV1().Roles(ns).Create(ctx, &rbacv1.Role{
        ObjectMeta: metav1.ObjectMeta{Name: team + "-role", Namespace: ns},
        Rules: []rbacv1.PolicyRule{{
            APIGroups: []string{"", "apps", "batch"},
            Resources: []string{"*"},
            Verbs:     []string{"*"},
        }},
    }, metav1.CreateOptions{})

    // 4. RoleBinding
    clientset.RbacV1().RoleBindings(ns).Create(ctx, &rbacv1.RoleBinding{
        ObjectMeta: metav1.ObjectMeta{Name: team + "-binding", Namespace: ns},
        Subjects:   []rbacv1.Subject{{Kind: "ServiceAccount", Name: team + "-sa", Namespace: ns}},
        RoleRef:    rbacv1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "Role", Name: team + "-role"},
    }, metav1.CreateOptions{})

    fmt.Printf("Namespace '%s' provisioned for team '%s'\n", ns, team)
    return nil
}
```

---

## List Roles and Their Rules

```go
roles, _ := clientset.RbacV1().Roles("default").List(ctx, metav1.ListOptions{})
for _, role := range roles.Items {
    fmt.Printf("Role: %s\n", role.Name)
    for _, rule := range role.Rules {
        fmt.Printf("  Resources: %v | Verbs: %v\n", rule.Resources, rule.Verbs)
    }
}
```
