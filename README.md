# matrix-market

[![GoDoc Reference](https://godoc.org/github.com/wamuir/matrix-market?status.svg)](http://godoc.org/github.com/wamuir/matrix-market)
[![Build Status](https://travis-ci.org/wamuir/matrix-market.svg?branch=master)](https://travis-ci.org/wamuir/matrix-market)
[![Go Report Card](https://goreportcard.com/badge/github.com/wamuir/matrix-market)](https://goreportcard.com/report/github.com/wamuir/matrix-market)

# Description

Go module to read matrices from files in the [NIST Matrix Market native exchange
format](https://math.nist.gov/MatrixMarket/formats.html#MMformat). The
[Matrix Market](https://math.nist.gov/MatrixMarket/) is a service of the
Mathematical and Computational Sciences Division of the Information
Technology Laboratory of the National Institute of Standards and Technology
(NIST), containing "test data for use in comparative studies of algorithms
for numerical linear algebra, featuring nearly 500 sparse matrices from a
variety of applications, as well as matrix generation tools and services."
The Matrix Market native exchange format has become a standard for
exchanging matrix data.

# Installation

  go get -u github.com/wamuir/matrix-market

# Usage

```go
  file, err := os.Open("matrix.mtx")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  mtx, err := market.Read(file)
  if err != nil {
      log.Fatal(err)
  }

  var dok sparse.DOK = mtx.ToSparse() // github.com/james-bowman/sparse

  var arr mat.Matrix = mtx.ToDense()  // gonum.org/v1/gonum/mat
```

# See also

- [github.com/gonum/gonum](https://github.com/gonum/gonum)
- [github.com/james-bowman/sparse](https://github.com/james-bowman/sparse)
- [Matrix Market Exchange Formats](https://math.nist.gov/MatrixMarket/formats.html#MMformat)
