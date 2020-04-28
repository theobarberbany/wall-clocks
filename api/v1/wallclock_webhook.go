/*

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

package v1

import (
	"fmt"
	"time"

	"github.com/tkuchiki/go-timezone"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var wallclocklog = logf.Log.WithName("wallclock-resource")

//SetupWebhookWithManager registers the webook with the manager
func (r *WallClock) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-wallclocks-ziglu-io-v1-wallclock,mutating=false,failurePolicy=fail,groups=wallclocks.ziglu.io,resources=wallclocks,versions=v1,name=vwallclock.kb.io

var _ webhook.Validator = &WallClock{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *WallClock) ValidateCreate() error {
	wallclocklog.Info("validate create", "name", r.Name)

	return r.validateWallClock()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered
// for the type
func (r *WallClock) ValidateUpdate(old runtime.Object) error {
	wallclocklog.Info("validate update", "name", r.Name)

	return r.validateWallClock()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered
// for the type. We don't care about the validity of a WallClock we are
// deleting, so always return nil
func (r *WallClock) ValidateDelete() error {
	wallclocklog.Info("validate delete", "name", r.Name)

	return nil
}

// validateWallClock checks that the requested timezone is valid
func (r *WallClock) validateWallClock() error {
	// First check, see if the built in time library can find the timezone
	// If the timezone is empty it'll assume UTC, which is OK I think?
	_, err := time.LoadLocation(r.Spec.Timezone)
	if err != nil {
		// Secondary check, is it possible the provided timezone is a short
		// code, e.g 'PDT'
		_, err2 := timezone.GetTimezones(r.Spec.Timezone)
		if err2 != nil {
			return fmt.Errorf("Invalid timezone: %s", r.Spec.Timezone)
		}

		return nil
	}

	return nil
}
