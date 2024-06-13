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
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	crdv1 "rudro.dev/kubebuilder-crd/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	_ "strings"
)

const (
	ourKind = "KubebuilderCrd"
)

// KubebuilderCrdReconciler reconciles a Kubebuilder Object
type KubebuilderCrdReconciler struct {
	client.Client
	Log            logr.Logger
	ctx            context.Context
	Scheme         *runtime.Scheme
	kubebuildercrd *crdv1.KubebuilderCrd
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
	r.Log = ctrl.Log.WithValues("KubebuilderCrd", req.NamespacedName)
	r.ctx = ctx
	//fmt.Println("INside reconcile")
	// TODO(user): your logic here

	var kubebuildercrd crdv1.KubebuilderCrd

	if err := r.Get(ctx, req.NamespacedName, &kubebuildercrd); err != nil {
		r.Log.Info("KubebuilderCrd not found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.kubebuildercrd = &kubebuildercrd

	if err := r.CheckDeployment(); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.CheckService(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
