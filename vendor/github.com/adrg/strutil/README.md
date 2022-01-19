<h1 align="center">strutil</h1>

<p align="center">
    <a href="https://github.com/adrg/strutil/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/adrg/strutil/workflows/CI/badge.svg" />
    </a>
    <a href="https://codecov.io/gh/adrg/strutil">
        <img alt="Code coverage" src="https://codecov.io/gh/adrg/strutil/branch/master/graphs/badge.svg?branch=master" />
    </a>
    <a href="https://pkg.go.dev/github.com/adrg/strutil">
        <img alt="pkg.go.dev documentation" src="https://pkg.go.dev/badge/github.com/adrg/strutil" />
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/adrg/strutil" />
    </a>
    <a href="https://goreportcard.com/report/github.com/adrg/strutil">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/adrg/strutil" />
    </a>
    <a href="https://github.com/adrg/strutil/issues">
        <img alt="GitHub issues" src="https://img.shields.io/github/issues/adrg/strutil" />
    </a>
    <a href="https://ko-fi.com/T6T72WATK">
        <img alt="Buy me a coffee" src="https://img.shields.io/static/v1.svg?label=%20&message=Buy%20me%20a%20coffee&color=579fbf&logo=buy%20me%20a%20coffee&logoColor=white" />
    </a>
</p>

strutil provides a collection of string metrics for calculating string similarity as well as
other string utility functions.  
Full documentation can be found at https://pkg.go.dev/github.com/adrg/strutil.

## Installation

```
go get github.com/adrg/strutil
```

## String metrics

