package logsecret

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

type Config struct{ Token string }

func bad(password string, cfg Config) {
	log.Printf("login %s", password)                 // want `possible secret "password"`
	slog.Info("auth", "token", cfg.Token)            // want `possible secret "token"` `possible secret "Token"`
	fmt.Println(cfg.Token)                           // want `possible secret "Token"`
	_ = fmt.Errorf("bad %s", password)               // want `possible secret "password"`
	log.Println(string(password))                    // want `possible secret "password"`
}

func good(user string, tokenType string) {
	log.Printf("login %s", user)
	log.Println("token expired")
	slog.Info("auth", "user", user, "type", tokenType)
	fmt.Fprintln(os.Stderr, "password rejected")
}
