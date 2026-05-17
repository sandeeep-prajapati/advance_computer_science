<?php
/**
 * Smoke test — verifies the PHP k8s client library loaded correctly.
 * Run: php test_connection.php
 * Requires a valid ~/.kube/config to actually connect.
 */

require __DIR__ . '/vendor/autoload.php';

use Maclof\Kubernetes\Client;

echo "=== k8s-client-php smoke test ===\n\n";

// 1. Check library loaded
echo "✓ maclof/kubernetes-client loaded (v" . \Composer\InstalledVersions::getVersion('maclof/kubernetes-client') . ")\n";

// 2. Try to build a client (won't connect until a real cluster is configured)
$kubeconfigPath = $_SERVER['HOME'] . '/.kube/config';

if (!file_exists($kubeconfigPath)) {
    echo "⚠  No ~/.kube/config found.\n";
    echo "   To connect to a cluster, run one of:\n";
    echo "     minikube start                          (local)\n";
    echo "     gcloud container clusters get-credentials CLUSTER --region REGION\n";
    echo "     aws eks update-kubeconfig --name CLUSTER --region REGION\n";
    echo "     az aks get-credentials --resource-group RG --name CLUSTER\n\n";
    echo "   Then re-run: php test_connection.php\n";
    exit(0);
}

echo "✓ kubeconfig found: $kubeconfigPath\n";

try {
    $client = Client::fromKubeConfig();
    echo "✓ Client built from kubeconfig\n";

    // Attempt real connection
    $namespaces = $client->namespaces()->find();
    echo "✓ Connected! Namespaces found: " . count($namespaces) . "\n\n";

    foreach ($namespaces as $ns) {
        $name   = $ns->getMetadata('name');
        $status = $ns->getStatus()['phase'] ?? 'Unknown';
        echo "  [$status] $name\n";
    }

    echo "\n✓ All good — you can now use the PHP k8s client.\n";

} catch (\Exception $e) {
    echo "✗ Connection failed: " . $e->getMessage() . "\n";
    echo "  Make sure your cluster is running and kubeconfig is valid.\n";
}
