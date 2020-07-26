package market

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type mtxArrayComplex struct {
	Header mmType
	M, N   int
	V      []complex128
}

type mtxCoordinateComplex struct {
	Header mmType
	M, N   int
	I, J   []int
	V      []complex128
}

func (m mtxArrayComplex) scanElement(k int, line string) error {

	var (
		a, b float64 // real, imaginary
	)

	_, err := fmt.Sscanf(line, "%f %f", &a, &b)
	if err != nil {
		return err
	}

	m.V[k] = complex(a, b)

	return nil
}

func (m mtxArrayComplex) ToDense() mat.CMatrix {

	return mat.NewCDense(m.M, m.N, m.V)
}

func (m mtxCoordinateComplex) scanElement(k int, line string) error {

	var (
		i, j int
		a, b float64 // real, imaginary
	)

	_, err := fmt.Sscanf(line, "%d %d %f %f", &i, &j, &a, &b)
	if err != nil {
		return err
	}

	m.I[k] = i - 1
	m.J[k] = j - 1
	m.V[k] = complex(a, b)

	return nil
}

func (m mtxCoordinateComplex) ToDense() mat.CMatrix {

	dense := mat.NewCDense(m.M, m.N, nil)
	for k, v := range m.V {
		dense.Set(m.I[k], m.J[k], v)
	}
	return dense
}
