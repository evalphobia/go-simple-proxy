go-simple-proxy
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Co decov Coverage][11]][12] [![Go Report Card][13]][14] [![Downloads][15]][16]

[1]: https://godoc.org/github.com/evalphobia/go-simple-proxy?status.svg
[2]: https://godoc.org/github.com/evalphobia/go-simple-proxy
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/go-simple-proxy.svg
[6]: https://github.com/evalphobia/go-simple-proxy/releases/latest
[7]: https://travis-ci.org/evalphobia/go-simple-proxy.svg?branch=master
[8]: https://travis-ci.org/evalphobia/go-simple-proxy
[9]: https://coveralls.io/repos/evalphobia/go-simple-proxy/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/go-simple-proxy?branch=master
[11]: https://codecov.io/github/evalphobia/go-simple-proxy/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/go-simple-proxy?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/go-simple-proxy
[14]: https://goreportcard.com/report/github.com/evalphobia/go-simple-proxy
[15]: https://img.shields.io/github/downloads/evalphobia/go-simple-proxy/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/go-simple-proxy/releases
[17]: https://img.shields.io/github/stars/evalphobia/go-simple-proxy.svg
[18]: https://github.com/evalphobia/go-simple-proxy/stargazers

`go-simple-proxy` is simple proxy to forwading port.
Suit for unit testing.

# Installation

Install go-simple-proxy using `go get` command:

```bash
$ go get github.com/evalphobia/go-simple-proxy
```

# Usage

```bash
$ go-simple-proxy
[PROXY ERROR] 2017/05/19 15:03:48 Argument for proxy is missing.

usage:

$ go-simple-proxy "<protocol>,<from_host:port>,<to_host:port>" ...
ex) go-simple-proxy "tcp,localhost:8080,example.com:80" "udp,localhost:8081,example2.com:80" "ws,localhost:8082,example3.com:80"

# hostnames are supplied by env var.
# The ports of redis, DynamoDB, SQS and Elasticsearch will be forwarded.
$ go-simple-proxy "tcp,localhost:6379,${REDIS_HOST}:6379" \
    "tcp,localhost:8000,${AWS_DYNAMODB_HOST}:8000" \
    "tcp,localhost:4568,${AWS_SQS_HOST}:4568" \
    "tcp,localhost:9200,${ELASTICSEARCH_HOST}:9200"
```

## Supported Protocol

- TCP
- ~~UDP~~ (not supported yet)
- ~~WebSocket~~ (not supported yet)

# Parameters

## arguments

Add forwarding setting.
Setting format is `<protocol>,<from_host:port>,<to_host:port>`.

- protocol
    - `tcp`
- from_host
    - localhost
    - 127.0.0.1
    - 0.0.0.0
    - x.x.x.x (network I/F IP address)
- to_host
    - set any IP address

## flag option

| name | description |
| ------- | ------- |
| `-v`  | output verbose logs |
| `-vv`  | output debug logs |
| `-timeout`  | request/response timeout, integer with time sign (e.g. `10s`, `500ms`, `1.5h`) |


# Whats For?

For the monolithic application, sometimes the configuration is almost static and difficult to change dynamically.
Then in unit testing on CI services or Docker containers, it's really hard to setting middleware hosts.

For my usecase, using Docker Compose with these containers

- x3 containers for Golang tests
- x1 Middlewares for each services
    - Redis
    - MySQL
    - Elasticsearch
    - etc...

Settings are like below;

```yml
version: "2"
services:
  # unittest container #1
  unittest-app-1:
    container_name: "unittest-app-1"
    image: "go:1.8"
    entrypoint:
      - make
      - unittest-1
    environment:
      CONFIG_TEST_PREFIX: test_1_
      REDIS_HOST: unittest-redis
      REDIS_TEST_DB: 1
      ELASTICSEARCH_HOST: unittest-elasticsearch
      MYSQL_TEST_HOST: unittest-database
      MYSQL_TEST_PORT: 3306
      MYSQL_TEST_DATABASE: test1
      MYSQL_TEST_USER: root
    links:
      - unittest-redis
      - unittest-elasticsearch
      - unittest-database
  # unittest container #2
  unittest-app-2:
    container_name: "unittest-app-2"
    image: "go:1.8"
    entrypoint:
      - make
      - unittest-2
    environment:
      CONFIG_TEST_PREFIX: test_2_
      REDIS_HOST: unittest-redis
      REDIS_TEST_DB: 2
      ELASTICSEARCH_HOST: unittest-elasticsearch
      MYSQL_TEST_HOST: unittest-database
      MYSQL_TEST_PORT: 3306
      MYSQL_TEST_DATABASE: test2
      MYSQL_TEST_USER: root
    links:
      - unittest-redis
      - unittest-elasticsearch
      - unittest-database

...
  # middleware container: Elasticsearch
  unittest-elasticsearch:
    image: evalphobia/elasticsearch:1.0.1
    restart: always
  # middleware container: Redis
  unittest-redis:
    image: redis:3.2-alpine
    restart: always
  # middleware container: MySQL
  unittest-database:
    container_name: "unittest-database"
    image: mysql:5.7
    environment:
      MYSQL_ALLOW_EMPTY_PASSWORD: "yes"
      MYSQL_USER: test_user
      MYSQL_PASSWORD: test_pass
    restart: always
```

And Makefile is like this,
(go-simple-proxy is used on the last part)

```bash
# test only dir1/, dir2/, dir3/
unittest-1: unittest-prepare
	go test -p 1 ./dir1/... ./dir2/... ./dir3/...

# test only dir4/, dir5/
unittest-2: unittest-prepare
	go test -p 1 ./dir4/... ./dir5/...

# setup  test data
unittest-prepare:
	# database migration
	while ! mysqladmin ping -h ${MYSQL_TEST_HOST} --silent; do \
	    sleep 1; \
	done
	mysql -h ${MYSQL_TEST_HOST} -e "DROP DATABASE IF EXISTS ${MYSQL_TEST_DBNAME};"
	mysql -h ${MYSQL_TEST_HOST} -e "CREATE DATABASE IF NOT EXISTS ${MYSQL_TEST_DBNAME} DEFAULT COLLATE utf8mb4_general_ci;"
    # In this case, test data is already set on seed DB (${MYSQL_TEST_DBNAME_SEED})
    # Here, copy data from seed DB to new DB for this container.
	mysqldump -h ${MYSQL_TEST_HOST} ${MYSQL_TEST_DBNAME_SEED} | mysql -h ${MYSQL_TEST_HOST} ${MYSQL_TEST_DBNAME}

	# change config for this container
	# (In this case, change prefix to this container, looks like: test_1_, test_2_, ...)
	sed -e "s@^\(prefix\s*= \)\".*@\1\"${CONFIG_TEST_PREFIX}\"@g" -i ${CONFIG_PATH}/config.tml

	# proxy to middlewares
	# (Our golang test code wants to connect middlewares on localhost, so this proxy foward to other container)
	go get github.com/evalphobia/go-simple-proxy
	go-simple-proxy \
	  "tcp,localhost:6379,${REDIS_HOST}:6379" \
	  "tcp,localhost:3306,${MYSQL_TEST_HOST}:3306" \
	  "tcp,localhost:9200,${ELASTICSEARCH_HOST}:9200" & # run as background process
```
