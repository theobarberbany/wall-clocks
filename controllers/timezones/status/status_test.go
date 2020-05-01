package status

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wallclocksv1 "github.com/ziglu/wallclocks/api/v1"
	"github.com/ziglu/wallclocks/test/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Timezones Status Suite", func() {
	var c client.Client
	var m utils.Matcher

	var timezones *wallclocksv1.Timezones
	var result *Result

	const timeout = time.Second * 5
	const consistentlyTimeout = time.Second
	const testStr = "test"

	BeforeEach(func() {
		var err error
		c, err = client.New(cfg, client.Options{})
		Expect(err).NotTo(HaveOccurred())
		m = utils.Matcher{Client: c}

		timezones = &wallclocksv1.Timezones{
			ObjectMeta: metav1.ObjectMeta{
				Name: "example",
			},
			Spec: wallclocksv1.TimezonesSpec{
				Clocks: []string{
					"PDT", "GMT",
				},
			},
		}
		m.Create(timezones).Should(Succeed())

		result = &Result{}
	})

	AfterEach(func() {
		utils.DeleteAll(cfg, timeout,
			&wallclocksv1.TimezonesList{},
		)
	})

	Context("UpdateStatus", func() {
		var updateErr error

		JustBeforeEach(func() {
			updateErr = UpdateStatus(c, timezones, result)
		})

		Context("when the phase is set in the Result", func() {
			var phase wallclocksv1.TimezonesPhase

			BeforeEach(func() {
				phase = wallclocksv1.TimezonesPhaseNew
				Expect(timezones.Status.Phase).ToNot(Equal(phase))
				result.Phase = &phase
			})

			It("updates the phase in the status", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.Phase", Equal(phase)))
			})

			It("does not cause an error", func() {
				Expect(updateErr).To(BeNil())
			})
		})

		Context("when no existing WallClocksCreated is set", func() {
			var wallClocksCreated []string

			BeforeEach(func() {
				wallClocksCreated = []string{"GMT", "PDT", "BST", "SGT"}
				Expect(timezones.Status.WallClocksCreated).To(BeEmpty())
				result.WallClocksCreated = wallClocksCreated
			})

			It("sets the WallClocksCreated field", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.WallClocksCreated", Equal(wallClocksCreated)))
			})

			It("sets the WallClocksCreatedCount field", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.WallClocksCreatedCount", Equal(len(wallClocksCreated))))
			})

			It("does not cause an error", func() {
				Expect(updateErr).To(BeNil())
			})
		})

		Context("when an existing WallClocksCreated is set", func() {
			var wallClocksCreated []string
			var existingWallClocksCreated []string
			var expectedWallClocksCreated []string

			BeforeEach(func() {
				// Set up the existing expected state
				existingWallClocksCreated = []string{"BST", "SGT"}
				m.Update(timezones, func(obj utils.Object) utils.Object {
					nr, _ := obj.(*wallclocksv1.Timezones)
					nr.Status.WallClocksCreated = existingWallClocksCreated
					nr.Status.WallClocksCreatedCount = len(existingWallClocksCreated)
					return nr
				}, timeout).Should(Succeed())

				// Introduce some duplication, this implicitly tests for de-duplication.
				wallClocksCreated = []string{"SGT", "GMT", "SGT"}
				result.WallClocksCreated = wallClocksCreated
				expectedWallClocksCreated = []string{"SGT", "GMT", "BST"}

			})

			It("joins the new and existing WallClocksCreated field", func() {
				m.Eventually(timezones, timeout).Should(
					utils.WithField("Status.WallClocksCreated", ConsistOf(expectedWallClocksCreated)),
				)
			})

			It("updates the WallClocksCreatedCount field", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.WallClocksCreatedCount", Equal(len(expectedWallClocksCreated))))
			})

			It("does not cause an error", func() {
				Expect(updateErr).To(BeNil())
			})
		})

		Context("when no existing WallClocksCreated is set", func() {
			var wallClocksFailed []string

			BeforeEach(func() {
				wallClocksFailed = []string{"ASDF", "QWERTY"}
				Expect(timezones.Status.WallClocksFailed).To(BeEmpty())
				result.WallClocksFailed = wallClocksFailed
			})

			It("sets the WallClocksFailed field", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.WallClocksFailed", Equal(wallClocksFailed)))
			})

			It("sets the WallClocksFailedCount field", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.WallClocksFailedCount", Equal(len(wallClocksFailed))))
			})

			It("does not cause an error", func() {
				Expect(updateErr).To(BeNil())
			})
		})

		Context("when an existing WallClocksFailed is set", func() {
			var wallClocksFailed []string
			var existingWallClocksFailed []string
			var expectedWallClocksFailed []string

			BeforeEach(func() {
				// Set up the existing expected state
				existingWallClocksFailed = []string{"QWERTY", "ASDF"}
				m.Update(timezones, func(obj utils.Object) utils.Object {
					nr, _ := obj.(*wallclocksv1.Timezones)
					nr.Status.WallClocksFailed = existingWallClocksFailed
					nr.Status.WallClocksFailedCount = len(existingWallClocksFailed)
					return nr
				}, timeout).Should(Succeed())

				// Introduce some duplication, this implicitly tests for de-duplication.
				wallClocksFailed = []string{"QWERTY", "ASDF", "ASDF", "LMNOP"}
				result.WallClocksFailed = wallClocksFailed
				expectedWallClocksFailed = []string{"QWERTY", "ASDF", "LMNOP"}

			})

			It("joins the new and existing WallClocksFailed field", func() {
				m.Eventually(timezones, timeout).Should(
					utils.WithField("Status.WallClocksFailed", ConsistOf(expectedWallClocksFailed)),
				)
			})

			It("updates the WallClocksFailedCount field", func() {
				m.Eventually(timezones, timeout).Should(utils.WithField("Status.WallClocksFailedCount", Equal(len(expectedWallClocksFailed))))
			})

			It("does not cause an error", func() {
				Expect(updateErr).To(BeNil())
			})
		})

		Context("when no existing CompletionTimestamp is set", func() {
			var completionTimestamp metav1.Time

			Context("and there is a CompletionTimestamp set in the Result", func() {
				BeforeEach(func() {
					completionTimestamp = metav1.Now()
					Expect(timezones.Status.CompletionTimestamp).To(BeNil())
					result.CompletionTimestamp = &completionTimestamp
				})

				It("sets the CompletionTimestamp field", func() {
					m.Eventually(timezones, timeout).Should(utils.WithField("Status.CompletionTimestamp", Equal(&completionTimestamp)))
				})

				It("does not cause an error", func() {
					Expect(updateErr).To(BeNil())
				})
			})

			Context("and there is not a CompletionTimestamp set in the Result", func() {
				BeforeEach(func() {
					Expect(timezones.Status.CompletionTimestamp).To(BeNil())
				})

				It("does not set the CompletionTimestamp", func() {
					m.Consistently(timezones, consistentlyTimeout).Should(utils.WithField("Status.CompletionTimestamp", BeNil()))
				})
			})
		})

		Context("when an existing CompletionTimestamp is set and CompletionTimestamp is set in  the Result", func() {
			var completionTimestamp metav1.Time
			var existingCompletionTimestamp metav1.Time

			BeforeEach(func() {
				// Set up the existing expected state
				existingCompletionTimestamp = metav1.NewTime(metav1.Now().Add(-time.Hour))
				m.Update(timezones, func(obj utils.Object) utils.Object {
					nr, _ := obj.(*wallclocksv1.Timezones)
					nr.Status.CompletionTimestamp = &existingCompletionTimestamp
					return nr
				}, timeout).Should(Succeed())

				completionTimestamp = metav1.Now()
				result.CompletionTimestamp = &completionTimestamp
			})

			It("does not update the CompletionTimestamp field", func() {
				m.Consistently(timezones, consistentlyTimeout).Should(utils.WithField("Status.CompletionTimestamp", Equal(&existingCompletionTimestamp)))
			})

			It("returns an error", func() {
				Expect(updateErr).ToNot(BeNil())
				Expect(updateErr.Error()).To(Equal("cannot update CompletionTimestamp, field is immutable once set"))
			})
		})

	})
})
