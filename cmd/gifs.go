package main

import (
	"fmt"

	"github.com/jirwin/gifs-quadlek/src"
)

func main() {
	g := gifs.NewGifs("", "PG-13")
	random, err := g.Search("the more you know")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(random)
}
