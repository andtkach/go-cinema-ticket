package main

import "os"

const DefaultServerVersion = "0.0.1"

func getServerVersion() string {
	if v := os.Getenv("SERVER_VERSION"); v != "" {
		return v
	}
	return DefaultServerVersion
}
