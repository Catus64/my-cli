package main

import (
	"encoding/hex"
	"fmt"
)

func main() {

	s := "hello \nworld"

	b := []byte(s)

	hexstr := hex.EncodeToString(b)

	fmt.Println("str: ", s)
	fmt.Println("binary:", b)
	fmt.Println("hex:", hexstr)

	for _, by := range b {
		fmt.Println("byte:", by, "Hex : ", hex.EncodeToString([]byte{by}), " str: ", string(by))
	}
}
