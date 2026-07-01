# Chapter 3: Setup Local Lab

## Tools You Need

Install:

- Docker
- Git
- Go
- `kubectl`
- `kind` or `minikube`

For contribution work, `kind` is a good default because it runs Kubernetes clusters inside Docker containers.

---

## Check Docker

```bash
docker version
docker ps
```

If Docker works, you already have the base runtime needed for a local Kubernetes lab.

---

## Install kubectl

`kubectl` is the Kubernetes command-line tool.

Check:

```bash
kubectl version --client
```

---

## Install kind

`kind` means **Kubernetes in Docker**.

Check:

```bash
kind version
```

Create a cluster:

```bash
kind create cluster --name learning
```

Check nodes:

```bash
kubectl get nodes
```

Delete when done:

```bash
kind delete cluster --name learning
```

---

## First Workload

```bash
kubectl create deployment hello --image=nginx:1.27
kubectl get deployments
kubectl get pods
kubectl expose deployment hello --port=80 --type=ClusterIP
kubectl port-forward service/hello 8080:80
```

Open:

```text
http://localhost:8080
```

Stop `port-forward` with `Ctrl+C`.

---

## Cleanup

```bash
kubectl delete service hello
kubectl delete deployment hello
```

---

## What You Should Understand Before Moving On

- `kind` creates a real Kubernetes cluster locally.
- `kubectl` talks to the API server.
- A Deployment creates Pods.
- A Service gives stable access to Pods.
- You can inspect everything with `kubectl get` and `kubectl describe`.

