# recv-context

"recv-context" is a CLI tool designed for listing Go language methods that do not include context.Context in their arguments. This tool is particularly useful for developers looking to ensure context propagation in their Go applications.

## Usage

```
$ go run . $GOPATH/src/github.com/takashabe/btcli/pkg/**
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/bigtable/bigtable_mock.go at line 80: mr.Get()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/bigtable/bigtable_mock.go at line 100: mr.GetRows()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/bigtable/bigtable_mock.go at line 116: mr.Count()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/bigtable/bigtable_mock.go at line 131: mr.Tables()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/cmd/interactive/completer.go at line 17: c.Do()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/cmd/interactive/executor.go at line 33: e.Do()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/cmd/interactive/interactive.go at line 37: c.Run()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/printer/printer.go at line 28: w.PrintRows()
/Users/takashabe/dev/src/github.com/takashabe/btcli/pkg/printer/printer.go at line 35: w.PrintRow()
```
