
[![GitHub release](https://img.shields.io/github/release/QuickBase/quickbase-cli.svg)](https://github.com/QuickBase/quickbase-cli/releases/)
[![Go Reference](https://pkg.go.dev/badge/github.com/QuickBase/quickbase-cli.svg)](https://pkg.go.dev/github.com/QuickBase/quickbase-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/QuickBase/quickbase-cli)](https://goreportcard.com/report/github.com/QuickBase/quickbase-cli)

# Quickbase Command Line Interface

The Quickbase Command Line Interface (CLI) is a tool to manage your Quickbase applications.

## Overview

The Quickbase CLI consumes Quickbase's APIs, so the commands should feel familiar to those versed in the [JSON](https://developer.quickbase.com/) and [XML](https://help.quickbase.com/api-guide/intro.html) API docs. In addition to being an easy way to consume the APIs, this tool is much more than a simple wrapper around them. Using the Quickbase CLI gives you the following benefits:

* Configuration/credential management with profiles for different realms and apps
* Resiliency through [backoff retries](https://en.wikipedia.org/wiki/Exponential_backoff)
* Security through tools that mask sensitive information such as Quickbase tokens
* Scriptability through parsable output, I/O redirection, and [JMESPath filtering](https://jmespath.org/)
* Observability through logging that implements [best practices](https://dev.splunk.com/enterprise/docs/developapps/logging/loggingbestpractices/)
* Debugging through request/response dumping to see what is sent/received over the wire
* Delighters that make it easier to work with data
  * Simple Query syntax (e.g., `--where 3=2`, `--where 2` which both equal `--where {'3'.EX.'2'}`)
  * Natural language processing for various data types, e.g., dates, durations, and addresses (planned).


## Support

The Quickbase CLI is an open source project supported by the community, and it does not fall under the purview of Quickbase Support entitlements. Please use [GitHub issues](https://docs.github.com/en/github/managing-your-work-on-github/about-issues) for bugs, questions, and enhancement requests.

## Installation

### Mac OSX

We recommend using [Homebrew](https://brew.sh/) to install the Quickbase CLI.

```
brew tap quickbase/tap
brew install quickbase-cli
```

Verify the installation:

```
quickbase-cli version
```

### All Platforms

Download and extract the latest release for your platform from the [Releases](https://github.com/QuickBase/quickbase-cli/releases) section. Copy the binary to a directory that is in your system's [PATH](https://en.wikipedia.org/wiki/PATH_(variable)).

### Build From Source

With a [correctly configured](https://golang.org/doc/install#install) Go toolchain, clone the repository to a directory of your choosing, change into it, and run `make`:

```
git clone https://github.com/QuickBase/quickbase-cli.git
cd ./quickbase-cli
make
```

Run `make` in favor of `go build` because the version is set through linker flags, which the default `make build` target does automatically.

## Configuration

Configuration is read from command-line options, environment variables, and a configuration file in that order of precedence. You are advised to set up a configuration file using the command below, which will prompt for your realm hostname, user token, and an optional application ID.

```
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

```
quickbase-cli config dump
```

```json
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

Commands follow the traditional `APP COMMANDS ARGS --FLAGS` pattern. See the exmaple command below that gets an app definition:

```
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

Example command that queries for records, returning fields 6 throguh 8 where Record #ID equals 2:

```
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

#### Field Ranges

In the examples above, `--select 6:8` is equivalent to `--select 6,7,8`. You can also combine the explicit fields and ranges, where `--select 1,3:5` is equal to `--select 1,3,4,5`.

#### Simplified Query Filters

You can also use simplified query syntax for basic queries. The following command queries for records where field 6 equals "Record One" and field 7 equals 2:

```
quickbase-cli records query --select 6:8 --from bqgruir7z --where '6="Record Two" 7=2'
```

Just passing a number will find a record by its ID:

```
quickbase-cli records query --select 6:8 --from bqgruir7z --where 2
```

#### Record Output Formatting

Passing `--format table` for commands that return records will render the output as a table instead of JSON.

```
+------------+--------+---------+
| TITLE      | NUMBER | LIST    |
+------------+--------+---------+
| Record Two | 2      | One,Two |
+------------+--------+---------+
```

Other valid options for `--format` are `csv`, `markdown`.

### Creating Records

Example command that creates a record where field 6 equals "Another Record" and field 7 equals 3:

```
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

### Importing / Exporting Records

Example commands that export data from one table and import it into another that has a similar structure:

```
quickbase-cli table export bq67er5pj | quickbase-cli table import bq72kz6p8
```

Use the import command's `--map` option to reconcile field label differences between the tables. The import/export commands batch the reads and writes by default. Set the `--batch-size` option to control the number of records in each batch. You can also set the `--delay` option to pause between batches, which can help when processing large amounts of data in an active app.

### Deleting Records

Example commmand that deletes the record created above:

```
quickbase-cli records delete --from bqgruir7z --where '6="Another Record"'
```

```json
{
    "numberDeleted": 1
}
```

### Creating Relationships

Example commmand that creates a relationship:

```
quickbase-cli relationship create --child-table-id bqgruir7z --parent-table-id bq6qbvfbv --lookup-field-ids 6,7
```

```json
{
    "childTableId": "bqgruir7z",
    "foreignKeyField": {
        "id": 9,
        "label": "Related Record",
        "type": "numeric"
    },
    "lookupFields": [
        {
            "id": 6,
            "label": "parent - text field",
            "type": "text"
        },
        {
            "id": 7,
            "label": "parent - numeric field",
            "type": "numeric"
        }
    ],
    "id": 9,
    "parentTableId": "bq6qbvfbv"
}
```

### Running Formulas

Example command that runs a formula:

```
quickbase-cli formula run bck7gp3q2 1 --formula "Sum([NumericField],20)"
```

Formulas can span multiple lines and get pretty large. In this instance, you can pass the formula via `STDIN` to this command. The following example assumes the `formula.qb` file contains a large formula:

```
cat ./formula.qb | quickbase-cli formula run bck7gp3q2 1
```

### Transforming Output

[JMESPath](https://jmespath.org/) is a powerful query language for JSON. You can apply JMESPath filters to transform the output of commands to make the data easier to work with. For example, let say you want to get only a list of table names in an app sorted alphabetically. To accomplish this, you can apply a JMESPath filter using the `--filter` option to the command below:

```
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

Valid log levels are `debug`, `info`, `notice`, `error`, `fatal`, and `none`. The default value is `none`.

#### -f, --log-file

Pass `--log-file ./qb.log` to write logs to the `./qb.log` file instead of STDERR.

#### -d, --dump-dir

Pass `--dump-dir ./dump` to write the requests and responses sent over the wire as text files in the directory. The filenames are prefixed with the timestamp and contain the transaction id that can be found in the `transid` context in log messages. All tokens are maked for security.

## Other Resources

The [./jq](https://stedolan.github.io/jq/) tool compliments the Quickbase CLI nicely and makes it easier to work with the output.
