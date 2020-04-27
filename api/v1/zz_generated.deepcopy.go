// +build !ignore_autogenerated

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

// autogenerated by controller-gen object, do not modify manually

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Timezones) DeepCopyInto(out *Timezones) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Timezones.
func (in *Timezones) DeepCopy() *Timezones {
	if in == nil {
		return nil
	}
	out := new(Timezones)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Timezones) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimezonesCondition) DeepCopyInto(out *TimezonesCondition) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimezonesCondition.
func (in *TimezonesCondition) DeepCopy() *TimezonesCondition {
	if in == nil {
		return nil
	}
	out := new(TimezonesCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimezonesList) DeepCopyInto(out *TimezonesList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Timezones, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimezonesList.
func (in *TimezonesList) DeepCopy() *TimezonesList {
	if in == nil {
		return nil
	}
	out := new(TimezonesList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TimezonesList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimezonesSpec) DeepCopyInto(out *TimezonesSpec) {
	*out = *in
	if in.Clocks != nil {
		in, out := &in.Clocks, &out.Clocks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimezonesSpec.
func (in *TimezonesSpec) DeepCopy() *TimezonesSpec {
	if in == nil {
		return nil
	}
	out := new(TimezonesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimezonesStatus) DeepCopyInto(out *TimezonesStatus) {
	*out = *in
	if in.WallClocksCreated != nil {
		in, out := &in.WallClocksCreated, &out.WallClocksCreated
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.CompletionTimestamp != nil {
		in, out := &in.CompletionTimestamp, &out.CompletionTimestamp
		*out = new(metav1.Time)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimezonesStatus.
func (in *TimezonesStatus) DeepCopy() *TimezonesStatus {
	if in == nil {
		return nil
	}
	out := new(TimezonesStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WallClock) DeepCopyInto(out *WallClock) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WallClock.
func (in *WallClock) DeepCopy() *WallClock {
	if in == nil {
		return nil
	}
	out := new(WallClock)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WallClock) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WallClockCondition) DeepCopyInto(out *WallClockCondition) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WallClockCondition.
func (in *WallClockCondition) DeepCopy() *WallClockCondition {
	if in == nil {
		return nil
	}
	out := new(WallClockCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WallClockList) DeepCopyInto(out *WallClockList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WallClock, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WallClockList.
func (in *WallClockList) DeepCopy() *WallClockList {
	if in == nil {
		return nil
	}
	out := new(WallClockList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WallClockList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WallClockSpec) DeepCopyInto(out *WallClockSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WallClockSpec.
func (in *WallClockSpec) DeepCopy() *WallClockSpec {
	if in == nil {
		return nil
	}
	out := new(WallClockSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WallClockStatus) DeepCopyInto(out *WallClockStatus) {
	*out = *in
	in.Time.DeepCopyInto(&out.Time)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]WallClockCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WallClockStatus.
func (in *WallClockStatus) DeepCopy() *WallClockStatus {
	if in == nil {
		return nil
	}
	out := new(WallClockStatus)
	in.DeepCopyInto(out)
	return out
}
