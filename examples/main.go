package main

type Interface interface {
	Foo(id int)
}

type impl struct{}

func (i *impl) Foo(id int) {}

func main() {
	i := &impl{}
	i.Foo(1)
}
