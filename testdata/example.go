// CodeOwner: @go_owner

// This is a single-line comment

/*
This is a multi-line
block comment in Go.
*/

package main

import "fmt"

var x = 42  // Inline comment

func greet(name string) string {
	// Another single-line comment
	return "Hello, " + name + "!"
}

func main() {
	fmt.Println(greet("world"))  /* Inline block comment */
}
