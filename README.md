# dryck
[![Build Status](https://travis-ci.com/fredericobormann/dryck.svg?branch=master)](https://travis-ci.com/fredericobormann/dryck)

dryck is a very simple web app written in Go using Gin to keep track of the amount of drinks taken by different people in a shared space.
At the moment it doesn't support user authentication, so you should trust everyone who has access to it.
On the http-basic-auth branch is a version available that uses http basic authentication (user: dryck), so you can limit the access in a simple way.
Also adding new drinks is only possible by editing the underlying Postgres database directly.

## Installation
The simplest way to start the application is by using docker compose.
1. Install docker and docker compose.
1. Pull the repo and make sure that you change the database password (and the http basic auth password) in the `docker-compose.yml`.
1. Then start everything by
```
docker-compose up -d
```
