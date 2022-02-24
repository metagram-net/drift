# Drift

Drift manages Postgresql database migrations.

There's no shortage of tools for running migrations, but this one embodies my
particular opinions:

1. Migrations should be written in database-specific SQL.
2. Migrations are not generally idempotent or reversible.
3. Migration dependencies form a tree, not a line.

## Installation

To install Drift as a command, use `go install`:

```bash
go install github.com/metagram-net/drift/cmd/drift
```

Since this tool is still very unstable, consider pinning the version in your
`go.mod` with the [tools pattern].

[tools pattern]: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

To use Drift as a library, use `go get`:

```bash
go get github.com/metagram-net/drift
```

## Usage

Run `drift help` to get usage information from each subcommand.

### First-time setup

Write the configuration file (`drift.toml`) or set the equivalent environment
variables. The environment variables take precedence over the file.

The environment variables are upper snake-case versions of the ones in the file
with `DRIFT_` prefixes. For example, `database-url` is `DRIFT_DATABASE_URL`.

```toml
# The connection string for the database to run migrations on.
#
# You might prefer to set this using an environment variable.
#
# Default: "" (default PostgreSQL server)
database-url = ""

# The directory used to store migration files.
#
# Default: "migrations"
migrations-dir = "migrations"

# The template to use for new migration files.
#
# Default: "" (use the embedded default migration template)
template-file = "migrations/_template.sql"

# How much info to log to stderr. Greater numbers mean more output, and 0 logs
# nothing.
#
# Default: 1
verbosity = 1
```

Then, generate the first migration that sets up Drift's requirements:

```bash
drift setup
```

That should have written `0-init.sql` to your migrations directory. Read it and
make any changes you want.

Finally, run the migration:

```bash
drift migrate
```

### Writing a new migration

Create a new empty migration file:

```bash
drift new --slug 'create_users_table'
```

Write your migration in the file. Then run it:

```bash
drift migrate
```

### Undoing a migration

For a migration that has already been run in production (or some other shared
environment), the best option is to write a new migration to undo the old one.

To be able to re-run a migration you're developing, call `_drift_unclaim_migration`
with the migration ID. You might find it useful to keep an `undo.sql` file
around (ignored by version control) to modify along with the migration.

```sql
-- undo.sql
select _drift_unclaim_migration(1645673864);

drop table if exists users;
```

## License

Source code and binaries are distributed under the terms of the MIT license.

## Contributing

I welcome contributions from anyone who finds a way to make this better.

In particular, I appreciate these things:
- Pull requests with clear motivation (tests are nice too!)
- Bug reports with reproducible setup instructions
- Ideas about things that work but can be improved

## Support

This is is hobby project, so I'll only work on it as much as I find it fun to
do so. That said, I find software maintenance techniques interesting, so feel
free to start a conversation about a stable v1 if you start relying on this for
something important.
