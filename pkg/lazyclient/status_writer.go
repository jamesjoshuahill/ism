package lazyclient

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type statusWriter struct {
	err error
}

func (l statusWriter) Update(ctx context.Context, obj runtime.Object) error {
	return l.err
}
