package deferinloop

func bad(items []int) {
	for range items {
		defer println() // want `defer inside a loop`
	}
	for i := 0; i < 3; i++ {
		defer println(i) // want `defer inside a loop`
	}
}

func good(items []int) {
	defer println()
	for range items {
		func() {
			defer println()
		}()
	}
}
