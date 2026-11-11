package layout

type Loose struct{} // want `type "Loose" must live in type.go`

const LOOSE = 2 // want `constants must live in defines.go`
