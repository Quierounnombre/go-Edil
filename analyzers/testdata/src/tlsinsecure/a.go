package tlsinsecure

import "crypto/tls"

func bad() {
	_ = &tls.Config{InsecureSkipVerify: true} // want `InsecureSkipVerify disables`
	cfg := &tls.Config{}
	cfg.InsecureSkipVerify = true // want `InsecureSkipVerify disables`
}

func good() {
	_ = &tls.Config{InsecureSkipVerify: false}
}
