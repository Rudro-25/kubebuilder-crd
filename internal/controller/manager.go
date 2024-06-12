package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	crdv1 "rudro.dev/kubebuilder-crd/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KubebuilderCrdReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.KubebuilderCrd{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
