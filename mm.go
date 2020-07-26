package market

import (
	"bufio"
	"fmt"
	"strings"

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

type mmType struct {
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

func (t *mmType) isMatrix() bool     { return t.Object == mtxObjectMatrix }
func (t *mmType) isArray() bool      { return t.Format == mtxFormatArray }
func (t *mmType) isCoordinate() bool { return t.Format == mtxFormatCoordinate }
func (t *mmType) isDense() bool      { return t.Format == mtxFormatDense }
func (t *mmType) isSparse() bool     { return t.Format == mtxFormatSparse }
func (t *mmType) isComplex() bool    { return t.Field == mtxFieldComplex }
func (t *mmType) isInteger() bool    { return t.Field == mtxFieldInteger }
func (t *mmType) isPattern() bool    { return t.Field == mtxFieldPattern }
func (t *mmType) isReal() bool       { return t.Field == mtxFieldReal }
func (t *mmType) isGeneral() bool    { return t.Symmetry == mtxSymmetryGeneral }
func (t *mmType) isHermitian() bool  { return t.Symmetry == mtxSymmetryHermitian }
func (t *mmType) isSkew() bool       { return t.Symmetry == mtxSymmetrySkew }
func (t *mmType) isSymmetric() bool  { return t.Symmetry == mtxSymmetrySymm }

// isMMType tests equality of two Matrix Market headers
func (t *mmType) isMMType(t2 *mmType) bool {

	if strings.ToLower(t.Object) != t2.Object {
		return false
	}
	if strings.ToLower(t.Format) != t2.Format {
		return false
	}
	if strings.ToLower(t.Field) != t2.Field {
		return false
	}
	if strings.ToLower(t.Symmetry) != t2.Symmetry {
		return false
	}
	return true
}

// isSupported reports if receiver is among supported Matrix Market types,
// based on comparison against object, format, field and symmetry.
func (t *mmType) isSupported() bool {

	for _, t2 := range supported {
		if t.isMMType(&t2) {
			return true
		}
	}

	return false
}

func scanHeader(scanner *bufio.Scanner) (*mmType, error) {

	var (
		banner string
		t      mmType
	)

	if ok := scanner.Scan(); !ok {
		return nil, ErrInputScanError
	}

	_, err := fmt.Sscan(scanner.Text(), &banner, &t.Object, &t.Format, &t.Field, &t.Symmetry)
	if err != nil {
		return nil, ErrPrematureEOF
	}

	if banner != matrixMktBanner {
		return nil, ErrNoHeader
	}

	if !(t.isSupported()) {
		return nil, ErrUnsupportedType
	}

	return &t, nil
}
