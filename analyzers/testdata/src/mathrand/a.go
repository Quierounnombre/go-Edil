package mathrand

import (
	"math/rand" // want `math/rand is predictable`
	_ "math/rand/v2" // want `math/rand/v2 is predictable`
)

var _ = rand.Int
