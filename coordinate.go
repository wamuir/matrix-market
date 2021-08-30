package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	// "strconv"
	"strings"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

// COO is a type embedding of sparse.COO, for reading and writing
// real-valued matrices in Matrix Market coordinate format.
type COO struct {
	Object   string
	Format   string
	Field    string
	Symmetry string
	mat      *sparse.COO
}

// NewCOO initializes a new COO sparse matrix from a sparse.COO matrix
func NewCOO(c *sparse.COO) *COO {
	return &COO{
		Object:   mtxObjectMatrix,
		Format:   mtxFormatCoordinate,
		Field:    mtxFieldReal,
		Symmetry: mtxSymmetryGeneral,
		mat:      c,
	}
}

func (m *COO) Do(fn func(i, j int, v float64)) {
	m.mat.DoNonZero(fn)
}

// ToCOO returns a sparse.COO matrix that shared underlying storage
// with the receiver.
func (m *COO) ToCOO() *sparse.COO { return m.mat }

// ToMatrix returns a mat.Matrix real matrix that shared underlying
// storage with the receiver.
func (m *COO) ToMatrix() mat.Matrix { return m.mat }

// MarshalText serializes the receiver to []byte in Matrix Market
// format and returns the result.
func (m *COO) MarshalText() ([]byte, error) {

	var b strings.Builder

	if _, err := m.MarshalTextTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

// MarshalTextTo serializes the receiver to w in Matrix Market format
// and returns the result.
func (m *COO) MarshalTextTo(w io.Writer) (int, error) {

	var total int

	t := mmType{m.Object, m.Format, m.Field, m.Symmetry}

	// Need additional checks on mmType
	if !(t.isMatrix() && t.isCoordinate()) {
		return total, ErrUnsupportedType
	}

	if n, err := w.Write(t.Bytes()); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	M, N := m.mat.Dims()
	if n, err := fmt.Fprintf(w, "%%\n %d  %d  %d\n", M, N, m.mat.NNZ()); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	var a floatTripletAligner
	m.Do(a.Fit('f', -1, 64))

	// entries in column major order
	var (
		buf = make([]byte, 0, 64)
		err error
		n   int
	)
	m.mat.DoNonZero(func(i, j int, v float64) {
		buf = a.Append(buf[:0], i, j, v, 'f', -1, 64)
		buf = append(buf, '\n')

		n, err = w.Write(buf)
		total += n
	})

	return total, err
}

// UnmarshalText deserializes []byte from Matrix Market format into
// the receiver.
func (m *COO) UnmarshalText(text []byte) error {

	r := bytes.NewReader(text)

	if _, err := m.UnmarshalTextFrom(r); err != nil {
		return err
	}

	return nil
}

// UnmarshalTextFrom deserializes r from Matrix Market format into the
// receiver.
func (m *COO) UnmarshalTextFrom(r io.Reader) (int, error) {

	var n counter

	r = io.TeeReader(r, &n)

	scanner := bufio.NewScanner(r)
	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	// read header
	t, err := scanHeader(scanner)
	if err != nil {
		return n.total, err
	}

	// apply header fields
	m.Object = t.Object
	m.Format = t.Format
	m.Field = t.Field
	m.Symmetry = t.Symmetry

	switch t.index() {

	case 1, 2, 3, 4, 5, 6, 21, 22:
		if err := m.scanCoordinateData(scanner); err != nil {
			return n.total, err
		}

		if err := scanner.Err(); err != nil {
			return n.total, err
		}

	default:
		return n.total, ErrUnsupportedType

	}

	return n.total, nil
}

func (m *COO) scanCoordinateData(scanner *bufio.Scanner) error {

	var M, N, L, k int

	for scanner.Scan() {

		line := scanner.Text()

		// blank line or comment (%, Unicode 37)
		if r := []rune(line); len(r) == 0 || r[0] == 37 {
			continue
		}

		if _, err := fmt.Sscanf(line, "%d %d %d", &M, &N, &L); err != nil {
			return ErrInputScanError
		}

		break

	}

	c := sparse.NewCOO(M, N, make([]int, L), make([]int, L), make([]float64, L))

	for scanner.Scan() {

		var (
			i, j int
			v    float64
		)

		line := scanner.Text()

		// blank lines are allowed in data per design spec
		if r := []rune(line); len(r) == 0 {
			continue
		}

		// error out if data rows exceed expected non-zero entries
		// (note that k is zero indexed)
		if k == L {
			return ErrInputScanError
		}

		switch m.Field {

		case mtxFieldInteger, mtxFieldReal:

			if _, err := fmt.Sscanf(line, "%d %d %f", &i, &j, &v); err != nil {
				return err
			}

		case mtxFieldPattern:

			v = 1.0

			if _, err := fmt.Sscanf(line, "%d %d", &i, &j); err != nil {
				return err
			}
		}

		switch m.Symmetry {

		case mtxSymmetrySymm:

			// if off diagonal, set value for symm element
			if i != j {
				c.Set(j-1, i-1, v)
			}

		case mtxSymmetrySkew:

			// if off diagonal, set skew value for symm element
			// (note. diagonal elements aren't allowed for skew mats)
			if i != j {
				c.Set(j-1, i-1, -v)
			}
		}

		c.Set(i-1, j-1, v)

		k++
	}

	// compare counter k against expected number of expected entries L
	if k != L {
		return ErrInputScanError
	}

	if err := scanner.Err(); err != nil {
		return ErrInputScanError
	}

	m.mat = c

	return nil
}
