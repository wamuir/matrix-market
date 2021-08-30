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

func (m *Dense) Do(fn func(i, j int, v float64)) {
	M, N := m.mat.Dims()
	for j := 0; j < N; j++ {
		for i := 0; i < M; i++ {
			fn(i, j, m.mat.At(i, j))
		}
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
	if n, err := fmt.Fprintf(w, "%%\n %d  %d\n", M, N); err == nil {
		total += n
	} else {
		return total, ErrUnwritable
	}

	var a floatAligner
	m.Do(a.Fit('f', -1, 64))

	// entries in column major order
	var buf = make([]byte, 0, 64)
	for j := 0; j < N; j++ {

		for i := 0; i < M; i++ {

			buf = a.Append(buf[:0], m.mat.At(i, j), 'f', -1, 64)
			buf = append(buf, '\n')

			n, err := w.Write(buf)
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

	// apply header fields
	m.Object = t.Object
	m.Format = t.Format
	m.Field = t.Field
	m.Symmetry = t.Symmetry

	switch t.index() {

	case 10, 11, 12, 13, 14, 15:
		if err := m.scanArrayData(scanner); err != nil {
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

func (m *Dense) scanArrayData(scanner *bufio.Scanner) error {

	var M, N, k int

	for scanner.Scan() {

		line := scanner.Text()

		// blank line or comment (%, Unicode 37)
		if r := []rune(line); len(r) == 0 || r[0] == 37 {
			continue
		}

		// M is number of rows, N is number of columns
		if _, err := fmt.Sscanf(line, "%d %d", &M, &N); err != nil {
			return ErrInputScanError
		}

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
		if k == M*N {
			return ErrInputScanError
		}

		if _, err := fmt.Sscanf(line, "%g", &v); err != nil {
			return ErrInputScanError
		}

		switch m.Symmetry {

		case mtxSymmetrySymm:

			// if above diagonal, move to diag
			for k%M < int(k/M) {
				k++
			}

			// if off diagonal, set value for symm element
			if int(k/M) != k%M {
				d.Set(int(k/M), k%M, v)
			}

		case mtxSymmetrySkew:

			// if on or above diagonal, move below diag
			for k%M <= int(k/M) {
				k++
			}

			// set skew value for symm element
			d.Set(int(k/M), k%M, -v)
		}

		d.Set(k%M, int(k/M), v)
		k++
	}

	// as skew-symmetric entries are below the diagonal in the matrix
	// market specification, there are no entries in the last column. The
	// counter k must must be advanced by one full column for
	// skew-symmetric matrices prior to the check which follows, otherwise
	// k == M * (N - 1) and the check will fail.
	if m.Symmetry == mtxSymmetrySkew {
		k += M
	}

	// compare counter k against expected number of entries (matrix size)
	if k != M*N {
		return ErrInputScanError
	}

	if err := scanner.Err(); err != nil {
		return ErrInputScanError
	}

	m.mat = d

	return nil
}
