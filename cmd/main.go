package main

import "go-proxy"

func main() {
	server := proxy.NewServer("https://dummyjson.com", 3000)
	server.Start()
}
