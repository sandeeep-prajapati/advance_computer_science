# k8s-client-go — kubectl vs Go Client Cheatsheet

Side-by-side reference.

---

## Pods

| kubectl | Go (client-go) |
|---|---|
| `kubectl get pods` | `cs.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})` |
| `kubectl get pods -A` | `cs.CoreV1().Pods("").List(ctx, metav1.ListOptions{})` |
| `kubectl get pods -n staging` | `cs.CoreV1().Pods("staging").List(ctx, metav1.ListOptions{})` |
| `kubectl get pods -l app=nginx` | `cs.CoreV1().Pods("ns").List(ctx, metav1.ListOptions{LabelSelector: "app=nginx"})` |
| `kubectl get pod my-pod` | `cs.CoreV1().Pods("ns").Get(ctx, "my-pod", metav1.GetOptions{})` |
| `kubectl logs my-pod` | `cs.CoreV1().Pods("ns").GetLogs("my-pod", &corev1.PodLogOptions{})` |
| `kubectl logs my-pod --tail=100` | `&corev1.PodLogOptions{TailLines: int64Ptr(100)}` |
| `kubectl logs my-pod -p` | `&corev1.PodLogOptions{Previous: true}` |
| `kubectl logs my-pod -c sidecar` | `&corev1.PodLogOptions{Container: "sidecar"}` |
| `kubectl logs -f my-pod` | `&corev1.PodLogOptions{Follow: true}` |
| `kubectl delete pod my-pod` | `cs.CoreV1().Pods("ns").Delete(ctx, "my-pod", metav1.DeleteOptions{})` |
| `kubectl delete pod my-pod --grace-period=0` | `metav1.DeleteOptions{GracePeriodSeconds: int64Ptr(0)}` |
| `kubectl get pods -w` | `cs.CoreV1().Pods("ns").Watch(ctx, metav1.ListOptions{})` |

---

## Deployments

| kubectl | Go (client-go) |
|---|---|
| `kubectl get deployments` | `cs.AppsV1().Deployments("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get deploy my-app` | `cs.AppsV1().Deployments("ns").Get(ctx, "my-app", metav1.GetOptions{})` |
| `kubectl scale deploy my-app --replicas=5` | `retry.RetryOnConflict(...) { d.Spec.Replicas = int32Ptr(5) }` |
| `kubectl set image deploy/my-app app=myapp:2.0` | `retry.RetryOnConflict(...) { d.Spec.Template.Spec.Containers[0].Image = "myapp:2.0" }` |
| `kubectl rollout restart deploy/my-app` | `d.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now()...` |
| `kubectl rollout status deploy/my-app` | `waitForDeployment(ctx, cs, "ns", "my-app")` |
| `kubectl delete deploy my-app` | `cs.AppsV1().Deployments("ns").Delete(ctx, "my-app", metav1.DeleteOptions{})` |
| `kubectl apply -f deploy.yaml` | `applyDeployment(ctx, cs, dep)` (create-or-update pattern) |

---

## Services

| kubectl | Go (client-go) |
|---|---|
| `kubectl get services` | `cs.CoreV1().Services("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get svc my-svc` | `cs.CoreV1().Services("ns").Get(ctx, "my-svc", metav1.GetOptions{})` |
| `kubectl delete svc my-svc` | `cs.CoreV1().Services("ns").Delete(ctx, "my-svc", metav1.DeleteOptions{})` |

---

## ConfigMaps & Secrets

| kubectl | Go (client-go) |
|---|---|
| `kubectl get configmaps` | `cs.CoreV1().ConfigMaps("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get cm app-config -o json` | `cs.CoreV1().ConfigMaps("ns").Get(ctx, "app-config", metav1.GetOptions{})` |
| `kubectl create cm app-config --from-literal=K=V` | `cs.CoreV1().ConfigMaps("ns").Create(ctx, &corev1.ConfigMap{Data: map[string]string{"K":"V"}}, ...)` |
| `kubectl delete cm app-config` | `cs.CoreV1().ConfigMaps("ns").Delete(ctx, "app-config", metav1.DeleteOptions{})` |
| `kubectl get secrets` | `cs.CoreV1().Secrets("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get secret db-creds -o jsonpath='{.data.password}'` | `string(secret.Data["password"])` (already decoded) |
| `kubectl create secret generic db --from-literal=p=x` | `cs.CoreV1().Secrets("ns").Create(ctx, &corev1.Secret{StringData: map[string]string{"p":"x"}}, ...)` |
| `kubectl delete secret db-creds` | `cs.CoreV1().Secrets("ns").Delete(ctx, "db-creds", metav1.DeleteOptions{})` |

