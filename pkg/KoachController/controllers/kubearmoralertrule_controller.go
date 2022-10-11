package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	securityv1 "github.com/kubearmor/koach/pkg/KoachController/api/v1"
)

// KubeArmorAlertRuleReconciler reconciles a KubeArmorAlertRule object
type KubeArmorAlertRuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=security.kubearmor.com,resources=kubearmoralertrules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=security.kubearmor.com,resources=kubearmoralertrules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=security.kubearmor.com,resources=kubearmoralertrules/finalizers,verbs=update

func (r *KubeArmorAlertRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeArmorAlertRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1.KubeArmorAlertRule{}).
		Complete(r)
}
