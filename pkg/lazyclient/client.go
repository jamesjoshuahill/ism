package lazyclient

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Factory func() (client.Client, error)

type LazyClient struct {
	Factory Factory
	client  client.Client
}

func (l *LazyClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	err := l.buildClient()
	if err != nil {
		return err
	}

	return l.client.Get(ctx, key, obj)
}

func (l *LazyClient) List(ctx context.Context, opts *client.ListOptions, list runtime.Object) error {
	err := l.buildClient()
	if err != nil {
		return err
	}

	return l.client.List(ctx, opts, list)
}

func (l *LazyClient) Create(ctx context.Context, obj runtime.Object) error {
	err := l.buildClient()
	if err != nil {
		return err
	}

	return l.client.Create(ctx, obj)
}

func (l *LazyClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOptionFunc) error {
	err := l.buildClient()
	if err != nil {
		return err
	}

	return l.client.Delete(ctx, obj, opts...)
}

func (l *LazyClient) Update(ctx context.Context, obj runtime.Object) error {
	err := l.buildClient()
	if err != nil {
		return err
	}

	return l.client.Update(ctx, obj)
}

func (l *LazyClient) Status() client.StatusWriter {
	err := l.buildClient()
	if err != nil {
		return statusWriter{err: err}
	}

	return l.client.Status()
}

func (l *LazyClient) buildClient() error {
	if l.client == nil {
		c, err := l.Factory()
		if err != nil {
			return fmt.Errorf("creating kubernetes api client: %s", err)
		}

		l.client = c
	}

	return nil
}
