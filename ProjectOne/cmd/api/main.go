package main

import (
	"appMove/pkg/config"
	"fmt"
)

func main() {

	// init config (viper)
	cfg := config.Load()

	fmt.Println(cfg)

	// init database (postgre)

	//init http - (gin)

}
