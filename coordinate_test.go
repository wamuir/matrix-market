package market

import (
	"bytes"
	"strings"
	"testing"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

func TestNewMMCOO(t *testing.T) {

	var (
		c          *sparse.COO
		m          *MMCOO
		mtx1, mtx2 mat.Matrix
	)

	c = sparse.NewCOO(4, 5, nil, nil, nil)
	c.Set(0, 2, 1)
	c.Set(1, 0, 1)
	c.Set(2, 1, 1)

	m = NewMMCOO(c)
	mtx1 = m.ToMatrix()
	mtx2 = mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	if !(mat.Equal(mtx1, mtx2)) {
		t.Errorf(
			"matrices differ\n \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(mtx1, mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(mtx2, mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}
}

func TestCOOMarshalText(t *testing.T) {

	var (
		b      strings.Builder
		c      *sparse.COO
		m1, m2 *MMCOO
	)

	c = sparse.NewCOO(4, 5, nil, nil, nil)
	c.Set(0, 2, 1)
	c.Set(1, 0, 1)
	c.Set(2, 1, 1)

	m1 = NewMMCOO(c)

	_, err := m1.MarshalTextTo(&b)
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	r := strings.NewReader(b.String())

	m2 = NewMMCOO(sparse.NewCOO(4, 5, nil, nil, nil))

	if _, err := m2.UnmarshalTextFrom(r); err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	if !(mat.Equal(m1.ToMatrix(), m2.ToMatrix())) {
		t.Errorf(
			"matrices differ\nexpected \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(m1.ToMatrix(), mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(m2.ToMatrix(), mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}

}

func TestCOOMarshalTextTo(t *testing.T) {

	var (
		c      *sparse.COO
		m1, m2 *MMCOO
		out    []byte
	)

	c = sparse.NewCOO(4, 5, nil, nil, nil)
	c.Set(0, 2, 1)
	c.Set(1, 0, 1)
	c.Set(2, 1, 1)

	m1 = NewMMCOO(c)

	out, err := m1.MarshalText()
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	m2 = NewMMCOO(sparse.NewCOO(4, 5, nil, nil, nil))

	if err := m2.UnmarshalText(out); err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	if !(mat.Equal(m1.ToMatrix(), m2.ToMatrix())) {
		t.Errorf(
			"matrices differ\nexpected \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(m1.ToMatrix(), mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(m2.ToMatrix(), mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}

}

func TestCOOUnmarshalText(t *testing.T) {

	var (
		in         []byte
		m          MMCOO
		mtx1, mtx2 mat.Matrix
	)

	// just a good matrix
	in = []byte(`%%MatrixMarket matrix coordinate real general
              4 5 3
              1 3 1
              2 1 1
              3 2 1`)

	if err := m.UnmarshalText(in); err == nil {
		mtx1 = m.ToMatrix()
	} else {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	mtx2 = mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	if !(mat.Equal(mtx1, mtx2)) {
		t.Errorf(
			"matrices differ\nexpected \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(mtx1, mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(mtx2, mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}

	// duplicate coordinate entries, should still pass
	in = []byte(`%%MatrixMarket matrix coordinate real general
              4 5 4
              1 3 0.5
              1 3 0.5
              2 1 1
              3 2 1`)

	if err := m.UnmarshalText(in); err == nil {
		mtx1 = m.ToMatrix()
	} else {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	mtx2 = mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	if !(mat.Equal(mtx1, mtx2)) {
		t.Errorf(
			"matrices differ\nexpected \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(mtx1, mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(mtx2, mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}
}

func TestCOOUnmarshalTextFrom(t *testing.T) {

	var (
		in         []byte
		m          MMCOO
		mtx1, mtx2 mat.Matrix
	)

	// just a good matrix
	in = []byte(`%%MatrixMarket matrix coordinate real general
              4 5 3
              1 3 1
              2 1 1
              3 2 1`)

	r := bytes.NewReader(in)

	n, err := m.UnmarshalTextFrom(r)
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	if n != len(in) {
		t.Errorf("Inconsistent number bytes read (%d), expected %d", n, len(in))
		return
	}

	mtx1 = m.ToMatrix()

	mtx2 = mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	if !(mat.Equal(mtx1, mtx2)) {
		t.Errorf(
			"matrices differ\nexpected \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(mtx1, mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(mtx2, mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}
}
