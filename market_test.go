package market

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sts(s string) *bufio.Scanner {

	r := strings.NewReader(s)
	return bufio.NewScanner(r)
}

func TestMmTypeIndex(t *testing.T) {

	var n mmType
	assert.Equal(t, n.index(), -1)
}

func TestScanHeader(t *testing.T) {

	var (
		h   *mmType
		err error
	)

	// example valid arry real general
	h, err = scanHeader(sts(`%%MatrixMarket matrix array real general`))
	assert.Nil(t, err)
	assert.True(t, (h.isArray() && h.isReal() && h.isGeneral()))

	// example valid coordinate-integer header
	h, err = scanHeader(sts(`%%MatrixMarket matrix coordinate integer skew-symmetric`))
	assert.Nil(t, err)
	assert.True(t, (h.isSparse() && h.isInteger() && h.isSkew()))

	// example valid coordinate-pattern header
	h, err = scanHeader(sts(`%%MatrixMarket matrix coordinate pattern symmetric`))
	assert.Nil(t, err)
	assert.True(t, (h.isSparse() && h.isPattern() && h.isSymmetric()))

	// example valid array-complex header
	h, err = scanHeader(sts(`%%MatrixMarket matrix array complex hermitian`))
	assert.Nil(t, err)
	assert.True(t, (h.isDense() && h.isComplex() && h.isHermitian()))

	// empty header
	_, err = scanHeader(sts(``))
	assert.EqualError(t, err, ErrInputScanError.Error())

	// too few fields in header
	_, err = scanHeader(sts(`%%MatrixMarket coordinate integer general`))
	assert.EqualError(t, err, ErrPrematureEOF.Error())

	// superfluous field(s) in header (expect to be discarded)
	_, err = scanHeader(sts(`%%MatrixMarket matrix coordinate integer general extra`))
	assert.Nil(t, err)

	// malformed banner
	_, err = scanHeader(sts(`MatrixMarket matrix coordinate integer general`))
	assert.EqualError(t, err, ErrNoHeader.Error())

	// unsupported object field
	_, err = scanHeader(sts(`%%MatrixMarket xirtam coordinate integer general`))
	assert.EqualError(t, err, ErrUnsupportedType.Error())

	// invalid field combination (real and hermitian)
	_, err = scanHeader(sts(`%%MatrixMarket matrix coordinate real hermitian`))
	assert.EqualError(t, err, ErrUnsupportedType.Error())

	// invalid field combination (array and pattern)
	_, err = scanHeader(sts(`%%MatrixMarket matrix array pattern general`))
	assert.EqualError(t, err, ErrUnsupportedType.Error())
}
