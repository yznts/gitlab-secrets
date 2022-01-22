
# GitLab Secrets

Tool to manipulate, pull and push GitLab secrets.  

## Installation

```bash
$ go install github.com/yuriizinets/gitlab-secrets@latest
```

## Usage

> Please note, all commands (except help messages) need to be run inside of a project directory with a GitLab repository upstream.

Help message

```bash
$ gitlab-secrets
Commands:
  auth     - Authenticate with a private token
  kv:list  - List keys or key/value pairs
  kv:get   - Get a value with a key
  kv:set   - Set a key/value pair
  kv:del   - Delete a key/value pair
  pull     - Pull secrets from GitLab to specified location
  push     - Push secrets from specified location to GitLab
Usage:
  gitlab-secrets auth -token <token>
  gitlab-secrets kv:get -key <key>
  gitlab-secrets kv:set -key <key> -value <value>
  gitlab-secrets kv:del -key <key>
  gitlab-secrets pull -file <filepath>
  gitlab-secrets push -file <filepath>
```

Set repository key

```bash
$ gitlab-secrets kv:set -key "TESTVAR" -value "test"
Variable was set.
```

Get repository key

```bash
$ gitlab-secrets kv:get -key "TESTVAR"
TESTVAR=test
```

List repository keys

```bash
$ gitlab-secrets kv:list
TESTVAR=test
```

Pull secrets from GitLab to specified file (`.env` by default)

```bash
$ gitlab-secrets pull -file .secrets
Variables pulled.
```
