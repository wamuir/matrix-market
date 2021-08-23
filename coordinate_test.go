package market

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/james-bowman/sparse"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

var mtx01 = sparse.NewCOO(
	4,
	5,
	[]int{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 3, 3, 3, 3},
	[]int{0, 1, 3, 4, 0, 2, 3, 4, 0, 1, 2, 0, 1, 3, 4},
	[]float64{
		+0.944853346337906500,
		-0.681501551465435000,
		+0.402696290353813800,
		+0.328601067704537230,
		+0.812079966562488300,
		+0.756536462138819500,
		+0.011573183932084063,
		-0.551271155078584300,
		+0.207552675260212820,
		+0.795981993703873700,
		-0.242048667319836990,
		+0.419371027177237500,
		+0.885613734373423400,
		+0.669421525553018500,
		+0.393836704188575100,
	},
)

var mtx02 = sparse.NewCOO(
	5,
	5,
	[]int{0, 1, 2, 2, 3, 3, 4, 4, 4, 1, 1, 0, 2},
	[]int{0, 1, 1, 2, 1, 3, 0, 2, 4, 2, 3, 4, 4},
	[]float64{
		+11.0,
		+22.0,
		+23.0,
		+33.0,
		+24.0,
		+44.0,
		+15.0,
		+35.0,
		+55.0,
		+23.0,
		+24.0,
		+15.0,
		+35.0,
	},
)

var mtx03 = sparse.NewCOO(
	5,
	5,
	[]int{3, 4, 4, 1, 0, 2},
	[]int{1, 0, 2, 3, 4, 4},
	[]float64{24.0, 15.0, 35.0, -24.0, -15.0, -35.0},
)

var mtx04 = sparse.NewCOO(
	4,
	5,
	[]int{0, 1, 2},
	[]int{2, 0, 1},
	[]float64{8, -2, 3},
)

var mtx05 = mtx02

var mtx06 = mtx03

var mtx21 = sparse.NewCOO(
	4,
	5,
	[]int{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 3, 3, 3, 3},
	[]int{0, 1, 3, 4, 0, 2, 3, 4, 0, 1, 2, 0, 1, 3, 4},
	[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
)

var mtx22 = sparse.NewCOO(
	5,
	5,
	[]int{0, 1, 2, 2, 3, 3, 4, 4, 4, 1, 1, 0, 2},
	[]int{0, 1, 1, 2, 1, 3, 0, 2, 4, 2, 3, 4, 4},
	[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
)

func TestNewCOO(t *testing.T) {

	m := NewCOO(mtx01)

	assert.True(t, mat.Equal(m.ToMatrix(), mtx01))
}

func TestCOOToCOO(t *testing.T) {

	m := NewCOO(mtx01)

	assert.True(t, mat.Equal(m.ToCOO(), mtx01))
}

func TestCOOToMatrix(t *testing.T) {

	m := NewCOO(mtx01)

	assert.True(t, mat.Equal(m.ToMatrix(), mtx01))
}

func TestCOOMarshalTextTo(t *testing.T) {

	var b strings.Builder

	m := NewCOO(mtx01)

	_, err := m.MarshalTextTo(&b)
	assert.Nil(t, err)

	mm, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-01.mtx"))
	assert.Nil(t, err)

	assert.Equal(t, b.String(), string(mm))
}

func TestCOOMarshalText(t *testing.T) {

	m := NewCOO(mtx01)

	mm1, err := m.MarshalText()
	assert.Nil(t, err)

	mm2, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-01.mtx"))
	assert.Nil(t, err)

	assert.Equal(t, string(mm1), string(mm2))
}

func TestCOOUnmarshalText(t *testing.T) {

	M, N := mtx01.Dims()

	c := sparse.NewCOO(M, N, nil, nil, nil)
	mm := NewCOO(c)

	b, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-01.mtx"))
	assert.Nil(t, err)

	assert.Nil(t, mm.UnmarshalText(b))

	assert.True(t, mat.Equal(mm.ToMatrix(), mtx01))
}

func TestCOOUnmarshalTextFrom(t *testing.T) {

	c := map[string]mat.Matrix{
		"mmtype-01.mtx": mtx01, // real general
		"mmtype-02.mtx": mtx02, // real symmetric
		"mmtype-03.mtx": mtx03, // real skew-symmetric
		"mmtype-04.mtx": mtx04, // integer general
		"mmtype-05.mtx": mtx05, // integer symmetric
		"mmtype-06.mtx": mtx06, // integer skew-symmetric
		"mmtype-21.mtx": mtx21, // pattern general
		"mmtype-22.mtx": mtx22, // pattern symmetric
	}

	for k, v := range c {

		f, _ := os.Open(filepath.Join("testdata", k))
		defer f.Close()

		var mm COO
		if _, err := mm.UnmarshalTextFrom(f); err != nil {
			t.Errorf("%v", err)
		}

		if !mat.Equal(mm.ToMatrix(), v) {
			t.Errorf(
				"\ngot:\n    %s\nwant:\n    %s\n",
				mat.Formatted(mm.ToCOO(), mat.Prefix("    "), mat.Squeeze()),
				mat.Formatted(v, mat.Prefix("    "), mat.Squeeze()),
			)
		}
	}

}

func BenchmarkCOOMarshalTextTo(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		a := sparse.NewCOO(i, i, nil, nil, nil)
		for j := 0; j < i*i; j++ {
			if j%10 < 8 {
				continue
			}
			a.Set(j%i, int(j/i), math.Sqrt(float64(j)))
		}
		m := NewCOO(a)
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				m.MarshalTextTo(io.Discard)
			}
		})
	}
}

func BenchmarkCOOUnmarshalTextFrom(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		a := sparse.NewCOO(i, i, nil, nil, nil)
		for j := 0; j < i*i; j++ {
			if j%10 < 8 {
				continue
			}
			a.Set(j%i, int(j/i), math.Sqrt(float64(j)))
		}
		m := NewCOO(a)
		t, _ := m.MarshalText()
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				m.UnmarshalText(t)
			}
		})
	}
}
