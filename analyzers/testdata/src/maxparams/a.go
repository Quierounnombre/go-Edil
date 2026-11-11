package maxparams

import "context"

func many(a, b, c, d, e, f int) {} // want `has 6 parameters`

func five(a, b, c, d, e int) {}

func withContext(ctx context.Context, a, b, c, d, e int) {}
