package mapt

import "fmt"

func Mapt() {
	a := make(map[string]int)
	a["hello"] = 1
	a["world"] = 2
	fmt.Println(a)

	b := make(map[string]map[string]int)
	bx := make(map[string]int)
	// var bx map[string]int
	bx["banana"] = 123
	b["apple"] = bx
	fmt.Println(b)
	fmt.Println(b["apple"])
	fmt.Println(b["apple"]["banana"])
}
