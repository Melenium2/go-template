package main

import (
	"log"

	"github.com/Melenium2/go-template/internal/container"
)

func main() {
	c := container.NewContainer()

	if err := c.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
