package market

import (
	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

const matrixMktBanner = `%%MatrixMarket`

const (
	MM_MTX_STR = "matrix"

	MM_ARRAY_STR      = "array"
	MM_COORDINATE_STR = "coordinate"
	MM_DENSE_STR      = "array"
	MM_SPARSE_STR     = "coordinate"

	MM_COMPLEX_STR = "complex"
	MM_INT_STR     = "integer"
	MM_PATTERN_STR = "pattern"
	MM_REAL_STR    = "real"

	MM_GENERAL_STR = "general"
	MM_HERM_STR    = "hermitian"
	MM_SKEW_STR    = "skew-symmetric"
	MM_SYMM_STR    = "symmetric"
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

func (h *header) isMatrix() bool     { return h.Object == MM_MTX_STR }
func (h *header) isArray() bool      { return h.Format == MM_ARRAY_STR }
func (h *header) isCoordinate() bool { return h.Format == MM_COORDINATE_STR }
func (h *header) isDense() bool      { return h.Format == MM_DENSE_STR }
func (h *header) isSparse() bool     { return h.Format == MM_SPARSE_STR }
func (h *header) isComplex() bool    { return h.Field == MM_COMPLEX_STR }
func (h *header) isInteger() bool    { return h.Field == MM_INT_STR }
func (h *header) isPattern() bool    { return h.Field == MM_PATTERN_STR }
func (h *header) isReal() bool       { return h.Field == MM_REAL_STR }
func (h *header) isGeneral() bool    { return h.Symmetry == MM_GENERAL_STR }
func (h *header) isHermitian() bool  { return h.Symmetry == MM_HERM_STR }
func (h *header) isSkew() bool       { return h.Symmetry == MM_SKEW_STR }
func (h *header) isSymmetric() bool  { return h.Symmetry == MM_SYMM_STR }

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
