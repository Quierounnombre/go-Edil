package noifinit

func bad(n int) {
	if v := n; v > 0 { // want `init statement in if`
		_ = v
	}
	switch x := n; x { // want `init statement in switch`
	}
}

func good(n int) {
	v := n
	if v > 0 {
		v--
	}
	for i := 0; i < n; i++ {
		v++
	}
}
