package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	crdv1 "rudro.dev/kubebuilder-crd/api/v1"
)

func (r *KubebuilderCrdReconciler) CheckDeployment() error {

	deploy := &appsv1.Deployment{}

	if err := r.Client.Get(r.ctx, types.NamespacedName{
		Name:      r.kubebuildercrd.DeploymentName(),
		Namespace: r.kubebuildercrd.Namespace,
	}, deploy); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Creating a new Deployment", "Namespace", r.kubebuildercrd.Namespace, "Name", r.kubebuildercrd.Name)
			deploy := r.NewDeployment()
			if err := r.Client.Create(r.ctx, deploy); err != nil {
				return err
			}
			r.Log.Info("Created Deployment", "Namespace", deploy.Namespace, "Name", deploy.Name)
			return nil
		}
		return err
	}
	if r.kubebuildercrd.Spec.Replicas != nil && *deploy.Spec.Replicas != *r.kubebuildercrd.Spec.Replicas {
		r.Log.Info("replica mismatch...")
		*deploy.Spec.Replicas = *r.kubebuildercrd.Spec.Replicas
		if err := r.Client.Update(r.ctx, deploy); err != nil {
			r.Log.Error(err, "Failed to update Deployment", "Namespace", deploy.Namespace, "Name", deploy.Name)
			return err
		}
	}

	return nil
}

func (r *KubebuilderCrdReconciler) NewDeployment() *appsv1.Deployment {
	r.Log.Info("New Deployment is called")
	labels := map[string]string{
		"controller": r.kubebuildercrd.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.kubebuildercrd.DeploymentName(),
			Namespace: r.kubebuildercrd.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(r.kubebuildercrd, crdv1.GroupVersion.WithKind(ourKind)),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: r.kubebuildercrd.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kubebuildercrd-container",
							Image: r.kubebuildercrd.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: r.kubebuildercrd.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}
