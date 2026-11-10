package main

type Rule struct {
	Name string
	Doc  string
}

type Page struct {
	Repo  string
	Rules []Rule
}
