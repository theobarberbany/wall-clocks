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

package timezones

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	wallclocksv1 "github.com/ziglu/wallclocks/api/v1"
	"github.com/ziglu/wallclocks/controllers/timezones/status"
)

// TimezonesReconciler reconciles a Timezones object
type TimezonesReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=wallclocks.ziglu.io,resources=timezones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wallclocks.ziglu.io,resources=wallclocks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wallclocks.ziglu.io,resources=timezones/status,verbs=get;update;patch
// Automatically generate RBAC rules to allow the Controller to read and write TimeZones and WallClocks

// Reconcile reads that state of the cluster for a Timezones object and makes changes based on the state read
// and what is in the Timezones.Spec.Clocks
func (r *TimezonesReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("timezones", req.NamespacedName)

	instance := &wallclocksv1.Timezones{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Error(err, "unable to fetch instance")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Starting to reconcile", "instance", instance.GetName())
	result, err := r.HandleTimezonesObject(instance)
	if err != nil {
		// Ensure we attempt to update the status even when the handler fails
		statusErr := status.UpdateStatus(r.Client, instance, result)
		if statusErr != nil {
			log.Error(statusErr, "error updating status")
		}

		return reconcile.Result{}, fmt.Errorf("error handling timezones %s: %+v", instance.GetName(), err)
	}
	log.Info("Updating status", "instance", instance.GetName())
	err = status.UpdateStatus(r.Client, instance, result)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating status: %v", err)
	}

	return ctrl.Result{}, nil
}

// HandleTimezonesObject is a top level handler for reconciling a
// wallclocksv1.Timezones Instance
func (r *TimezonesReconciler) HandleTimezonesObject(instance *wallclocksv1.Timezones) (*status.Result, error) {
	switch instance.Status.Phase {
	case wallclocksv1.TimezonesPhaseNew:
		// For each entry of .Spec.Clocks
		// Create a WallClock
		// If any fail, set in the result otherwise set succeded
		// Update status
		return r.HandleNew(instance)
	case wallclocksv1.TimezonesPhaseCompleted:
		// NoOp. Maybe clean up after a certain time?
		return &status.Result{}, nil
	default:
		return r.HandleNew(instance)
	}
}

// HandleNew handles an instance of a timezones object in the new phase
func (r *TimezonesReconciler) HandleNew(instance *wallclocksv1.Timezones) (*status.Result, error) {
	// initialise result
	result := &status.Result{}

	// try creating a WallClock for every timezone, agregate failures and successes
	failures := map[string]error{}
	created := []string{}
	for _, timeZone := range instance.Spec.Clocks {
		ctx := context.Background()
		clock := newWallClock(instance.GetName(), timeZone, instance)
		createErr := r.Create(ctx, clock)
		if createErr != nil {
			// if it already exists, ignore
			// I know this is inefficient
			if apierrors.IsAlreadyExists(createErr) {
				continue
			}

			failures[timeZone] = createErr
			continue
		}
		created = append(created, timeZone)
	}

	var failedTimeZones []string
	var failedTimezonesError []error
	for k, v := range failures {
		failedTimeZones = append(failedTimeZones, k)

		failedTimezonesError = append(failedTimezonesError, fmt.Errorf("Failed to create WallClock for timezone %s: %v", k, v))

	}

	// new aggregate error
	aggregate := utilerrors.NewAggregate(failedTimezonesError)

	// update result
	result.WallClocksCreated = created
	result.WallClocksFailed = failedTimeZones
	result.WallClocksCreatedError = aggregate

	completedPhase := wallclocksv1.TimezonesPhaseCompleted
	result.Phase = &completedPhase

	return result, aggregate

}

func newWallClock(name string, timezone string, instance *wallclocksv1.Timezones) *wallclocksv1.WallClock {
	wallClock := &wallclocksv1.WallClock{
		TypeMeta: metav1.TypeMeta{
			APIVersion: wallclocksv1.GroupVersion.String(),
			Kind:       "WallClock",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", name, timezone),
			OwnerReferences: []metav1.OwnerReference{
				newOwnerRef(instance, instance.GroupVersionKind(), true, true),
			},
		},
		Spec: wallclocksv1.WallClockSpec{
			Timezone: timezone,
		},
		Status: wallclocksv1.WallClockStatus{},
	}
	return wallClock.DeepCopy()
}

// newOwnerRef creates an OwnerReference pointing to the given owner
func newOwnerRef(owner metav1.Object, gvk schema.GroupVersionKind, isController bool, blockOwnerDeletion bool) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: &blockOwnerDeletion,
		Controller:         &isController,
	}
}

// SetupWithManager sets up the Timezones controller with the manager
func (r *TimezonesReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&wallclocksv1.Timezones{}).
		Complete(r)
}
