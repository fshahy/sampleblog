# Sample Blog [![Build Status](https://travis-ci.org/fshahy/sampleblog.svg?branch=master)](https://travis-ci.org/fshahy/sampleblog)
This repository contains the code for a sample blog which contains 3 RESTFul endpoints:
- Create blog article
- Get blog article by ID
- Get all blog articles

## Technical Overview
This blog is written in pure Go language (Golang). Blog articles are saved in a PostgreSQL database. 

It has two main components:
- API Server: Which runs in a docker container based on **golang:1.12.7** image.
- PostgreSQL Serve: It also runs in a seperate docker container based on **postgres:latest** image.
Using a docker compose file, the two containers are set up and the API server communicates with database server in docker bridged network.
All database installtion and initialization is done in the **start.sh** file.

## Docker Installation
You need to install Docker and Docker Compose before starting.
From a terminal change directory to the root of project and run this:

```docker-compose up -d```

All containers will start up and run in as daemon.
The API server exposes port **8080** and will be available the host through:

**http://localhost:8080/articles** 

which lists all articles. If you want access a specific articles use the following url:

**http://localhost:8080/articles/{id}** 

in which **{id}** is replaced with the articles ID.
For example 

**http://localhost:8080/1** 

will retrieve article with **ID=1**.

## Local Installation
If youdo not have docker installed, you can install this blog locally.
- Install Golang.
- Install PostgreSQL server.
- Install PostgreSQL driver for Golang: ```go get github.com/lib/pq```
- Create blog database: switch user to postgres and run **start.sh** file to create blog database.
- Set the following environment variables. Off course you can change the values, these are the values used in docker installation.
    - $POSTGRES_HOST: localhost
    - $POSTGRES_DB: blogdb
    - $POSTGRES_USER: postgres
    - $POSTGRES_PASSWORD: secret
    - $POSTGRES_TEST_DB: blogb_test

The last environment variable is meant to be used by unit/integration tests. We use a seperate database for testing the blog server.

To start your server run:

```go run main.go```

Open your browser and goto **http://localhost:8080/articles**
