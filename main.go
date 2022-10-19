package main

import "github.com/wenlaizhou/middleware"

func main() {
	// 1. 包名, 2. dist路径
	middleware.StartServer("", 8080)
}
