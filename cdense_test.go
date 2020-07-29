package market

import (
	"bytes"
	"strings"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// tol is the tolerance when testing matrices for equality
const tol = 1e-4

var cdense = mat.NewCDense(4, 5, []complex128{
	complex(+0.944853346337906500, -0.154091238677780850),
	complex(-0.681501551465435000, +0.594570321595631100),
	complex(-0.658745773257358300, +0.897566664045815500),
	complex(+0.402696290353813800, +0.009438983689089353),
	complex(+0.328601067704537230, +0.753843618074761200),
	complex(+0.812079966562488300, -0.274796067563821470),
	complex(+0.266121460291257600, -0.446018383861926500),
	complex(+0.756536462138819500, -0.429721939760935760),
	complex(+0.011573183932084063, +0.247960163711064440),
	complex(-0.551271155078584300, +0.157755862192646700),
	complex(+0.207552675260212820, -0.421555728398867800),
	complex(+0.795981993703873700, -0.288601857746140670),
	complex(-0.242048667319836990, -0.654258502990059600),
	complex(-0.247056369395660220, -0.190607085297800800),
	complex(-0.432441064387707700, +0.950877547679289700),
	complex(+0.419371027177237500, -0.664032247260985200),
	complex(+0.885613734373423400, +0.697886250370502100),
	complex(+0.593696988465424400, -0.223046160442398330),
	complex(+0.669421525553018500, +0.634515494429762400),
	complex(+0.393836704188575100, +0.061366273144705996),
})

func TestNewCDense(t *testing.T) {

	mtx1 := mat.NewCDense(4, 5, nil)
	_, _ = mtx1.Copy(cdense)

	m := NewCDense(mtx1)
	mtx2 := m.ToCMatrix()

	if !(mat.CEqual(mtx1, mtx2)) {
		t.Errorf("matrices differ")
		return
	}
}

func TestCDenseMarshalTextTo(t *testing.T) {

	var b strings.Builder

	mtx1 := mat.NewCDense(4, 5, nil)
	_, _ = mtx1.Copy(cdense)

	m1 := NewCDense(mtx1)

	_, err := m1.MarshalTextTo(&b)
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	r := strings.NewReader(b.String())

	m2 := NewCDense(mat.NewCDense(4, 5, nil))

	if _, err := m2.UnmarshalTextFrom(r); err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	if !(mat.CEqualApprox(m1.ToCMatrix(), m2.ToCMatrix(), tol)) {
		t.Errorf("matrices differ")
		return
	}

}

func TestCDenseMarshalText(t *testing.T) {

	mtx1 := mat.NewCDense(4, 5, nil)
	_, _ = mtx1.Copy(cdense)

	m1 := NewCDense(mtx1)

	out, err := m1.MarshalText()
	if err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	m2 := NewCDense(mat.NewCDense(4, 5, nil))

	if err := m2.UnmarshalText(out); err != nil {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	if !(mat.CEqualApprox(m1.ToCMatrix(), m2.ToCMatrix(), tol)) {
		t.Errorf("matrices differ")
		return
	}

}

func TestCDenseUnmarshalText(t *testing.T) {

	var (
		in   []byte
		mtx1 mat.CMatrix
	)

	// just a good matrix
	in = []byte(`%%MatrixMarket matrix array complex general
              4 5
	       0.944853346337906500 -0.154091238677780850
	      -0.681501551465435000  0.594570321595631100
	      -0.658745773257358300  0.897566664045815500
	       0.402696290353813800  0.009438983689089353
	       0.328601067704537230  0.753843618074761200
	       0.812079966562488300 -0.274796067563821470
	       0.266121460291257600 -0.446018383861926500
	       0.756536462138819500 -0.429721939760935760
	       0.011573183932084063  0.247960163711064440
	      -0.551271155078584300  0.157755862192646700
	       0.207552675260212820 -0.421555728398867800
	       0.795981993703873700 -0.288601857746140670
	      -0.242048667319836990 -0.654258502990059600
	      -0.247056369395660220 -0.190607085297800800
	      -0.432441064387707700  0.950877547679289700
	       0.419371027177237500 -0.664032247260985200
	       0.885613734373423400  0.697886250370502100
	       0.593696988465424400 -0.223046160442398330
	       0.669421525553018500  0.634515494429762400
	       0.393836704188575100  0.061366273144705996`,
	)

	m := NewCDense(mat.NewCDense(4, 5, nil))

	if err := m.UnmarshalText(in); err == nil {
		mtx1 = m.ToCMatrix()
	} else {
		t.Errorf("Received unexpected error: %v", err.Error())
		return
	}

	mtx2 := mat.NewCDense(4, 5, nil)
	_, _ = mtx2.Copy(cdense)

	if !(mat.CEqualApprox(mtx1, mtx2, tol)) {
		t.Errorf("matrices differ")
		return
	}
}

func TestCDenseUnmarshalTextFrom(t *testing.T) {

	var (
		in   []byte
		m    CDense
		mtx1 mat.CMatrix
	)

	// just a good matrix
	in = []byte(`%%MatrixMarket matrix array complex general
              4 5
	       0.944853346337906500 -0.154091238677780850
	      -0.681501551465435000  0.594570321595631100
	      -0.658745773257358300  0.897566664045815500
	       0.402696290353813800  0.009438983689089353
	       0.328601067704537230  0.753843618074761200
	       0.812079966562488300 -0.274796067563821470
	       0.266121460291257600 -0.446018383861926500
	       0.756536462138819500 -0.429721939760935760
	       0.011573183932084063  0.247960163711064440
	      -0.551271155078584300  0.157755862192646700
	       0.207552675260212820 -0.421555728398867800
	       0.795981993703873700 -0.288601857746140670
	      -0.242048667319836990 -0.654258502990059600
	      -0.247056369395660220 -0.190607085297800800
	      -0.432441064387707700  0.950877547679289700
	       0.419371027177237500 -0.664032247260985200
	       0.885613734373423400  0.697886250370502100
	       0.593696988465424400 -0.223046160442398330
	       0.669421525553018500  0.634515494429762400
	       0.393836704188575100  0.061366273144705996`,
	)

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

	mtx1 = m.ToCMatrix()

	mtx2 := mat.NewCDense(4, 5, nil)
	_, _ = mtx2.Copy(cdense)

	if !(mat.CEqualApprox(mtx1, mtx2, tol)) {
		t.Errorf("matrices differ")
		return
	}
}
