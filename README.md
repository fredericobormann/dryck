# dryck
[![Build Status](https://travis-ci.com/fredericobormann/dryck.svg?branch=master)](https://travis-ci.com/fredericobormann/dryck)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredericobormann/dryck)](https://goreportcard.com/report/github.com/fredericobormann/dryck)

dryck is a very simple web app written in Go using Gin to keep track of the amount of drinks taken by different people in a shared space.
At the moment it doesn't support user authentication, so you should trust everyone who has access to it.
If you set the environment variable `HTTP_PASSWORD`, you can limit the access in a very simple way with HTTP Basic Auth (username: dryck).
If `HTTP_PASSWORD` is not set, dryck won't require authentication.
Also adding new drinks is only possible by editing the underlying Postgres database directly.

## Installation
The simplest way to start the application is by using docker compose.
1. Install docker and docker compose.
2. Create a `docker-compose.yml` with the following content and change the variables according to your needs:
```yaml
version: '2.0'
services:
  web:
    image: fredericobormann/dryck
    restart: always
    environment:
      POSTGRES_HOST: postgres
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: mysecretpassword
      POSTGRES_DATABASE: postgres
      GIN_MODE: release
      HTTP_PASSWORD: mysecretpassword
    ports:
      - "8089:8080"
    depends_on:
      - postgres
  postgres:
    image: postgres
    restart: always
    environment:
      POSTGRES_PASSWORD: mysecretpassword
    ports:
      - "5432:5432"
    volumes:
      - db-data:/var/lib/postgresql/data
volumes:
  db-data:
    driver: local
```
3. Then start everything by
```
docker-compose up -d
```
