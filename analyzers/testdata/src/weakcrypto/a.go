package weakcrypto

import (
	"crypto/md5" // want `crypto/md5 is broken`
	"crypto/sha256"
)

var _ = md5.New
var _ = sha256.New
