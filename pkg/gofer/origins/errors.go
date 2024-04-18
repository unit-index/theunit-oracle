package origins

import (
	"errors"
	"fmt"
)

var ErrEmptyOriginResponse = fmt.Errorf("empty origin response received")
var ErrMissingResponseForPair = fmt.Errorf("no response for pair from origin")
var ErrInvalidResponseStatus = fmt.Errorf("invalid response status from origin")
var ErrInvalidResponse = fmt.Errorf("invalid response from origin")
var ErrInvalidPrice = fmt.Errorf("invalid price from origin")
var ErrUnknownOrigin = errors.New("unknown origin")
