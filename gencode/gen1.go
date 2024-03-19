package gencode

import (
	_ "embed"
	"fmt"
)

//go:embed hello.txt
var s []byte

func Gen1() {
	fmt.Println(s)
}
