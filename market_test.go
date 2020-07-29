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

func TestScanHeader(t *testing.T) {

	// example valid coordinate-integer header
	if h, err := scanHeader(sts(`%%MatrixMarket matrix coordinate integer skew-symmetric`)); err != nil {
		t.Errorf(err.Error())
	} else if !(h.isSparse() && h.isInteger() && h.isSkew()) {
		t.Errorf("Header elements failed to evaluate")
	}

	// example valid coordinate-pattern header
	if h, err := scanHeader(sts(`%%MatrixMarket matrix coordinate pattern symmetric`)); err != nil {
		t.Errorf(err.Error())
	} else if !(h.isSparse() && h.isPattern() && h.isSymmetric()) {
		t.Errorf("Header elements failed to evaluate")
	}

	// example valid array-complex header
	if h, err := scanHeader(sts(`%%MatrixMarket matrix array complex hermitian`)); err != nil {
		t.Errorf(err.Error())
	} else if !(h.isDense() && h.isComplex() && h.isHermitian()) {
		t.Errorf("Header elements failed to evaluate")
	}

	// empty header
	if _, err := scanHeader(sts(``)); err == nil {
		t.Errorf("Expected EOF error, received: %v", err)
	}

	// too few fields in header
	if _, err := scanHeader(sts(`%%MatrixMarket coordinate integer general`)); err == nil {
		t.Errorf("Expected EOF error, received: %v", err)
	}

	// superfluous field(s) in header (expect to be discarded)
	if _, err := scanHeader(sts(`%%MatrixMarket matrix coordinate integer general extra`)); err != nil {
		t.Errorf("Expected nil error, received: %v", err)
	}

	// malformed banner
	if _, err := scanHeader(sts(`MatrixMarket matrix coordinate integer general`)); err != ErrNoHeader {
		t.Errorf("Expected NO_HEADER error, received: %v", err)
	}

	// unsupported object field
	if _, err := scanHeader(sts(`%%MatrixMarket xirtam coordinate integer general`)); err == nil {
		t.Errorf("Expected EOF error, received: %v", err)
	}

	// invalid field combination (real and hermitian)
	if _, err := scanHeader(sts(`%%MatrixMarket matrix coordinate real hermitian`)); err != ErrUnsupportedType {
		t.Errorf("Expected UNSUPPORTED_TYPE error, received: %v", err)
	}

	// invalid field combination (array and pattern)
	if _, err := scanHeader(sts(`%%MatrixMarket matrix array pattern general`)); err != ErrUnsupportedType {
		t.Errorf("Expected UNSUPPORTED_TYPE error, received: %v", err)
	}

}
