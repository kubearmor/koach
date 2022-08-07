package server

import (
	"context"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kl "github.com/kubearmor/koach/koach/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"
)

// ================= //
// == K8s Handler == //
// ================= //

// K8s Handler
var K8s *K8sHandler

// init Function
func init() {
	K8s = NewK8sHandler()
}

// K8sHandler Structure
type K8sHandler struct {
	K8sToken string
	K8sHost  string
	K8sPort  string

	client dynamic.Interface
}

// NewK8sHandler Function
func NewK8sHandler() *K8sHandler {
	kh := &K8sHandler{}

	if val, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		kh.K8sHost = val
	} else {
		kh.K8sHost = "127.0.0.1"
	}

	if val, ok := os.LookupEnv("KUBERNETES_PORT_443_TCP_PORT"); ok {
		kh.K8sPort = val
	} else {
		kh.K8sPort = "8001" // kube-proxy
	}

	return kh
}

// ================ //
// == K8s Client == //
// ================ //

// InitK8sClient Function
func (kh *K8sHandler) InitK8sClient() bool {
	if !kl.IsK8sEnv() { // not Kubernetes
		return false
	}

	if kh.client == nil {
		if kl.IsInK8sCluster() {
			return kh.InitInclusterAPIClient()
		}
		if kl.IsK8sLocal() {
			return kh.InitLocalAPIClient()
		}

		return false
	}

	return true
}

// InitLocalAPIClient Function
func (kh *K8sHandler) InitLocalAPIClient() bool {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
		if _, err := os.Stat(filepath.Clean(kubeconfig)); err != nil {
			return false
		}
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return false
	}

	// creates the clientset
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return false
	}
	kh.client = client

	return true
}

// InitInclusterAPIClient Function
func (kh *K8sHandler) InitInclusterAPIClient() bool {
	read, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return false
	}
	kh.K8sToken = string(read)

	// create the configuration by token
	kubeConfig := &rest.Config{
		Host:        "https://" + kh.K8sHost + ":" + kh.K8sPort,
		BearerToken: kh.K8sToken,
		// #nosec
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return false
	}
	kh.client = client

	return true
}

// ========================== //
// == KubeArmor Alert Rule == //
// ========================== //

// WatchK8sAlertRules Function
func (kh K8sHandler) WatchKubeArmorAlertRules() watch.Interface {
	if !kl.IsK8sEnv() { // not Kubernetes
		return nil
	}

	watcher, err := kh.client.Resource(schema.GroupVersionResource{
		Group:    "security.kubearmor.com",
		Version:  "v1",
		Resource: "kubearmoralertrules",
	}).Namespace(metav1.NamespaceAll).Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil
	}

	return watcher
}
