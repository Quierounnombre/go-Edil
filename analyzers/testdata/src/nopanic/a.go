package main

import (
	"log"
	"os"
)

func bad() {
	panic("x")     // want `panic is forbidden`
	log.Panic("x") // want `Panic is forbidden`
	os.Exit(1)     // want `Exit is only allowed in main`
	log.Fatal("x") // want `Fatal is only allowed in main`
}

func main() {
	os.Exit(1)
}
