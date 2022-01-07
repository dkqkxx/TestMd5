package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	fmt.Println("MD5 hacking...")
	msg := []byte("0123456789ABCDEF0123456789ABCDEF")
	sum := md5.Sum(msg)
	fmt.Println(fmt.Sprintf("%02X", sum))
}
