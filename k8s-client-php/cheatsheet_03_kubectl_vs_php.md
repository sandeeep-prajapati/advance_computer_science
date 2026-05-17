# k8s-client-php — kubectl vs PHP Client Cheatsheet

Side-by-side reference: what you'd type in kubectl vs the PHP equivalent.

---

## Pods

| kubectl | PHP |
|---|---|
| `kubectl get pods` | `$client->pods()->find()` |
| `kubectl get pods -n staging` | `$client->setNamespace('staging'); $client->pods()->find()` |
| `kubectl get pods -A` | `$client->setNamespace(''); $client->pods()->find()` |
| `kubectl get pod my-pod` | `$client->pods()->findByName('my-pod')` |
| `kubectl get pods -l app=nginx` | `$client->pods()->setLabelSelector(['app'=>'nginx'])->find()` |
| `kubectl describe pod my-pod` | `$pod = $client->pods()->findByName('my-pod'); $pod->getStatus(); $pod->getSpec()` |
| `kubectl logs my-pod` | `$client->pods()->logs('my-pod')` |
| `kubectl logs my-pod --tail=100` | `$client->pods()->logs('my-pod', ['tailLines'=>100])` |
| `kubectl logs my-pod -p` | `$client->pods()->logs('my-pod', ['previous'=>true])` |
| `kubectl logs my-pod -c sidecar` | `$client->pods()->logs('my-pod', ['container'=>'sidecar'])` |
| `kubectl delete pod my-pod` | `$client->pods()->deleteByName('my-pod')` |
| `kubectl delete pod my-pod --grace-period=0` | `$client->pods()->delete($pod, ['gracePeriodSeconds'=>0])` |

---

## Deployments

| kubectl | PHP |
|---|---|
| `kubectl get deployments` | `$client->deployments()->find()` |
| `kubectl get deploy my-app` | `$client->deployments()->findByName('my-app')` |
| `kubectl scale deploy my-app --replicas=5` | `$client->deployments()->patch($dep, ['spec'=>['replicas'=>5]])` |
| `kubectl set image deploy/my-app app=myapp:2.0` | `$client->deployments()->patch($dep, ['spec'=>['template'=>['spec'=>['containers'=>[['name'=>'app','image'=>'myapp:2.0']]]]]])` |
| `kubectl rollout status deploy/my-app` | `waitForDeployment($client, 'my-app')` |
| `kubectl delete deploy my-app` | `$client->deployments()->deleteByName('my-app')` |
| `kubectl apply -f deploy.yaml` | `$client->deployments()->create($model)` or `update($model)` |

---

## Services

| kubectl | PHP |
|---|---|
| `kubectl get services` | `$client->services()->find()` |
| `kubectl get svc my-svc` | `$client->services()->findByName('my-svc')` |
| `kubectl describe svc my-svc` | `$svc->getSpec(); $svc->getStatus()` |
| `kubectl delete svc my-svc` | `$client->services()->deleteByName('my-svc')` |

---

## ConfigMaps & Secrets

| kubectl | PHP |
|---|---|
| `kubectl get configmaps` | `$client->configMaps()->find()` |
| `kubectl get cm app-config -o json` | `$client->configMaps()->findByName('app-config')` |
| `kubectl create configmap app-config --from-literal=KEY=val` | `$client->configMaps()->create(new ConfigMap([...'data'=>['KEY'=>'val']...]))` |
| `kubectl delete cm app-config` | `$client->configMaps()->deleteByName('app-config')` |
| `kubectl get secrets` | `$client->secrets()->find()` |
| `kubectl get secret db-creds -o jsonpath='{.data.password}'` | `base64_decode($secret->getAttribute('data')['password'])` |
| `kubectl create secret generic db-creds --from-literal=password=x` | `$client->secrets()->create(new Secret([...'data'=>['password'=>base64_encode('x')]...]))` |
| `kubectl delete secret db-creds` | `$client->secrets()->deleteByName('db-creds')` |

---

## Namespaces

| kubectl | PHP |
|---|---|
| `kubectl get namespaces` | `$client->namespaces()->find()` |
| `kubectl create namespace team-alpha` | `$client->namespaces()->create(new NamespaceModel(['metadata'=>['name'=>'team-alpha']]))` |
| `kubectl delete namespace team-alpha` | `$client->namespaces()->deleteByName('team-alpha')` |

---

## Storage

