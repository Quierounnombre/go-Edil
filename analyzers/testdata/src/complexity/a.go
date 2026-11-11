package complexity

func branches(n int) int { // want `cognitive complexity 16`
	if n == 00 {
		n++
	}
	if n == 10 {
		n++
	}
	if n == 20 {
		n++
	}
	if n == 30 {
		n++
	}
	if n == 40 {
		n++
	}
	if n == 50 {
		n++
	}
	if n == 60 {
		n++
	}
	if n == 70 {
		n++
	}
	if n == 80 {
		n++
	}
	if n == 90 {
		n++
	}
	if n == 100 {
		n++
	}
	if n == 110 {
		n++
	}
	if n == 120 {
		n++
	}
	if n == 130 {
		n++
	}
	if n == 140 {
		n++
	}
	if n == 150 {
		n++
	}
	return n
}

func simple(n int) int {
	if n > 0 {
		return n
	}
	return 0
}
