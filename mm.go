package market

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

const maxScanTokenSize = 64 * 1024

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

// Errors returned by failures to read a matrix
var (
	ErrInputScanError  = fmt.Errorf("error while scanning matrix input")
	ErrLineTooLong     = fmt.Errorf("input line exceed maximum length ")
	ErrPrematureEOF    = fmt.Errorf("required header items are missing")
	ErrNoHeader        = fmt.Errorf("missing matrix market header line")
	ErrNotMTX          = fmt.Errorf("input is not a matrix market")
	ErrUnsupportedType = fmt.Errorf("unrecognizable matrix description")
	ErrUnwritable      = fmt.Errorf("unable to write matrix to file")
)

var supported = []mmType{
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldReal, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldReal, mtxSymmetrySymm},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldReal, mtxSymmetrySkew},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldInteger, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldInteger, mtxSymmetrySymm},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldInteger, mtxSymmetrySkew},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetrySymm},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetrySkew},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldReal, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldReal, mtxSymmetrySymm},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldReal, mtxSymmetrySkew},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldInteger, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldInteger, mtxSymmetrySymm},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldInteger, mtxSymmetrySkew},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetrySymm},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetrySkew},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetryHermitian},
	{mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetryHermitian},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldPattern, mtxSymmetryGeneral},
	{mtxObjectMatrix, mtxFormatCoordinate, mtxFieldPattern, mtxSymmetrySymm},
}

type index struct {
	M int
	N int
	L int
}

// MM is a basic matrix interface type for real valued Matrix Market matrices.
type MM interface {
	// mat.Matrix
	scanElement(int, string) error
	ToDense() mat.Matrix
	ToSparse() *sparse.DOK
}

// CMM is a basic matrix interface type for complex valued Matrix Market matrices.
type CMM interface {
	// mat.CMatrix
	scanElement(int, string) error
	ToDense() mat.CMatrix
}

type mmType struct {
	Object   string
	Format   string
	Field    string
	Symmetry string
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

func (t *mmType) Bytes() []byte {

	s := fmt.Sprintf(
		"%s %s %s %s %s\n",
		matrixMktBanner,
		t.Object,
		t.Format,
		t.Field,
		t.Symmetry,
	)
	return []byte(s)
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

// counter tallies the number of bytes written to it,
type counter struct {
	total int
}

// Write implements the io.Writer interface.
func (c *counter) Write(p []byte) (int, error) {
	var n int = len(p)
	c.total += n
	return n, nil
}
