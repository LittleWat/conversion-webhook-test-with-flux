/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cwtestv1alpha2 "github.com/littlewat/conversion-webhook-test-with-flux/api/v1alpha2"
)

// TestResourceReconciler reconciles a TestResource object
type TestResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cwtest.littlewat.github.io,resources=testresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cwtest.littlewat.github.io,resources=testresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cwtest.littlewat.github.io,resources=testresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TestResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	rlog := log.Log.WithName("TestResourceReconciler")
	rlog.Info("----------------------------------")
	rlog.Info("Reconcile is called!")

	testResource := &cwtestv1alpha2.TestResource{}
	// Check if the Kafka object is defined, throw an error and requeue if not defined yet
	err := r.Get(ctx, req.NamespacedName, testResource)
	if err != nil {
		rlog.Info("Failed in r.Get", "testResource", testResource)
		return ctrl.Result{}, err
	}
	rlog.Info("Found CR spec", "testResource", testResource)

	testResource.Status.State = testResource.Spec.Foo + "-OK"
	err = r.Status().Update(ctx, testResource)
	if err != nil {
		rlog.Info("Failed in r.Status().Update", "testResource", testResource)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cwtestv1alpha2.TestResource{}).
		Complete(r)
}
