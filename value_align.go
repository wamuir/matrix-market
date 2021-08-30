package market

import (
	"math"
	"strconv"
)

type cmplxAligner struct {
	r floatAligner
	i floatAligner
}

func (a *cmplxAligner) Append(dst []byte, v complex128, fmt byte, p int, bitSize int) []byte {
	dst = a.r.PaddedAppend(dst, real(v), fmt, p, bitSize/2)
	dst = append(dst, ' ')
	dst = a.i.Append(dst, imag(v), fmt, p, bitSize/2)
	return dst
}

func (a *cmplxAligner) PaddedAppend(dst []byte, v complex128, fmt byte, p int, bitSize int) []byte {
	dst = a.r.PaddedAppend(dst, real(v), fmt, p, bitSize/2)
	dst = append(dst, ' ')
	dst = a.i.PaddedAppend(dst, imag(v), fmt, p, bitSize/2)
	return dst
}

func (a *cmplxAligner) Fit(fmt byte, p int, bitSize int) func(i, j int, v complex128) {
	buf := make([]byte, 0, 64)
	return func(_, _ int, v complex128) {
		a.fit(buf, v, fmt, p, bitSize)
	}
}

func (a *cmplxAligner) fit(buf []byte, v complex128, fmt byte, p int, bitSize int) {
	a.r.fit(buf, real(v), fmt, p, bitSize/2)
	a.i.fit(buf, imag(v), fmt, p, bitSize/2)
}

type cmplxTripletAligner struct {
	row intAligner
	col intAligner
	val cmplxAligner
}

func (a cmplxTripletAligner) Append(dst []byte, i, j int, v complex128, fmt byte, p int, bitSize int) []byte {
	dst = a.row.Append(dst, i+1, 10)
	dst = append(dst, ' ')
	dst = a.col.Append(dst, j+1, 10)
	dst = append(dst, ' ')
	dst = a.val.Append(dst, v, fmt, p, bitSize)
	return dst
}

func (a *cmplxTripletAligner) Fit(fmt byte, p int, bitSize int) func(i, j int, v complex128) {
	var buf = make([]byte, 0, bitSize)
	return func(i, j int, v complex128) {
		a.row.fit(i+1, 10)
		a.col.fit(j+1, 10)
		a.val.fit(buf, v, fmt, p, bitSize)
	}
}

type floatAligner [2]int

func (a *floatAligner) Fit(fmt byte, p int, bitSize int) func(i, j int, v float64) {
	buf := make([]byte, 0, bitSize)
	return func(i, j int, v float64) {
		a.fit(buf, v, fmt, p, bitSize)
	}
}

func (a *floatAligner) fit(buf []byte, v float64, fmt byte, p int, bitSize int) {
	buf = buf[:0]
	if v >= 0 || math.IsNaN(v) {
		buf = append(buf, ' ')
	}
	buf = strconv.AppendFloat(buf, v, fmt, p, bitSize)

	a[0] = max(a[0], characteristic(v))
	a[1] = max(a[1], len(buf)-characteristic(v))
}

func (a floatAligner) Append(dst []byte, v float64, fmt byte, p int, bitSize int) []byte {
	for i := 0; i < a[0]-characteristic(v); i++ {
		dst = append(dst, ' ')
	}

	if v >= 0 || math.IsNaN(v) {
		dst = append(dst, ' ')
	}

	return strconv.AppendFloat(dst, v, fmt, p, bitSize)
}

func (a floatAligner) PaddedAppend(dst []byte, v float64, fmt byte, p int, bitSize int) []byte {
	var l int

	l = len(dst)
	for i := 0; i < a[0]-characteristic(v); i++ {
		dst = append(dst, ' ')
	}

	if v >= 0 || math.IsNaN(v) {
		dst = append(dst, ' ')
	}

	dst = strconv.AppendFloat(dst, v, fmt, p, bitSize)

	for (a[0] + a[1]) > len(dst)-l {
		dst = append(dst, ' ')
	}

	return dst
}

type floatTripletAligner struct {
	row intAligner
	col intAligner
	val floatAligner
}

func (a floatTripletAligner) Append(dst []byte, i, j int, v float64, fmt byte, p int, bitSize int) []byte {
	dst = a.row.Append(dst, i+1, 10)
	dst = append(dst, ' ')
	dst = a.col.Append(dst, j+1, 10)
	dst = append(dst, ' ')
	dst = a.val.Append(dst, v, fmt, p, bitSize)
	return dst
}

func (a *floatTripletAligner) Fit(fmt byte, p int, bitSize int) func(i, j int, v float64) {
	var buf = make([]byte, 0, bitSize)
	return func(i, j int, v float64) {
		a.row.fit(i+1, 10)
		a.col.fit(j+1, 10)
		a.val.fit(buf, v, fmt, p, bitSize)
	}
}

type intAligner int

func (a intAligner) Append(dst []byte, v int, base int) []byte {
	if v >= 0 {
		dst = append(dst, ' ')
	}

	for i := 0; i < (int(a) - characteristic(float64(v))); i++ {
		dst = append(dst, ' ')
	}

	return strconv.AppendInt(dst, int64(v), base)
}

func (a *intAligner) Fit(base int) func(i, j int, v int) {
	return func(i, j int, v int) {
		a.fit(v, base)
	}
}

func (a *intAligner) fit(v int, base int) {
	if v == 0 {
		*a = intAligner(max(int(*a), 2))
	}
	*a = intAligner(
		max(
			int(*a),
			int(2+math.Log10(math.Abs(float64(v)))),
		),
	)
}

type intTripletAligner struct {
	row intAligner
	col intAligner
	val intAligner
}

func (a intTripletAligner) Append(dst []byte, i, j, v int, base int) []byte {
	dst = a.row.Append(dst, i+1, 10)
	dst = append(dst, ' ')
	dst = a.col.Append(dst, j+1, 10)
	dst = append(dst, ' ')
	dst = a.val.Append(dst, v, base)
	return dst
}

func (a *intTripletAligner) Fit(base int) func(i, j, v int) {
	return func(i, j, v int) {
		a.row.fit(i+1, 10)
		a.col.fit(j+1, 10)
		a.val.fit(v, 10)
	}
}

// characteristic counts the number of characters to the left of the decimal,
// always adding one to account for a potential sign.  This function is only
// useful when formatting in decimal point notation (%f/%F); i.e., will return
// incorrect counts for scientific or hex notations.  NaN and Inf are assumed
// left of the decimal.
func characteristic(v float64) int {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 4
	}
	if math.Abs(v) < 1 {
		return 2
	}
	return int(2 + math.Log10(math.Abs(v)))
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
