---
layout: post
title:  "Application Tuning in go: Benchmarks"
medium: https://medium.com
---

One of the best parts about writing and maintaining software with `go` is that it comes with "batteries included" for developing scalable applications. Beyond just verifying the correctness of your code using go's `testing` package, you can just as easily measure and track the _performance_ of your code, using more or less the same paradigms and constructs. To illustrate the power and simplicity of benchmarks, this post walks through an example, real-world, tuning scenario.

## Our sample program

For our exercise we will be tuning the performance of a [`statsd`](https://github.com/etsy/statsd) client that our hypothetical web server relies on. Our server calls out to this client multiple times when serving each request, making it a very "hot" region of code. By using [`pprof`](https://golang.org/pkg/runtime/pprof/), we have identified that we are spending a **lot** of CPU time in one function - a small utility for adding tags to a stat string as per the [DogStatsD](https://docs.datadoghq.com/developers/dogstatsd/) spec.

```golang
// fmtStatString takes a preformatted statting string and appends
// a flattened map of tags to it.
func fmtStatStr(stat string, tags map[string]string) string {
  parts := make([]string, 0, len(tags))
  for k, v := range tags {
    if v != "" {
      parts = append(parts, fmt.Sprintf("%s:%s", clean(k), clean(v)))
     }
  }

  return fmt.Sprintf("%s|%s", stat, strings.Join(parts, ","))
}
```

## Writing A Benchmark

The process of writing benchmarks should be very familiar to those who have written a functional test in `go` before. The key difference is a different parameter to your testing function - `*testing.B` instead of `*testing.T`, as well as the inclusion of a small looping construct in the body of your test.

Most benchmarks ultimately end up resembling this pattern:

```golang
func exLog() {
  log.Println("bench")
}

func BenchmarkLogging(b *testing.B) {
    for i := 0; i < b.N; i++ {
      exLog()
    }
}
```

Under the hood, the `go test` tool will the code inside of the benchmark loop until a certain time limit is reached, at which point it is simple to derive performance characteristics from the number of iterations that finished and the resource consumption of your application during that time.

In order to get the absolute tightest bound on your application's true performance, you should design your benchmark so that you have only the code you wish to measure inside the benchmark "`for`" loop. By doing any setup outside of the benchmark loop and then calling [`b.ResetTimer`](https://golang.org/pkg/testing/#B.ResetTimer) you will ensure that no allocations or cycles spent during setup will pollute your results.

For our `fmtStatStr` function, we don't need a lot of prep-work and can write a straightforward benchmark:

```golang
func BenchmarkLogging(b *testing.B) {
    for i := 0; i < b.N; i++ {
      fmtStatStr("handler.sample", map[string]string{"os": "osx", "locale": "en-US"})
    }
}
```

## Running your benchmark

You can run your benchmarks via the `go test` tool simply by adding a few flags, like so:

```bash
go test -memprofile -bench=.
```

After running this command and waiting a few seconds you should see something like this output to your terminal:

```bash
BenchmarkLogging-4     1000000       1462 ns/op      584 B/op       14 allocs/op
```

The output can be read from right to left as

- The name of the benchmark function
- The number of iterations the loop was run for.
- The number of nanoseconds of system time each iteration took.
- The number of bytes allocated per iteration.
- The number of dintinct allocations that took place in each operation.

## Analyzing and Iterating

Now that we have a good baseline of our program's performance, we can begin refactoring our code to see how much work we can shave off. There is no single process or magical remedy for all performance issues, but a few quick Google searches and a dive into the relevant `go` docs is probably a good place to start.

After a little research we discovered that using `fmt.Printf` to concatenate and join strings is a fairly CPU and memory intensive process and that using a [`bytes.Buffer`](https://golang.org/pkg/bytes/#NewBufferString) is much more efficient means of achieving the same result. Armed with this knowledge we refactor our function to look like this:

```golang
func fmtStatStrOpt(stat string, tags map[string]string) string {
  b := bytes.NewBufferString(stat)
  b.WriteString("|")

  first := true
  for k, v := range tags {
    if v != "" {
      if !first {
        // do not append a ',' at the
        // beginning of the first iteration.
        b.WriteString(",")
      }

      b.WriteString(k)
      b.WriteString(":")
      b.WriteString(v)

      // automatically set to falsey after first iteration
      first = false
    }
  }

  return b.String()
}
```

Now that we have two identical functions, we can compare their performance side to side and see which is better! To do so we add a new benchmark function and supply the same parameters to both tests.

```golang
var test := struct {
  stat string
  tags map[string]string
}{
  "handler.sample",
  map[string]string{
    "os": "ios",
    "locale": "en-US",
  }
}

func Benchmark_fmtStatStrBefore(b *testing.B) {
  for i := 0; i < b.N; i++ {
    fmtStatStr(test.stat, test.tags)
  }
}
```

To run both of our benchmarks we can run the same command as before:

```bash
go test -benchmem -bench=.
```

Which should output something like:

```bash
Benchmark_fmtStatStrBefore-4     1000000              1333 ns/op             248 B/op         12 allocs/op
Benchmark_fmtStatStrAfter-4      3000000               492 ns/op             224 B/op          4 allocs/op
PASS
```

At first glance our code immediately looks a lot better! Not only is our code faster (lower ns/op), but we are also doing fewer allocations per operation. Not too bad for a few minutes of work!


#### Aside:

Notice how much longer and difficult to read the optimized function is? Now is a good time to meditate on one of the core tradeoffs of optimization - while your code will be faster, it will almost always be more brittle, harder to read, and much more difficult to maintain by other engineers. **This is why it is imperative that you only tune code that needs to be optimized and write tests to verify the original behavior of your code before refactoring it and adding additional complexity for the sake of performance.**

## Caveat Emptor - Benchmarks and the Go compiler

Although writing benchmarks is simple and straightforward for 95% of all use cases, there are some common traps you have to look out for.

### Escape Analysis Tricks

In some situations, your code may be subject to some [compile-time optimizations](https://github.com/golang/go/wiki/CompilerOptimizations) which render your benchmarks invalid.

Consider the following scenario:

```golang
func concat(a, b string) string {
  return a + b
}

func Benchmark_concat(b *testing.B) {
  for i := 0; i < b.N; i++ {
    concat(a, b)
  }
}
```

Since `concat` is a _very_ short function and its return value isn't captured at its callsite, the `go` compiler may happily remove this function call for you - turning your benchmark into a really lame game of for-loop musical chairs. To save yourself from this benign mistake you can trick the compiler into thinking your code is more important than it really is by capturing return values in various scopes - like so:

```golang
// package level scope
var concatRes string

func Benchmark_concat(b *testing.B) {
  // function level scope
  var res string

  for i := 0; i < b.N; i++ {
    // block level scope

    // assign variable outside of block scope so that this
    // loop is not escaped.
    res = concat(a, b)
  }

  // assign final result to a package-level variable to prevent
  // closure level escaping.
  concatRes = res
}
```

If you notice that your benchmarks are running _suspiciously_ fast, you may want to to obfuscate your code and make it a little harder for the compiler to eliminate your slow code for you :)

### Noisy Neighbors

If you are running benchmarks on your local machine it is important that you try and minimize the amount of other work that your computer is doing during the benchmark. You should shut down Slack, XCode, your Bitcoin rig etc.. if you want to have fair and reproducible benchmarks.

## Recommended Reading

Want to learn more? Here are a few great links:

- [`testing` package documentation](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [subtests and subbenchmarks (go blog)](https://blog.golang.org/subtests)
- [go language benchmarks](https://github.com/golang/benchmarks)
