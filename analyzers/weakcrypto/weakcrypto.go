package weakcrypto

import (
	"strconv"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "weakcrypto",
	Doc:  "no md5, sha1, des or rc4; use sha256+, AES-GCM or ChaCha20",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		for _, spec := range file.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return nil, err
			}
			switch path {
			case MD5, SHA1, DES, RC4:
				kit.Reportf(pass, spec.Pos(), "%s is broken for security use; pick a modern primitive", path)
			}
		}
	}
	return nil, nil
}
