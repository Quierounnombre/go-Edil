package noglobals

import _ "embed"

var counter = 0 // want `package-level var "counter"`

var _ int = 1

//go:embed fixture.txt
var Embedded string

func local() {
	var x int
	_ = x
}
