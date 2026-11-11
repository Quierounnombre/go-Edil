package main

import "context"

func helper() context.Context {
	return context.Background() // want `root context created outside main`
}

func todo() context.Context {
	return context.TODO() // want `root context created outside main`
}

func main() {
	_ = context.Background()
}
