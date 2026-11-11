package errreturn

import (
	"errors"
	"log"
)

func fail() error { return errors.New("x") }

func swallowed() error {
	err := fail()
	if err != nil { // want `checked but neither returned`
		log.Println(err)
	}
	return nil
}

func blank() {
	_ = fail() // want `error assigned to _`
}

func returned() error {
	err := fail()
	if err != nil {
		return err
	}
	return nil
}

func reassigned() error {
	err := fail()
	if err != nil {
		err = errors.New("wrapped")
	}
	return err
}

func deferred() {
	defer func() { _ = fail() }()
}
