package market

import (
	"fmt"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

type mm_array_real struct {
	Header header
	M, N   int
	V      []float64
}

type mm_coordinate_real struct {
	Header header
	M, N   int
	I, J   []int
	V      []float64
}

func (m mm_array_real) scan_element(k int, line string) error {

	_, err := fmt.Sscanf(line, "%f", &m.V[k])
	if err != nil {
		return err
	}

	return nil
}

func (m mm_array_real) ToDense() mat.Matrix {

	dense := mat.NewDense(m.M, m.N, nil)

	for k, v := range m.V {
		dense.Set(int(k/m.N), k%m.N, v)
	}

	return dense
}

func (m mm_array_real) ToSparse() *sparse.DOK {

	dok := sparse.NewDOK(m.M, m.N)

	for k, v := range m.V {
		if v != 0 {
			dok.Set(int(k/m.N), k%m.N, v)
		}
	}

	return dok
}

func (m mm_coordinate_real) scan_element(k int, line string) error {

	var i, j int

	_, err := fmt.Sscanf(line, "%d %d %f", &i, &j, &m.V[k])
	if err != nil {
		return err
	}

	m.I[k] = i - 1
	m.J[k] = j - 1

	return nil
}

func (m mm_coordinate_real) ToDense() mat.Matrix {

	sparse := m.ToSparse()
	return sparse.ToDense()
}

func (m mm_coordinate_real) ToSparse() *sparse.DOK {

	dok := sparse.NewDOK(m.M, m.N)

	for k, v := range m.V {
		dok.Set(m.I[k], m.J[k], v)
	}

	return dok
}
