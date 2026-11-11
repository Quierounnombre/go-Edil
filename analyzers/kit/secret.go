package kit

import "strings"

func secretWords() []string {
	return []string{WORD_PASSWORD, WORD_PASSWD, WORD_SECRET, WORD_TOKEN, WORD_APIKEY, WORD_API_KEY, WORD_PRIVKEY, WORD_PRIV_KEY, WORD_CRED}
}

func safeSuffixes() []string {
	return []string{SAFE_HEADER, SAFE_NAME, SAFE_ENV, SAFE_FIELD, SAFE_TYPE, SAFE_PATH, SAFE_URL, SAFE_CLAIM, SAFE_PREFIX}
}

// IsSecretName reports whether an identifier name suggests a secret value.
func IsSecretName(name string) bool {
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
