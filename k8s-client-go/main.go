// Smoke test — verifies the Go k8s client library works.
// Run: go run main.go
// Requires a valid ~/.kube/config to actually connect.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildConfig() (*rest.Config, error) {
	// Try in-cluster first (when running as a Pod)
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}
	// Fall back to kubeconfig
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func main() {
	fmt.Println("=== k8s-client-go smoke test ===\n")

	// Check kubeconfig
	home, _ := os.UserHomeDir()
	kubeconfigPath := filepath.Join(home, ".kube", "config")

	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		fmt.Println("⚠  No ~/.kube/config found.")
		fmt.Println("   To connect to a cluster, run one of:")
		fmt.Println("     minikube start")
		fmt.Println("     gcloud container clusters get-credentials CLUSTER --region REGION")
		fmt.Println("     aws eks update-kubeconfig --name CLUSTER --region REGION")
		fmt.Println("     az aks get-credentials --resource-group RG --name CLUSTER")
		fmt.Println("\n   Then re-run: go run main.go")
		os.Exit(0)
	}
	fmt.Printf("✓ kubeconfig found: %s\n", kubeconfigPath)

	// Build REST config
	config, err := buildConfig()
	if err != nil {
		log.Fatalf("✗ Failed to build config: %v", err)
	}
	fmt.Printf("✓ Config built — API server: %s\n", config.Host)

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("✗ Failed to create clientset: %v", err)
	}
	fmt.Println("✓ Clientset created")

	// Verify server version
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		fmt.Printf("✗ Cannot reach API server: %v\n", err)
		fmt.Println("  Make sure your cluster is running.")
		os.Exit(1)
	}
	fmt.Printf("✓ Connected! Kubernetes version: %s\n\n", version.GitVersion)

	// List namespaces
	ctx := context.TODO()
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("✗ Failed to list namespaces: %v", err)
	}

	fmt.Printf("Namespaces (%d):\n", len(namespaces.Items))
	for _, ns := range namespaces.Items {
		fmt.Printf("  [%s] %s\n", ns.Status.Phase, ns.Name)
	}

	// List pods in default namespace
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("✗ Failed to list pods: %v", err)
	}

	fmt.Printf("\nPods in 'default' namespace (%d):\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf("  [%-10s] %s  node=%s\n", pod.Status.Phase, pod.Name, pod.Spec.NodeName)
	}

	fmt.Println("\n✓ All good — you can now use the Go k8s client.")
}
