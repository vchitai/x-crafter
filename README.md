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



[comment]: <> (## Stargazers over time)

[comment]: <> ([![Stargazers over time]&#40;https://starchart.cc/vchitai/x-crafter.svg&#41;]&#40;https://starchart.cc/vchitai/x-crafter&#41;)
