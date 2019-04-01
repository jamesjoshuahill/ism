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

var _ = Describe("Binding Actor", func() {

	var (
		fakeBindingRepository *actorsfakes.FakeBindingRepository
		bindingsActor         *BindingsActor
	)

	BeforeEach(func() {
		fakeBindingRepository = &actorsfakes.FakeBindingRepository{}

		bindingsActor = &BindingsActor{
			Repository: fakeBindingRepository,
		}
	})

	Describe("GetBindings", func() {
		var (
			bindings []*osbapi.Binding
			err      error
		)

		BeforeEach(func() {
			fakeBindingRepository.FindAllReturns([]*osbapi.Binding{
				{Name: "binding-1"},
				{Name: "binding-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			bindings, err = bindingsActor.GetBindings()
		})

		It("finds all bindings from the repository", func() {
			Expect(fakeBindingRepository.FindAllCallCount()).NotTo(BeZero())
			Expect(bindings).To(ConsistOf(
				&osbapi.Binding{Name: "binding-1"},
				&osbapi.Binding{Name: "binding-2"},
			))
		})

		When("finding all bindings returns an error", func() {
			BeforeEach(func() {
				fakeBindingRepository.FindAllReturns([]*osbapi.Binding{}, errors.New("error-finding-bindings"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-bindings"))
			})
		})
	})

	Describe("GetBindingByName", func() {
		var (
			binding *osbapi.Binding
			err     error
		)

		BeforeEach(func() {
			fakeBindingRepository.FindByNameReturns(&osbapi.Binding{Name: "my-binding"}, nil)
		})

		JustBeforeEach(func() {
			binding, err = bindingsActor.GetBindingByName("my-binding")
		})

		It("finds the binding from the repository", func() {
			Expect(fakeBindingRepository.FindByNameCallCount()).To(Equal(1))
			Expect(*binding).To(Equal(osbapi.Binding{Name: "my-binding"}))
		})

		When("find a binding returns an error", func() {
			BeforeEach(func() {
				fakeBindingRepository.FindByNameReturns(&osbapi.Binding{}, errors.New("error-finding-binding"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-binding"))
			})
		})
	})

	Describe("Create", func() {
		var err error

		JustBeforeEach(func() {
			err = bindingsActor.Create(&osbapi.Binding{
				Name: "binding-1",
			})
		})

		It("create the binding", func() {
			Expect(fakeBindingRepository.CreateArgsForCall(0)).To(Equal(&osbapi.Binding{
				Name: "binding-1",
			}))
		})

		When("creating the binding fails", func() {
			BeforeEach(func() {
				fakeBindingRepository.CreateReturns(errors.New("error-creating-binding"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-creating-binding"))
			})
		})
	})

	Describe("Delete", func() {
		var err error

		JustBeforeEach(func() {
			err = bindingsActor.Delete("binding-1")
		})

		It("deletes the binding", func() {
			Expect(fakeBindingRepository.DeleteArgsForCall(0)).To(Equal("binding-1"))
		})

		When("deleting the binding fails", func() {
			BeforeEach(func() {
				fakeBindingRepository.DeleteReturns(errors.New("error-deleting-binding"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-deleting-binding"))
			})
		})
	})
})
