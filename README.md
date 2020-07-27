# matrix-market

[![GoDoc Reference](https://godoc.org/github.com/wamuir/matrix-market?status.svg)](http://godoc.org/github.com/wamuir/matrix-market)
[![Build Status](https://travis-ci.org/wamuir/matrix-market.svg?branch=master)](https://travis-ci.org/wamuir/matrix-market)
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

## Sparse Matrices (Coordinate Format)

#### Sparse Real-Valued Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Coordinate | Real    | General        |
| Matrix | Coordinate | Real    | Skew-Symmetric |
| Matrix | Coordinate | Real    | Symmetric      |

#### Sparse Integer Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Coordinate | Integer | General        |
| Matrix | Coordinate | Integer | Skew-Symmetric |
| Matrix | Coordinate | Integer | Symmetric      |

#### Sparse Complex-Valued Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Coordinate | Complex | General        |
| Matrix | Coordinate | Complex | Hermitian      |
| Matrix | Coordinate | Complex | Skew-Symmetric |
| Matrix | Coordinate | Complex | Symmetric      |

#### Sparse Pattern Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Coordinate | Pattern | General        |
| Matrix | Coordinate | Pattern | Symmetric      |


## Dense Matrices (Array Format)

#### Dense Real-Valued Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Array      | Real    | General        |
| Matrix | Array      | Real    | Skew-Symmetric |
| Matrix | Array      | Real    | Symmetric      |

#### Dense Integer Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Array      | Integer | General        |
| Matrix | Array      | Integer | Skew-Symmetric |
| Matrix | Array      | Integer | Symmetric      |

#### Dense Complex-Valued Matrices
| Object | Format     | Field   | Symmetry       |
| ------ | ---------- | ------- | -------------- |
| Matrix | Array      | Complex | General        |
| Matrix | Array      | Complex | Hermitian      |
| Matrix | Array      | Complex | Skew-Symmetric |
| Matrix | Array      | Complex | Symmetric      |


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

# See also

- [github.com/gonum/gonum](https://github.com/gonum/gonum)
- [github.com/james-bowman/sparse](https://github.com/james-bowman/sparse)
- [Matrix Market Exchange Formats](https://math.nist.gov/MatrixMarket/formats.html#MMformat)
