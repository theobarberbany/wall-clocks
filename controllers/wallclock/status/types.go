package status

// Result is used as the basis to updating the status of the WallClock object.
// It contains information gathered during a single run of the reconcile loop.
type Result struct {
	// This is the time on the clock
	Time *string
}
