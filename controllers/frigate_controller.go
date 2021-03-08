/*
Copyright 2021.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	shipv1beta1 "github.com/openshift/route-monitor-operator/api/v1beta1"
)

// FrigateReconciler reconciles a Frigate object
type FrigateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ship.my.domain,resources=frigates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ship.my.domain,resources=frigates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ship.my.domain,resources=frigates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Frigate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *FrigateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("frigate", req.NamespacedName)
	res := &v1.Pod{}
	if err := r.Client.Get(ctx, req.NamespacedName, res); err != nil {
		if !k8serrors.IsNotFound(err) {
			return ctrl.Result{}, errors.Wrap(err, "G: resource pre-getting failed")
		}
	} else if res != new(v1.Pod) {
		// stop
		return ctrl.Result{Requeue: false}, nil

	}

	// it's not found, we can create it
	resource := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: req.NamespacedName.Name, Namespace: req.NamespacedName.Namespace},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "bob",
					Image: "nginx",
				}},
		},
	}
	if err := r.Client.Create(ctx, &resource); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "U: resouce could not be created")
	}
	resource.Spec.Containers[0].Image = "redis"

	if err := r.Client.Update(ctx, &resource); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "C: resouce could not be updated")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FrigateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipv1beta1.Frigate{}).
		Complete(r)
}
