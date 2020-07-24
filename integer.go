package market

import (
	"fmt"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

type mtxArrayInt struct {
	Header header
	M, N   int
	V      []int
}

type mtxCoordinateInt struct {
	Header header
	M, N   int
	I, J   []int
	V      []int
}

func (m mtxArrayInt) scanElement(k int, line string) error {

	_, err := fmt.Sscanf(line, "%d", &m.V[k])
	if err != nil {
		return err
	}

	return nil
}

func (m mtxArrayInt) ToDense() mat.Matrix {

	dense := mat.NewDense(m.M, m.N, nil)

	for k, v := range m.V {
		dense.Set(int(k/m.N), k%m.N, float64(v))
	}

	return dense
}

func (m mtxArrayInt) ToSparse() *sparse.DOK {

	dok := sparse.NewDOK(m.M, m.N)

	for k, v := range m.V {
		if v != 0 {
			dok.Set(int(k/m.N), k%m.N, float64(v))
		}
	}

	return dok
}

func (m mtxCoordinateInt) scanElement(k int, line string) error {

	var i, j int

	_, err := fmt.Sscanf(line, "%d %d %d", &i, &j, &m.V[k])
	if err != nil {
		return err
	}

	m.I[k] = i - 1
	m.J[k] = j - 1

	return nil
}

func (m mtxCoordinateInt) ToDense() mat.Matrix {

	sparse := m.ToSparse()
	return sparse.ToDense()
}

func (m mtxCoordinateInt) ToSparse() *sparse.DOK {

	dok := sparse.NewDOK(m.M, m.N)

	for k, v := range m.V {
		dok.Set(m.I[k], m.J[k], float64(v))
	}

	return dok
}
