package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TimezonesSpec defines the desired state of Timezones
type TimezonesSpec struct {
	// Clocks define the WallClocks to be created from this set of 'Timezones'.
	// It should be populated with short codes (GMT) or locations as defined by
	// the go Time package ("Europe/London")
	Clocks []string `json:"clocks,omitempty"`
}

// TimezonesPhase determines the phase in which the Timezones currently is
type TimezonesPhase string

// The following TimezonesPhases enumerate all possible TimezonesPhases
const (
	TimezonesPhaseNew       TimezonesPhase = "New"
	TimezonesPhaseCompleted TimezonesPhase = "Completed"
)

// TimezonesStatus defines the observed state of Timezones
type TimezonesStatus struct {
	// Phase is used to determine which phase of the creation cycle a Timezones
	// is currently in.
	Phase TimezonesPhase `json:"phase"`

	// WallClocksCreated lists the names of all WallClocks created by the
	// controller for the given Timezones.
	WallClocksCreated []string `json:"wallClocksCreated,omitempty"`

	// WallClocksCreatedCount is the count of WallClocksCreated.
	// This is used for printing in kubectl.
	WallClocksCreatedCount int `json:"wallClocksCreatedCount,omitempty"`

	// WallClocksFailed lists the names of all WallClocks that the controller failed to create
	// for the given Timezones.
	WallClocksFailed []string `json:"wallClocksFailed,omitempty"`

	// WallClocksFailedCount is the count of WallClocksFailed.
	// This is used for printing in kubectl.
	WallClocksFailedCount int `json:"wallClocksFailedCount,omitempty"`

	// CompletionTimestamp is a timestamp for when the creation of WallClocks
	// has completed
	CompletionTimestamp *metav1.Time `json:"completionTimestamp,omitempty"`
}

// Timezones is the Schema for the timezones API
// +kubebuilder:object:root=true
// +genclient:nonNamespaced
// +kubebuilder:resource:path=timezones,scope=Cluster,shortName=tzone;tz;tzs
// +kubebuilder:printcolumn:name="WallClocks created",type="integer",JSONPath=".status.wallClocksCreatedCount",description="Number of WallClocks created"
// +kubebuilder:printcolumn:name="WallClocks failed",type="integer",JSONPath=".status.wallClocksFailedCount",description="Number of WallClocks that failed to create"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Completed",type="date",JSONPath=".status.completionTimestamp",description="The time since the creation of WallClocks completed"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Timezones struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TimezonesSpec   `json:"spec,omitempty"`
	Status TimezonesStatus `json:"status,omitempty"`
}

// TimezonesList contains a list of Timezones
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=timezonesList,scope=Cluster
type TimezonesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Timezones `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Timezones{}, &TimezonesList{})
}
