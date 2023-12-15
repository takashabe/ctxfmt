## Overview
`ctxfmt` is a tool designed to automate the insertion of context.Context in Go method declarations and calls. It's especially useful for existing codebases, allowing for seamless integration of context handling.

## Features

- Method Definition Completion: Automatically adds ctx context.Context to method declarations in interfaces and existing method definitions.
- Method Call Completion: Inserts context.TODO() in method calls where arguments are insufficient.
- Dry-Run Mode: Acts as a linter for code not yet using context, allowing you to preview changes without applying them.

## Installation

```bash
go intall github.com/takashabe/ctxfmt@latest
```

## Usage

### Method Definition Completion

```bash
$ ctxfmt signature $GOPATH/src/github.com/takashabe/ctxfmt/examples/**
```

This command modifies the code as follows:
- ctx context.Context is added to method definitions
- The import "context" statement is included automatically

```diff
@@ -1,12 +1,14 @@
 package main

+import "context"
+
 type Interface interface {
-       Foo(id int)
+       Foo(ctx context.Context, id int)
 }

 type impl struct{}

-func (i *impl) Foo(id int) {}
+func (i *impl) Foo(ctx context.Context, id int) {}

 func main() {
        i := &impl{}
```

### Method Call Completion

```bash
$ ctxfmt args --pkg 'github.com/takashabe/ctxfmt/examples' $GOPATH/src/github.com/takashabe/ctxfmt/examples
processed /Users/takashabe/dev/src/github.com/takashabe/ctxfmt/examples/main.go
```

This command modifies the code as follows:
- add context.TODO() to method calls that require a context.Context argument

```diff
@@ -12,5 +12,5 @@ func (i *impl) Foo(ctx context.Context, id int) {}

 func main() {
        i := &impl{}
-       i.Foo(1)
+       i.Foo(context.TODO(), 1)
 }
```

#### :warning: Prerequisites

Ensure the file has a compilation error like below, indicating missing context.Context in method calls:

```bash

$ go vet .
# github.com/takashabe/ctxfmt/examples
vet: ./main.go:15:9: not enough arguments in call to i.Foo
        have (number)
        want (context.Context, int)
```
