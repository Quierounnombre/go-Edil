package mathrand

import (
	"math/rand" // want `math/rand is predictable`
	//edil:ignore mathrand jitter only, not security
	_ "math/rand/v2"
	//edil:ignore mathrand
	nr "math/rand" // want `math/rand is predictable`
)

var _ = rand.Int
var _ = nr.Int
