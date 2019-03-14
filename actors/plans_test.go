/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms
of the under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package actors_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/actors/actorsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Plans Actor", func() {

	var (
		fakePlanRepository *actorsfakes.FakePlanRepository

		plansActor *PlansActor
	)

	BeforeEach(func() {
		fakePlanRepository = &actorsfakes.FakePlanRepository{}

		plansActor = &PlansActor{
			Repository: fakePlanRepository,
		}
	})

	Describe("GetPlan", func() {
		var (
			plan *osbapi.Plan
			err  error
		)

		BeforeEach(func() {
			fakePlanRepository.FindReturns(&osbapi.Plan{ID: "plan-1", Name: "my-plan-1"}, nil)
		})

		JustBeforeEach(func() {
			plan, err = plansActor.GetPlan("plan-1")
		})

		It("finds the plan by plan id", func() {
			Expect(fakePlanRepository.FindArgsForCall(0)).To(Equal("plan-1"))
			Expect(plan).To(Equal(&osbapi.Plan{ID: "plan-1", Name: "my-plan-1"}))
		})

		When("finding the plan returns an error", func() {
			BeforeEach(func() {
				fakePlanRepository.FindReturns(nil, errors.New("error-finding-plan"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-plan"))
			})
		})
	})

	Describe("GetPlans", func() {
		var (
			plans []*osbapi.Plan
			err   error
		)

		BeforeEach(func() {
			fakePlanRepository.FindByServiceReturns([]*osbapi.Plan{
				{Name: "plan-1"},
				{Name: "plan-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			plans, err = plansActor.GetPlans("service-1")
		})

		It("finds plans by service id", func() {
			Expect(fakePlanRepository.FindByServiceArgsForCall(0)).To(Equal("service-1"))

			Expect(plans).To(Equal([]*osbapi.Plan{
				{Name: "plan-1"},
				{Name: "plan-2"},
			}))
		})

		When("finding plans returns an error", func() {
			BeforeEach(func() {
				fakePlanRepository.FindByServiceReturns([]*osbapi.Plan{}, errors.New("error-finding-plans"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-plans"))
			})
		})
	})
})
