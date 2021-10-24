# XCrafter

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
$ go get -u github.com/vchitai/x-crafter
```

## Usage

Use --help to see more option
```bash 
X-Crafter is used to quickly make a go code prototype quickly become reproducible

Usage:
  x-crafter [command]

Available Commands:
  break       Convert your prototype into recraftable golang templates
  build       Using your broken pile to once again craft your thing again
  completion  generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help   help for x-crafter

Use "x-crafter [command] --help" for more information about a command.
```

## What is this tool for?
Remember the last time you set a new project up? Normally, we just clone an old project into a new one.
Later, we will replace some old names with new ones. And do some setup jobs.
This may have been done from time to time.
Especially, when you are working on a microservice ecosystem, where the need to initialize a new project comes up occasionally. I have initialized many projects.
And I have also been doing that job from time to time for a long time.

To resolve this problem, some organizations may come up with a prototype project.
You can just fork it and do some setup.
But this also requires you to know what to replace and make some correct replacements to make it work.

Some company will develop their own SDK for their own use.
That may cost them a lot of time like I have.
Also, the template, in some cases, is compressed into assets so that it may be re-parsable. You may set all the projects up for the first time. But as the time flies, we always need to upgrade the legacy. That's where the problems come from. The language used to write templates like gotemplate or jinja2 is different from the main usage language. We can just validate the output only when all the files have been parsed. To make that work, we may need to write again and test from time to time, limiting our proeffiency.

So there's this tool that comes to help, to combine all the advantages of both approaches. XCrafter comes with a batch of tooling that enables developers to deal with template creation and parsing them much faster. We will use the base prototype project to generate the template, which will keep the process easy to approach.
Let us go through a quick example to demonstrate the way it may help you achieve your goal.

## Creating k8s chart

If you've ever written a K8s application, you've probably struggled with the deploy configuration, service configuration, secret binding, and configmap.As a result, many organizations may use helm to manage the packages in their k8s ecosystem.The Helm template is useful to help developers create deployments with a minimum knowledge of the K8s system. But it is written in Golang template syntax. A yaml with a golang template syntax can not be parsed and be well formatted by any yaml-supported IDE. Hence, you may make many mistakes in creating those templates, like I did.

Let's take a look at an example project:

```yaml
{{- if .Values.rbac.create }}
apiVersion: rbac.authorization.k8s.io/v1
{{- if .Values.rbac.clusterWide }}
kind: ClusterRole
{{- else }}
kind: Role 
{{- end }}
metadata:
  name: {{ template "telegraf.fullname" . }}
  namespace: {{ .Release.Namespace }}
  labels:
  {{- include "telegraf.labels" . | nindent 4 }}
rules:
  {{ toYaml .Values.rbac.rules | indent 2 }}
{{- end }}
```

This syntax is not supported by any yaml-supported IDE.So the formatter may not be used. This creates a lot of problems, especially if you're writing for a heavily indented language.Any errors will come up and you will need to adjust and regenerate them to test them.

With XCrafter, we take another approach that takes the native comments of each language to mark a syntax to Golang template conversion:

```yaml
#X- if .Values.rbac.create
apiVersion: rbac.authorization.k8s.io/v1
#X- if .Values.rbac.clusterWide
kind: ClusterRole
#X- else
#Okind: Role
#X- end
metadata:
  name: #X template "telegraf.fullname" .
  namespace: #X .Release.Namespace
  labels:
  #X- include "telegraf.labels" . | nindent 4
rules:
  #X toYaml .Values.rbac.rules | indent 2
#X- end
```

The syntax is marked without any errors happening to the formatter. So we ensure the last output is valid. The comment should be converted to golang template syntax so that it can be parsed by the golang template engine. In order to do that, just provide the tool a descriptor for your own syntax.

```yaml
syntax:
  ".yaml":
    inline: "#X"
    optional: "#O"

layers:
  charts:
    - role.yaml
```

With the `layers` option, we will put the parsed components into different baskets. That will be useful for analyzing a large project later.
The `syntax` option defines what we should do with the comments in the file with the provided extension. In this example, for each `yaml` file, we will try to turn the inline comment, which is marked with `#X` at the start of the line, into a golang syntax block. `#O` stands for an optional code block, that may be available with the condition above. It may cause some problems if it comes up in the main code base, but it should be available for parsing as an optional choice. So we will escape it in order to make it work like that.

The result is as below:

```yaml
{{- if .Values.rbac.create }}
apiVersion: rbac.authorization.k8s.io/v1
{{- if .Values.rbac.clusterWide }}
kind: ClusterRole
{{- else }}
kind: Role 
{{- end }}
metadata:
  name: {{ template "telegraf.fullname" . }}
  namespace: {{ .Release.Namespace }}
  labels:
  {{- include "telegraf.labels" . | nindent 4 }}
rules:
  {{ toYaml .Values.rbac.rules | indent 2 }}
{{- end }}
```

This can be used with any Golang template engine and the Helm template also. The template can be quickly updated by changing the protype yaml above and hitting the break button again.

## Break a whole project

Setting up a whole project is another common case in the real world. As aforementioned, we can fork a template repository into a new one. We can also copy an old project into a new one and do some setup. But this approach may meet some limits.

So we are going to turn a whole project into templates.

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

As in the above example, we should also provide a break guide for the breaker in order to create the templates.

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

We define the syntax for Go files and HTML files in the `syntax` option. Then the `layers` option is used to collect the scrap buckets for later use. With the `args` and `kwargs options, the `breaker` will work fine and mark all the provided vocabulary sets as arguments, which are ready to be replaced in the future.

After all those works, just run the break command

```bash
$ x-crafter break example/go-proj
```

And get this result.

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

The broken components have been well arranged into a structure that is ready to be used to recraft a new project.

In order to use the builder, you must provide a build guide. You can provide a custom path with `guide` option like `--guide=example/build/xbuild.yml`. If you don't provide one, the builder will look in your template folder for "xbuild.yaml."

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

In this guide, we provide the builder with the parameters that are used to replace the arguments that were marked in the break step with `params` option. `template_root specifies the template folder root if you have another template project structure. The `steps` option specifies the steps that need to be taken in order to craft your next project. `name annotates the name of the step. `parse` option tells the builder to parse the layer, relative to the `template_root` into the location specified in `on` option. The module may be parsed multiple times using the `repeat` option using the params.

Parsing may not be the only thing needed to set your project up. Sometimes you may need to run additional setups with command lines. So use the `run` option. That will run the following command from the location annotated in `on` option. Extra environment variables may be passed into the command using the `env` option. You can also use the environment that is created through the `params` option. For example, `${PROJECT_PKG_PATH}` is created through `project_pkg_path`

Then just use the build command

```bash
$ x-crafter build example/go-proj_broken example/go-proj_broken
```

And this is be what we will get 

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


## As an SDK

The components is exported and can be used in any other go project if you want to make use of it. Take a look at 
[here](https://pkg.go.dev/github.com/vchitai/x-crafter) for the document.

[comment]: <> (## Stargazers over time)

[comment]: <> ([![Stargazers over time]&#40;https://starchart.cc/vchitai/x-crafter.svg&#41;]&#40;https://starchart.cc/vchitai/x-crafter&#41;)
