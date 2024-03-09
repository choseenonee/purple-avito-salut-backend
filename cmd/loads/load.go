package main

import "template/internal/fixtures"

func main() {
	err := fixtures.LoadRegions(fixtures.Regions)
	if err != nil {
		panic(err)
	}
}
