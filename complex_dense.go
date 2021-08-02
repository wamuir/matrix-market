package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// CDense is a type embedding of mat.CDense, for reading and writing
// complex-valued matrices in Matrix Market array format.
type CDense struct {
	Object   string
	Format   string
	Field    string
	Symmetry string
	mat      *mat.CDense
}

// NewCDense initializes a new CDense dense matrix from a mat.CDense
// matrix.
func NewCDense(d *mat.CDense) *CDense {
	return &CDense{
		Object:   mtxObjectMatrix,
		Format:   mtxFormatArray,
		Field:    mtxFieldComplex,
		Symmetry: mtxSymmetryGeneral,
		mat:      d,
	}
}

// ToCDense returns a mat.CDense matrix that shares underlying storage
// with the receiver.
func (m *CDense) ToCDense() *mat.CDense { return m.mat }

// ToCMatrix returns a mat.CMatrix complex matrix that shares underlying
// storage with the receiver.
func (m *CDense) ToCMatrix() mat.CMatrix { return m.mat }

// MarshalText serializes the receiver to []byte in Matrix Market format
// and returns the result.
func (m *CDense) MarshalText() ([]byte, error) {

	var b strings.Builder

	if _, err := m.MarshalTextTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

// MarshalTextTo serializes the receiver to w in Matrix Market format
// and returns the result.
func (m *CDense) MarshalTextTo(w io.Writer) (int, error) {

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

			v := m.mat.At(i, j)

			n, err := fmt.Fprintf(w, "%g %g\n", real(v), imag(v))
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
func (m *CDense) UnmarshalText(text []byte) error {

	r := bytes.NewReader(text)

	if _, err := m.UnmarshalTextFrom(r); err != nil {
		return err
	}

	return nil
}

// UnmarshalTextFrom deserializes r from Matrix Market format
// into the receiver.
func (m *CDense) UnmarshalTextFrom(r io.Reader) (int, error) {

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

	if !(t.isMatrix() && t.isArray() && t.isComplex() && t.isGeneral()) {
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

func (m *CDense) scanArrayData(scanner *bufio.Scanner) error {

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

	d := mat.NewCDense(M, N, nil)

	for scanner.Scan() {

		var vr, vi float64

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

		_, err := fmt.Sscanf(line, "%f %f", &vr, &vi)
		if err != nil {
			return ErrInputScanError
		}

		d.Set(int(k/N), k%N, complex(vr, vi))

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
