package market

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const MaxScanTokenSize = 64 * 1024

var (
	INPUT_SCAN_ERROR = fmt.Errorf("error while scanning matrix input") // 11
	LINE_TOO_LONG    = fmt.Errorf("input line exceed maximum length ") // 16
	PREMATURE_EOF    = fmt.Errorf("required header items are missing") // 12
	NO_HEADER        = fmt.Errorf("missing matrix market header line") // 14
	NOT_MTX          = fmt.Errorf("input is not a matrix market")      // 13
	UNSUPPORTED_TYPE = fmt.Errorf("unrecognizable matrix description") // 15
	UNWRITABLE       = fmt.Errorf("unable to write matrix to file")    // 17
)

var supported = []header{
	{MM_MTX_STR, MM_COORDINATE_STR, MM_REAL_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_REAL_STR, MM_SYMM_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_REAL_STR, MM_SKEW_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_INT_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_INT_STR, MM_SYMM_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_INT_STR, MM_SKEW_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_COMPLEX_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_COMPLEX_STR, MM_SYMM_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_COMPLEX_STR, MM_SKEW_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_REAL_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_REAL_STR, MM_SYMM_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_REAL_STR, MM_SKEW_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_INT_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_INT_STR, MM_SYMM_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_INT_STR, MM_SKEW_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_COMPLEX_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_COMPLEX_STR, MM_SYMM_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_COMPLEX_STR, MM_SKEW_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_COMPLEX_STR, MM_HERM_STR},
	{MM_MTX_STR, MM_ARRAY_STR, MM_COMPLEX_STR, MM_HERM_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_PATTERN_STR, MM_GENERAL_STR},
	{MM_MTX_STR, MM_COORDINATE_STR, MM_PATTERN_STR, MM_SYMM_STR},
}

func mm_scan_header(scanner *bufio.Scanner) (*header, error) {

	if ok := scanner.Scan(); !ok {
		return nil, INPUT_SCAN_ERROR
	}

	var banner, object, format, field, symm string

	n, err := fmt.Sscan(scanner.Text(), &banner, &object, &format, &field, &symm)
	if err != nil {
		return nil, err
	}

	if n != 5 {
		return nil, PREMATURE_EOF
	}

	if banner != MatrixMktBanner {
		return nil, NO_HEADER
	}

	h := header{
		Object:   strings.ToLower(object),
		Format:   strings.ToLower(format),
		Field:    strings.ToLower(field),
		Symmetry: strings.ToLower(symm),
	}

	if !h.isValid() {
		return nil, UNSUPPORTED_TYPE
	}

	return &h, nil
}

func mm_scan_index(scanner *bufio.Scanner, hdr *header) (*index, error) {

	switch {

	case hdr.isArray():

		return mm_scan_array_index(scanner, hdr)

	case hdr.isCoordinate():

		return mm_scan_coordinate_index(scanner, hdr)

	default:

		return nil, UNSUPPORTED_TYPE

	}
}

func mm_scan_array_index(scanner *bufio.Scanner, hdr *header) (*index, error) {

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
			return nil, INPUT_SCAN_ERROR
		}

		idx.L = idx.M * idx.N

		break

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &idx, nil

}

func mm_scan_coordinate_index(scanner *bufio.Scanner, hdr *header) (*index, error) {

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
			return nil, INPUT_SCAN_ERROR
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

	buf := make([]byte, MaxScanTokenSize)
	scanner.Buffer(buf, MaxScanTokenSize)

	return scanner

}

// Read reads Matrix Market inputs from an io.Reader
func Read(r io.Reader) (Matrix, error) {

	var matrix Matrix

	scanner := makeScanner(r)

	// read header
	hdr, err := mm_scan_header(scanner)
	if err != nil {
		return nil, err
	}

	// read index
	idx, err := mm_scan_index(scanner, hdr)
	if err != nil {
		return nil, err
	}

	// read data
	switch {

	case hdr.isComplex():

		return nil, UNSUPPORTED_TYPE

	case hdr.isArray() && hdr.isInteger():

		matrix = &mm_array_int{
			Header: *hdr,
			M:      idx.M,
			N:      idx.N,
			V:      make([]int, idx.L),
		}

	case hdr.isArray() && hdr.isReal():

		matrix = &mm_array_real{
			Header: *hdr,
			M:      idx.M,
			N:      idx.N,
			V:      make([]float64, idx.L),
		}

	case hdr.isCoordinate() && hdr.isInteger():

		matrix = &mm_coordinate_int{
			Header: *hdr,
			M:      idx.M,
			N:      idx.N,
			I:      make([]int, idx.L),
			J:      make([]int, idx.L),
			V:      make([]int, idx.L),
		}

	case hdr.isCoordinate() && hdr.isPattern():

		matrix = &mm_coordinate_pattern{
			Header: *hdr,
			M:      idx.M,
			N:      idx.N,
			I:      make([]int, idx.L),
			J:      make([]int, idx.L),
		}

	case hdr.isCoordinate() && hdr.isReal():

		matrix = &mm_coordinate_real{
			Header: *hdr,
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
			return nil, INPUT_SCAN_ERROR
		}

		if err := matrix.scan_element(k, line); err != nil {
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
		return nil, INPUT_SCAN_ERROR
	}

	return matrix, nil
}

// ReadComplex reads Matrix Market inputs for complex type from an io.Reader
func ReadComplex(r io.Reader) (CMatrix, error) {

	var matrix CMatrix

	scanner := makeScanner(r)

	// read header
	hdr, err := mm_scan_header(scanner)
	if err != nil {
		return nil, err
	}

	// read index
	idx, err := mm_scan_index(scanner, hdr)
	if err != nil {
		return nil, err
	}

	// read data
	switch {

	case !(hdr.isComplex()):

		return nil, UNSUPPORTED_TYPE

	case hdr.isArray():

		matrix = &mm_array_complex{
			Header: *hdr,
			M:      idx.M,
			N:      idx.N,
			V:      make([]complex128, idx.L),
		}

	case hdr.isCoordinate():

		matrix = &mm_coordinate_complex{
			Header: *hdr,
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
