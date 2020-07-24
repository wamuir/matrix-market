package market

import (
	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

const matrixMktBanner = `%%MatrixMarket`

const (
	mtxObjectMatrix = "matrix"

	mtxFormatArray      = "array"
	mtxFormatCoordinate = "coordinate"
	mtxFormatDense      = "array"
	mtxFormatSparse     = "coordinate"

	mtxFieldComplex = "complex"
	mtxFieldInteger = "integer"
	mtxFieldPattern = "pattern"
	mtxFieldReal    = "real"

	mtxSymmetryGeneral   = "general"
	mtxSymmetryHermitian = "hermitian"
	mtxSymmetrySkew      = "skew-symmetric"
	mtxSymmetrySymm      = "symmetric"
)

type header struct {
	Object   string
	Format   string
	Field    string
	Symmetry string
}

type index struct {
	M int
	N int
	L int
}

// CMatrix is a basic matrix interface type for complex matrices.
type CMatrix interface {
	scanElement(int, string) error
	ToDense() mat.CMatrix
}

// Matrix is a basic matrix interface type for matrices.
type Matrix interface {
	scanElement(int, string) error
	ToDense() mat.Matrix
	ToSparse() *sparse.DOK
}

func (h *header) isMatrix() bool     { return h.Object == mtxObjectMatrix }
func (h *header) isArray() bool      { return h.Format == mtxFormatArray }
func (h *header) isCoordinate() bool { return h.Format == mtxFormatCoordinate }
func (h *header) isDense() bool      { return h.Format == mtxFormatDense }
func (h *header) isSparse() bool     { return h.Format == mtxFormatSparse }
func (h *header) isComplex() bool    { return h.Field == mtxFieldComplex }
func (h *header) isInteger() bool    { return h.Field == mtxFieldInteger }
func (h *header) isPattern() bool    { return h.Field == mtxFieldPattern }
func (h *header) isReal() bool       { return h.Field == mtxFieldReal }
func (h *header) isGeneral() bool    { return h.Symmetry == mtxSymmetryGeneral }
func (h *header) isHermitian() bool  { return h.Symmetry == mtxSymmetryHermitian }
func (h *header) isSkew() bool       { return h.Symmetry == mtxSymmetrySkew }
func (h *header) isSymmetric() bool  { return h.Symmetry == mtxSymmetrySymm }

// Equals tests equality of two matrix market headers
func (h *header) equals(j header) bool {

	if j.Object != h.Object {
		return false
	}

	if j.Format != h.Format {
		return false
	}

	if j.Field != h.Field {
		return false
	}

	if j.Symmetry != h.Symmetry {
		return false
	}

	return true
}

func (h *header) isValid() bool {

	for _, j := range supported {
		if h.equals(j) {
			return true
		}
	}

	return false
}
