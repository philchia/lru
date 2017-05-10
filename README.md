# lru is a in memory lru cache made by go [![CircleCI](https://circleci.com/gh/philchia/lru.svg?style=svg)](https://circleci.com/gh/philchia/lru)

[![Go Report Card](https://goreportcard.com/badge/github.com/philchia/lru)](https://goreportcard.com/report/github.com/philchia/lru)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/philchia/lru?status.svg)](https://godoc.org/github.com/philchia/lru)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://opensource.org/licenses/MIT)

## Installation

```shell
    go get -u github.com/philchia/lru
```

## Usage

### Create a cache

```go
import "github.com/philchia/lru"
cache := lru.New(10)
```

### Create a thread safe cache

```go
import "github.com/philchia/lru"
cache := lru.NewLockCache(10)
```

### Store a k/v

```go
cache.Set("key", "value")
```

### Store a k/v with expired time

```go
cache.Set("key", "value", time.Now())
```

### Get a value with a key

```go
if value :=cache.Get("key"); value != nil {
    print(value)
}
```

### Delete a key

```go
cache := lru.New(10)
cache.Del("key")
```

## License

lru is published under the MIT license