# dubins
> Pure Golang implementation of Dubins Path, based on C++ [Dubins-Curves](https://github.com/AndrewWalker/Dubins-Curves)

## Installation
```bash
go get github.com/dohyeunglee/dubins
```

## API
Visit https://pkg.go.dev/github.com/dohyeunglee/dubins

## Example
An example is available in [example/main.go](example/main.go)

## Demo
A demo with plot is available in [demo](demo/README.md)

## Test
```bash
go test
```

## Benchmark
```bash
go test -bench .
```
### Result
Run on M1 Macbook Pro
```bash
BenchmarkMinLengthPath-8   	 1000000	      1126 ns/op
```

## Note
If you want ReedsSheppPath, check https://github.com/dohyeunglee/reedsshepp.

