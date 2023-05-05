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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mondepllog = logf.Log.WithName("mondepl-resource")

func (r *MonDepl) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kally-io-kally-io-v1alpha1-mondepl,mutating=true,failurePolicy=fail,sideEffects=None,groups=kally.io.kally.io,resources=mondepls,verbs=create;update,versions=v1alpha1,name=mmondepl.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MonDepl{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MonDepl) Default() {
	mondepllog.Info("Defaulting a little bit...")

	if r.Spec.Reason == "" {
		r.Spec.Reason = "love"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-kally-io-kally-io-v1alpha1-mondepl,mutating=false,failurePolicy=fail,sideEffects=None,groups=kally.io.kally.io,resources=mondepls,verbs=create;update,versions=v1alpha1,name=vmondepl.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MonDepl{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MonDepl) ValidateCreate() error {
	mondepllog.Info("Validating creation a little bit...")
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MonDepl) ValidateUpdate(old runtime.Object) error {
	mondepllog.Info("Validating update a little bit...")
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MonDepl) ValidateDelete() error {
	mondepllog.Info("Validating deletion a little bit...")
	return nil
}
