package main

import (
	"APIGetway/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustNew()
	fmt.Printf("%+v", cfg)

}
