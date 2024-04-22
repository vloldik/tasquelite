# TasqueLite: A [tasque](https://github.com/vloldik/tasque) Storage Package with GORM and SQLite

[![GoDoc](https://godoc.org/github.com/vloldik/tasquelite?status.svg)](https://godoc.org/github.com/vloldik/tasquelite)
[![Go Report Card](https://goreportcard.com/badge/github.com/vloldik/tasquelite)](https://goreportcard.com/report/github.com/vloldik/tasquelite)

TasqueLite is a Go package that provides a [tasque](https://github.com/vloldik/tasque) storage for concurrent processing of tasks. It uses GORM, a powerful ORM library for Go, to interact with SQLite, a lightweight disk-based database, for task storage.

## Features

- **Task Queue**: Process tasks concurrently using a task queue.
- **GORM Integration**: Use GORM to interact with SQLite for task storage.
- **SQLite Database**: Store tasks in a lightweight disk-based database.

## Installation

To install TasqueLite, use `go get`:
```sh
go get github.com/vloldik/tasquelite
```

## Usage

```go
taskStorage, err := tasquelite.NewGormTaskStorageManager[EmailTask](DatabaseName, &Task{})
```
