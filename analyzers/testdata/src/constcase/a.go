package constcase

const maxRetries = 3 // want `"maxRetries" must be SCREAMING_SNAKE_CASE`
const EVEN_IF_IT_IS_A_2 = 2
const _ = 0

const (
	Foo      = iota // want `"Foo" must be`
	BAR__BAZ        // want `"BAR__BAZ" must be`
	TRAIL_          // want `"TRAIL_" must be`
	GOOD_ONE
)

func local() int {
	const inner = 4 // want `"inner" must be`
	return inner
}
