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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WallClockSpec defines the desired state of WallClock
type WallClockSpec struct {
	// Timezone is the timezone for the clock
	Timezone string `json:"timezone,omitempty"`
}

// WallClockStatus defines the observed state of a WallClock
type WallClockStatus struct {
	// Time is the time on the WallClock
	Time metav1.Time `json:"time"`
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
