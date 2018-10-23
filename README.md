# Tyk-Git

## What is it?

Tyk-Git is a command line tool and library to manage and synchronise a Tyk installation with a Git repository.

## Features

- Update APIs and policies on remote Tyk Dashboards
- Update APIs on remote Tyk CE Gateways
- Publish APIS/Policies to remote Tyk Dashboards
- Publish APIs to remote Tyk CE Gateways
- Synchronise a Tyk Dashboard's APIs and Policies with those stored in a repository (one-way, Git writes to Dashboard)
- Synchronise a Tyk CE Gateway's APIs with those stored in a repository (one-way, Git writes to the Gateway)
- Dump Policies and APIs in a transportable format from a Dashboard to a directory
- Support for importing, converting and publishing Swagger (Open API Spec) files to Tyk.

### Sync

Tyk-Git tries to be clever about what APIs and Policies to update and which to create, it will actually base all
ID matching on the API ID and the masked Policy ID, so it can identify the same object across installations. Tyk has
a tendency to generate fresh IDs for all new Objects, so Tyk-Git gets around this by using portable IDs and ensuring
the necessary portable IDs are set when using the `dump` command.

This means that Tyk-Git can be used to back-up your most important API Gateway configurations as code, and to deploy
those configurations to any target and ensure that API IDs and Policy IDs will remain consistent, ensuring that any
dependent tokens continue to have access to your services.

### Prerequisites:

- Tyk-Git was built using Go 1.10. The minimum Go version required to install is 1.7.
- In order for policy ID matching to work correctly, your Gateway must have `policies.allow_explicit_policy_id: true`.
- It is assumed you have a Tyk CE or Tyk Pro installation.

## Installation

Currently the application is only available via Go, so to install you must have Go installed and run:

```
go install github.com/TykTechnologies/tyk-git
```

This should make the `tyk-git` command available to your console.

## Usage

```
Usage:
  tyk-git [flags]
  tyk-git [command]

Available Commands:
  dump        Dump will extract policies and APIs from a target (dashboard)
  help        Help about any command
  publish     publish API definitions from a Git repo to a gateway or dashboard
  sync        Synchronise a github repo with a gateway
  update      A brief description of your command

Flags:
  -h, --help   help for tyk-git

Use "tyk-git [command] --help" for more information about a command.
```

## Example: Transfer from one Tyk Dashboard to another

First, we need to extract the data from our Tyk Dashboard, here we `dump` into ./tmp, let's assume this is a git-enabled
directory

```
./tyk-git dump -d="http://localhost:3000" -s="b2d420ca5302442b6f20100f76de7d83" -t="./tmp"
Extracting APIs and Policies from http://localhost:3000
> Fetching policies
--> Identified 1 policies
--> Fetching and cleaning policy objects
> Fetching APIs
--> Fetched 3 APIs
> Creating spec file in: tmp/.tyk.json
Done.
```

Next, let's push those changes back to the Git repo on the branch `my-test-branch`:

```
cd tmp
git add .
git commit -m "My dashboard dump"
git push -u origin my-test-branch
```

Now to restore this data directly from GitHub:

```
./tyk sync -d="http://localhost:3010" -s="b2d420ca5302442b6f20100f76de7d83" -b="refs/heads/my-test-branch" https://github.com/myname/my-test.git
Using publisher: Dashboard Publisher
Fetched 3 definitions
Fetched 1 policies
Processing APIs...
Deleting: 0
Updating: 3
Creating: 0
SYNC Updating: 598ec94f9695f201730d835b
SYNC Updating: 598ec9589695f201730d835c
SYNC Updating: 5990cfee9695f201730d836e
Processing Policies...
Deleting policies: 0
Updating policies: 1
Creating policies: 0
SYNC Updating Policy: Test policy 1
--> Found policy using explicit ID, substituting remote ID for update
```

The command provides output to identify which actions have been taken. If using a Tyk Gateway, the Gateway will be
automatically hot-reloaded.

