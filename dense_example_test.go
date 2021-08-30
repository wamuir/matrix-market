package market

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func ExampleDense_MarshalText() {

	// mtx is a dense matrix representation
	mtx := mat.NewDense(4, 2, []float64{
		+0.944853346337906500,
		-0.154091238677780850,
		-0.681501551465435000,
		+0.594570321595631100,
		-0.658745773257358300,
		+0.897566664045815500,
		+0.402696290353813800,
		+0.009438983689089353,
	})

	// m is a Dense matrix initialized with mtx
	m := NewDense(mtx)

	// serialize m into []byte (mm)
	mm, err := m.MarshalText()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(mm))
	// output:
	// %%MatrixMarket matrix array real general
	// %
	//  4  2
	//  0.9448533463379065
	// -0.681501551465435
	// -0.6587457732573583
	//  0.4026962903538138
	// -0.15409123867778085
	//  0.5945703215956311
	//  0.8975666640458155
	//  0.009438983689089353
}

func ExampleDense_UnmarshalText() {

	// mm is a complex-valued matrix in Matrix Market format
	mm := []byte(
		`%%MatrixMarket matrix array real general
		 4 2
		 0.9448533463379065
		 -0.15409123867778085
		 -0.681501551465435
		 0.5945703215956311
		 -0.6587457732573583
		 0.8975666640458155
		 0.4026962903538138
		 0.009438983689089353`,
	)

	// mtx is a dense matrix representation
	mtx := mat.NewDense(4, 2, nil)

	// m is a Dense matrix initialized with mtx
	m := NewDense(mtx)

	// deserialize mm into m
	err := m.UnmarshalText(mm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", mat.Formatted(m.ToMatrix()))
	// output:
	//⎡  0.9448533463379065   -0.6587457732573583⎤
	//⎢-0.15409123867778085    0.8975666640458155⎥
	//⎢  -0.681501551465435    0.4026962903538138⎥
	//⎣  0.5945703215956311  0.009438983689089353⎦
}
