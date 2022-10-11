package common

import (
	"os"
	"path/filepath"
)

// ================ //
// == Kubernetes == //
// ================ //

// IsK8sLocal Function
func IsK8sLocal() bool {
	k8sConfig := os.Getenv("KUBECONFIG")
	if k8sConfig != "" {
		if _, err := os.Stat(filepath.Clean(k8sConfig)); err == nil {
			return true
		}
	}

	home := os.Getenv("HOME")
	if _, err := os.Stat(filepath.Clean(home + "/.kube/config")); err == nil {
		return true
	}

	return false
}

// IsInK8sCluster Function
func IsInK8sCluster() bool {
	if _, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		return true
	}

	if _, err := os.Stat(filepath.Clean("/run/secrets/kubernetes.io")); err == nil {
		return true
	}

	return false
}

// IsK8sEnv Function
func IsK8sEnv() bool {
	// local
	if IsK8sLocal() {
		return true
	}

	// in-cluster
	if IsInK8sCluster() {
		return true
	}

	return false
}
