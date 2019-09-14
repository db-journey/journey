# Journey Changelog

## 2.2.1 - 2019-09-14

* Add go modules support

## 2.2.0 - 2018-01-22

* Update to migrate [v2.0.0](https://github.com/db-journey/migrate/releases/tag/v2.0.0)
* [Feature] apply specific version (run up migration)
* [Feature] rollback specific version (run down migration)

## 2.1.1 - 2017-12-12

* Fix docker image building

## 2.1.0 - 2017-07-07

* [migrate] Add templating support in migration files
* [scheduler] Add jobs duration to log

## 2.0.0 - 2017-04-06

This is a breaking change release.
The cli now expose two main commands:

* "migrate", the old behaviour
* "cronjobs", added in this release

ie: old params and commands are working if using `journey migrate [...]`

* Add cronjobs support (https://github.com/db-journey/cronjobs)
* Switch to urfave cli
* Provide commands as a package

## 1.4.2 - 2017-02-07

* Split repos from gemnasium/migrate (a fork of mattes/migrate)

## v1.4.1 - 2016-12-16

* [cassandra] Add [disable_init_host_lookup](https://github.com/gocql/gocql/blob/master/cluster.go#L92) url param (@GeorgeMac / #17)

## v1.4.0 - 2016-11-22

* [crate] Add [Crate](https://crate.io) database support, based on the Crate sql driver by [herenow](https://github.com/herenow/go-crate) (@dereulenspiegel / #16)

## v1.3.2 - 2016-11-11

* [sqlite] Allow multiple statements per migration (dklimkin / #11)

## v1.3.1 - 2016-08-16

* Make MySQL driver aware of SSL certificates for TLS connection by scanning ENV variables (https://github.com/mattes/migrate/pull/117/files)

## v1.3.0 - 2016-08-15

* Initial changelog release
* Timestamp migration, instead of increments (https://github.com/mattes/migrate/issues/102)
* Versions will now be tagged
* Added consistency parameter to cassandra connection string (https://github.com/mattes/migrate/pull/114)