| kubectl | PHP |
|---|---|
| `kubectl get pvc` | `$client->persistentVolumeClaims()->find()` |
| `kubectl get pv` | `$client->persistentVolumes()->find()` |
| `kubectl get storageclass` | `$client->storageClasses()->find()` |
| `kubectl delete pvc my-pvc` | `$client->persistentVolumeClaims()->deleteByName('my-pvc')` |

---

## Jobs & CronJobs

| kubectl | PHP |
|---|---|
| `kubectl get jobs` | `$client->jobs()->find()` |
| `kubectl get job my-job` | `$client->jobs()->findByName('my-job')` |
| `kubectl delete job my-job` | `$client->jobs()->deleteByName('my-job', ['propagationPolicy'=>'Foreground'])` |
| `kubectl get cronjobs` | `$client->cronJobs()->find()` |
| `kubectl patch cronjob my-cj -p '{"spec":{"suspend":true}}'` | `$client->cronJobs()->patch($cj, ['spec'=>['suspend'=>true]])` |

---

## StatefulSets & DaemonSets

| kubectl | PHP |
|---|---|
| `kubectl get statefulsets` | `$client->statefulSets()->find()` |
| `kubectl scale sts mysql --replicas=5` | `$client->statefulSets()->patch($sts, ['spec'=>['replicas'=>5]])` |
| `kubectl get daemonsets` | `$client->daemonSets()->find()` |

---

## Ingress

| kubectl | PHP |
|---|---|
| `kubectl get ingress` | `$client->ingresses()->find()` |
| `kubectl get ing my-ingress` | `$client->ingresses()->findByName('my-ingress')` |
| `kubectl delete ingress my-ingress` | `$client->ingresses()->deleteByName('my-ingress')` |

---

## RBAC

| kubectl | PHP |
|---|---|
| `kubectl get roles` | `$client->roles()->find()` |
| `kubectl get rolebindings` | `$client->roleBindings()->find()` |
| `kubectl get clusterroles` | `$client->clusterRoles()->find()` |
| `kubectl get clusterrolebindings` | `$client->clusterRoleBindings()->find()` |
| `kubectl get serviceaccounts` | `$client->serviceAccounts()->find()` |
| `kubectl create serviceaccount ci-sa` | `$client->serviceAccounts()->create(new ServiceAccount(['metadata'=>['name'=>'ci-sa']]))` |

---

## Nodes

| kubectl | PHP |
|---|---|
| `kubectl get nodes` | `$client->nodes()->find()` |
| `kubectl get node node-1` | `$client->nodes()->findByName('node-1')` |
| `kubectl describe node node-1` | `$node->getStatus()` |

---

## Events

| kubectl | PHP |
|---|---|
| `kubectl get events` | `$client->events()->find()` |
| `kubectl get events --field-selector type=Warning` | `$client->events()->setFieldSelector(['type'=>'Warning'])->find()` |
| `kubectl get events -w` | `$client->events()->watch(function($t, $e){ return true; })` |

---

## Watch

| kubectl | PHP |
|---|---|
| `kubectl get pods -w` | `$client->pods()->watch(function($type, $pod){ return true; })` |
| `kubectl get deployments -w` | `$client->deployments()->watch(function($type, $dep){ return true; })` |
| `kubectl get events -w` | `$client->events()->watch(function($type, $evt){ return true; })` |

---

## Label Selectors

| kubectl | PHP |
|---|---|
| `kubectl get pods -l app=nginx` | `$client->pods()->setLabelSelector(['app'=>'nginx'])->find()` |
| `kubectl get pods -l app=nginx,env=prod` | `$client->pods()->setLabelSelector(['app'=>'nginx','env'=>'prod'])->find()` |
| `kubectl get pods --field-selector status.phase=Running` | `$client->pods()->setFieldSelector(['status.phase'=>'Running'])->find()` |

---

## Cluster Info

| kubectl | PHP |
|---|---|
| `kubectl version` | `$client->getVersion()` |
| `kubectl cluster-info` | Uses the `master` URL configured in Client |
| `kubectl api-resources` | `$client->customResourceDefinitions()->find()` |

---

## HTTP Status Code Reference

| Code | kubectl equivalent | Meaning |
|---|---|---|
| `200` | Success | Resource returned/updated |
| `201` | Created | Resource created |
| `401` | error: Unauthorized | Bad/expired token |
| `403` | Error from server (Forbidden) | RBAC missing |
| `404` | Error from server (NotFound) | Resource doesn't exist |
| `409` | Error from server (Conflict) | Resource already exists |
| `422` | error: Invalid | Spec validation failed |
| `429` | — | Rate limited — add retry |
| `500` | — | K8s API bug |
