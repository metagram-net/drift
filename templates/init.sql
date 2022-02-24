/*
Set up the Drift framework requirements. Naturally, this first migration is
going to break a few rules ;)

First, this includes a drift:no-transaction directive, which tells Drift to
skip two steps it would normally take:
1. Opening a transaction around the migration file. In Postgres, DDL can be
   done in a transaction. This can make some migrations safer, so Drift assumes
transactions as the default.
2. Calling _drift_claim_migration(id, slug) before running the file. Since this
   claim would fail on a duplicate ID, this ensures we never run a migration
twice (since it's normally part of a transaction).

It doesn't make sense to call _drift_claim_migration yet, because this is the
migration that defines it!

You can modify the _drift_claim_migration function if you want to. The only
expectation Drift has of it (besides the signature) is that it writes the
migration ID to the table and fails if that ID is already recorded.

You can also modify _drift_require_migration; Drift doesn't use it. It's useful
to call within a migration when it would only make sense to run after some
earlier one has completed.

You can also modify the schema_migrations table, but (at least for now) Drift
assumes that the migration records table has exactly that name and has the
integer primary key id column.
*/
--drift:no-transaction

begin;

create table schema_migrations (
    id integer primary key,
    slug text not null,
    run_at timestamp not null default current_timestamp
);

-- _drift_claim_migration registers a migration in the schema_migrations table.
-- It will fail if the migration ID has already been claimed.
--
-- Drift will call this at the start of every migration transaction. For
-- migrations that cannot be run within transactions, it is the migration's
-- responsibility to call this.
create function _drift_claim_migration(mid integer, mslug text) returns void as $$
    insert into schema_migrations (id, slug) values (mid, mslug);
$$ language sql;

-- _drift_unclaim_migration removes a migration from the schema_migrations
-- table.
--
-- When iterating on a migration in development, it's useful to have an "undo"
-- or "down" migration to reset back to the previous schema. Call this to undo
-- the automatic _drift_claim_migration to be able to re-run the "up"
-- migration.
create function _drift_unclaim_migration(mid integer) returns void as $$
    delete from schema_migrations where id = mid;
$$ language sql;

-- _drift_require_migration asserts that the migration ID has already been
-- claimed in the schema_migrations table.
--
-- Call this from within a migration to ensure that another migration has
-- already run to completion.
create function _drift_require_migration(mid integer) returns void as $$
declare
    mrow schema_migrations%rowtype;
begin
    select * into mrow from schema_migrations where id = mid;
    if not found then
        raise exception 'Required migration has not been run: %', mid;
    end if;
end;
$$ language plpgsql;

-- Normally, this would be the first thing in the migration, but we had to
-- create the schema_migrations table first!
select _drift_claim_migration(0, 'init');

commit;