---

## Namespaces

| kubectl | Go (client-go) |
|---|---|
| `kubectl get namespaces` | `cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})` |
| `kubectl create namespace team` | `cs.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "team"}}, ...)` |
| `kubectl delete namespace team` | `cs.CoreV1().Namespaces().Delete(ctx, "team", metav1.DeleteOptions{})` |

---

## Jobs & CronJobs

| kubectl | Go (client-go) |
|---|---|
| `kubectl get jobs` | `cs.BatchV1().Jobs("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl delete job my-job` | `cs.BatchV1().Jobs("ns").Delete(ctx, "my-job", metav1.DeleteOptions{PropagationPolicy: &foreground})` |
| `kubectl get cronjobs` | `cs.BatchV1().CronJobs("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl patch cj my-cj -p '{"spec":{"suspend":true}}'` | `cs.BatchV1().CronJobs("ns").Patch(ctx, "name", types.MergePatchType, suspendJSON, metav1.PatchOptions{})` |

---

## StatefulSets & DaemonSets

| kubectl | Go (client-go) |
|---|---|
| `kubectl get statefulsets` | `cs.AppsV1().StatefulSets("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl scale sts mysql --replicas=5` | `retry.RetryOnConflict(...) { sts.Spec.Replicas = int32Ptr(5) }` |
| `kubectl get daemonsets` | `cs.AppsV1().DaemonSets("ns").List(ctx, metav1.ListOptions{})` |

---

## Ingress

| kubectl | Go (client-go) |
|---|---|
| `kubectl get ingress` | `cs.NetworkingV1().Ingresses("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl delete ingress my-ing` | `cs.NetworkingV1().Ingresses("ns").Delete(ctx, "my-ing", metav1.DeleteOptions{})` |

---

## RBAC

| kubectl | Go (client-go) |
|---|---|
| `kubectl get roles` | `cs.RbacV1().Roles("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get rolebindings` | `cs.RbacV1().RoleBindings("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get clusterroles` | `cs.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})` |
| `kubectl get clusterrolebindings` | `cs.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})` |
| `kubectl get serviceaccounts` | `cs.CoreV1().ServiceAccounts("ns").List(ctx, metav1.ListOptions{})` |

---

## Storage

| kubectl | Go (client-go) |
|---|---|
| `kubectl get pvc` | `cs.CoreV1().PersistentVolumeClaims("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get pv` | `cs.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})` |
| `kubectl get storageclass` | `cs.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})` |

---

## Nodes

| kubectl | Go (client-go) |
|---|---|
| `kubectl get nodes` | `cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{})` |
| `kubectl get node node-1` | `cs.CoreV1().Nodes().Get(ctx, "node-1", metav1.GetOptions{})` |

---

## Events

| kubectl | Go (client-go) |
|---|---|
| `kubectl get events` | `cs.CoreV1().Events("ns").List(ctx, metav1.ListOptions{})` |
| `kubectl get events --field-selector type=Warning` | `cs.CoreV1().Events("ns").List(ctx, metav1.ListOptions{FieldSelector: "type=Warning"})` |
| `kubectl get events -w` | Informer on Events resource |

---

## Watch

| kubectl | Go (client-go) |
|---|---|
| `kubectl get pods -w` | `cs.CoreV1().Pods("ns").Watch(ctx, metav1.ListOptions{})` |
| `kubectl get deployments -w` | `cs.AppsV1().Deployments("ns").Watch(ctx, metav1.ListOptions{})` |
| Efficient + cached watching | `informers.NewSharedInformerFactory(cs, 30*time.Second)` |

---

## Error Code Reference

| HTTP Code | kubectl output | Go check |
|---|---|---|
| 200 | Success | `err == nil` |
| 401 | Unauthorized | `errors.IsUnauthorized(err)` |
| 403 | Forbidden | `errors.IsForbidden(err)` |
| 404 | NotFound | `errors.IsNotFound(err)` |
| 409 | Conflict / AlreadyExists | `errors.IsConflict(err)` or `errors.IsAlreadyExists(err)` |
| 422 | Invalid | `errors.IsInvalid(err)` |
| 429 | TooManyRequests | `errors.IsTooManyRequests(err)` |
| 503 | ServiceUnavailable | `errors.IsServiceUnavailable(err)` |
