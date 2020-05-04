package utils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega"
	gtypes "github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Matcher has Gomega Matchers that use the controller-runtime client
type Matcher struct {
	Client client.Client
}

// Object is the combination of two interfaces as a helper for passing
// Kubernetes objects between methods
type Object interface {
	runtime.Object
	metav1.Object
}

// UpdateFunc modifies the object fetched from the API server before sending
// the update
type UpdateFunc func(Object) Object

// Create creates the object on the API server
func (m *Matcher) Create(obj Object, extras ...interface{}) gomega.GomegaAssertion {
	err := m.Client.Create(context.TODO(), obj)
	return gomega.Expect(err, extras)
}

// Delete deletes the object from the API server
func (m *Matcher) Delete(obj Object, extras ...interface{}) gomega.GomegaAssertion {
	err := m.Client.Delete(context.TODO(), obj)
	return gomega.Expect(err, extras)
}

// Update udpates the object on the API server by fetching the object
// and applying a mutating UpdateFunc before sending the update
func (m *Matcher) Update(obj Object, fn UpdateFunc, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	key := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
	update := func() error {
		err := m.Client.Get(context.TODO(), key, obj)
		if err != nil {
			return err
		}
		return m.Client.Update(context.TODO(), fn(obj))
	}
	return gomega.Eventually(update, intervals...)
}

// UpdateStatus udpates the object's status subresource on the API server by
// fetching the object and applying a mutating UpdateFunc before sending the
// update
func (m *Matcher) UpdateStatus(obj Object, fn UpdateFunc, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	key := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
	update := func() error {
		err := m.Client.Get(context.TODO(), key, obj)
		if err != nil {
			return err
		}
		return m.Client.Status().Update(context.TODO(), fn(obj))
	}
	return gomega.Eventually(update, intervals...)
}

// Get gets the object from the API server
func (m *Matcher) Get(obj Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	key := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
	get := func() error {
		return m.Client.Get(context.TODO(), key, obj)
	}
	return gomega.Eventually(get, intervals...)
}

// List gets the list object from the API server
func (m *Matcher) List(obj runtime.Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	list := func() error {
		return m.Client.List(context.TODO(), obj)
	}
	return gomega.Eventually(list, intervals...)
}

// Consistently continually gets the object from the API for comparison
func (m *Matcher) Consistently(obj runtime.Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	// If the object is a list, return a list
	if meta.IsListType(obj) {
		return m.consistentlyList(obj, intervals...)
	}
	if o, ok := obj.(Object); ok {
		return m.consistentlyObject(o, intervals...)
	}
	//Should not get here
	panic("Unknown object.")
}

// consistentlyObject gets an individual object from the API server
func (m *Matcher) consistentlyObject(obj Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	key := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
	get := func() Object {
		err := m.Client.Get(context.TODO(), key, obj)
		if err != nil {
			panic(err)
		}
		return obj
	}
	return gomega.Consistently(get, intervals...)
}

// consistentlyObject gets an individual object from the API server
func (m *Matcher) consistentlyList(obj runtime.Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	list := func() runtime.Object {
		err := m.Client.List(context.TODO(), obj)
		if err != nil {
			panic(err)
		}
		return obj
	}
	return gomega.Consistently(list, intervals...)
}

// Eventually continually gets the object from the API for comparison
func (m *Matcher) Eventually(obj runtime.Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	// If the object is a list, return a list
	if meta.IsListType(obj) {
		return m.eventuallyList(obj, intervals...)
	}
	if o, ok := obj.(Object); ok {
		return m.eventuallyObject(o, intervals...)
	}
	//Should not get here
	panic("Unknown object.")
}

// eventuallyObject gets an individual object from the API server
func (m *Matcher) eventuallyObject(obj Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	key := types.NamespacedName{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
	get := func() Object {
		err := m.Client.Get(context.TODO(), key, obj)
		if err != nil {
			panic(err)
		}
		return obj
	}
	return gomega.Eventually(get, intervals...)
}

// eventuallyList gets a list type  from the API server
func (m *Matcher) eventuallyList(obj runtime.Object, intervals ...interface{}) gomega.GomegaAsyncAssertion {
	list := func() runtime.Object {
		err := m.Client.List(context.TODO(), obj)
		if err != nil {
			panic(err)
		}
		return obj
	}
	return gomega.Eventually(list, intervals...)
}

// WithListItems returns the items of the list
func WithListItems(matcher gtypes.GomegaMatcher) gtypes.GomegaMatcher {
	return gomega.WithTransform(func(obj runtime.Object) []runtime.Object {
		items, err := meta.ExtractList(obj)
		if err != nil {
			panic(err)
		}
		return items
	}, matcher)
}

// WithField gets the value of the named field from the object
func WithField(field string, matcher gtypes.GomegaMatcher) gtypes.GomegaMatcher {
	// Addressing Field by <struct>.<field> can be recursed
	fields := strings.SplitN(field, ".", 2)
	if len(fields) == 2 {
		matcher = WithField(fields[1], matcher)
	}

	return gomega.WithTransform(func(obj interface{}) interface{} {
		r := reflect.ValueOf(obj)
		f := reflect.Indirect(r).FieldByName(fields[0])
		if !f.IsValid() {
			panic(fmt.Sprintf("Object '%s' does not have a field '%s'", reflect.TypeOf(obj), fields[0]))
		}
		return f.Interface()
	}, matcher)
}
