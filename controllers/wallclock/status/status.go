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
func UpdateStatus(c client.Client, instance *wallclocksv1.WallClock, result *Result) error {
	status := instance.Status

	setTime(&status, result)

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

// setTime sets the phase when it is set in the Result
func setTime(status *wallclocksv1.WallClockStatus, result *Result) {
	if result.Time != nil {
		status.Time = result.Time
	}
}
