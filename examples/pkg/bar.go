package pkg

import (
	"github.com/takashabe/ctxfmt/examples"
)

type implInterface struct{}

func NewInterface() examples.Interface {
	return &implInterface{}
}

func (i *implInterface) Foo() {
}

// infucient context.Context argument
func (i *implInterface) FooCtx() {
}

// infucient context.Context argument
var _ = examples.BarCtx()
