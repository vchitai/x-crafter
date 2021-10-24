# x-crafter

:smile: x-crafter is used to quickly create templates from your prototype, also come with a builder to quickly regenerate your code.

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/vchitai/x-crafter)
[![License](https://img.shields.io/badge/license-MIT-%2397ca00.svg)](https://github.com/vchitai/x-crafter/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/vchitai/x-crafter.svg)](https://github.com/vchitai/x-crafter/releases)
[![Made by vchitai](https://img.shields.io/badge/made%20by-vchitai-blue.svg?style=flat)](https://vchitai.github.io/)

[![GolangCI](https://golangci.com/badges/github.com/vchitai/x-crafter.svg)](https://golangci.com/r/github.com/vchitai/x-crafter)
[![codecov](https://codecov.io/gh/vchitai/x-crafter/branch/main/graph/badge.svg?token=6QWOopYRPD)](https://codecov.io/gh/vchitai/x-crafter)
[![Go Report Card](https://goreportcard.com/badge/github.com/vchitai/x-crafter)](https://goreportcard.com/report/github.com/vchitai/x-crafter)
[![CodeFactor](https://www.codefactor.io/repository/github/vchitai/x-crafter/badge)](https://www.codefactor.io/repository/github/vchitai/x-crafter)


## Install

### Using go

```console
$ go get -u github.com/vchitai/x-crafter/cmd
```

## What is this tool for?

Remember the last time you set a new project up? Normally, we just clone an old project to a new one. 
Later, we replace some old names with new ones. And do some setup jobs.
This may have been done from time to time. 
Especially, when you are working on a microservice ecosystem, where the need of initializing a new project come up occasionally.

I have been through many projects, thourgh many companies.
And I have also been doing that job from time to time for a long time. 

Some organization may come up with a prototype project.
You can just fork from it and do some setup. 
But this also require you to know what to replace, and do some correct replacements to make it work.

Some will create their SDK for their own organization. 
That may cost them a lot of time like I have. 
Also the template in some case is compressed into assets so that it may be re-parsable. You may set all the project up for the first time. But as the time flies, we always need to upgrade the legacy. That's where the problems come frome. The language used to write template like golang template or jinja2 is different than the main using language, we can just validate the output only when all the files have parsed. To make that work, we may need to write again and test from time to time, limited our proeffiency.

So there's this tool come to help, to combine all the advantages from both the approaches. XCrafter come with a batch of tooling that enable developers to deal with the template creation and parsing them much faster. We will use the base prototype project to generate the template, that keep the process easy to approach.

Let us go through a quick example to demonstrate the way it may help you achieve your goal.

# Break a whole project

Set up a whole project is a common case in a new world. As aforementioned, we can fork a template repository into a new one. We can also copy an old project into a new one and do some setup. But this approach may meet some limits. 

So we are going to use the turn a whole project into templates.

```bash
example/go-proj
├── cmd
│   └── main
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── module_name
│   │   └── module.go
│   └── test
│       ├── embed.go
│       └── x.html
└── xbreak.yaml
```

Same as the above example, we should also provide a break guide for the breaker in order to create the templates.

```yaml
syntax:
  ".go":
    start: "/*X"
    end: "X*/"
    inline: "//X"
  ".html":
    start: "<!--X"
    end: "X-->"

layers:
  base:
    - cmd/*
    - internal/test/*
  module:
    - internal/module_name/*

args:
    - module_name

kwargs:
  package_name: package_name
  var_customer: var_customer
```

We define the syntax for go files and html files in the `syntax` option. Then the `layers` to start collect the scraps buckets for later use. With the `args` and `kwargs` option, the `breaker` will fine and mark all the provided vocabulary set as arguments, which are ready to be replaced in the future.

Run the break command
```bash
$ x-crafter break example/go-proj
```

And as result,

```bash
example/go-proj_broken
├── layers
│   ├── base
│   │   ├── cmd
│   │   │   └── main
│   │   │       └── main.go.tmpl
│   │   └── internal
│   │       └── test
│   │           ├── embed.go.tmpl
│   │           └── x.html.tmpl
│   └── module
│       └── internal
│           └── {{module_name}}
│               └── module.go.tmpl
└── version
```

## Rebuild the project again

The broken components have been arranged into a structure that is ready to be used to recrafted a new project.

In order to use the builder, you must provide a build guide. You can provide a custom path with `guide` option like `--guide=example/build/xbuild.yml`. If you did not provide one, the builder will look for `xbuild.yaml` in your template folder for the guide.

```yaml
params:
  module:
    - module_name: a_secret_module
  package_name: package_x
  var_customer: customer A
  project_pkg_path: github.com/vchitai/x-crafter/example/go-proj_rebuilt

template_root: layers

steps:
  - name: parse base
    parse: base
    on: .
  - name: parse module
    parse: module
    repeat:
      for: module
  - name: init go mod
    run: ["go", "mod", "init", "${PROJECT_PKG_PATH}"]
  - name: tidy
    run: ["go", "mod", "tidy"]
    on: .
    env:
        - GOSUMDB=off
```

In this guide, we provide the builder the parameter that used to replace for the arguments was marked in the break step with `params` option. `template_root` specifiy the template folder root if you have another template project structure. The `steps` option specify the steps that need to be done in order to craft you next project. `name` annotate the name of the step. `parse` option tell the builder to parse the layer, relative to the `template_root` into the the location specify in `on` option. The module may be parse multiple times using the `repeat` option using the params.

Parsing may not the only thing needed to set you project up. Sometimes you may need to run additional setup with command lines. So use the `run` option, that will run the following command from the location annotated in `on` option. Extra environment variable may be passed into the comman using the `env` option. You can also use the environment that created through the `params` option, for example `${PROJECT_PKG_PATH}` that is provided through `project_pkg_path`

Run the build command
```bash
$ x-crafter build example/go-proj_broken example/go-proj_broken
```

and what we will get 

```bash
example/go-proj_rebuilt
├── cmd
│   └── main
│       └── main.go
├── go.mod
└── internal
    ├── module_name
    │   └── module.go
    └── test
        ├── embed.go
        └── x.html
```

The project is all set up and ready to be run. 

[comment]: <> (## Stargazers over time)

[comment]: <> ([![Stargazers over time]&#40;https://starchart.cc/vchitai/x-crafter.svg&#41;]&#40;https://starchart.cc/vchitai/x-crafter&#41;)
