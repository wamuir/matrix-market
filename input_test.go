package market

import (
	"bufio"
	"strings"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func sts(s string) *bufio.Scanner {

	r := strings.NewReader(s)
	return makeScanner(r)
}

func Test_mm_read_header(t *testing.T) {

	// valid header
	if _, err := mm_scan_header(sts(`%%MatrixMarket matrix coordinate integer general`)); err != nil {
		t.Errorf(err.Error())
	}

	// empty header
	if h, err := mm_scan_header(sts(``)); err == nil {
		t.Errorf("Expected EOF error, received %v, %v", h, err)
	}

	// object field missing from header
	if h, err := mm_scan_header(sts(`%%MatrixMarket coordinate integer general`)); err == nil {
		t.Errorf("Expected EOF error, received %v, %v", h, err)
	}

	// malformed banner
	if h, err := mm_scan_header(sts(`MatrixMarket matrix coordinate integer general`)); err != NO_HEADER {
		t.Errorf("Expected NO_HEADER error, received %v, %v", h, err)
	}

	// invalid field combination (real and hermitian)
	if h, err := mm_scan_header(sts(`%%MatrixMarket matrix coordinate real hermitian`)); err != UNSUPPORTED_TYPE {
		t.Errorf("Expected UNSUPPORTED_TYPE error, received %v, %v", h, err)
	}

	// invalid field combination (array and pattern)
	if h, err := mm_scan_header(sts(`%%MatrixMarket matrix array pattern general`)); err != UNSUPPORTED_TYPE {
		t.Errorf("Expected UNSUPPORTED_TYPE error, received %v, %v", h, err)
	}

	// superfluous field(s) in header (expect to be discarded)
	if h, err := mm_scan_header(sts(`%%MatrixMarket matrix coordinate integer general extra`)); err != nil {
		t.Errorf("Expected nil error, received %v, %v", h, err)
	}

}

func TestRead(t *testing.T) {

	mmx := `%%MatrixMarket matrix coordinate integer general
 4 5 3
 1 3 1
 2 1 1
 3 2 1
`

	m, err := Read(strings.NewReader(mmx))
	if err != nil {
		t.Errorf("Received unexpected error %v", err.Error())
		return
	}

	d := mat.NewDense(4, 5, []float64{
		0, 0, 1, 0, 0,
		1, 0, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	})

	if !(mat.Equal(d, m.ToDense())) {
		t.Errorf(
			"matrices differ\nexpected \nm = %v\n\nreceived \nm = %v\n\n",
			mat.Formatted(m.ToDense(), mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(d, mat.Prefix("    "), mat.Squeeze()),
		)
	}
}
