package main

import "strings"

func parseTesters(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	return lines
}
