package market

import (
	"fmt"
)

// Error codes returned by failures to write a matrix
var (
	ErrUnwritable = fmt.Errorf("unable to write matrix to file")
)
