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

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

var mtx10 = mat.NewDense(4, 5, []float64{
	+0.944853346337906500,
	+0.328601067704537230,
	+0.011573183932084063,
	-0.242048667319836990,
	+0.885613734373423400,
	-0.681501551465435000,
	+0.812079966562488300,
	-0.551271155078584300,
	-0.000000000000000000,
	-0.000000000000000000,
	-0.000000000000000000,
	-0.000000000000000000,
	+0.207552675260212820,
	-0.000000000000000000,
	+0.669421525553018500,
	+0.402696290353813800,
	+0.756536462138819500,
	+0.795981993703873700,
	+0.419371027177237500,
	+0.393836704188575100,
})

var mtx11 = mat.NewDense(5, 5, []float64{11, 0, 0, 0, 15, 0, 22, 23, 24, 0, 0, 23, 33, 0, 35, 0, 24, 0, 44, 0, 15, 0, 35, 0, 55})
var mtx12 = mat.NewDense(5, 5, []float64{0, 0, 0, 0, -15, 0, 0, -23, -24, 0, 0, 23, 0, 0, -35, 0, 24, 0, 0, 0, 15, 0, 35, 0, 0})
var mtx13 = mat.NewDense(4, 5, []float64{0, 0, 0, 0, 0, 0, -2, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 3, 0, 0})
var mtx14 = mtx11
var mtx15 = mtx12

func TestNewDense(t *testing.T) {

	m := NewDense(mtx10)

	assert.True(t, mat.Equal(m.ToMatrix(), mtx10))
}

func TestDenseToDense(t *testing.T) {

	m := NewDense(mtx10)

	assert.True(t, mat.Equal(m.ToDense(), mtx10))
}

func TestDenseToMatrix(t *testing.T) {

	m := NewDense(mtx10)

	assert.True(t, mat.Equal(m.ToMatrix(), mtx10))
}

func TestDenseMarshalTextTo(t *testing.T) {

	var b strings.Builder

	m := NewDense(mtx10)

	_, err := m.MarshalTextTo(&b)
	assert.Nil(t, err)

	mm, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-10.mtx"))
	assert.Nil(t, err)

	assert.Equal(t, b.String(), string(mm))
}

func TestDenseMarshalText(t *testing.T) {

	m := NewDense(mtx10)

	mm1, err := m.MarshalText()
	assert.Nil(t, err)

	mm2, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-10.mtx"))
	assert.Nil(t, err)

	assert.Equal(t, string(mm1), string(mm2))
}

func TestDenseUnmarshalText(t *testing.T) {

	b, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-10.mtx"))
	assert.Nil(t, err)

	var mm Dense
	assert.Nil(t, mm.UnmarshalText(b))

	if !mat.Equal(mm.ToMatrix(), mtx10) {
		t.Errorf(
			"\ngot:\n    %v\nwant:\n    %v\n",
			mat.Formatted(mm.ToMatrix(), mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(mtx10, mat.Prefix("    "), mat.Squeeze()),
		)
	}
}

func TestDenseUnmarshalTextFrom(t *testing.T) {

	c := map[string]mat.Matrix{
		"mmtype-10.mtx": mtx10, // real general
		"mmtype-11.mtx": mtx11, // real symmetric
		"mmtype-12.mtx": mtx12, // real skew-symmetric
		"mmtype-13.mtx": mtx13, // integer general
		"mmtype-14.mtx": mtx14, // integer symmetric
		"mmtype-15.mtx": mtx15, // integer skew-symmetric
	}

	for k, v := range c {

		f, _ := os.Open(filepath.Join("testdata", k))
		defer f.Close()

		var mm Dense
		if _, err := mm.UnmarshalTextFrom(f); err != nil {
			t.Errorf("%v", err)
		}

		if !mat.Equal(mm.ToMatrix(), v) {
			t.Errorf(
				"\ngot:\n    %v\nwant:\n    %v\n",
				mat.Formatted(mm.ToMatrix(), mat.Prefix("    "), mat.Squeeze()),
				mat.Formatted(v, mat.Prefix("    "), mat.Squeeze()),
			)
		}
	}

}

func BenchmarkDenseMarshalTextTo(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		a := mat.NewDense(i, i, nil)
		for j := 0; j < i*i; j++ {
			a.Set(j%i, int(j/i), math.Sqrt(float64(j)))
		}
		m := NewDense(a)
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				m.MarshalTextTo(io.Discard)
			}
		})
	}
}

func BenchmarkDenseUnmarshalTextFrom(b *testing.B) {
	for i := 1; i <= 1000; i *= 10 {
		a := mat.NewDense(i, i, nil)
		for j := 0; j < i*i; j++ {
			a.Set(j%i, int(j/i), math.Sqrt(float64(j)))
		}
		m := NewDense(a)
		t, _ := m.MarshalText()
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				m.UnmarshalText(t)
			}
		})
	}
}
