package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// Dense is a type embedding of mat.Dense, for reading and writing
// complex-valued matrices in Matrix Market array format.
type Dense struct {
	Object   string
	Format   string
	Field    string
	Symmetry string
	mat      *mat.Dense
}

// NewDense initializes a new CDense dense matrix from a mat.Dense
// matrix.
func NewDense(d *mat.Dense) *Dense {
	return &Dense{
		Object:   mtxObjectMatrix,
		Format:   mtxFormatArray,
		Field:    mtxFieldReal,
		Symmetry: mtxSymmetryGeneral,
		mat:      d,
	}
}

// ToDense returns a mat.Dense matrix that shares underlying storage
// with the receiver.
func (m *Dense) ToDense() *mat.Dense { return m.mat }

// ToMatrix returns a mat.Matrix complex matrix that shares underlying
// storage with the receiver.
func (m *Dense) ToMatrix() mat.Matrix { return m.mat }

// MarshalText serializes the receiver to []byte in Matrix Market format
// and returns the result.
func (m *Dense) MarshalText() ([]byte, error) {

	var b strings.Builder

	if _, err := m.MarshalTextTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

// MarshalTextTo serializes the receiver to w in Matrix Market format
// and returns the result.
func (m *Dense) MarshalTextTo(w io.Writer) (int, error) {

	var total int

	t := mmType{m.Object, m.Format, m.Field, m.Symmetry}

	if n, err := w.Write(t.Bytes()); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	M, N := m.mat.Dims()
	if n, err := fmt.Fprintf(w, "%d %d\n", M, N); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	for i := 0; i < M; i++ {

		for j := 0; j < N; j++ {

			n, err := fmt.Fprintf(w, "%f\n", m.mat.At(i, j))
			if err != nil {
				return total, ErrUnwritable
			}

			total += n
		}

	}

	return total, nil
}

// UnmarshalText deserializes []byte from Matrix Market format
// into the receiver.
func (m *Dense) UnmarshalText(text []byte) error {

	r := bytes.NewReader(text)

	if _, err := m.UnmarshalTextFrom(r); err != nil {
		return err
	}

	return nil
}

// UnmarshalTextFrom deserializes r from Matrix Market format
// into the receiver.
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

	m.mat = d

	return nil
}
