package market

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

var mtx10 = mat.NewDense(4, 5, []float64{
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

	M, N := mtx10.Dims()

	mtx := mat.NewDense(M, N, nil)
	mm := NewDense(mtx)

	b, err := ioutil.ReadFile(filepath.Join("testdata", "mmtype-10.mtx"))
	assert.Nil(t, err)

	assert.Nil(t, mm.UnmarshalText(b))

	assert.True(t, mat.Equal(mm.ToMatrix(), mtx10))
}

func TestDenseUnmarshalTextFrom(t *testing.T) {

	M, N := mtx10.Dims()

	mtx := mat.NewDense(M, N, nil)
	mm := NewDense(mtx)

	r, err := os.Open(filepath.Join("testdata", "mmtype-10.mtx"))
	assert.Nil(t, err)

	_, err = mm.UnmarshalTextFrom(r)
	assert.Nil(t, err)

	assert.True(t, mat.Equal(mm.ToMatrix(), mtx10))
}
