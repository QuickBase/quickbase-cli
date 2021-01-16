
[![Go Reference](https://pkg.go.dev/badge/github.com/QuickBase/quickbase-cli.svg)](https://pkg.go.dev/github.com/QuickBase/quickbase-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/QuickBase/quickbase-cli)](https://goreportcard.com/report/github.com/QuickBase/quickbase-cli)

# Quickbase Command Line Interface

The Quickbase Command Line Interface (CLI) is a tool to manage your Quickbase applications.

## Overview

The Quickbase CLI consumes Quickbase's RESTful API, so the commands should feel familiar to those versed in the [API docs](https://developer.quickbase.com/). In addition to being an easy way to consume the RESTful API, this tool is much more than a simple wrapper around it. Using the Quickbase CLI gives you the following benefits:

* Configuration/credential management with profiles for different realms and apps
* Resiliency through [backoff retries](https://en.wikipedia.org/wiki/Exponential_backoff)
* Security through tools that mask sensitive information such as Quickbase tokens
* Scriptability through parsable output, I/O redirection, and [JMESPath filtering](https://jmespath.org/)
* Observability through logging that implements [best practices](https://dev.splunk.com/enterprise/docs/developapps/logging/loggingbestpractices/)
* Debugging through request/response dumping to see what is sent/received over the wire
* Delighters that make it easier to work with data
  * Simple Query syntax (e.g., `--where 3=2`, `--where 2` which both equal `--where {'3'.EX.'2'}`)
  * Natural language processing for various data types, e.g., dates, durations, and addresses (planned).

## Installation

Download the latest binary for your platform from the [Releases](https://github.com/QuickBase/quiickbase-cli/releases) section to a directory that is in your system's [PATH](https://en.wikipedia.org/wiki/PATH_(variable)).

### Building From Source

With a [correctly configured](https://golang.org/doc/install#install) Go toolchain:

```sh
go get github.com/QuickBase/quickbase-cli
```

## Configuration

Configuration is read from command-line options, environment variables, and a configuration file in that order of precedence. You are advised to set up a configuration file using the command below, which will prompt for your realm hostname, user token, and an optional application ID.

```sh
quickbase-cli config setup
```

The configuration is written to a file named `.config/quickbase/config.yml` under your home directory, which you can edit to add additional configuration sets called **profiles**:

```yml
default:
  realm_hostname: example1.quickbase.com
  user_token: b3b6se_mzif_dy36********************hi7b
  app_id: bqgruir3g

another_realm:
  realm_hostname: example2.quickbase.com
  user_token: b3b6se_uyp_iybv********************js2k
```

The `default` profile is used unless the `QUICKBASE_PROFILE` environment variable or `--profile` command line option specify another value, such as `another_realm`.

Run the following command to dump the configuration values for the active profile:

```sh
quickbase-cli config dump
{
    "realm_hostname": "example1.quickbase.com",
    "user_token": "b3b6se_mzif_dy36********************hi7b",
    "app_id": "bqgruir3g"
}
```

You can also set environment variables for common options, e.g., app IDs, table IDs, and field IDs. This makes it easy to chain together a string of commands that act on the same resource:

```sh
export QUICKBASE_TABLE_ID=bqgruir7z

# The commands below use "bqgruir7z" for the --to and --from options.
quickbase-cli records insert --data '6="Another Record" 7=3'
quickbase-cli records query --select 6 --where '6="Another Record"'
quickbase-cli records delete --where '6="Another Record"'
```

## Usage

### Command Format

Exmaple command that gets an app definition:

```sh
quickbase-cli app get --app-id bqgruir3g
```

```json
{
    "id": "bqgruir3g",
    "name": "New API Test",
    "timeZone": "(UTC-08:00) Pacific Time (US \u0026 Canada)",
    "dateFormat": "MM-DD-YYYY",
    "created": "2020-11-03T19:33:01Z",
    "updated": "2020-11-03T19:33:01Z",
    "variables": [
        {
            "name": "var1",
            "value": "Test variable value"
        }
    ]
}
```

### Querying For Records

Example command that queries for records (where Record #ID is 2):

```sh
quickbase-cli records query --select 6,7,8 --from bqgruir7z --where '{3.EX.2}'
```

```json
{
    "data": [
        {
            "6": {
                "value": "Record Two"
            },
            "7": {
                "value": 2
            },
            "8": {
                "value": [
                    "One",
                    "Two"
                ]
            }
        }
    ],
    "fields": [
        {
            "id": 6,
            "label": "Title",
            "type": "text"
        },
        {
            "id": 7,
            "label": "Number",
            "type": "numeric"
        },
        {
            "id": 8,
            "label": "List",
            "type": "multitext"
        }
    ],
    "metadata": {
        "totalRecords": 1,
        "numRecords": 1,
        "numFields": 3,
        "skip": 0
    }
}
```

You can also use simplified query syntax for basic queries. The following command queries for records where field 6 equals "Record One" and field 7 equals 2:

```sh
quickbase-cli records query --select 6,7,8 --from bqgruir7z --where '6="Record Two" 7=2'
```

Just passing a number will find a record by its ID:

```sh
quickbase-cli records query --select 6,7,8 --from bqgruir7z --where 2
```

### Creating Records

Example command that creates a record where field 6 equals "Another Record" and field 7 equals 3:

```sh
quickbase-cli records insert --to bqgruir7z --data '6="Another Record" 7=3'
```

```json
{
    "metadata": {
        "createdRecordIds": [
            7
        ],
        "totalNumberOfRecordsProcessed": 1,
        "unchangedRecordIds": [],
        "updatedRecordIds": []
    }
}
```

### Transforming Output

[JMESPath](https://jmespath.org/) is a powerful query language for JSON. You can apply JMESPath filters to transform the output of commands to make the data easier to work with. For example, let say you want to get only a list of table names in an app sorted alphabetically. To accomplish this, you can apply a JMESPath filter using the `--filter` option to the command below:

```sh
quickbase-cli table list --app-id bqgruir3g --filter "tables[].name | sort(@) | {Tables: join(', ', @)}"
```

```json
{
    "Tables": "Fields, Records, Tasks"
}
```

### Navigation Helpers

The CLI tool has navigation helpers via `open` commands that make it easy to jump to specific pages in the UI. The commands below assume a default application is confgured, which is why the `--app-id` option is omitted, and open your browser when run:

```sh

# Navigate to the app's homepage.
quickbase-cli app open

# Navigate to the app's settings page.
quickbase-cli app open --settings all

# Navigate to the app's "Roles" settings page.
quickbase-cli app open --settings roles

# Valid values for the --settings option are all, branding, management, pages,
# properties, roles, tables, and variables.

# Navigate the the table's home page.
quickbase-cli table open bq4w73asu

# Navigate to the table's settings page.
quickbase-cli table open bq4w73asu --settings all

# Valid values for the --settings option are actions, access, advanced, all,
# forms, fields, notifications, relationships, reports, webhooks.

```

### Global Options

#### -h, --help

Returns help for commands.

#### -q, --quiet

Suppress output written to STDOUT.

#### -l, --log-level

Pass `--log-level debug` to get information useful for debugging. Log messages are written to STDERR, so you can redirect the logs using `2>` without disrupting the normal output.

Valid log levels are `debug`, `info`, `error`, `notice`, `fatal`, and `none`. The default value is `none`.

#### -f, --log-file

Pass `--log-file ./qb.log` to write logs to the `./qb.log` file instead of STDERR.

#### -d, --dump-dir

Pass `--dump-dir ./dump` to write the requests and responses sent over the wire as text files in the directory. The filenames are prefixed with the timestamp and contain the transaction id that can be found in the `transid` context in log messages. All tokens are maked for security.

## Other Resources

The [./jq](https://stedolan.github.io/jq/) tool compliments the Quickbase CLI nicely and makes it easier to work with the output.
