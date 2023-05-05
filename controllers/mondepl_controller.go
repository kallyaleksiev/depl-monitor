/*
Copyright 2023 Kaloyan Aleksiev.

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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/kallyaleksiev/depl-monitor/api/v1alpha1"
)

// MonDeplReconciler reconciles a MonDepl object
type MonDeplReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kally.io.kally.io,resources=mondepls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kally.io.kally.io,resources=mondepls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kally.io.kally.io,resources=mondepls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MonDeplReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// return ctrl.Result{Requeue: true}, fmt.Errorf("pls redo")
	logger := log.FromContext(ctx)

	var monDepl v1alpha1.MonDepl
	if err := r.Get(ctx, req.NamespacedName, &monDepl); err != nil {

		if errors.IsNotFound(err) {
			logger.Info("MonDepl cannot be found, seems to be deleted, will attempt to delete underlying...")
			var underlying appsv1.Deployment
			var underName = types.NamespacedName{
				Name:      GetUnderlyingName(req.Name),
				Namespace: req.Namespace,
			}

			if err := r.Get(ctx, underName, &underlying); err != nil {
				if errors.IsNotFound(err) {
					logger.Info("Underlying already deleted...")
					return ctrl.Result{}, nil
				} else {
					logger.Error(err, "Could not get underlying...")
					return ctrl.Result{}, err
				}
			}

			logger.Info("Trying to delete underlying...")

			if err := r.Delete(ctx, &underlying); err != nil {
				logger.Error(err, "Fail to delete underlying...")
				return ctrl.Result{Requeue: true}, err

			}
			return ctrl.Result{}, nil

		} else {
			logger.Error(err, "Could not get MonDepl...")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	var underlying appsv1.Deployment
	var underName = types.NamespacedName{
		Name:      GetUnderlyingName(req.Name),
		Namespace: req.Namespace,
	}
	var isDeleting = !monDepl.DeletionTimestamp.IsZero()

	logger.Info("Trying to get underlying...")

	if err := r.Get(ctx, underName, &underlying); err != nil {
		if errors.IsNotFound(err) && isDeleting {
			logger.Info("Underlying already deleted...")
			return ctrl.Result{}, nil
		} else if !errors.IsNotFound(err) {
			logger.Error(err, "Could not get underlying...")
			return ctrl.Result{Requeue: true}, err
		} else {
			underlying = appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      GetUnderlyingName(req.Name),
					Namespace: req.Namespace,
				},
				Spec: monDepl.Spec.Underlying,
			}

			if err := controllerutil.SetControllerReference(&monDepl, &underlying, r.Scheme); err != nil {
				logger.Error(err, "Could not set owner reference...")
				return ctrl.Result{Requeue: true}, err
			}

			if err := r.Create(ctx, &underlying); err != nil {
				logger.Error(err, "Could not create underlying...")
				return ctrl.Result{Requeue: true}, err
			}

			logger.Info("Updating MonDepl status...")
			underRef, err := ref.GetReference(r.Scheme, &underlying)
			if err != nil {
				logger.Error(err, "Could not get object reference for underlying...")
				return ctrl.Result{Requeue: true}, err
			}

			monDepl.Status.Active = underRef
			if err := r.Status().Update(ctx, &monDepl); err != nil {
				logger.Error(err, "Could not update status of MonDepl...")

				return ctrl.Result{Requeue: true}, err
			}

			return ctrl.Result{}, nil
		}
	}

	if isDeleting {

		logger.Info("Trying to delete underlying...")

		if err := r.Delete(ctx, &underlying); err != nil {
			logger.Error(err, "Fail to delete underlying...")
			return ctrl.Result{Requeue: true}, err

		}
		return ctrl.Result{}, nil
	}

	logger.Info("Trying to set owner reference...")

	if err := controllerutil.SetControllerReference(&monDepl, &underlying, r.Scheme); err != nil {
		logger.Error(err, "Could not set owner reference...")
		return ctrl.Result{Requeue: true}, err
	}

	logger.Info("Trying to update underlying...")

	underlying.Spec = monDepl.Spec.Underlying

	if err := r.Update(ctx, &underlying); err != nil {
		logger.Error(err, "Could not update underlying...")
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonDeplReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.MonDepl{}).
		Complete(r)
}
