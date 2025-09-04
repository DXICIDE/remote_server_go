# remote_server_go

## What is it
It's a copy of cheap copy of twitter(X) without a name, that has basic functions of Authentication and Authorization, making chirps(tweets) and some basic API and token functions.  

This is a guided project on [Boot.dev](https://boot.dev). Yes it is not a very presentable project for a job, nevertheless it's something I've written with a little bit of guidance and taught me a lot about api, routing, databases, webhooks, tokens, refresh tokens and so much more. It just shows that I can follow instructions and deliver.

## Installation

### 1) Clone
git clone github.com/DXICIDE/remote_server_go

### 2) Configure your environment
you can check the main.go which .env variable I use

### 3) initialize the DB
go into db/schema and create also create account for postgres
goose postgres "postgres://*username*:*password*@localhost:5432/chirpy" up

### 4) Start the root
go run . in the root

### 5) Enjoy
Open another console and make your requests

Example of requests:
* curl -X POST http://localhost:8080/admin/reset
* curl -X http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"email":"a@example.com","password":"pass"}'
* curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"email":"a@example.com","password":"pass"}'
* curl -X POST http://localhost:8080/api/chirps -H "Authorization: Bearer $TOKEN_A" -H "Content-Type: application/json" -d '{"body":"First from A"}'
* curl -X POST http://localhost:8080/api/chirps -H "Authorization: Bearer $TOKEN_A" -H "Content-Type: application/json" -d '{"body":"Second from A"}'
* curl http://localhost:8080/api/chirps
* curl http://localhost:8080/api/chirps?author_id=<A_id>