- [Hamming](#hamming)
- [Levenshtein](#levenshtein)
- [Jaro](#jaro)
- [Jaro-Winkler](#jaro-winkler)
- [Smith-Waterman-Gotoh](#smith-waterman-gotoh)
- [Sorensen-Dice](#sorensen-dice)
- [Jaccard](#jaccard)
- [Overlap Coefficient](#overlap-coefficient)

The package defines the `StringMetric` interface, which is implemented by all
the string metrics. The interface is used with the `Similarity` function, which
calculates the similarity between the specified strings, using the provided
string metric.

```go
type StringMetric interface {
    Compare(a, b string) float64
}

func Similarity(a, b string, metric StringMetric) float64 {
}
```

All defined string metrics can be found in the
[metrics](https://pkg.go.dev/github.com/adrg/strutil/metrics) package.

#### Hamming

Calculate similarity.
```go
similarity := strutil.Similarity("text", "test", metrics.NewHamming())
fmt.Printf("%.2f\n", similarity) // Output: 0.75
```

Calculate distance.
```go
ham := metrics.NewHamming()
fmt.Printf("%d\n", ham.Distance("one", "once")) // Output: 2
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#Hamming).

#### Levenshtein

Calculate similarity using default options.
```go
similarity := strutil.Similarity("graph", "giraffe", metrics.NewLevenshtein())
fmt.Printf("%.2f\n", similarity) // Output: 0.43
```

Configure edit operation costs.
```go
lev := metrics.NewLevenshtein()
lev.CaseSensitive = false
lev.InsertCost = 1
lev.ReplaceCost = 2
lev.DeleteCost = 1

similarity := strutil.Similarity("make", "Cake", lev)
fmt.Printf("%.2f\n", similarity) // Output: 0.50
```

Calculate distance.
```go
lev := metrics.NewLevenshtein()
fmt.Printf("%d\n", lev.Distance("graph", "giraffe")) // Output: 4
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#Levenshtein).

#### Jaro

```go
similarity := strutil.Similarity("think", "tank", metrics.NewJaro())
fmt.Printf("%.2f\n", similarity) // Output: 0.78
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#Jaro).

#### Jaro-Winkler

```go
similarity := strutil.Similarity("think", "tank", metrics.NewJaroWinkler())
fmt.Printf("%.2f\n", similarity) // Output: 0.80
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#JaroWinkler).

#### Smith-Waterman-Gotoh

Calculate similarity using default options.
```go
swg := metrics.NewSmithWatermanGotoh()
similarity := strutil.Similarity("times roman", "times new roman", swg)
fmt.Printf("%.2f\n", similarity) // Output: 0.82
```

Customize gap penalty and substitution function.
```go
swg := metrics.NewSmithWatermanGotoh()
swg.CaseSensitive = false
swg.GapPenalty = -0.1
swg.Substitution = metrics.MatchMismatch {
    Match:    1,
    Mismatch: -0.5,
}

similarity := strutil.Similarity("Times Roman", "times new roman", swg)
fmt.Printf("%.2f\n", similarity) // Output: 0.96
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#SmithWatermanGotoh).

#### Sorensen-Dice

Calculate similarity using default options.
```go
sd := metrics.NewSorensenDice()
similarity := strutil.Similarity("time to make haste", "no time to waste", sd)
fmt.Printf("%.2f\n", similarity) // Output: 0.62
```

Customize n-gram size.
```go
sd := metrics.NewSorensenDice()
sd.CaseSensitive = false
sd.NgramSize = 3

similarity := strutil.Similarity("Time to make haste", "no time to waste", sd)
fmt.Printf("%.2f\n", similarity) // Output: 0.53
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#SorensenDice).

#### Jaccard

Calculate similarity using default options.
```go
j := metrics.NewJaccard()
similarity := strutil.Similarity("time to make haste", "no time to waste", j)
fmt.Printf("%.2f\n", similarity) // Output: 0.45
```

Customize n-gram size.
```go
j := metrics.NewJaccard()
j.CaseSensitive = false
j.NgramSize = 3

similarity := strutil.Similarity("Time to make haste", "no time to waste", j)
fmt.Printf("%.2f\n", similarity) // Output: 0.36
```

The input of the Sorensen-Dice example is the same as the one of Jaccard
because the metrics bear a resemblance to each other. In fact, each of the
coefficients can be used to calculate the other one.

Sorensen-Dice to Jaccard.
```
J = SD/(2-SD)

where SD is the Sorensen-Dice coefficient and J is the Jaccard index.
```

Jaccard to Sorensen-Dice.
```
SD = 2*J/(1+J)

where SD is the Sorensen-Dice coefficient and J is the Jaccard index.
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#Jaccard).

#### Overlap Coefficient

Calculate similarity using default options.
```go
oc := metrics.NewOverlapCoefficient()
similarity := strutil.Similarity("time to make haste", "no time to waste", oc)
fmt.Printf("%.2f\n", similarity) // Output: 0.67
```

Customize n-gram size.
```go
oc := metrics.NewOverlapCoefficient()
oc.CaseSensitive = false
oc.NgramSize = 3

similarity := strutil.Similarity("Time to make haste", "no time to waste", oc)
fmt.Printf("%.2f\n", similarity) // Output: 0.57
```

More information and additional examples can be found on
[pkg.go.dev](https://pkg.go.dev/github.com/adrg/strutil/metrics#OverlapCoefficient).

## References

For more information see:
- [Hamming distance](https://en.wikipedia.org/wiki/Hamming_distance)
- [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
- [Jaro-Winkler distance](https://en.wikipedia.org/wiki/Jaro-Winkler_distance)
- [Smith-Waterman algorithm](https://en.wikipedia.org/wiki/Smith-Waterman_algorithm)
- [Sorensen-Dice coefficient](https://en.wikipedia.org/wiki/Sorensen–Dice_coefficient)
- [Jaccard index](https://en.wikipedia.org/wiki/Jaccard_index)
- [Overlap coefficient](https://en.wikipedia.org/wiki/Overlap_coefficient)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/adrg/strutil.svg)](https://starchart.cc/adrg/strutil)

## Contributing

Contributions in the form of pull requests, issues or just general feedback,
are always welcome.  
See [CONTRIBUTING.MD](CONTRIBUTING.md).

## License

Copyright (c) 2019 Adrian-George Bostan.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
