package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {

	data := []byte("hello world")

	hash := sha1.Sum(data)

	hash_hex := fmt.Sprintf("%x", hash[:])

	fmt.Printf("%x\n", hash)
	fmt.Println(hash_hex)
}
