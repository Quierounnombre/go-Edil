package main

type Rule struct {
	Name    string
	Doc     string
	Example string
}

type Page struct {
	Repo  string
	Rules []Rule
}

// output is one generated file: where it goes and the template that makes it.
type output struct {
	path     string
	template string
}
