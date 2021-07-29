package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	if n, err := fmt.Fprintf(w, "%d %d %d\n", M, N, m.mat.NNZ()); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	var err error
	m.mat.DoNonZero(func(i, j int, v float64) {

		if err != nil {
			return
		}

		n, e := fmt.Fprintf(w, "%d %d %g\n", i+1, j+1, v)
		if e != nil {
			err = ErrUnwritable
			return
		}

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

	switch t.index() {

	case 1, 2, 3, 4, 5, 6:
		if err := m.scanCoordinateData(scanner); err != nil {
			return n.total, err
		}

		if err := scanner.Err(); err != nil {
			return n.total, err
		}

	case 21, 22:
		if err := m.scanPatternData(scanner); err != nil {
			return n.total, err
		}

		if err := scanner.Err(); err != nil {
			return n.total, err
		}

	default:
		return n.total, ErrUnsupportedType

	}

	if t.isSymmetric() {
		m.mat.DoNonZero(func(i, j int, v float64) {
			if i == j {
				return
			}
			m.mat.Set(j, i, v)
		})
	}

	if t.isSkew() {
		m.mat.DoNonZero(func(i, j int, v float64) {
			if i == j {
				return
			}
			m.mat.Set(j, i, -v)
		})
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

		_, err := fmt.Sscanf(line, "%d %d %d", &M, &N, &L)
		if err != nil {
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

		_, err := fmt.Sscanf(line, "%d %d %f", &i, &j, &v)
		if err != nil {
			return err
		}

		c.Set(i-1, j-1, v)

		k++
	}

	// check if number of non-empty rows read is equal to expected
	// count of non-zero rows
	if k != L {
		return ErrInputScanError
	}

	if err := scanner.Err(); err != nil {
		return ErrInputScanError
	}

	m.mat = c

	return nil
}

func (m *COO) scanPatternData(scanner *bufio.Scanner) error {

	var M, N, L, k int

	for scanner.Scan() {

		line := scanner.Text()

		// blank line or comment (%, Unicode 37)
		if r := []rune(line); len(r) == 0 || r[0] == 37 {
			continue
		}

		_, err := fmt.Sscanf(line, "%d %d %d", &M, &N, &L)
		if err != nil {
			return ErrInputScanError
		}

		break

	}

	c := sparse.NewCOO(M, N, make([]int, L), make([]int, L), make([]float64, L))

	for scanner.Scan() {

		var i, j int

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

		_, err := fmt.Sscanf(line, "%d %d", &i, &j)
		if err != nil {
			return err
		}

		c.Set(i-1, j-1, 1)

		k++
	}

	// check if number of non-empty rows read is equal to expected
	// count of non-zero rows
	if k != L {
		return ErrInputScanError
	}

	if err := scanner.Err(); err != nil {
		return ErrInputScanError
	}

	m.mat = c

	return nil
}
