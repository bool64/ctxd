// Package internal is something special.
package internal

// Foo foes.
type Foo string

type foo string

func doSomething(f Foo) {
	f = "something"
}
