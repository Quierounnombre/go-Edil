package httpdefault

import "net/http"

func bad() {
	http.Get("http://x")   // want `Get has no timeout`
	_ = http.DefaultClient // want `DefaultClient has no timeout`
	_ = &http.Client{}     // want `without Timeout`
}

func good() {
	client := &http.Client{Timeout: 5}
	client.Get("http://x")
}
