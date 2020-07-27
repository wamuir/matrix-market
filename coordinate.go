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

// MMCOO is a type embedding of sparse.COO
type MMCOO struct{ Matrix *sparse.COO }

// NewMMCOO creates a new MMCOO from a sparse.COO
func NewMMCOO(c *sparse.COO) *MMCOO { return &MMCOO{c} }

// ToCOO shares data with the receiver
func (m *MMCOO) ToCOO() *sparse.COO { return m.Matrix }

func (m *MMCOO) ToMatrix() mat.Matrix { return m.Matrix }

func (m *MMCOO) MarshalText() ([]byte, error) {

	var b strings.Builder

	if _, err := m.MarshalTextTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

func (m *MMCOO) MarshalTextTo(w io.Writer) (int, error) {

	var total int

	t := mmType{
		mtxObjectMatrix,
		mtxFormatCoordinate,
		mtxFieldReal,
		mtxSymmetryGeneral,
	}

	if n, err := w.Write(t.Bytes()); err == nil {
		total += n
	} else {
		return total, err
	}

	M, N := m.Matrix.Dims()
	if n, err := fmt.Fprintf(w, "%d %d %d\n", M, N, m.Matrix.NNZ()); err == nil {
		total += n
	} else {
		return total, err
	}

	m.Matrix.DoNonZero(func(i, j int, v float64) {
		if n, err := fmt.Fprintf(w, "%d %d %f\n", i+1, j+1, v); err == nil {
			total += n
		} else {
			panic(err)
		}
	})

	return total, nil
}

// Should the receiver not be a pointer?
func (m *MMCOO) UnmarshalText(text []byte) error {

	r := bytes.NewReader(text)

	if _, err := m.UnmarshalTextFrom(r); err != nil {
		return err
	}

	return nil
}

func (m *MMCOO) UnmarshalTextFrom(r io.Reader) (int, error) {

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

	if t.isComplex() {
		return n.total, ErrUnsupportedType
	}

	if err := m.scanCoordinateData(scanner); err != nil {
		return n.total, err
	}

	if err := scanner.Err(); err != nil {
		return n.total, err
	}

	return n.total, nil
}

func (m *MMCOO) scanCoordinateData(scanner *bufio.Scanner) error {

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

	m.Matrix = c

	return nil
}
