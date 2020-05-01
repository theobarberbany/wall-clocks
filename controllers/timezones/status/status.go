package status

import (
	"context"
	"fmt"
	"reflect"

	wallclocksv1 "github.com/ziglu/wallclocks/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateStatus merges the status in the existing instance with the information
// provided in the Result and then updates the instance if there is any
// difference between the new and updated status
func UpdateStatus(c client.Client, instance *wallclocksv1.Timezones, result *Result) error {
	status := instance.Status

	setPhase(&status, result)

	err := setWallClocksCreated(&status, result)
	if err != nil {
		return err
	}

	setWallClocksFailed(&status, result)

	err = setCompletionTimestamp(&status, result)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(status, instance.Status) {
		copy := instance.DeepCopy()
		copy.Status = status

		err := c.Update(context.TODO(), copy)
		if err != nil {
			return fmt.Errorf("error updating status: %v", err)
		}
	}

	return nil
}

// setPhase sets the phase when it is set in the Result
func setPhase(status *wallclocksv1.TimezonesStatus, result *Result) {
	if result.Phase != nil {
		status.Phase = *result.Phase
	}
}

// setWallClocksCreated sets the WallClocksCreated field,if it has not been
// set before it is added. If it has been set before the two are appended
func setWallClocksCreated(status *wallclocksv1.TimezonesStatus, result *Result) error {
	if status.WallClocksCreated != nil && result.WallClocksCreated != nil {
		status.WallClocksCreated = appendIfMissingStr(status.WallClocksCreated, result.WallClocksCreated...)
		status.WallClocksCreatedCount = len(status.WallClocksCreated)
	}

	if status.WallClocksCreated == nil && result.WallClocksCreated != nil {
		status.WallClocksCreated = result.WallClocksCreated
		status.WallClocksCreatedCount = len(result.WallClocksCreated)
	}

	return nil
}

// setWallClocksFailed sets the WallClocksFailed field, if it has not been
// set before it is added. If it has been set before the two are appended
func setWallClocksFailed(status *wallclocksv1.TimezonesStatus, result *Result) {
	if status.WallClocksFailed != nil && result.WallClocksFailed != nil {
		status.WallClocksFailed = appendIfMissingStr(status.WallClocksFailed, result.WallClocksFailed...)
		status.WallClocksFailedCount = len(status.WallClocksFailed)
	}

	if status.WallClocksFailed == nil && result.WallClocksFailed != nil {
		status.WallClocksFailed = result.WallClocksFailed
		status.WallClocksFailedCount = len(result.WallClocksFailed)
	}

}

// setCompletionTimestamp sets the setCompletionTimestamp field. If it has not
// been set before it is added. If it has been set before an error is returned
func setCompletionTimestamp(status *wallclocksv1.TimezonesStatus, result *Result) error {
	if status.CompletionTimestamp != nil && result.CompletionTimestamp != nil {
		return fmt.Errorf("cannot update CompletionTimestamp, field is immutable once set")
	}

	if status.CompletionTimestamp == nil && result.CompletionTimestamp != nil {
		status.CompletionTimestamp = result.CompletionTimestamp
	}
	return nil
}

// appendIfMissingStr will append two []string(s) dropping duplicate elements
func appendIfMissingStr(slice []string, str ...string) []string {
	merged := slice
	for _, ele := range str {
		merged = appendIfMissingElement(merged, ele)
	}
	return merged
}

// appendIfMissingElement will append a string to a []string only if it is
// unique
func appendIfMissingElement(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
