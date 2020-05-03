---
layout: post
title: "Advanced Go Testing Techniques"
tags: golang testing
redirect_from: /golang/testing/2018/03/23/advanced-go-testing-techniques
medium: https://medium.com
---

Writing tests in `go` is extremely simple and straightforward, mostly due to the extremely powerful testing library and test runner (`go test`) which come bundled with the standard library. While the [`testing` ](https://golang.org/pkg/testing/) docs do a very good job of outlining basic applications - there are many tricks and design patterns that you can employ to make your tests even faster to write, easier to debug, and simpler to maintain. This article outlines some of these more "advanced" testing techniques and provides examples inspired by real-world use cases.

## Table Driven Tests and Sub-Tests

Table driven tests are one of the most expressive design patterns in `golang`. By leveraging composite literals and anonymous `struct`s, table driven tests separate test inputs from their logic and allow you to modify or extend your test suite with ease.

You can gain a little more control over the execution of your tests cases by splitting the test into multiple "sub-tests" via the [`T.Run`](https://golang.org/pkg/testing/#T.Run) method. For easier debugging and better reporting, each sub-test is given its own name and identifier. You can use these identifiers when running individual sub-tests via the `go test` tool - useful for hunting down that one failing or flaky test case. As an added bonus, you can run each scenario of a test suite in parallel - a huge win for suites of large integration tests.

An example of this technique is as follows:

```golang
func TestfmtStatStr(t *testing.T){
	tests := []struct{
		msg string
		stat string
		tags map[string]string
		exp string
	}{
		{
			"empty string and tags",
			"",
			nil,
			"",
		},
		{
			"empty tags",
			"handler.latency",
			nil,
			"foo",
		},
		{
			"with tags",
			"handler.latency",
			map[string]string{
				"host": "aws_",
				"service": "users",
			},
			"foo",
		},
	}

	for _  , tt := range tests {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
		    t.Parallel() // run sub-tests in parallel
		    res := fmtStatStr(tt.st, tt.tags)
		    if tt.exp != res {
		        t.Error("exp:", tt.exp, "got:", res)
		    }
		})
	}
}
```

## External Test Packages

When developing an API designed to be consumed by other engineers, it is important that you exercise the interfaces and behaviors of your package in the same ways that you anticipate your end user to. `Go`'s package-based file hierarchy can make it too easy to tests for public methods alongside their private counterparts makes it hard to distinguish between the two and allows you to take shortcuts that those who vendor your library won't be able to.

By creating a separate `_test` package which sits alongside a package, you are able to separate your tests for private functionality from those for public methods and interfaces. Although this pattern is a little more work upfront and adds additional complexity to your project's layout - it will ultimately result in friendlier more testable APIs.

In practice, your file structure will look something like this:

```
- src
    - client
        - client.go
        - client_test.go
    - client_test <-- test package for public interfaces in 'client'
        - client_test.go
```

## Time Multipliers

Managing clocks and timeouts is one of the hardest parts of any language, and `go` is no exception. A popular choice is using one of [several](https://github.com/uber-go/ratelimit/tree/master/internal/clock) [excellent](https://github.com/benbjohnson/clock) open source libraries which provide a "clock" interface for programmatically managing and manipulating time in tests. This tactic gives you a high degree of accuracy and control, but adds an additional API to your code and muddles your business logic with testing constructs.

A simple way of solving this issue is by adding a "time multiplier" variable to each of your critical timeouts and intervals. By manipulating this variable, you are able to effectively control the flow of time in your tests and make them run faster or slower at will.

The pattern is particularly useful when testing the `time.Ticker` object - as illustrated below:

```golang
var (
	timeMultiplier = 1.0
	batchMaxAge = time.Second
	batchSize = 25
)

// flushEvents listens for new events on the `events` channel and sends
// them to the events client in batches of size `batchSize`.
func flushEvents(events <- chan string) {

	var batch []string
	select {

	case e := <- events:
		batch = append(batch, e)

		if len(batch) >= batchSize {
			log.Println("flushing batch of events")
			batch = nil
		}

	case <-time.Tick(time.Duration(timeMultiplier) * batchMaxAge):
		log.Println("flushing batch of stale events")
		batch = nil
	}
}
```

## Embedding Types

When testing data structures which have to manage a lot of internal state, you may find yourself adding "helper methods" or additional functionality to make your tests easier to run or debug. Instead of polluting your business logic with artifacts meant only for testing, you can instead add test-only "wrapper types" which provide additional functionality without changing the underlying interface. By keeping these types within `_test.go` files you can have all the benefits of accessing private structs and methods without any of the interface bloat.

A common use case for this pattern is "setting up" and "tearing down" datastore resources like thus:

```golang
// User represents a user in our application.
type User struct {
	Id int64
	FirstName string
	LastName string
}

// UserMutator defines a type which is able to mutate and persist User data.
type UserMutator interface {
	Fetch(int64) (*User, error)
	Create(*User) (*User, error)
	Update(int64, *User) (*User, error)
	Delete(int64) error
}

// UserManager handles datastore operations for 'user' data.
type UserManager struct {
	// database tables, caches, etc...
}

// Fetch retrieves a single user model keyed by id.
func (*UserManager) Fetch(id int64) (*User, error) {
	return nil, nil
}

// Create instantiates a new user model
func (*UserManager) Create(u *User) (*User, error) {
	return nil, nil
}

// Update updates the user keyed by the provided id to a new state.
func (*UserManager) Update(id int64, u *User) (*User, error) {
	return nil, nil
}

// Delete deletes the user keyed by the provided id.
func (*UserManager) Delete(id int64) error {
	return nil
}

// UserManagerCreator wraps a UserManager and provides additional
// functionality for setting up and tearing down this manager's underlying
// datastores.
type UserManagerCreator struct {
	UserManager
}

// Setup instantiates resources used by this manager.
func (u *UserManagerCreator) Setup() error {
	return nil
}

// Teardown cleans up resources used by this manager.
func (u *UserManagerCreator) Teardown() error {
	return nil
}
```

## Further Readings

Want to learn more? Here are a few great links:

- [`testing` package documentation](https://golang.org/pkg/testing/)
- [Testing Techniques - Andrew Gerrand](https://talks.golang.org/2014/testing.slide)
- [Writing table driven tests](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go)
