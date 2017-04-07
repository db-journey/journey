# Journey

[![Build Status](https://travis-ci.org/db-journey/journey.svg?branch=master)](https://travis-ci.org/db-journey/journey)
[![GoDoc](https://godoc.org/github.com/db-journey/journey?status.svg)](https://godoc.org/github.com/db-journey/journey)

Journey is based on the work of @mattes on his tool "migrate": https://github.com/mattes/migrate/

## __Features__

* Super easy to implement [Driver interface](http://godoc.org/github.com/db-journey/migrate/driver#Driver).
* Gracefully quit running migrations on ``^C``.
* No magic search paths routines, no hard-coded config files.
* CLI is build on top of the ``migrate`` package.


## Available Drivers

 * [PostgreSQL](https://github.com/db-journey/postgresql-driver)
 * [Cassandra](https://github.com/db-journey/cassandra-driver)
 * [SQLite](https://github.com/db-journey/sqlite3-driver)
 * [MySQL](https://github.com/db-journey/mysql-driver) ([experimental](https://github.com/mattes/migrate/issues/1#issuecomment-58728186))
 * Bash (planned)

Need another driver? Just implement the [Driver interface](http://godoc.org/github.com/db-journey/migrate/driver#Driver) and open a PR.

## Usage from Terminal

```bash
# install
go get github.com/db-journey/journey

# create new migration file in path
journey --url driver://url --path ./migrations migrate create migration_file_xyz

# apply all available migrations
journey --url driver://url --path ./migrations migrate up

# roll back all migrations
journey --url driver://url --path ./migrations migrate down

# roll back the most recently applied migration, then run it again.
journey --url driver://url --path ./migrations migrate redo

# run down and then up command
journey --url driver://url --path ./migrations migrate reset

# show the current migration version
journey --url driver://url --path ./migrations migrate version

# apply the next n migrations
journey --url driver://url --path ./migrations migrate migrate +1
journey --url driver://url --path ./migrations migrate migrate +2
journey --url driver://url --path ./migrations migrate migrate +n

# roll back the previous n migrations
journey --url driver://url --path ./migrations migrate migrate -1
journey --url driver://url --path ./migrations migrate migrate -2
journey --url driver://url --path ./migrations migrate migrate -n

# go to specific migration
journey --url driver://url --path ./migrations migrate goto 1
journey --url driver://url --path ./migrations migrate goto 10
journey --url driver://url --path ./migrations migrate goto v
```

## CronJobs

Journey also provides a command to run scheduled jobs on databases:


```bash
journey --url driver://url --path ./cronjobs scheduler start
```
