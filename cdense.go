package market

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// CDense is a type embedding of sparse.COO
type CDense struct{ CMatrix *mat.CDense }

// NewCDense creates a new MMCOO from a sparse.COO
func NewCDense(d *mat.CDense) *CDense { return &CDense{d} }

// ToCDense shares data with the receiver
func (m *CDense) ToCDense() *mat.CDense { return m.CMatrix }

func (m *CDense) ToCMatrix() mat.CMatrix { return m.CMatrix }

func (m *CDense) MarshalText() ([]byte, error) {

	var b strings.Builder

	if _, err := m.MarshalTextTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

func (m *CDense) MarshalTextTo(w io.Writer) (int, error) {

	var total int

	t := mmType{
		mtxObjectMatrix,
		mtxFormatArray,
		mtxFieldComplex,
		mtxSymmetryGeneral,
	}

	if n, err := w.Write(t.Bytes()); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	M, N := m.CMatrix.Dims()
	if n, err := fmt.Fprintf(w, "%d %d\n", M, N); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	for i := 0; i < M; i++ {

		for j := 0; j < N; j++ {

			v := m.CMatrix.At(i, j)

			n, err := fmt.Fprintf(w, "%f %f\n", real(v), imag(v))
			if err != nil {
				return total, ErrUnwritable
			}

			total += n
		}

	}

	return total, nil
}

// Should the receiver not be a pointer?
func (m *CDense) UnmarshalText(text []byte) error {

	r := bytes.NewReader(text)

	if _, err := m.UnmarshalTextFrom(r); err != nil {
		return err
	}

	return nil
}

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

	if t.isComplex() {
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

	m.CMatrix = d

	return nil
}
