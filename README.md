# matrix-market

[![Go Reference](https://pkg.go.dev/badge/github.com/wamuir/matrix-market.svg)](https://pkg.go.dev/github.com/wamuir/matrix-market)
[![Build Status](https://travis-ci.com/wamuir/matrix-market.svg?branch=master)](https://travis-ci.com/wamuir/matrix-market)
[![codecov](https://codecov.io/gh/wamuir/matrix-market/branch/master/graph/badge.svg)](https://codecov.io/gh/wamuir/matrix-market)
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

  var m market.COO

  file, err := os.Open("sparse.mtx")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  _, err := m.UnmarshalTextFrom(file)
  if err != nil {
      log.Fatal(err)
  }

  var c *sparse.COO = m.ToCOO()  // github.com/james-bowman/sparse

```

# Supported Formats

## Sparse Matrices (Coordinate Format)

#### Sparse Real-Valued Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Coordinate | Real    | General        | *Yes*     | [market.COO](https://godoc.org/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://godoc.org/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Real    | Skew-Symmetric | Planned   |                                                                           |                                                                          |
| Matrix | Coordinate | Real    | Symmetric      | Planned   |                                                                           |                                                                          |

#### Sparse Integer-Valued Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Coordinate | Integer | General        | *Yes*     | [market.COO](https://godoc.org/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://godoc.org/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Integer | Skew-Symmetric | Planned   |
| Matrix | Coordinate | Integer | Symmetric      | Planned   |

#### Sparse Complex-Valued Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Coordinate | Complex | General        | *Yes*     | [market.CDense](https://godoc.org/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://godoc.org/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Coordinate | Complex | Hermitian      | Planned   |                                                                           |                                                                          |
| Matrix | Coordinate | Complex | Skew-Symmetric | Planned   |                                                                           |                                                                          |
| Matrix | Coordinate | Complex | Symmetric      | Planned   |                                                                           |                                                                          |

#### Sparse Pattern Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Coordinate | Pattern | General        | *Yes*     | [market.COO](https://godoc.org/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://godoc.org/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Pattern | Symmetric      | Planned   |                                                                           |                                                                          |


## Dense Matrices (Array Format)

#### Dense Real-Valued Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Array      | Real    | General        | *Yes*     | [market.Dense](https://godoc.org/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://godoc.org/gonum.org/v1/gonum/mat#Dense)              |
| Matrix | Array      | Real    | Skew-Symmetric | Planned   |                                                                           |                                                                          |
| Matrix | Array      | Real    | Symmetric      | Planned   |                                                                           |                                                                          |

#### Dense Integer-Valued Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Array      | Integer | General        | *Yes*     | [market.Dense](https://godoc.org/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://godoc.org/gonum.org/v1/gonum/mat#Dense)              |
| Matrix | Array      | Integer | Skew-Symmetric | Planned   |                                                                           |                                                                          |
| Matrix | Array      | Integer | Symmetric      | Planned   |                                                                           |                                                                          |

#### Dense Complex-Valued Matrices
| Object | Format     | Field   | Symmetry       | Supported | Concrete Type                                                             | Storage                                                                  |
| ------ | ---------- | ------- | -------------- | :-------: | :-----------------------------------------------------------------------: | :----------------------------------------------------------------------: |
| Matrix | Array      | Complex | General        | *Yes*     | [market.CDense](https://godoc.org/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://godoc.org/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Array      | Complex | Hermitian      | Planned   |                                                                           |                                                                          |
| Matrix | Array      | Complex | Skew-Symmetric | Planned   |                                                                           |                                                                          |
| Matrix | Array      | Complex | Symmetric      | Planned   |                                                                           |                                                                          |




# See also

- [github.com/gonum/gonum](https://github.com/gonum/gonum)
- [github.com/james-bowman/sparse](https://github.com/james-bowman/sparse)
- [Matrix Market Exchange Formats](https://math.nist.gov/MatrixMarket/formats.html#MMformat)
