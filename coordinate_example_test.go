package market

import (
	"fmt"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

func ExampleCOO_MarshalText() {

	// mtx is a sparse matrix representation in coordinate format
	mtx := sparse.NewCOO(4, 5, nil, nil, nil)
	mtx.Set(0, 0, 0.944853346337906500)
	mtx.Set(1, 1, 0.897566664045815500)
	mtx.Set(2, 2, 0.402696290353813800)

	// m is a COO matrix initialized with mtx
	m := NewCOO(mtx)

	// serialized m into []byte (mm)
	mm, err := m.MarshalText()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(mm))
	// output:
	// %%MatrixMarket matrix coordinate real general
	// %
	//  4  5  3
	//  1  1  0.9448533463379065
	//  2  2  0.8975666640458155
	//  3  3  0.4026962903538138
}

func ExampleCOO_UnmarshalText() {

	// mm is a real-valued sparse matrix in Matrix Market coordinate format
	mm := []byte(
		`%%MatrixMarket matrix coordinate real general
		  3  3  3
		  1  1  0.9448533463379065
		  2  2  0.8975666640458155
		  3  3  0.4026962903538138`,
	)

	// mtx is a coo matrix representation
	mtx := sparse.NewCOO(3, 3, nil, nil, nil)

	// m is a COO matrix initialized with mtx
	m := NewCOO(mtx)

	// deserialize mm into m
	err := m.UnmarshalText(mm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", mat.Formatted(m.ToMatrix()))
	// output:
	// ⎡0.9448533463379065                   0                   0⎤
	// ⎢                 0  0.8975666640458155                   0⎥
	// ⎣                 0                   0  0.4026962903538138⎦
}
