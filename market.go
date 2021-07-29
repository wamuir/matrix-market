package market

import (
	"bufio"
	"fmt"
	"strings"
)

const maxScanTokenSize = 64 * 1024

const matrixMktBanner = `%%MatrixMarket`

const (
	// object
	mtxObjectMatrix = "matrix"

	// format
	mtxFormatArray      = "array"
	mtxFormatCoordinate = "coordinate"
	mtxFormatDense      = "array"
	mtxFormatSparse     = "coordinate"

	// field
	mtxFieldComplex = "complex"
	mtxFieldInteger = "integer"
	mtxFieldPattern = "pattern"
	mtxFieldReal    = "real"

	// symmetry
	mtxSymmetryGeneral   = "general"
	mtxSymmetryHermitian = "hermitian"
	mtxSymmetrySkew      = "skew-symmetric"
	mtxSymmetrySymm      = "symmetric"
)

// Errors returned by failures to read a matrix
var (
	ErrInputScanError  = fmt.Errorf("error while scanning matrix input")
	ErrLineTooLong     = fmt.Errorf("input line exceeds maximum length")
	ErrPrematureEOF    = fmt.Errorf("required header items are missing")
	ErrNoHeader        = fmt.Errorf("missing matrix market header line")
	ErrNotMTX          = fmt.Errorf("input is not a matrix market file")
	ErrUnsupportedType = fmt.Errorf("unrecognizable matrix description")
	ErrUnwritable      = fmt.Errorf("error writing matrix to io writer")
)

var supported = map[int]mmType{
	1:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldReal, mtxSymmetryGeneral},
	2:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldReal, mtxSymmetrySymm},
	3:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldReal, mtxSymmetrySkew},
	4:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldInteger, mtxSymmetryGeneral},
	5:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldInteger, mtxSymmetrySymm},
	6:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldInteger, mtxSymmetrySkew},
	7:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetryGeneral},
	8:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetrySymm},
	9:  {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetrySkew},
	10: {mtxObjectMatrix, mtxFormatArray, mtxFieldReal, mtxSymmetryGeneral},
	11: {mtxObjectMatrix, mtxFormatArray, mtxFieldReal, mtxSymmetrySymm},
	12: {mtxObjectMatrix, mtxFormatArray, mtxFieldReal, mtxSymmetrySkew},
	13: {mtxObjectMatrix, mtxFormatArray, mtxFieldInteger, mtxSymmetryGeneral},
	14: {mtxObjectMatrix, mtxFormatArray, mtxFieldInteger, mtxSymmetrySymm},
	15: {mtxObjectMatrix, mtxFormatArray, mtxFieldInteger, mtxSymmetrySkew},
	16: {mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetryGeneral},
	17: {mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetrySymm},
	18: {mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetrySkew},
	19: {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldComplex, mtxSymmetryHermitian},
	20: {mtxObjectMatrix, mtxFormatArray, mtxFieldComplex, mtxSymmetryHermitian},
	21: {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldPattern, mtxSymmetryGeneral},
	22: {mtxObjectMatrix, mtxFormatCoordinate, mtxFieldPattern, mtxSymmetrySymm},
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

// Bytes returns a formatted Matrix Market headers
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

// index returns the (one-indexed) index of the Matrix Market type or -1
func (t *mmType) index() int {

	for i, t2 := range supported {
		if t.isMMType(&t2) {
			return i
		}
	}

	return -1
}

// scanHeader scans one line from a scanner and attempts to parse as a
// Matrix Market header
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

// counter tallies the number of bytes written to it
type counter struct {
	total int
}

// Write implements the io.Writer interface.
func (c *counter) Write(p []byte) (int, error) {
	var n int = len(p)
	c.total += n
	return n, nil
}
