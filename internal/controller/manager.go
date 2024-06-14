/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	crdv1 "rudro.dev/kubebuilder-crd/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	deployOwnerKey = ".metadata.controller"
	svcOwnerKey    = ".metadata.controller"
	ourApiGVStr    = crdv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *KubebuilderCrdReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// Extra part
	// Indexing our Owns resource.This will allow for quickly answer the question:
	// If Owns Resource x is updated, which KubebuilderCrd are affected?
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, deployOwnerKey, func(object client.Object) []string {
		// grab the deployment object, extract the owner.
		deployment := object.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deployment)
		if owner == nil {
			return nil
		}
		// make sure it's a KubebuilderCrd
		if owner.APIVersion != ourApiGVStr || owner.Kind != ourKind {
			return nil
		}
		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Service{}, svcOwnerKey, func(object client.Object) []string {
		svc := object.(*corev1.Service)
		owner := metav1.GetControllerOf(svc)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != ourApiGVStr || owner.Kind != ourKind {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	// Extra part
	// Implementation with watches and custom eventHandler
	// if someone edit the resources(here example given for deployment resource) by kubectl
	handlerFunc := handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		//fmt.Println("111111111111111111111111")
		// List all the CR
		customRs := &crdv1.KubebuilderCrdList{}
		if err := r.List(context.Background(), customRs); err != nil {
			return nil
		}
		// This func return a reconcile request array
		var req []reconcile.Request
		for _, c := range customRs.Items {
			//fmt.Println("crrrrrrrrrrrrrr", c.Name)
			deploymentName := c.DeploymentName()
			// Find the deployment owned by the CR
			if deploymentName == obj.GetName() && c.Namespace == obj.GetNamespace() {
				deploy := &appsv1.Deployment{}
				if err := r.Get(context.Background(), types.NamespacedName{
					Namespace: obj.GetNamespace(),
					Name:      obj.GetName(),
				}, deploy); err != nil {
					//fmt.Println("2222222222222222222222", obj.GetName(), c.Name)
					// This case can happen if somehow deployment gets deleted by
					// Kubectl command. We need to append new reconcile request to array
					// to create desired number of deployment again.
					if errors.IsNotFound(err) {
						req = append(req, reconcile.Request{
							NamespacedName: types.NamespacedName{
								Namespace: c.Namespace,
								Name:      c.Name,
							},
						})
						continue
					} else {
						//fmt.Println("errrrrrrrrrrrr", err)
						return nil
					}
				}
				// Only append to the reconcile request array if replica count miss match.
				if deploy.Spec.Replicas != nil && c.Spec.Replicas != nil {
					//fmt.Println(*deploy.Spec.Replicas, *c.Spec.Replicas, deploy.Name, c.Name)
					if *deploy.Spec.Replicas != *c.Spec.Replicas {
						//fmt.Println("checking", *deploy.Spec.Replicas, *c.Spec.Replicas)
						req = append(req, reconcile.Request{
							NamespacedName: types.NamespacedName{
								Namespace: c.Namespace,
								Name:      c.Name,
							},
						})
					}
				}
			}
		}
		return req
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.KubebuilderCrd{}).
		Watches(
			&appsv1.Deployment{}, handlerFunc,
		).
		Owns(&corev1.Service{}).
		Complete(r)

	//
	//Main part. Simplified. Comment above for simplification.
	//return ctrl.NewControllerManagedBy(mgr).
	//	For(&crdv1.KubebuilderCrd{}).
	//	Owns(&appsv1.Deployment{}).
	//	Owns(&corev1.Service{}).
	//	Complete(r)
}
