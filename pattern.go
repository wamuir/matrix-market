package market

import (
	"fmt"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

type mtxCoordinatePattern struct {
	Header header
	M, N   int
	I, J   []int
}

func (m mtxCoordinatePattern) scan_element(k int, line string) error {

	var i, j int

	_, err := fmt.Sscanf(line, "%d %d", &i, &j)
	if err != nil {
		return err
	}

	m.I[k] = i - 1
	m.J[k] = j - 1

	return nil
}

func (m mtxCoordinatePattern) ToDense() mat.Matrix {

	sparse := m.ToSparse()
	return sparse.ToDense()
}

func (m mtxCoordinatePattern) ToSparse() *sparse.DOK {

	dok := sparse.NewDOK(m.M, m.N)

	for k := range m.I {
		dok.Set(m.I[k], m.J[k], 1)
	}

	return dok
}
