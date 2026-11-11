package main

type group struct{}

func (group) Go(f func()) { f() }

func work() {}

func bad() {
	go work()      // want `go statement`
	go func() {}() // want `go statement`
}

func good(g group) {
	g.Go(work)
}

func main() {
	go work()
}
