package market

import (
	"bufio"
	"fmt"
	"io"
)

func mtxScanIndex(scanner *bufio.Scanner, t *mmType) (*index, error) {

	switch {

	case t.isArray():

		return mtxScanArrayIndex(scanner, t)

	case t.isCoordinate():

		return mtxScanCoordinateIndex(scanner, t)

	default:

		return nil, ErrUnsupportedType

	}
}

func mtxScanArrayIndex(scanner *bufio.Scanner, t *mmType) (*index, error) {

	var idx index

	for scanner.Scan() {

		line := scanner.Text()

		// blank line or comment (%, Unicode 37)
		if r := []rune(line); len(r) == 0 || r[0] == 37 {
			continue
		}

		n, err := fmt.Sscanf(line, "%d %d", &idx.M, &idx.N)
		if err != nil {
			return nil, err
		}

		if n != 2 {
			return nil, ErrInputScanError
		}

		idx.L = idx.M * idx.N

		break

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &idx, nil

}

func mtxScanCoordinateIndex(scanner *bufio.Scanner, t *mmType) (*index, error) {

	var idx index

	for scanner.Scan() {

		line := scanner.Text()

		// blank line or comment (%, Unicode 37)
		if r := []rune(line); len(r) == 0 || r[0] == 37 {
			continue
		}

		n, err := fmt.Sscanf(line, "%d %d %d", &idx.M, &idx.N, &idx.L)
		if err != nil {
			return nil, err
		}

		if n != 3 {
			return nil, ErrInputScanError
		}

		break

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &idx, nil
}

func makeScanner(r io.Reader) *bufio.Scanner {

	scanner := bufio.NewScanner(r)

	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	return scanner

}

// Read reads Matrix Market inputs from an io.Reader
func Read(r io.Reader) (Matrix, error) {

	var matrix Matrix

	scanner := makeScanner(r)

	// read header
	t, err := scanHeader(scanner)
	if err != nil {
		return nil, err
	}

	// read index
	idx, err := mtxScanIndex(scanner, t)
	if err != nil {
		return nil, err
	}

	// read data
	switch {

	case t.isComplex():

		return nil, ErrUnsupportedType

	case t.isArray() && t.isInteger():

		matrix = &mtxArrayInt{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			V:      make([]int, idx.L),
		}

	case t.isArray() && t.isReal():

		matrix = &mtxArrayReal{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			V:      make([]float64, idx.L),
		}

	case t.isCoordinate() && t.isInteger():

		matrix = &mtxCoordinateInt{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			I:      make([]int, idx.L),
			J:      make([]int, idx.L),
			V:      make([]int, idx.L),
		}

	case t.isCoordinate() && t.isPattern():

		matrix = &mtxCoordinatePattern{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			I:      make([]int, idx.L),
			J:      make([]int, idx.L),
		}

	case t.isCoordinate() && t.isReal():

		matrix = &mtxCoordinateReal{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			I:      make([]int, idx.L),
			J:      make([]int, idx.L),
			V:      make([]float64, idx.L),
		}
	}

	var k int
	for scanner.Scan() {

		line := scanner.Text()

		// blank lines are allowed in data per design spec
		if r := []rune(line); len(r) == 0 {
			continue
		}

		// error out if data rows exceed expected non-zero entries
		// (note that k is zero indexed)
		if k == idx.L {
			return nil, ErrInputScanError
		}

		if err := matrix.scanElement(k, line); err != nil {
			return nil, err
		}

		k++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// check if number of non-empty rows read is equal to expected
	// count of non-zero rows
	if k != idx.L {
		return nil, ErrInputScanError
	}

	return matrix, nil
}

// ReadComplex reads Matrix Market inputs for complex type from an io.Reader
func ReadComplex(r io.Reader) (CMatrix, error) {

	var matrix CMatrix

	scanner := makeScanner(r)

	// read header
	t, err := scanHeader(scanner)
	if err != nil {
		return nil, err
	}

	// read index
	idx, err := mtxScanIndex(scanner, t)
	if err != nil {
		return nil, err
	}

	// read data
	switch {

	case !(t.isComplex()):

		return nil, ErrUnsupportedType

	case t.isArray():

		matrix = &mtxArrayComplex{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			V:      make([]complex128, idx.L),
		}

	case t.isCoordinate():

		matrix = &mtxCoordinateComplex{
			Header: *t,
			M:      idx.M,
			N:      idx.N,
			I:      make([]int, idx.L),
			J:      make([]int, idx.L),
			V:      make([]complex128, idx.L),
		}
	}

	// return mm_scan_complex(scanner, hdx, idx)
	return matrix, nil
}
