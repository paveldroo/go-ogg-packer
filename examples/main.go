package main

import (
	"fmt"
	"log"

	"github.com/paveldroo/go-ogg-packer"
)

func main() {
	p, err := oggpacker.NewPacker(8000, 2)
	if err != nil {
		log.Fatalf("create ogg packer: %s", err)
	}
	fmt.Println(p)
}
