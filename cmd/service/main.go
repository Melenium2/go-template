package main

import (
	"log"

	"github.com/Melenium2/go-tempalte/internal/container"
)

func main() {
	c := container.NewContainer()

	if err := c.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
