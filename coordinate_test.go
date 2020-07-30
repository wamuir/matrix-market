package market

import (
	"io/ioutil"
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

func TestNewCOO(t *testing.T) {

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

	M, N := mtx01.Dims()

	c := sparse.NewCOO(M, N, nil, nil, nil)
	mm := NewCOO(c)

	r, err := os.Open(filepath.Join("testdata", "mmtype-01.mtx"))
	assert.Nil(t, err)

	_, err = mm.UnmarshalTextFrom(r)
	assert.Nil(t, err)

	assert.True(t, mat.Equal(mm.ToMatrix(), mtx01))
}
