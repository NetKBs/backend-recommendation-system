package main

import "example/config"

func init() {
	config.ConnectDB()
}

func main() {
	defer config.SESSION.Close()
}
