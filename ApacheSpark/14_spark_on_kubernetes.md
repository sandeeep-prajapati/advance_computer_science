# Chapter 14: Spark on Kubernetes

## Why Spark on Kubernetes?

Running Spark on Kubernetes gives you:
- **Unified infrastructure** — same cluster as the rest of your apps
- **Dynamic resource allocation** — no dedicated Spark cluster
- **Container isolation** — reproducible environments
- **Cloud-native** — works on EKS, GKE, AKS, on-premise

---

## How It Works

```
spark-submit → Driver Pod created → Requests Executor Pods
                      ↓
            K8s API Server allocates pods
                      ↓
            Executor Pods start, run tasks, terminate
                      ↓
            Driver Pod completes, terminates
```

No persistent cluster — pods are ephemeral. Each job gets its own pods.

---

## Basic spark-submit on Kubernetes

```bash
spark-submit \
  --master k8s://https://k8s-api-server:6443 \
  --deploy-mode cluster \
  --name my-spark-job \
  --conf spark.kubernetes.container.image=apache/spark:3.5.0 \
  --conf spark.kubernetes.namespace=spark-jobs \
  --conf spark.executor.instances=5 \
  --conf spark.executor.memory=4g \
  --conf spark.executor.cores=2 \
  --conf spark.driver.memory=2g \
  local:///opt/spark/examples/jars/spark-examples.jar \
  10000000
```

For a PySpark job:

```bash
spark-submit \
  --master k8s://https://k8s-api-server:6443 \
  --deploy-mode cluster \
  --conf spark.kubernetes.container.image=my-registry/pyspark:3.5 \
  --conf spark.kubernetes.namespace=spark-jobs \
  --conf spark.kubernetes.authenticate.driver.serviceAccountName=spark \
  --conf spark.executor.instances=10 \
  local:///opt/spark/jobs/my_job.py
```

---

## Custom Docker Image

```dockerfile
FROM apache/spark:3.5.0

# Install Python dependencies
USER root
COPY requirements.txt /opt/spark/requirements.txt
RUN pip3 install --no-cache-dir -r /opt/spark/requirements.txt

# Copy job files
COPY src/ /opt/spark/jobs/

USER spark
```

```bash
docker build -t my-registry/pyspark:3.5 .
docker push my-registry/pyspark:3.5
```

---

## RBAC Setup

Spark needs permission to create/delete pods:

```yaml
# spark-rbac.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: spark
  namespace: spark-jobs
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: spark-role
  namespace: spark-jobs
rules:
  - apiGroups: [""]
    resources: ["pods", "services", "configmaps", "persistentvolumeclaims"]
    verbs: ["create", "get", "watch", "list", "delete", "patch"]
  - apiGroups: [""]
    resources: ["pods/log"]
    verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: spark-rolebinding
  namespace: spark-jobs
subjects:
  - kind: ServiceAccount
    name: spark
    namespace: spark-jobs
roleRef:
  kind: Role
  name: spark-role
  apiGroup: rbac.authorization.k8s.io
```

```bash
kubectl apply -f spark-rbac.yaml
```

---

## Dynamic Resource Allocation on K8s

Scale executors automatically without pre-allocating:

```bash
spark-submit \
  --conf spark.dynamicAllocation.enabled=true \
  --conf spark.dynamicAllocation.shuffleTracking.enabled=true \
  --conf spark.dynamicAllocation.minExecutors=1 \
  --conf spark.dynamicAllocation.maxExecutors=50 \
  --conf spark.dynamicAllocation.initialExecutors=5 \
  --conf spark.dynamicAllocation.executorIdleTimeout=60s \
  ...
```

> Note: On Kubernetes, use `shuffleTracking` instead of External Shuffle Service.

---

## Spark Operator (Kubernetes-Native)

The **Spark Operator** (Google's open-source controller) manages Spark jobs as Kubernetes CRDs:

```bash
helm install spark-operator spark-operator/spark-operator \
  --namespace spark-operator \
  --create-namespace
```

```yaml
# spark-job.yaml
apiVersion: sparkoperator.k8s.io/v1beta2
kind: SparkApplication
metadata:
  name: my-spark-job
  namespace: spark-jobs
spec:
  type: Python
  pythonVersion: "3"
  mode: cluster
  image: "my-registry/pyspark:3.5"
  imagePullPolicy: Always
  mainApplicationFile: "local:///opt/spark/jobs/etl.py"
  sparkVersion: "3.5.0"

  driver:
    cores: 2
    memory: "4g"
    serviceAccount: spark
    labels:
      app: my-spark-job

  executor:
    cores: 4
    instances: 10
    memory: "8g"
    labels:
      app: my-spark-job

  restartPolicy:
    type: OnFailure
    onFailureRetries: 3
```

```bash
kubectl apply -f spark-job.yaml
kubectl get sparkapplication -n spark-jobs
kubectl describe sparkapplication my-spark-job -n spark-jobs
```

---

## Accessing Data from S3

```bash
# Add S3 credentials as env vars in the image, or use IAM roles for service accounts
spark-submit \
  --conf spark.hadoop.fs.s3a.access.key=YOUR_KEY \
  --conf spark.hadoop.fs.s3a.secret.key=YOUR_SECRET \
  --conf spark.hadoop.fs.s3a.impl=org.apache.hadoop.fs.s3a.S3AFileSystem \
  ...
```

For EKS with IAM Roles for Service Accounts (IRSA):

```yaml
metadata:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789:role/spark-s3-role
```

---

## Spark History Server on K8s

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: spark-history-server
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: history-server
          image: apache/spark:3.5.0
          command: ["/opt/spark/bin/spark-class"]
          args: ["org.apache.spark.deploy.history.HistoryServer"]
          env:
            - name: SPARK_HISTORY_OPTS
              value: "-Dspark.history.fs.logDirectory=s3a://bucket/spark-logs/"
          ports:
            - containerPort: 18080
```

---

## Resource Management Tips

```bash
# Limit resource usage with resource quotas
--conf spark.kubernetes.driver.request.cores=1
--conf spark.kubernetes.driver.limit.cores=2
--conf spark.kubernetes.executor.request.cores=2
--conf spark.kubernetes.executor.limit.cores=4

# Node selectors (run on specific node pool)
--conf spark.kubernetes.node.selector.workload=spark

# Tolerations (run on tainted nodes)
--conf spark.kubernetes.driver.tolerations=[{"key":"spark","operator":"Equal","value":"true","effect":"NoSchedule"}]

# Pod priority
--conf spark.kubernetes.driver.priorityClassName=high-priority
--conf spark.kubernetes.executor.priorityClassName=medium-priority
```
