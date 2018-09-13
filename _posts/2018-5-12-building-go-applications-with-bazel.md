---
layout: post
title:  "Building Go Applications with Bazel"
categories: golang bazel
medium: https://medium.com
---

Let's admit it - managing dependencies and building binaries is quite possibly the most frustrating and least fulfilling part of software development. To make matters worse, these frustrations only compound as your application grows - resulting in a hornet's nest of bash scripts and fatalist debug instructions like "just clean the entire project and re-download all dependencies."

In the last few years, a variety of tools which aim to break down this complexity have been developed and open sourced. Facebook's [`buck`](https://github.com/facebook/buck) and Foursquare's [`pants`](https://www.pantsbuild.org/) are two popular derivatives of Google's internal `blaze` build tool, components of which have been open sourced under the [`bazel`](https://bazel.build/) project. These tools require a pretty hefty upfront investment to get setup, but will quickly pay dividends time and time over again as your builds are faster, easier to debug, and consistent across operating systems and architectures.

This post presents `bazel` as a viable alternative to the native `go` toolchains and walks through the process of setting up and using `bazel` to build a real-world application.

## The Benefits of `Bazel`

### Fast and Reproducible Builds

The core selling point of `bazel` is that, if set up correctly, your application's build process is guaranteed to be completely reproducible and consistent - meaning no more afternoons wasted trying to figure out why the code you wrote behaves differently in CI / your boss's laptop / prod (!). On top of this ambitious promise, `bazel` also takes strides to make your builds faster, spreading work across all of your machine's processing power, and ensuring that only the necessary files are rebuilt between runs.

### Language Agnostic, Extensible Tooling

`bazel` can be used for more than just `go` projects, and is configured via the powerful [`skylark`](https://docs.bazel.build/versions/master/skylark/language.html) language, a breath of fresh air for accustomed to hacking together bespoke `bash` scripts for every repository. Beyond just being able to build code, `bazel` can also be used to manage more complex workflows, such as [building and pushing docker containers](https://github.com/bazelbuild/rules_docker) and even [integration tests](https://github.com/bazelbuild/rules_webtesting).

### Consistent UX

One of the hardest parts of developing and maintaining a suite of projects spanning multiple programming languages is the constant burden of context switching between the different frameworks and toolchains. `bazel` attempts to solve this issue by providing a consistent and familiar user experience and workflow, no matter if you are building a Javascript web app or a fleet of Scala microservices. For many projects, developers can hit the ground running with only two commands, `bazel build` and `bazel test`.

## Setting up Bazel for Go

This post details setting up `bazel` for the popular [`groupcache`](github.com/golang/groupcache) project. If you want to follow along or reference this project later, you can check out the code on [github](https://github.com/brendanjryan/groupcache-bazel).

The first step to setting up a `bazel` repo is creating what is known as a `WORKSPACE` file. This file contains a manifest of all of you external dependencies and `bazel` libraries. Our project will build the `groupcache` binary and then package it into a Docker container for other developers to use. As such, our `WORKSPACE` file will look something like this:

```python
# download go bazel tools
http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.11.0/rules_go-0.11.0.tar.gz",
    sha256 = "f70c35a8c779bb92f7521ecb5a1c6604e9c3edd431e50b6376d7497abc8ad3c1",
)
# download the gazelle tool
http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.11.0/bazel-gazelle-0.11.0.tar.gz",
    sha256 = "92a3c59734dad2ef85dc731dbcb2bc23c4568cded79d4b87ebccd787eb89e8d0",
)

# load go rules
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()

# load gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

# load go docker rules
load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)
_go_image_repos()

# external dependencies

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    tag = "v1.0.0",
)
```

The syntax of this file should be very familiar with those who have written `python` before - the `skylark` language is essentially a pared down version of the `python` language.

In addition to creating a `WORKSPACE` file, each directory or "package" in a `bazel` project needs to have a `BUILD.bazel` file. These files declare how to build and test each package, along with any dependencies and additional tasks.

`BUILD` files are the lifeblood of `bazel` - but they are also a huge pain to initially write and then keep up to date as dependencies change - especially for `go` programmers who are used to just declaring `import ()` blocks and having the compiler figure out all of the semantics for you. Luckily, the `bazel` team has recognized this pain point and has written a nifty tool called `gazelle` which can completely automate this process for you! For the sake of brevity (and sanity), the rest of this walkthrough will use the `gazelle` tool, something which I strongly recommend you adopt in your own projects as well.

## Scaffolding Dependencies with `gazelle`

To use `gazelle` and generate `BUILD` files for your project you must first create `BUILD.bazel` file in the root of your repo and configure the `gazelle` tool.

```python
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    # you project name here!
    prefix = "github.com/brendanjryan/groupcache-bazel",
)
```

After this brief setup, invoking `gazelle` is simple straightforward - just `"run"` the job via `bazel`.

```bash
$ bazel run //:gazelle
```
{:.code-bg-yellow}

That's it! You should now see `BUILD` files in each package of your project. Take a few minutes to check these out and bask in the power of `gazelle`.

## Building your application

Now that we have set up our `BUILD` files, the process of building our application is extremely straightforward. By running commands of the form `bazel build <target>`, you can build any package or target declared in your project.

```bash
bazel build //lru/...
INFO: Analysed 2 targets (3 packages loaded).
INFO: Found 2 targets...
INFO: Elapsed time: 0.628s, Critical Path: 0.04s
INFO: Build completed successfully, 1 total action
````

**N.B. In `bazel`'s vernacular `//` denotes the "root" of your project and `...` denotes all "child" packages of the specified package. For example, the command `bazel build //lr/...` will build the `lru` package and all sub-packages underneath it.**

Note that subsequent builds of the same target should be significantly faster:

```bash
bazel build //lru/...
INFO: Analysed 2 targets (0 packages loaded).
INFO: Found 2 targets...
INFO: Elapsed time: 0.268s, Critical Path: 0.01s
INFO: Build completed successfully, 1 total action
```

If you want to build the _entire_ project, you can run the following command - note the _significant_ speedups gained from using `bazel`.

```bash
bazel build //...
INFO: Analysed 19 targets (64 packages loaded).
INFO: Found 19 targets...
INFO: Elapsed time: 8.206s, Critical Path: 3.24s
INFO: Build completed successfully, 35 total actions

bazel build //...
zsh: correct '//...' to '//..' [nyae]? n
INFO: Analysed 19 targets (0 packages loaded).
INFO: Found 19 targets...
INFO: Elapsed time: 0.382s, Critical Path: 0.00s
INFO: Build completed successfully, 1 total action
```

## Testing your applications

Under the hood `bazel` runs tests using the same `go test` tools that you should be familiar with but exposes them under the same `bazel <command> <taget>` pattern used by the `build` process.

For example, to test the `consistenthash` package you would run:

```bash
bazel test //consistenthash/...
INFO: Analysed 2 targets (0 packages loaded).
INFO: Found 1 target and 1 test target...
INFO: Elapsed time: 0.502s, Critical Path: 0.15s
INFO: Build completed successfully, 2 total actions

Executed 1 out of 1 test: 1 test passes.
```

And to test the entire project:

```bash
bazel test //...
INFO: Analysed 19 targets (0 packages loaded).
INFO: Found 15 targets and 4 test targets...
INFO: Elapsed time: 1.733s, Critical Path: 0.91s
INFO: Build completed successfully, 4 total actions

Executed 4 out of 4 tests: 4 tests pass.
```

Note that we get the same benefits of cached results as we do with `bazel build`

```bash
bazel test //...
INFO: Analysed 19 targets (0 packages loaded).
INFO: Found 15 targets and 4 test targets...
INFO: Elapsed time: 0.381s, Critical Path: 0.00s
INFO: Build completed successfully, 1 total action

Executed 0 out of 4 tests: 4 tests pass.
```

The `bazel` testrunner also provides additional functionality on top of `go test` - for instance you can pass the `--runs-per-test` flag to run your suite multiple times in parallel -- very useful for catching flaky tests and data races between test runs.

```bash
bazel test --runs_per_test=10 //...
INFO: Analysed 19 targets (0 packages loaded).
INFO: Found 15 targets and 4 test targets...
INFO: Elapsed time: 7.456s, Critical Path: 1.10s
INFO: Build completed successfully, 41 total actions

Executed 4 out of 4 tests: 4 tests pass.
```

## Packaging your application

Now that we've gotten our project building with `bazel` - publishing the final binary as a `docker` container is surprisingly little work. To do so, we just declare each of the layers of the final image and then how and where the image will be published, like so:

```python
# load bazel rules for docker images
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push", "container_image")

# declare the base `go` image - this is the same format as the standard
# `go_binary` rule.
go_image(
    name = "groupcache_image_base",
    embed = [":go_default_library"],
)

# wrapper image used to expose ports to the underlying go_image
container_image(
    name = "groupcache_image",
    base = ":groupcache_image_base",
    ports = ["8080"],
)

# declare where and how the image will be published
container_push(
    name = "push",
    format = "Docker",
    image = ":groupcache_image",
    registry = "index.docker.io",
    repository = "brendanjryan/groupcache-bazel",
    tag = "master",  # don't use this on production image :)
)
```

One of the strengths of this process over the standard docker workflow is that no `Dockerfiles` are required and you can easily build and publish multiple images to multiple repositories - all in parallel!

In our case, pushing our image up to `Dockerhub` is as simple as:

```bash
$ bazel run //example:push
```
{:.code-bg-yellow}

_Caveat: I do not recommend pushing images from your local workstation. This step should be part of your CI workflow_.

## Final Words

Hopefully this walkthrough gives you enough to start integrating `bazel` into one of your own `go` projects - or conversely know that you never want to :)

Feel free to reach out on [Twitter](https://twitter.com/Brendan_J_Ryan) or [Github](https://github.com/brendanjryan/groupcache-bazel) if you have any questions!

## Further Readings

Want to learn more? Here are a few great links:

* [Official `Bazel` documentation](https://golang.org/pkg/testing/)
* [`Bazel` rules for `go`](https://github.com/bazelbuild/rules_go)
* [`bazel-gazelle` - used for generating BUILD files](https://github.com/bazelbuild/bazel-gazelle)
* [Golang UK - Building `Go` with `Bazel`](https://www.youtube.com/watch?v=2TKxuERTnks)
