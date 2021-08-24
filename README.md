# matrix-market

![Project Stability: Experimental](https://img.shields.io/badge/stability-experimental-critical.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/wamuir/matrix-market.svg)](https://pkg.go.dev/github.com/wamuir/matrix-market)
[![Build Status](https://github.com/wamuir/matrix-market/actions/workflows/go.yml/badge.svg?branch=master&event=push)](https://github.com/wamuir/matrix-market/actions/workflows/go.yml?query=event%3Apush+branch%3Amaster)
[![codecov](https://codecov.io/gh/wamuir/matrix-market/branch/master/graph/badge.svg)](https://codecov.io/gh/wamuir/matrix-market)
[![Go Report Card](https://goreportcard.com/badge/github.com/wamuir/matrix-market)](https://goreportcard.com/report/github.com/wamuir/matrix-market)


# Installation

```sh
  go get -u github.com/wamuir/matrix-market
```

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

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Coordinate | Real    | General        | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Real    | Skew-Symmetric | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Real    | Symmetric      | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |

#### Sparse Integer-Valued Matrices

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Coordinate | Integer | General        | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Integer | Skew-Symmetric | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Integer | Symmetric      | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |

#### Sparse Complex-Valued Matrices

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Coordinate | Complex | General        | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Coordinate | Complex | Hermitian      | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Coordinate | Complex | Skew-Symmetric | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Coordinate | Complex | Symmetric      | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |

#### Sparse Pattern Matrices

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Coordinate | Pattern | General        | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |
| Matrix | Coordinate | Pattern | Symmetric      | [market.COO](https://pkg.go.dev/github.com/wamuir/matrix-market#COO)       | [sparse.COO](https://pkg.go.dev/github.com/james-bowman/sparse#COO)       |


## Dense Matrices (Array Format)

#### Dense Real-Valued Matrices

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Array      | Real    | General        | [market.Dense](https://pkg.go.dev/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://pkg.go.dev/gonum.org/v1/gonum/mat#Dense)              |
| Matrix | Array      | Real    | Skew-Symmetric | [market.Dense](https://pkg.go.dev/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://pkg.go.dev/gonum.org/v1/gonum/mat#Dense)              |
| Matrix | Array      | Real    | Symmetric      | [market.Dense](https://pkg.go.dev/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://pkg.go.dev/gonum.org/v1/gonum/mat#Dense)              |

#### Dense Integer-Valued Matrices

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Array      | Integer | General        | [market.Dense](https://pkg.go.dev/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://pkg.go.dev/gonum.org/v1/gonum/mat#Dense)              |
| Matrix | Array      | Integer | Skew-Symmetric | [market.Dense](https://pkg.go.dev/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://pkg.go.dev/gonum.org/v1/gonum/mat#Dense)              |
| Matrix | Array      | Integer | Symmetric      | [market.Dense](https://pkg.go.dev/github.com/wamuir/matrix-market#Dense)   | [mat.Dense](https://pkg.go.dev/gonum.org/v1/gonum/mat#Dense)              |

#### Dense Complex-Valued Matrices

| Object | Format     | Field   | Symmetry       | Concrete Type                                                              | Storage                                                                   |
| ------ | ---------- | ------- | -------------- | :------------------------------------------------------------------------: | :-----------------------------------------------------------------------: |
| Matrix | Array      | Complex | General        | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Array      | Complex | Hermitian      | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Array      | Complex | Skew-Symmetric | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
| Matrix | Array      | Complex | Symmetric      | [market.CDense](https://pkg.go.dev/github.com/wamuir/matrix-market#CDense) | [mat.CDense](https://pkg.go.dev/gonum.org/v1/gonum/mat#CDense)            |
