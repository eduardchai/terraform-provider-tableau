```
NOTES: This provider has been archived and moved to https://github.com/traveloka/terraform-provider-tableau or https://registry.terraform.io/providers/traveloka/tableau/latest/docs
```

# Terraform Provider Tableau
Terraform provider for Tableau Cloud resources
The Tableau Cloud provider provides resources to interact with Tableau Rest API using Tableau Personal Access Token (PAT).

### Known Limitations

Tableau PAT can only generate 1 authentication token at a time. This technically will limit the capability of the provider to be run in parallel or concurrently since one process will invalidate the auth token that is used in the other process. Different PAT can be used for the other process if running plan/apply in parallel is required.

## Usage

The provider is published to the Terraform registry and can be used in the same way as any other provider. For detailed documentation with usage examples [view the generated docs in the Terraform registry](https://registry.terraform.io/providers/traveloka/tableau/latest/docs).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 1.x.x
-	[Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies   

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
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

```sh
$ export TABLEAU_SERVER_URL="xxx"
$ export TABLEAU_API_VERSION="x.x"
$ export TABLEAU_SITE="xxx"
$ export TABLEAU_PAT_NAME="xxx"
$ export TABLEAU_PAT_SECRET="xxx"
$ TF_ACC=1 go test -count=1 -v -cover ./internal/provider/
```
