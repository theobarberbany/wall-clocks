package v1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Timezones", func() {
	key := types.NamespacedName{
		Name: "foo",
	}
	created := &Timezones{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}

	It("can create, update and delete the object", func() {
		// Test Create
		fetched := &Timezones{}
		Expect(c.Create(context.TODO(), created)).NotTo(HaveOccurred())

		Expect(c.Get(context.TODO(), key, fetched)).NotTo(HaveOccurred())
		Expect(fetched).To(Equal(created))

		// Test Updating the Labels
		updated := fetched.DeepCopy()
		updated.Labels = map[string]string{"hello": "world"}
		Expect(c.Update(context.TODO(), updated)).NotTo(HaveOccurred())

		Expect(c.Get(context.TODO(), key, fetched)).NotTo(HaveOccurred())
		Expect(fetched).To(Equal(updated))

		// Test Delete
		Expect(c.Delete(context.TODO(), fetched)).NotTo(HaveOccurred())
		Expect(c.Get(context.TODO(), key, fetched)).To(HaveOccurred())
	})
})
