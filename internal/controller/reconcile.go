package controller

import (
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ourkind = "KubebuilderCrd"
)

type KubebuilderCrdReconciler struct {
	Client client.Client
	Log    logr.Logger
	ctx    context.Context
	Scheme *runtime.Scheme
	kubebuildercrd *
}

// +kubebuilder:rbac:groups=rudro.dev.rudro.dev,resources=kubebuildercrds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rudro.dev.rudro.dev,resources=kubebuildercrds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=rudro.dev.rudro.dev,resources=kubebuildercrds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KubebuilderCrd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile

func (r *KubebuilderCrdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	return ctrl.Result{}, nil
}
