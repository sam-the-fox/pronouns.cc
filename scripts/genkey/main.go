package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	b := make([]byte, 64)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	s := base64.URLEncoding.EncodeToString(b)

	fmt.Println(s)
}
