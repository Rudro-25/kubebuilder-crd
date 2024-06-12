/*
Copyright 2024.

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
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	rudrodevv1 "rudro.dev/kubebuilder-crd/api/v1"
)

// KubebuilderCrdReconciler reconciles a KubebuilderCrd object
type KubebuilderCrdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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
	log := log.FromContext(ctx)
	log.WithValues("ReqName", req.Name, "ReqNamespace", req.Namespace)
	fmt.Println("Reconcile started")

	// TODO(user): your logic here

	// customR have all data of CustomR Resources
	var customR rudrov1.customR
	if err := r.Get(ctx, req.NamespaceName, &customR); err != nil {
		fmt.Println("unable to fetch customR")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, nil
	}
	fmt.Println("Custom fetched", req.NamespacedName)

	// deploymentInstance have all data of deployment in specific namespace and name under the controller
	var deploymentInstance appsv1.Deployment
	// making deployment name
	deploymentName := func() string {
		if customR, Spec.DeploymentName == "" {
			return customR.Name + "-" + "randomName"
		} else {
			return customR.Name + "-" + customR, Spec.DeploymentName
		}
	}()
	// Creating NamespaceName for deploymentInstance
	obk := client.ObjectKey{
		Name:      deploymentName,
		Namespace: req.Namespace,
	}

	if err := r.Get(ctx, obk, &deploymentInstance); err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("could not find existing Deployment for ", customR.Name, ", creating one...")
			err := r.Client.Create(ctx, newDeployment(&customR, deploymentName))
			if err != nil {
				fmt.Errorf("error")
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubebuilderCrdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rudrodevv1.KubebuilderCrd{}).
		Complete(r)
}
