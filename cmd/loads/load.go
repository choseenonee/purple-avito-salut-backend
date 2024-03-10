package main

import "template/internal/fixtures"

func main() {
	var err error

	//err = fixtures.LoadRegions(fixtures.Regions)
	//if err != nil {
	//	panic(err)
	//}

	err = fixtures.Loader(fixtures.UsersSegments)
	if err != nil {
		panic(err)
	}
}
