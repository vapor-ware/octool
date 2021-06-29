package main

import (
	"math/rand"
	"time"

	"github.com/vapor-ware/octool/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
