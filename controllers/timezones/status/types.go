package status

import (
	wallclocksv1 "github.com/ziglu/wallclocks/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Result is used as the basis to updating the status of the Timezones object.
// It contains information gathered during a single run of the reconcile loop.
type Result struct {
	// This represents the Phase of the Timezones object that the status should be set
	// to when updating the status.
	// If Phase == nil, don't update the Phase, else, overwrite it.
	Phase *wallclocksv1.TimezonesPhase

	// This should contain any errors related to the creation of the WallClocks
	WallClocksCreatedError error

	// This should list all WallClocks created
	WallClocksCreated []string

	// This should be a list of any WallClocks that failed to create
	WallClocksFailed []string

	// CompletionTimestamp is a timestamp for when the creation of WallClocks
	// has completed
	CompletionTimestamp *metav1.Time
}
