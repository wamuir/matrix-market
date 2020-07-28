package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// Dense is a type embedding of sparse.COO
type Dense struct{ Matrix *mat.Dense }

// NewDense creates a new MMCOO from a sparse.COO
func NewDense(d *mat.Dense) *Dense { return &Dense{d} }

// ToDense shares data with the receiver
func (m *Dense) ToDense() *mat.Dense { return m.Matrix }

func (m *Dense) ToMatrix() mat.Matrix { return m.Matrix }

func (m *Dense) MarshalText() ([]byte, error) {

	var b strings.Builder

	if _, err := m.MarshalTextTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

func (m *Dense) MarshalTextTo(w io.Writer) (int, error) {

	var total int

	t := mmType{
		mtxObjectMatrix,
		mtxFormatArray,
		mtxFieldReal,
		mtxSymmetryGeneral,
	}

	if n, err := w.Write(t.Bytes()); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	M, N := m.Matrix.Dims()
	if n, err := fmt.Fprintf(w, "%d %d\n", M, N); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	for i := 0; i < M; i++ {

		for j := 0; j < N; j++ {

			n, err := fmt.Fprintf(w, "%f\n", m.Matrix.At(i, j))
			if err != nil {
				return total, ErrUnwritable
			}

			total += n
		}

	}

	return total, nil
}

// Should the receiver not be a pointer?
func (m *Dense) UnmarshalText(text []byte) error {

	r := bytes.NewReader(text)

	if _, err := m.UnmarshalTextFrom(r); err != nil {
		return err
	}

	return nil
}

func (m *Dense) UnmarshalTextFrom(r io.Reader) (int, error) {

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

	if !(t.isMatrix() && t.isArray() && (t.isReal() || t.isInteger()) && t.isGeneral()) {
		return n.total, ErrUnsupportedType
	}

	if err := m.scanArrayData(scanner); err != nil {
		return n.total, err
	}

	if err := scanner.Err(); err != nil {
		return n.total, err
	}

	return n.total, nil
}

func (m *Dense) scanArrayData(scanner *bufio.Scanner) error {

	var M, N, L, k int

	for scanner.Scan() {

		line := scanner.Text()

		// blank line or comment (%, Unicode 37)
		if r := []rune(line); len(r) == 0 || r[0] == 37 {
			continue
		}

		_, err := fmt.Sscanf(line, "%d %d", &M, &N)
		if err != nil {
			return ErrInputScanError
		}

		L = M * N

		break
	}

	d := mat.NewDense(M, N, nil)

	for scanner.Scan() {

		var v float64

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

		_, err := fmt.Sscanf(line, "%f", &v)
		if err != nil {
			return ErrInputScanError
		}

		d.Set(int(k/N), k%N, v)

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

	m.Matrix = d

	return nil
}
