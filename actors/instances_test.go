package actors_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/actors/actorsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Instance Actor", func() {

	var (
		fakeInstanceRepository *actorsfakes.FakeInstanceRepository
		instancesActor         *InstancesActor
	)

	BeforeEach(func() {
		fakeInstanceRepository = &actorsfakes.FakeInstanceRepository{}

		instancesActor = &InstancesActor{
			Repository: fakeInstanceRepository,
		}
	})

	Describe("Create", func() {
		var err error

		JustBeforeEach(func() {
			err = instancesActor.Create(&osbapi.Instance{
				Name: "instance-1",
			})
		})

		It("create the instance", func() {
			Expect(fakeInstanceRepository.CreateArgsForCall(0)).To(Equal(&osbapi.Instance{
				Name: "instance-1",
			}))
		})

		When("creating the instance fails", func() {
			BeforeEach(func() {
				fakeInstanceRepository.CreateReturns(errors.New("error-creating-instance"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-creating-instance"))
			})
		})
	})
})
