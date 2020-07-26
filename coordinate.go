package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

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

	return nil, nil
}

func (m *MMCOO) MarshalTextTo(w io.Writer) (int, error) {

	return 0, nil
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

	scanner := bufio.NewScanner(r)
	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	// read header
	t, err := scanHeader(scanner)
	if err != nil {
		return 0, err
	}

	if t.isComplex() {
		return 0, ErrUnsupportedType
	}

	if err := m.scanCoordinateData(scanner); err != nil {
		return 0, err
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, nil
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
