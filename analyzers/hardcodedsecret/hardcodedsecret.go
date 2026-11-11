package hardcodedsecret

import (
	"go/ast"
	"strconv"
	"strings"
	"unicode"

	"edil/analyzers/kit"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "hardcodedsecret",
	Doc:  "no secret-looking string literals assigned to password, token, secret or key names; load them from the environment",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if kit.Skip(pass, file) {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch typed := node.(type) {
			case *ast.AssignStmt:
				checkPairs(pass, typed.Lhs, typed.Rhs)
			case *ast.ValueSpec:
				names := make([]ast.Expr, 0, len(typed.Names))
				for _, name := range typed.Names {
					names = append(names, name)
				}
				checkPairs(pass, names, typed.Values)
			case *ast.KeyValueExpr:
				checkPair(pass, typed.Key, typed.Value)
			}
			return true
		})
	}
	return nil, nil
}

func checkPairs(pass *analysis.Pass, targets []ast.Expr, values []ast.Expr) {
	if len(targets) != len(values) {
		return
	}
	for index, target := range targets {
		checkPair(pass, target, values[index])
	}
}

func checkPair(pass *analysis.Pass, target ast.Expr, value ast.Expr) {
	name := nameOf(target)
	secret, isConst := kit.ConstString(pass, value)
	if name == "" || !isConst || !isSecretName(name) || !looksLikeSecret(secret) {
		return
	}
	kit.Reportf(pass, target.Pos(), "hardcoded secret in %q; load it from the environment or a secret store", name)
}

func nameOf(expr ast.Expr) string {
	switch typed := expr.(type) {
	case *ast.Ident:
		return typed.Name
	case *ast.SelectorExpr:
		return typed.Sel.Name
	case *ast.BasicLit:
		unquoted, err := strconv.Unquote(typed.Value)
		if err != nil {
			return ""
		}
		return unquoted
	}
	return ""
}

func secretWords() []string {
	return []string{WORD_PASSWORD, WORD_PASSWD, WORD_SECRET, WORD_TOKEN, WORD_APIKEY, WORD_API_KEY, WORD_PRIVKEY, WORD_PRIV_KEY, WORD_CRED}
}

func safeSuffixes() []string {
	return []string{SAFE_HEADER, SAFE_NAME, SAFE_ENV, SAFE_FIELD, SAFE_TYPE, SAFE_PATH, SAFE_URL, SAFE_CLAIM, SAFE_PREFIX}
}

func isSecretName(name string) bool {
	lower := strings.ToLower(name)
	for _, suffix := range safeSuffixes() {
		if strings.HasSuffix(lower, suffix) {
			return false
		}
	}
	for _, word := range secretWords() {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

func looksLikeSecret(value string) bool {
	if len(value) < MIN_SECRET_LEN || strings.ContainsAny(value, WHITESPACE) || isEnvName(value) {
		return false
	}
	for _, char := range value {
		if !unicode.IsLetter(char) {
			return true
		}
	}
	return false
}

func isEnvName(value string) bool {
	for _, char := range value {
		if !unicode.IsUpper(char) && char != UNDERSCORE {
			return false
		}
	}
	return true
}
