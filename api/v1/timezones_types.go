/*

Copyright 2020 Theo Barber-Bany

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
	corev1 "k8s.io/api/core/v1"
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
	TimezonesPhaseNew        TimezonesPhase = "New"
	TimezonesPhaseInProgress TimezonesPhase = "InProgress"
	TimezonesPhaseCompleted  TimezonesPhase = "Completed"
)

// TimezonesStatus defines the observed state of Timezones
type TimezonesStatus struct {
	// Phase is used to determine which phase of the creation cycle a Timezones
	// is currently in.
	Phase TimezonesPhase `json:"phase"`

	// WallClocksCreated lists the names of all WallClocks created by the
	// controller for this Timezones.
	WallClocksCreated []string `json:"wallClocksCreated,omitempty"`

	// WallClocksCreatedCount is the count of WallClocksCreated.
	// This is used for printing in kubectl.
	WallClocksCreatedCount int `json:"wallClocksCreatedCount,omitempty"`

	// CompletionTimestamp is a timestamp for when the creation of WallClocks
	// has completed
	CompletionTimestamp *metav1.Time `json:"completionTimestamp,omitempty"`

	// Conditions gives detailed condition information about the Timezones
	//Conditions []TimezonesCondition `json:"conditions,omitempty"`
}

// TimezonesConditionType is the type of a TimezonesCondition
type TimezonesConditionType string

const (
	// WallClocksCreatedType refers to whether the controller successfully
	// created all of the required Timezoness
	WallClocksCreatedType TimezonesConditionType = "TimezonesCreated"

	// WallClocksInProgressType refers to whether the controller is currently
	// processing WallClocks
	WallClocksInProgressType TimezonesConditionType = "TimezonesInProgress"
)

// TimezonesConditionReason represents a valid condition reason for a Timezones
type TimezonesConditionReason string

// TimezonesCondition is a status condition for a Timezones
type TimezonesCondition struct {
	// Type of this condition
	Type TimezonesConditionType `json:"type"`

	// Status of this condition
	Status corev1.ConditionStatus `json:"status"`

	// this makes code gen blow up, I can't be bothered to hunt down why given
	// time constraints

	// LastUpdateTime of this condition
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`

	// LastTransitionTime of this condition
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason for the current status of this condition
	Reason TimezonesConditionReason `json:"reason,omitempty"`

	// Message associated with this condition
	Message string `json:"message,omitempty"`
}

// Timezones is the Schema for the timezones API
// +kubebuilder:object:root=true
// +genclient:nonNamespaced
// +kubebuilder:resource:path=timezones,scope=Cluster,shortName=tzone;tz;tzs
// +kubebuilder:printcolumn:name="WallClocks created",type="integer",JSONPath=".status.wallClocksCreatedCount",description="Number of WallClocks created"
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
