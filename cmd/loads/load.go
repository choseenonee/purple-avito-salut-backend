package main

import "template/internal/fixtures"

func main() {
	err := fixtures.Loader(fixtures.Microcategories)
	if err != nil {
		panic(err)
	}
}
