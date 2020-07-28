package market

import (
	"bufio"
	"strings"
	"testing"
)

func sts(s string) *bufio.Scanner {

	r := strings.NewReader(s)
	return bufio.NewScanner(r)
}

func TestMMTypeIsMMType(t *testing.T) {

	for i := range supported {
		for j := range supported {
			if i == j {
				continue
			}
			if supported[i].isMMType(&supported[j]) != (i == j) {
				t.Errorf(
					"Expected equality of types to evaluate %v, received: %v",
					(i == j),
					!(i == j),
				)
			}

		}
	}

}

func TestScanHeader(t *testing.T) {

	// valid header
	if _, err := scanHeader(sts(`%%MatrixMarket matrix coordinate integer general`)); err != nil {
		t.Errorf(err.Error())
	}

	// empty header
	if _, err := scanHeader(sts(``)); err == nil {
		t.Errorf("Expected EOF error, received: %v", err)
	}

	// object field missing from header
	if _, err := scanHeader(sts(`%%MatrixMarket coordinate integer general`)); err == nil {
		t.Errorf("Expected EOF error, received: %v", err)
	}

	// malformed banner
	if _, err := scanHeader(sts(`MatrixMarket matrix coordinate integer general`)); err != ErrNoHeader {
		t.Errorf("Expected NO_HEADER error, received: %v", err)
	}

	// invalid field combination (real and hermitian)
	if _, err := scanHeader(sts(`%%MatrixMarket matrix coordinate real hermitian`)); err != ErrUnsupportedType {
		t.Errorf("Expected UNSUPPORTED_TYPE error, received: %v", err)
	}

	// invalid field combination (array and pattern)
	if _, err := scanHeader(sts(`%%MatrixMarket matrix array pattern general`)); err != ErrUnsupportedType {
		t.Errorf("Expected UNSUPPORTED_TYPE error, received: %v", err)
	}

	// superfluous field(s) in header (expect to be discarded)
	if _, err := scanHeader(sts(`%%MatrixMarket matrix coordinate integer general extra`)); err != nil {
		t.Errorf("Expected nil error, received: %v", err)
	}
}
