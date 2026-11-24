package main

import "fmt"

func main() {
	salarios := map[string]int{
		"João":  1000,
		"Maria": 2000,
		"Pedro": 3000,
	}
	for nome, salario := range salarios {
		fmt.Printf("O salário do %s é %d\n", nome, salario)
	}

	//salario2 := make(map[string]int)
	//salario3 := map[string]int{}
}
