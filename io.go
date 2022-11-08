package main

import (
	"os"
)

func readFile(s string) (string, error) {
	data, err := os.ReadFile(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
