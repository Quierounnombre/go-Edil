package hardcodedsecret

import "os"

const (
	apiKey       = "AKIAIOSFODNN7EXAMPLE1" // want `hardcoded secret in "apiKey"`
	TOKEN_HEADER = "Authorization"
	API_TOKEN    = "API_TOKEN_ENV"
)

type Config struct{ Secret string }

var dbPassword = "hunter2-hunter2" // want `hardcoded secret in "dbPassword"`

func bad() Config {
	return Config{Secret: "s3cr3t-value!"} // want `hardcoded secret in "Secret"`
}

func good() {
	password := os.Getenv("PASSWORD")
	tokenType := "Bearer"
	message := "invalid password supplied"
	_, _, _ = password, tokenType, message
}
