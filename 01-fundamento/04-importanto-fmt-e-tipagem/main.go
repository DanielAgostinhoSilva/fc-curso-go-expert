package main

import "fmt"

const a = "Hello, World!"

type ID int

var (
	b bool    = true
	c int     = 10
	d string  = "Daniel"
	e float64 = 1.2
	f ID      = 1
)

func main() {
	fmt.Printf(" O tipo da variavel A é %T\n", a)
	fmt.Printf(" O tipo da variavel B é %T\n", b)
	fmt.Printf(" O tipo da variavel C é %T\n", c)
	fmt.Printf(" O tipo da variavel D é %T\n", d)
	fmt.Printf(" O tipo da variavel E é %T\n", e)
	fmt.Printf(" O tipo da variavel F é %T\n", f)
}
