package controller

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	crdv1 "rudro.dev/kubebuilder-crd/api/v1"
)

func (r *KubebuilderCrdReconciler) CheckService() error {
	srv := &corev1.Service{}
	if err := r.Client.Get(r.ctx, types.NamespacedName{
		Name:      r.kubebuildercrd.ServiceName(),
		Namespace: r.kubebuildercrd.Namespace,
	}, srv); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("KubebuilderCrd service not found")
			if err := r.Client.Create(r.ctx, r.newService()); err != nil {
				r.Log.Error(err, "failed to create service")
				return err
			}
			r.Log.Info("created service")
			return nil
		}
		return err
	}
	return nil
}

func (r *KubebuilderCrdReconciler) newService() *corev1.Service {
	fmt.Println("New Service is called")
	labels := map[string]string{
		"controller": r.kubebuildercrd.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.kubebuildercrd.ServiceName(),
			Namespace: r.kubebuildercrd.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(r.kubebuildercrd, crdv1.GroupVersion.WithKind(ourKind)),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     getServiceType(r.kubebuildercrd.Spec.Service.ServiceName),
			Ports: []corev1.ServicePort{
				{
					Port:       r.kubebuildercrd.Spec.Container.Port,
					NodePort:   r.kubebuildercrd.Spec.Service.ServiceNodePort,
					TargetPort: intstr.FromInt32(r.kubebuildercrd.Spec.Container.Port),
				},
			},
		},
	}
}
