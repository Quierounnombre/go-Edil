package maxdepth

func deep(a int) {
	for a > 0 {
		if a > 1 {
			for a > 2 {
				if a > 3 { // want `nesting deeper than 3`
					a--
				}
			}
		}
	}
}

func shallow(a int) {
	if a > 0 {
		for a > 1 {
			a--
		}
	}
	if a > 5 {
		a--
	} else if a > 4 {
		a++
	} else if a > 3 {
		a--
	}
}
