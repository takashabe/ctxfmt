package examples

import "context"

type Interface interface {
	Foo()
	FooCtx(ctx context.Context)
}

func BarCtx(ctx context.Context) error {
	return nil
}
