package nakedreturn

func named() (n int, err error) {
	n = 1
	return // want `bare return`
}

func explicit() (n int, err error) {
	return 1, nil
}

func unnamed() int {
	return 1
}
