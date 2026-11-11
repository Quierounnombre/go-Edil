package errwrap

// verbAt is a formatting verb and its byte offset inside the format string.
type verbAt struct {
	verb   rune
	offset int
}
