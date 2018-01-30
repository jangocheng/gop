# GOP

[简体中文](README_ZH.md)

[![CircleCI](https://circleci.com/gh/lunny/gop.svg?style=shield)](https://circleci.com/gh/lunny/gop) [![codecov](https://codecov.io/gh/lunny/gop/branch/master/graph/badge.svg)](https://codecov.io/gh/lunny/gop)
[![](https://goreportcard.com/badge/github.com/lunny/gop)](https://goreportcard.com/report/github.com/lunny/gop) 

GOP is a project manangement tool for building your golang applications out of global GOPATH. In fact gop will keep both global GOPATH and every project GOPATH. But that means your project will  **not** go-getable. Of course, GOP itself is go-getable. GOP copy all denpendencies from global GOPATH to your project's `src/vendor` directory and all application's source is also in `src` directory.

A normal process using gop is below:

```
git clone xxx@mydata.com:bac/aaa.git
cd aaa
gop ensure -g
gop build
gop test
```

## Features

* GOPATH compitable
* Multiple build targets support
* Put your projects out of global GOPATH

## Installation

Please ensure you have installed the `go` command, GOP will invoke it on building or testing

```
go get github.com/lunny/gop
```

## Directory structure

Every project should have a GOPATH directory structure and put a `gop.yml` int the root directory. This is an example project's directory tree.

```
<project root>
├── gop.yml
├── bin
├── doc
└── src
    ├── main
    │   └── main.go
    ├── models
    │   └── models.go
    ├── routes
    │   └── routes.go
    └── vendor
        └── github.com
            ├── go-xorm
            │   ├── builder
            │   ├── core
            │   └── xorm
            └── lunny
                ├── log
                └── tango
```

## Gop.yml

Gop will recognize a gop project which has `gop.yml`. The file is also a project configuration file. Below is an example. If you didn't define any target, the default target is src/main and the target name is the project name.

```yml
targets:
- name: myproject1
  dir: main
  assets:
  - templates
  - public
  - config.ini
  - key.pem
  - cert.pem
- name: myproject2
  dir: web
  assets:
  - templates
  - public
  - config.ini
```

## Command

### init

Create the default directory structure tree.

```
gop init
```

### ensure

Automatically copy dependencies from $GOPATH to local project directory. `-g` will let you automatically call `go get <package>` when the package is missing on `GOPATH`. `-u` will always `go get <package>` on all the dependencies and copy them to `vendor`.

```
gop ensure [-g|-u] [target_name]
```

### status

List all dependencies of this project and show the status.

```
gop status [target_name]
```

### add

Add a package to this project.

```
gop add <package>
```

### update

Update a package to this project.

```
gop update <package>
```

### rm

Remove a package from this project.

```
gop rm <package>
```

### build

Run `go build` on the src directory.

```
gop build [target_name]
```

### run

Run `go run` on the src directory.

```
gop run [target_name]
```

### test

Run `go test` on the src directory.

```
gop test [target_name]
```

### release

Run `go release` on the src directory.

```
gop release [target_name]
```

## TODO

* [ ] Versions support, specify a dependency package verison
* [ ] Improve binary package building support
* [ ] Live support for `gop run`
* [ ] Support run `gop` in `GOPATH`