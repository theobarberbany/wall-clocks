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

package wallclock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"github.com/tkuchiki/go-timezone"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	wallclocksv1 "github.com/ziglu/wallclocks/api/v1"
	"github.com/ziglu/wallclocks/controllers/wallclock/status"
)

// WallClockReconciler reconciles a WallClock object
type WallClockReconciler struct {
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=wallclocks.ziglu.io,resources=wallclocks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wallclocks.ziglu.io,resources=wallclocks/status,verbs=get;update;patch

// Reconcile reads that state of the cluster for a WallClock object and makes changes based on the state read
// and what is in the WallClock.Status.Time
func (r *WallClockReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("wallclock", req.NamespacedName)

	instance := &wallclocksv1.WallClock{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Error(err, "unable to fetch instance")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Starting to reconcile", "instance", instance.GetName())
	result, err := r.HandleWallClockObject(instance)
	if err != nil {
		// Ensure we attempt to update the status even when the handler fails
		statusErr := status.UpdateStatus(r.Client, instance, result)
		if statusErr != nil {
			log.Error(statusErr, "error updating status")
		}

		return reconcile.Result{}, fmt.Errorf("error handling wallclock %s: %+v", instance.GetName(), err)
	}
	log.Info("Updating status", "instance", instance.GetName())
	err = status.UpdateStatus(r.Client, instance, result)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating status: %v", err)
	}

	return ctrl.Result{}, nil
}

// HandleWallClockObject is a top level handler for reconciling a
// wallclocksv1.WallClock Instance. It sets the time in the result. If the
// TimeZone fails to parse it will return an error that is logged. This should
// not happen, as the validating admission webhook should prevent any WallClock
// with an invalid timezone set from being created.
func (r *WallClockReconciler) HandleWallClockObject(instance *wallclocksv1.WallClock) (*status.Result, error) {
	// initialise result
	result := &status.Result{}

	var time string
	time, err := TimeInWrapper(instance.Spec.Timezone)
	if err != nil {
		locations, err := timezone.GetTimezones(instance.Spec.Timezone)
		if err != nil {
			return result, fmt.Errorf("failed to get time for timezone %s: %v", instance.Spec.Timezone, err)
		}

		time, err = TimeInWrapper(locations[0])
		if err != nil {
			return result, fmt.Errorf("failed to get time for timezone %s: %v", instance.Spec.Timezone, err)
		}

		result.Time = &time
		return result, nil

	}

	result.Time = &time
	return result, nil
}

// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// TimeInWrapper wraps TimeIn to make checking a potentially invalid timezone
// twice neater. If the time is succesfully calculated, it is returned as a
// string.
func TimeInWrapper(name string) (string, error) {
	t, err := TimeIn(time.Now(), name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(t.Format("15:04:01")), nil
}

func (r *WallClockReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wallclocksv1.WallClock{}).
		Complete(r)
}
