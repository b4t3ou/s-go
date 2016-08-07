package main

import (
	"fmt"
	"github.com/b4t3ou/s-go/sainsbury"
)

func main() {
	products := &sainsbury.Products{}
	err := products.GetList("http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html")

	if err != nil {
		fmt.Println(err)
		return
	}

	jsonString, err := products.ToJSON()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonString))
}