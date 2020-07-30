package market

import (
	"bytes"
	"strings"
	"testing"

	"gonum.org/v1/gonum/mat"
)

var dense = mat.NewDense(4, 5, []float64{
	+0.944853346337906500,
	-0.681501551465435000,
	-0.000000000000000000,
	+0.402696290353813800,
	+0.328601067704537230,
	+0.812079966562488300,
	-0.000000000000000000,
	+0.756536462138819500,
	+0.011573183932084063,
	-0.551271155078584300,
	+0.207552675260212820,
	+0.795981993703873700,
	-0.242048667319836990,
	-0.000000000000000000,
	-0.000000000000000000,
	+0.419371027177237500,
	+0.885613734373423400,
	-0.000000000000000000,
	+0.669421525553018500,
	+0.393836704188575100,
})

func TestNewDense(t *testing.T) {

	mtx1 := mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	m := NewDense(mtx1)
	mtx2 := m.ToMatrix()

	if !(mat.Equal(mtx1, mtx2)) {
		t.Errorf(
			"matrices differ\n \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(mtx1, mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(mtx2, mat.Prefix("    "), mat.Squeeze()),
		)
		return
	}
}

func TestDenseMarshalTextTo(t *testing.T) {

	var b strings.Builder

	mtx := mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	m1 := NewDense(mtx)

	_, err := m1.MarshalTextTo(&b)
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	r := strings.NewReader(b.String())

	m2 := NewDense(mat.NewDense(4, 5, nil))

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

func TestDenseMarshalText(t *testing.T) {

	mtx := mat.NewDense(
		4, 5, []float64{0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
	)

	m1 := NewDense(mtx)

	out, err := m1.MarshalText()
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	m2 := NewDense(mat.NewDense(4, 5, nil))

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

func TestDenseUnmarshalText(t *testing.T) {

	var (
		in         []byte
		mtx1, mtx2 mat.Matrix
	)

	// just a good matrix
	in = []byte(`%%MatrixMarket matrix array real general
              4 5
	      0
	      0
	      1
	      0
	      0
	      1
	      0
	      0
	      0
	      0
	      0
	      1
	      0
	      0
	      0
	      0
	      0
	      0
	      0
	      0`)

	m := NewDense(mat.NewDense(4, 5, nil))

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

	// number of entries != M * N, should fail
	in = []byte(`%%MatrixMarket matrix array real general
              4 5
	      0
	      0
	      1
	      0
	      0
	      1
	      0
	      0
	      0
	      0
	      0
	      1
	      0
	      0
	      0
	      0
	      0
	      0
	      0
	      0
	      999`)

	m2 := NewDense(mat.NewDense(4, 5, nil))

	if err := m2.UnmarshalText(in); err != ErrInputScanError {
		t.Errorf("Expected ErrInputScanError; received: %v", err)
		return
	}
}

func TestDenseUnmarshalTextFrom(t *testing.T) {

	var (
		in         []byte
		m          Dense
		mtx1, mtx2 mat.Matrix
	)

	// just a good matrix
	in = []byte(`%%MatrixMarket matrix array real general
              4 5
	      0
	      0
	      1
	      0
	      0
	      1
	      0
	      0
	      0
	      0
	      0
	      1
	      0
	      0
	      0
	      0
	      0
	      0
	      0
	      0`)

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
