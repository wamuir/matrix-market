package market

import (
	"fmt"
	"testing"

	"github.com/james-bowman/sparse"
)

func TestNewCOOFrom(t *testing.T) {

	c := sparse.NewCOO(2, 2, nil, nil, nil)
	c.Set(0, 1, 100)

	m := NewMMCOO(c)

	//buf := make([]byte, 0)

	//c.TextUnmarshaler(buf)

	fmt.Printf("%v\n", m.ToCOO())
}

func TestCOOTextUnmarshaler(t *testing.T) {

	var m MMCOO

	m.TextUnmarshaler(make([]byte, 0))

	fmt.Printf("%v\n", m.ToCOO())
}
