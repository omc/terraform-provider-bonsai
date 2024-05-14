# Contributing

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Pre-commit

This project uses [pre-commit](https://pre-commit.com/) to lint and store 3rd-party dependency licenses.
Installation instructions are available on the [pre-commit](https://pre-commit.com/) website!

To verify your installation, run this project's pre-commit hooks against all files:

```shell
pre-commit run --all-files
```

#### Hook dependencies

- [Terraform](https://developer.hashicorp.com/terraform/install)
- [Terraform-docs](https://terraform-docs.io/user-guide/installation/)
