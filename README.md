# manrs-observatory-go

[![Build](https://github.com/manrs-tools/manrs-observatory-go/actions/workflows/build.yml/badge.svg)](https://github.com/manrs-tools/manrs-observatory-go/actions/workflows/build.yml)
[![Lint](https://github.com/manrs-tools/manrs-observatory-go/actions/workflows/lint.yml/badge.svg)](https://github.com/manrs-tools/manrs-observatory-go/actions/workflows/lint.yml)
[![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8?logo=go)](go.mod)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

Go SDK and CLI for the [MANRS Observatory](https://observatory.manrs.org) API, the platform that tracks [MANRS](https://www.manrs.org) (Mutually Agreed Norms for Routing Security) participant conformance, ASN scoring, RPKI/ROA validation, and routing security incidents.

This is the public SDK: it covers every publicly available Observatory endpoint (autonomous systems, participants, ROAs, stats, tickets, and more). It's used the same way whether you're building tooling around your own participant data or just querying public Observatory data.

## Requirements

- An Observatory API key ([`client login`](#authentication), or generate one from the Observatory dashboard).
- Go 1.25 or later (only for the SDK, or when building the CLI from source).

## Installation

```bash
go get github.com/manrs-tools/manrs-observatory-go
```

## Quick start

```go
package main

import (
	"context"
	"log"

	manrs "github.com/manrs-tools/manrs-observatory-go"
	"github.com/manrs-tools/manrs-observatory-go/model"
)

func main() {
	client := manrs.New("your-api-key")

	participants, err := client.FindParticipants(context.Background(), &model.ParticipantFilter{
		Status: "approved",
		Limit:  10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range participants {
		log.Printf("%s (%s): ASNs=%v", p.Name, p.ID, p.ASNs)
	}
}
```

## SDK guide

The SDK is the root `client` package (module `github.com/manrs-tools/manrs-observatory-go`), conventionally imported as `manrs`. Request/response types live in the `model` subpackage.

### Client options

`manrs.New` accepts functional options:

| Option                             | Description                                                                                         |
| ---------------------------------- | --------------------------------------------------------------------------------------------------- |
| `manrs.WithClient(*http.Client)`   | Swap in a custom underlying HTTP client (proxying, custom transport, etc).                          |
| `manrs.WithDebug(bool)`            | Print every outgoing request as an equivalent `curl` command, plus response status/body, to stderr. |
| `manrs.WithTimeout(time.Duration)` | Override the default 30s request timeout.                                                           |

`Client.SetAPIKey` and `Client.SetTimeout` allow adjusting a client after it's been constructed (used by the CLI's `login` flow, for example). Requests are retried automatically via [`hashicorp/go-retryablehttp`](https://github.com/hashicorp/go-retryablehttp), and gzip-encoded responses are decoded transparently.

### Filtering

Every `FindX`/`FindXPaged` method takes a `*model.XFilter` struct. Fields are translated into query parameters automatically based on `query` struct tags: slices are comma-joined, zero values are omitted where the tag specifies `omitempty`/`omitzero`, and `time.Time` fields are encoded as RFC 3339. Passing `nil` sends no query parameters.

```go
ases, err := client.FindASes(ctx, &model.ASFilter{Economies: []string{"US", "CA"}})
```

### Pagination

Paginated list calls (`FindXPaged`) return a generic `*Iterator[T]`:

```go
it := client.FindParticipantsPaged(ctx, &model.ParticipantFilter{Status: "approved"})
for it.Next() {
	participant := it.Value()
	log.Println(participant.Name)
}
if err := it.Err(); err != nil {
	log.Fatal(err)
}
```

`Iterator` fetches pages lazily as `Next()` walks past the end of the current page, using the API's `nextPageToken` cursor.

### Error handling

Non-2xx responses are parsed into the API's problem-detail format and surfaced as an `error` whose message already includes the HTTP status and any field-level error messages returned by the server, so most callers need nothing further:

```go
if _, err := client.FindParticipant(ctx, "does-not-exist"); err != nil {
	log.Println(err) // e.g. "find participant: Participant not found (HTTP 400)"
}
```

The error is a `client.ResponseError` carrying the full problem-detail payload, so a caller that needs the structured detail behind the message — which field failed validation, what value was rejected — recovers it with `errors.As`:

```go
var apiErr client.ResponseError
if errors.As(err, &apiErr) {
	for _, detail := range apiErr.Errors {
		log.Printf("%s: %s (got %v)", detail.Location, detail.Message, detail.Value)
	}
}
```

### Resource coverage

**Core resources** (available to any API key):

| File                   | Covers                                                  |
| ---------------------- | ------------------------------------------------------- |
| `as.go`                | Autonomous systems, ASN scores, ASN incidents           |
| `api_key.go`           | API key management                                      |
| `auth.go`              | Login, registration, password reset, email verification |
| `email.go`             | Sent email lookup                                       |
| `favorite.go`          | Favorited participants                                  |
| `incident_feedback.go` | Incident feedback and summaries                         |
| `invite.go`            | Invite acceptance                                       |
| `participant.go`       | Participant lookup/search                               |
| `roa.go`               | ROA/RPKI validation records and stats                   |
| `stats.go`             | Conformance, summary, history, and overview stats       |
| `testimonial.go`       | Testimonials                                            |
| `ticket.go`            | Support tickets                                         |
| `user.go`              | Session user                                            |
| `version.go`           | API version                                             |

**Administrative resources** (`*_admin.go`, require an admin-role API key):

`application_admin.go`, `as_admin.go`, `auth_admin.go`, `bogon_admin.go`, `contact_admin.go`, `email_admin.go`, `gap_admin.go`, `grip_admin.go`, `incident_feedback_admin.go`, `invite_admin.go`, `job_admin.go`, `log_admin.go`, `nro_admin.go`, `participant_admin.go`, `request_admin.go`, `ris_announcement_admin.go`, `ro_admin.go`, `rov_admin.go`, `spoofer_admin.go`, `stats_admin.go`, `testimonial_admin.go`, `ticket_admin.go`, `user_admin.go`.

## CLI guide

The `client` binary (`cmd/client`) wraps the entire SDK in a [Cobra](https://github.com/spf13/cobra)-based CLI.

### Installing the CLI

Prebuilt binaries for macOS, Linux, FreeBSD, and Windows are attached to every [GitHub release](https://github.com/manrs-tools/manrs-observatory-go/releases). Download the zip for your platform, extract it, and place the binary on your `PATH`:

```bash
unzip manrs-observatory-darwin-arm64.zip
mv manrs-observatory-darwin-arm64 /usr/local/bin/client
```

Alternatively, build from source with Go:

```bash
go install github.com/manrs-tools/manrs-observatory-go/cmd/client@latest
# or, from a checkout of this repo:
make dev            # builds ./bin/client
```

### Authentication

```bash
# Using an existing API key
client login 2fbf890c-aa01-4f5b-907f-a9c2ba634ec2

# Using an email/password
client login --email you@example.com --password '...'
```

Your API key, email, and role are cached in `~/.config/manrs-observatory-client/config.yml` (your password is never persisted). Inspect or edit that cache directly with `client config get|set|show`.

### Global flags

These apply to every subcommand:

| Flag                | Default | Description                                                                   |
| ------------------- | ------- | ----------------------------------------------------------------------------- |
| `-d, --debug`       | `false` | Print outgoing requests/responses (equivalent curl command + body) to stderr. |
| `-f, --format`      | `yaml`  | Output format: `json`, `jsonp` (pretty JSON), or `yaml`.                      |
| `-l, --limit`       | `0`     | Limit the number of results returned by `find`-style commands.                |
| `--logLevel`        | `info`  | Log level: `debug`, `info`, `warn`, `error`, `fatal`, `panic`.                |
| `--prettyLog`       | `true`  | Pretty-print logs to the console.                                             |
| `-t, --timeout`     | `60s`   | HTTP request timeout.                                                         |
| `-w, --writetofile` | `false` | Write command output to a timestamped file instead of stdout.                 |

### Command reference

Every command group supports `--help` for its full set of subcommands and flags (e.g. `client participants find --help`). Top-level groups:

| Command                      | Description                                       |
| ---------------------------- | ------------------------------------------------- |
| `apiKeys`                    | Manage your API keys                              |
| `as`                         | Autonomous systems, ASN scores, ASN incidents     |
| `config`                     | Get/set/show local client configuration           |
| `docs`                       | Open the OpenAPI docs in your browser             |
| `emails`                     | Sent email lookup                                 |
| `favorites`                  | Favorited participants                            |
| `feedback`                   | Incident feedback and summaries                   |
| `login`                      | Authenticate and cache credentials                |
| `participants`               | Participant search and lookup                     |
| `roa`                        | ROA/RPKI validation records and stats             |
| `stats`                      | Conformance, summary, history, and overview stats |
| `testimonials`               | Browse testimonials                               |
| `tickets` / `ticket-replies` | Support tickets and replies                       |
| `user` / `users`             | Session user / user lookup                        |
| `version`                    | Print API version information                     |

A couple of concrete examples:

```bash
client as scores 15169 --format json
client participants find --status approved --region "North America" --limit 20
client roa find --asns 15169
```

## Running with Docker

```bash
docker build -t ghcr.io/manrs-tools/manrs-observatory-go:latest .
docker run --rm ghcr.io/manrs-tools/manrs-observatory-go:latest participants find --limit 5
```

Tagged images are published to `ghcr.io/manrs-tools/manrs-observatory-go` on release (see `.github/workflows/release.yml`).

## Development

```bash
make setup   # install toolchain dependencies (betteralign, gofumpt, nilaway, golangci-lint)
make dev     # build ./bin/client
make lint    # run golangci-lint
make test    # run tests
make format  # run gofumpt
make nil     # run nilaway
```

`golangci-lint` is configured (`.golangci.yml`) with a strict, opinionated rule set; CI runs it with `--fix` on every pull request and commits any resulting formatting fixes automatically.

## License

[Apache-2.0](LICENSE) © Global Cyber Alliance
