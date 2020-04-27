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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WallClockSpec defines the desired state of WallClock
type WallClockSpec struct {
	// Timezone is the timezone for the clock
	Timezone string `json:"timezone,omitempty"`
}

// WallClockPhase determines the phase in which the WallClock currently is in
type WallClockPhase string

// The following WallClockPhases enumerate all possible WallClockPhases
const (
	WallClockPhaseNew      WallClockPhase = "New"
	WallClockPhaseUpdating WallClockPhase = "Updating"
	WallClockPhaseUpdated  WallClockPhase = "Updated"
	WallClockPhaseFailed   WallClockPhase = "Failed"
)

// WallClockStatus defines the observed state of a WallClock
type WallClockStatus struct {
	// Time is the time on the WallClock
	Time metav1.Time `json:"time"`

	// Phase is used to determine which phase of the clock cycle a clock
	// is currently in.
	Phase WallClockPhase `json:"phase"`

	// Conditions gives detailed condition information about the clock
	Conditions []WallClockCondition `json:"conditions,omitempty"`
}

// WallClockConditionType is the type of a WallClockCondition
type WallClockConditionType string

const (
	// TimezoneParsedType refers to the type of condition where the controller
	// successfully managed parse a Timezone
	TimezoneParsedType WallClockConditionType = "TimezoneParsed"
)

const (
	// ReasonGotTimezone refers to whether the controller successfully managed to
	// locate the desired timezone
	ReasonGotTimezone WallClockConditionReason = "GotTimezone"

	// ReasonErrorGettingTimezone is a clock condition for an error finding the
	// time zone
	ReasonErrorGettingTimezone WallClockConditionReason = "ErrorGettingTimezone"
)

// WallClockConditionReason represents a valid condition reason for a WallClock
type WallClockConditionReason string

// WallClockCondition is a status condition for a WallClock
type WallClockCondition struct {
	// Type of this condition
	Type WallClockConditionType `json:"type"`

	// Status of this condition
	Status corev1.ConditionStatus `json:"status"`

	// LastUpdateTime of this condition
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`

	// LastTransitionTime of this condition
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason for the current status of this condition
	Reason WallClockConditionReason `json:"reason,omitempty"`

	// Message associated with this condition
	Message string `json:"message,omitempty"`
}

// WallClock is the Schema for the wallclocks API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=wallclock,scope=Cluster,shortName=wclock;wc;wcs
type WallClock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WallClockSpec   `json:"spec,omitempty"`
	Status WallClockStatus `json:"status,omitempty"`
}

// WallClockList contains a list of WallClock
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=wallclockList,scope=Cluster
type WallClockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WallClock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WallClock{}, &WallClockList{})
}
