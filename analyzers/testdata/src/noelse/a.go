package noelse

func bad(n int) int {
	if n > 0 {
		return 1
	} else { // want `else after a block that exits`
		return 2
	}
}

func good(n int) int {
	if n > 0 {
		n++
	} else {
		n--
	}
	return n
}
