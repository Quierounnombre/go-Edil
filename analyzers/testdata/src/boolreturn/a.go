package boolreturn

func isOdd(n int) int { // want `returns only 0 and 1`
	if n%2 == 0 {
		return 0
	}
	return 1
}

func isPositive(n int) bool {
	return n > 0
}

func count(n int) int {
	return n + 1
}

func compare(ok bool) bool {
	if ok == true { // want `comparison with true/false`
		return true
	}
	return false
}
