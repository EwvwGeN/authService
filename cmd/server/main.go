package main

import (
	"flag"
	"fmt"

	"github.com/EwvwGeN/authService/internal/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file")
}

func main() {
	flag.Parse()
	c, e := config.LoadConfig(configPath)
	fmt.Printf("%v         %v", c, e)
}