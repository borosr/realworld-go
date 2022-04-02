# ![RealWorld Example App](logo.png)

> ### Golang net/http (with generics) + BadgerDB codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with Golang net/http (with generics) + BadgerDB including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the Golang community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

I tried to create a wrapper library on the top of the built-in net/http package to support generics in http handler methods.

# Getting started

Project is using Go version 1.18 and BadgerDB. Because BadgerDB is a general key-value store, no need for migration.

- Make sure Go 1.18 version installed on your machine
- To install dependencies run `go mod vendor`
- Then start the server with `go run main.go` from the project root
- NOTE: the badgerDB will create its own files under `/tmp/badger`, to flush the database, delete this directory
