## Overview

ctxfmt is an linter tool for Go developers, designed to enhance the readability and maintainability of Go codebases by ensuring the consistent use of `context.Context` in interface method definitions and implementations. This tool automatically adds `context.Context` as the first parameter in methods where it's missing.

## Features

- **Automatic Injection**: ContextLint automatically injects `context.Context` as the first parameter in interface methods and their implementations if it's missing.
- **Dry Run Mode**: Offers a dry run option to report where changes would be made without actually modifying the code.

## Installation

```bash
go get -u github.com/takashabe/ctxfmt
```

## Usage

ContextLint can be run in two modes: normal and dry run.

- **Automatic Injection**: This will modify your Go files directly, adding `context.Context` where necessary.

  ```bash
  contextlint ./...
  ```

- **Dry Run Mode**: Use this to see what changes would be made without applying them.

  ```bash
  contextlint --dry-run ./...
  ```

## Configuration

TODO
