package main

import "template/internal/fixtures"

func main() {
	err := fixtures.LoadRegions(fixtures.Users)
	if err != nil {
		panic(err)
	}
}
