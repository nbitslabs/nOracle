# nOracle

[![Go Report Card](https://goreportcard.com/badge/github.com/nbitslabs/nOracle)](https://goreportcard.com/report/github.com/nbitslabs/nOracle)

## 📦 Requirements

- [Go](https://golang.org/doc/install)
- [Bazelisk](https://github.com/bazelbuild/bazelisk)

## 🚀 Getting Started

### Install Go dependencies

Use Bazel to tidy the Go module and fetch external dependencies.

```bash
bazel mod tidy
```

### Build all targets

This ensures all Bazel targets in the repo are building correctly.

```bash
bazel build //...
```

### Generate or update Bazel build files with Gazelle

Gazelle scans your Go code and automatically generates or updates `BUILD.bazel` files.

```bash
bazel run //:gazelle
```

> ✅ Re-run Gazelle whenever:
>
> - You add or remove `.go` files
> - You rename Go packages or files
> - You change import paths

### Run tests

This ensures all Bazel targets in the repo are tested correctly.

```bash
bazel test //...
```
