package main

import (
	"fmt"
	"github.com/paveldroo/go-ogg-packer/oggpacker"
	"log"
)

func main() {
	p, err := oggpacker.NewPacker(8000, 2)
	if err != nil {
		log.Fatalf("%w", err)
	}
	fmt.Println(p)
}
