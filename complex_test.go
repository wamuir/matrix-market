package market

import (
	"fmt"
	"io"
	"math/cmplx"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

var mtx16 = mat.NewCDense(4, 5, []complex128{
	complex(+0.944853346337906500, -0.154091238677780850),
	complex(+0.328601067704537230, +0.753843618074761200),
	complex(+0.011573183932084063, +0.247960163711064440),
	complex(+0.242048667319836990, -0.654258502990059600),
	complex(+0.885613734373423400, +0.697886250370502100),
	complex(-0.681501551465435000, +0.594570321595631100),
	complex(+0.812079966562488300, -0.274796067563821470),
	complex(-0.551271155078584300, +0.157755862192646700),
	complex(-0.247056369395660220, -0.190607085297800800),
	complex(+0.593696988465424400, -0.223046160442398330),
	complex(-0.658745773257358300, +0.897566664045815500),
	complex(+0.266121460291257600, -0.446018383861926500),
	complex(+0.207552675260212820, -0.421555728398867800),
	complex(-0.432441064387707700, +0.950877547679289700),
	complex(+0.669421525553018500, +0.634515494429762400),
	complex(+0.402696290353813800, +0.009438983689089353),
	complex(+0.756536462138819500, -0.429721939760935760),
	complex(+0.795981993703873700, -0.288601857746140670),
	complex(+0.419371027177237500, -0.664032247260985200),
	complex(+0.393836704188575100, +0.061366273144705996),
})

var mtx17 = mat.NewCDense(5, 5, []complex128{
	complex(+0.944853346337906500, -0.154091238677780850),
	complex(-0.681501551465435000, +0.594570321595631100),
	complex(-0.658745773257358300, +0.897566664045815500),
	complex(+0.402696290353813800, +0.009438983689089353),
	complex(+0.328601067704537230, +0.753843618074761200),
	complex(-0.681501551465435000, +0.594570321595631100),
	complex(+0.812079966562488300, -0.274796067563821470),
	complex(+0.266121460291257600, -0.446018383861926500),
	complex(+0.756536462138819500, -0.429721939760935760),
	complex(+0.011573183932084063, +0.247960163711064440),
	complex(-0.658745773257358300, +0.897566664045815500),
	complex(+0.266121460291257600, -0.446018383861926500),
	complex(-0.551271155078584300, +0.157755862192646700),
	complex(+0.207552675260212820, -0.421555728398867800),
	complex(+0.795981993703873700, -0.288601857746140670),
	complex(+0.402696290353813800, +0.009438983689089353),
	complex(+0.756536462138819500, -0.429721939760935760),
	complex(+0.207552675260212820, -0.421555728398867800),
	complex(+0.242048667319836990, -0.654258502990059600),
	complex(-0.247056369395660220, -0.190607085297800800),
	complex(+0.328601067704537230, +0.753843618074761200),
	complex(+0.011573183932084063, +0.247960163711064440),
	complex(+0.795981993703873700, -0.288601857746140670),
	complex(-0.247056369395660220, -0.190607085297800800),
	complex(-0.432441064387707700, +0.950877547679289700),
})

var mtx18 = mat.NewCDense(5, 5, []complex128{
	complex(+0.000000000000000000, +0.000000000000000000),
	complex(+0.681501551465435000, -0.594570321595631100),
	complex(+0.658745773257358300, -0.897566664045815500),
	complex(-0.402696290353813800, -0.009438983689089353),
	complex(-0.328601067704537230, -0.753843618074761200),
	complex(-0.681501551465435000, +0.594570321595631100),
	complex(+0.000000000000000000, +0.000000000000000000),
	complex(-0.266121460291257600, +0.446018383861926500),
	complex(-0.756536462138819500, +0.429721939760935760),
	complex(-0.011573183932084063, -0.247960163711064440),
	complex(-0.658745773257358300, +0.897566664045815500),
	complex(+0.266121460291257600, -0.446018383861926500),
	complex(+0.000000000000000000, +0.000000000000000000),
	complex(-0.207552675260212820, +0.421555728398867800),
	complex(-0.795981993703873700, +0.288601857746140670),
	complex(+0.402696290353813800, +0.009438983689089353),
	complex(+0.756536462138819500, -0.429721939760935760),
	complex(+0.207552675260212820, -0.421555728398867800),
	complex(+0.000000000000000000, +0.000000000000000000),
	complex(+0.247056369395660220, +0.190607085297800800),
	complex(+0.328601067704537230, +0.753843618074761200),
	complex(+0.011573183932084063, +0.247960163711064440),
	complex(+0.795981993703873700, -0.288601857746140670),
	complex(-0.247056369395660220, -0.190607085297800800),
	complex(+0.000000000000000000, +0.000000000000000000),
})

var mtx20 = mat.NewCDense(5, 5, []complex128{
	complex(+0.944853346337906500, -0.154091238677780850),
	complex(-0.681501551465435000, -0.594570321595631100),
	complex(-0.658745773257358300, -0.897566664045815500),
	complex(+0.402696290353813800, -0.009438983689089353),
	complex(+0.328601067704537230, -0.753843618074761200),
	complex(-0.681501551465435000, +0.594570321595631100),
	complex(+0.812079966562488300, -0.274796067563821470),
	complex(+0.266121460291257600, +0.446018383861926500),
	complex(+0.756536462138819500, +0.429721939760935760),
	complex(+0.011573183932084063, -0.247960163711064440),
	complex(-0.658745773257358300, +0.897566664045815500),
	complex(+0.266121460291257600, -0.446018383861926500),
	complex(-0.551271155078584300, +0.157755862192646700),
	complex(+0.207552675260212820, +0.421555728398867800),
	complex(+0.795981993703873700, +0.288601857746140670),
	complex(+0.402696290353813800, +0.009438983689089353),
	complex(+0.756536462138819500, -0.429721939760935760),
	complex(+0.207552675260212820, -0.421555728398867800),
	complex(+0.242048667319836990, -0.654258502990059600),
	complex(-0.247056369395660220, +0.190607085297800800),
	complex(+0.328601067704537230, +0.753843618074761200),
	complex(+0.011573183932084063, +0.247960163711064440),
	complex(+0.795981993703873700, -0.288601857746140670),
	complex(-0.247056369395660220, -0.190607085297800800),
	complex(-0.432441064387707700, +0.950877547679289700),
})

func TestNewCDense(t *testing.T) {

	m := NewCDense(mtx16)

	assert.True(t, mat.CEqual(m.ToCMatrix(), mtx16))
}

func TestCDenseToCDense(t *testing.T) {

	m := NewCDense(mtx16)

	assert.True(t, mat.CEqual(m.ToCDense(), mtx16))
}

func TestCDenseToCMatrix(t *testing.T) {

	m := NewCDense(mtx16)

	assert.True(t, mat.CEqual(m.ToCDense(), mtx16))
}

func TestCDenseMarshalTextTo(t *testing.T) {

	var b strings.Builder

	m := NewCDense(mtx16)

	_, err := m.MarshalTextTo(&b)
	assert.Nil(t, err)

	mm, err := os.ReadFile(filepath.Join("testdata", "mmtype-16.mtx"))
	assert.Nil(t, err)

	assert.Equal(t, b.String(), string(mm))
}

func TestCDenseMarshalText(t *testing.T) {

	m := NewCDense(mtx16)

	mm1, err := m.MarshalText()
	assert.Nil(t, err)

	mm2, err := os.ReadFile(filepath.Join("testdata", "mmtype-16.mtx"))
	assert.Nil(t, err)

	assert.Equal(t, string(mm1), string(mm2))
}

func TestCDenseUnmarshalText(t *testing.T) {

	var mm CDense

	b, err := os.ReadFile(filepath.Join("testdata", "mmtype-16.mtx"))
	assert.Nil(t, err)

	assert.Nil(t, mm.UnmarshalText(b))

	assert.True(t, mat.CEqual(mm.ToCMatrix(), mtx16))
}

func TestCDenseUnmarshalTextFrom(t *testing.T) {

	c := map[string]mat.CMatrix{
		"mmtype-07.mtx": mtx16, // coordinate complex general
		"mmtype-08.mtx": mtx17, // coordinate complex symmetric
		"mmtype-09.mtx": mtx18, // coordinate complex skew-symmetric
		"mmtype-19.mtx": mtx20, // coordinate complex hermitian
		"mmtype-16.mtx": mtx16, // array complex general
		"mmtype-17.mtx": mtx17, // array complex symmetric
		"mmtype-18.mtx": mtx18, // array complex skew-symmetric
		"mmtype-20.mtx": mtx20, // array complex hermitian
	}

	for k, v := range c {

		f, _ := os.Open(filepath.Join("testdata", k))
		defer f.Close()

		var mm CDense
		if _, err := mm.UnmarshalTextFrom(f); err != nil {
			t.Errorf("%v", err)
		}

		if !mat.CEqualApprox(mm.ToCMatrix(), v, 1e-14) {
			t.Errorf(
				"\ngot:\n    %v\nwant:\n    %v\n",
				mm.ToCMatrix(),
				v,
			)
		}
	}

}

func BenchmarkCDenseMarshalTextTo(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		a := mat.NewCDense(i, i, nil)
		for j := 0; j < i*i; j++ {
			a.Set(j%i, int(j/i), cmplx.Sqrt(complex(float64(j), float64(j))))
		}
		m := NewCDense(a)
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				_, err := m.MarshalTextTo(io.Discard)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCDenseUnmarshalTextFrom(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		a := mat.NewCDense(i, i, nil)
		for j := 0; j < i*i; j++ {
			a.Set(j%i, int(j/i), cmplx.Sqrt(complex(float64(j), float64(j))))
		}
		m := NewCDense(a)
		t, _ := m.MarshalText()
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				if err := m.UnmarshalText(t); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
