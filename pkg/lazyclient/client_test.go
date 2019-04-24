package lazyclient_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/ism/pkg/lazyclient"
	"github.com/pivotal-cf/ism/pkg/lazyclient/fakes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate counterfeiter -o fakes/fake_client.go ../../vendor/sigs.k8s.io/controller-runtime/pkg/client/interfaces.go Client
//go:generate counterfeiter -o fakes/fake_status_writer.go ../../vendor/sigs.k8s.io/controller-runtime/pkg/client/interfaces.go StatusWriter
//go:generate counterfeiter -o fakes/fake_object.go ../../vendor/k8s.io/apimachinery/pkg/runtime/interfaces.go Object

var _ = Describe("Lazy Client", func() {
	var (
		fakeClient *fakes.FakeClient

		factory func() (client.Client, error)
		ctx     context.Context
		obj     *fakes.FakeObject
	)

	BeforeEach(func() {
		fakeClient = new(fakes.FakeClient)

		factory = func() (client.Client, error) {
			return fakeClient, nil
		}
		ctx = context.Background()
		obj = new(fakes.FakeObject)
	})

	var ItIsLazy = func(action func(c client.Client) error) {
		It("builds the client on the first action", func() {
			factoryCalled := 0
			c := &lazyclient.LazyClient{Factory: func() (client.Client, error) {
				factoryCalled++
				return fakeClient, nil
			}}
			Expect(factoryCalled).To(BeZero())

			_ = action(c)
			_ = action(c)
			Expect(factoryCalled).To(Equal(1))
		})

		It("returns the client factory error", func() {
			c := &lazyclient.LazyClient{Factory: func() (client.Client, error) {
				return nil, errors.New("factory error")
			}}

			err := action(c)

			Expect(err).To(MatchError("creating kubernetes api client: factory error"))
		})
	}

	When("Get is called", func() {
		ItIsLazy(func(c client.Client) error {
			return c.Get(ctx, client.ObjectKey{}, obj)
		})

		It("passes the args to the client", func() {
			c := &lazyclient.LazyClient{Factory: factory}
			key := client.ObjectKey{}

			_ = c.Get(ctx, key, obj)

			Expect(fakeClient.GetCallCount()).To(Equal(1))
			actualCtx, actualKey, actualObj := fakeClient.GetArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualKey).To(Equal(key))
			Expect(actualObj).To(Equal(obj))
		})
	})

	When("List is called", func() {
		ItIsLazy(func(c client.Client) error {
			return c.List(ctx, &client.ListOptions{}, obj)
		})

		It("passes the args to the client", func() {
			c := &lazyclient.LazyClient{Factory: factory}
			opts := &client.ListOptions{}

			_ = c.List(ctx, opts, obj)

			Expect(fakeClient.ListCallCount()).To(Equal(1))
			actualCtx, actualOpts, actualObj := fakeClient.ListArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualOpts).To(Equal(opts))
			Expect(actualObj).To(Equal(obj))
		})
	})

	When("Create is called", func() {
		ItIsLazy(func(c client.Client) error {
			return c.Create(ctx, obj)
		})

		It("passes the args to the client", func() {
			c := &lazyclient.LazyClient{Factory: factory}
			opts := &client.ListOptions{}

			_ = c.List(ctx, opts, obj)

			Expect(fakeClient.ListCallCount()).To(Equal(1))
			actualCtx, actualOpts, actualObj := fakeClient.ListArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualOpts).To(Equal(opts))
			Expect(actualObj).To(Equal(obj))
		})
	})

	When("Delete is called", func() {
		ItIsLazy(func(c client.Client) error {
			var opts []client.DeleteOptionFunc
			return c.Delete(ctx, obj, opts...)
		})

		It("passes the args to the client", func() {
			c := &lazyclient.LazyClient{Factory: factory}
			var opts []client.DeleteOptionFunc

			_ = c.Delete(ctx, obj, opts...)

			Expect(fakeClient.DeleteCallCount()).To(Equal(1))
			actualCtx, actualObj, actualOpts := fakeClient.DeleteArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualObj).To(Equal(obj))
			Expect(actualOpts).To(Equal(opts))
		})
	})

	When("Update is called", func() {
		ItIsLazy(func(c client.Client) error {
			return c.Update(ctx, obj)
		})

		It("passes the args to the client", func() {
			c := &lazyclient.LazyClient{Factory: factory}

			_ = c.Update(ctx, obj)

			Expect(fakeClient.UpdateCallCount()).To(Equal(1))
			actualCtx, actualObj := fakeClient.UpdateArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualObj).To(Equal(obj))
		})
	})

	When("Status().Update() is called", func() {
		BeforeEach(func() {
			fakeStatusWriter := new(fakes.FakeStatusWriter)
			fakeClient.StatusReturns(fakeStatusWriter)
		})

		ItIsLazy(func(c client.Client) error {
			statusWriter := c.Status()
			return statusWriter.Update(ctx, obj)
		})
	})
})
