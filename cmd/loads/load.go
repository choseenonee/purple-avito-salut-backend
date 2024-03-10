package main

import "template/internal/fixtures"

func main() {
	err := fixtures.Loader(fixtures.Regions, fixtures.Users, fixtures.UsersSegments)
	if err != nil {
		panic(err)
	}
}
