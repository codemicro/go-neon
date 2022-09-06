# The Neon intelligent template compiler

Neon compiles your HTML templates into native Go code, much like [valyala/QuickTemplate](https://github.com/valyala/quicktemplate).

However, unlike QuickTemplate, Neon does its own typechecking so you don't have to annotate your templates with type information. Coupled with changes to the syntax and style of templates, Neon templates are much more ergonomic to use compared to QuickTemplate templates.

## Features
* Automatic type inference
* HTML string escaping by default
* Much faster than standard library HTML templating

## Getting started

### Installing `neontc`, the Neon Template Compiler

```
go install github.com/codemicro/go-neon/neontc@latest
```

### Your first template

Now you've got `neontc` installed, check out the [basic example](examples/basic) to get started.

## Benchmarks

The following are the benchmark results from [Neon](examples/benchmark/templates/bench.ntc), [QuickTemplate](https://github.com/valyala/quicktemplate/blob/master/testdata/templates/bench.qtpl) and [`html/template`](https://github.com/valyala/quicktemplate/blob/master/testdata/templates/bench.tpl), all generating the same output.

QuickTemplate is more optimised than Neon, but Neon is still many times faster than the standard library templating tools.

```
BenchmarkNeonTemplate1-8         7974867               145.3 ns/op           560 B/op          3 allocs/op
BenchmarkNeonTemplate10-8        3526574               326.1 ns/op          1200 B/op          4 allocs/op
BenchmarkNeonTemplate100-8        290163              4146 ns/op           18608 B/op         10 allocs/op

BenchmarkQuickTemplate1-8       20372656                51.00 ns/op            0 B/op          0 allocs/op
BenchmarkQuickTemplate10-8       8364355               145.1 ns/op             0 B/op          0 allocs/op
BenchmarkQuickTemplate100-8       747440              1559 ns/op               0 B/op          0 allocs/op

BenchmarkHTMLTemplate1-8         1278934               946.0 ns/op           440 B/op         21 allocs/op
BenchmarkHTMLTemplate10-8         261238              4585 ns/op            1953 B/op        102 allocs/op
BenchmarkHTMLTemplate100-8         26552             45779 ns/op           19243 B/op       1047 allocs/op
```

## License

The `go-neon` project is licensed under the MIT License. See [`LICENSE`](LICENSE) for more information.