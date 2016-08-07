package sainsbury

import "testing"

func TestGetProductList(t *testing.T) {
	products := &Products{}
	err := products.GetList("http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html")

	if err != nil {
		t.Error("HTML cannot be parsed")
	}

	if len(products.Products) == 0 {
		t.Error("Products list cannot be empty")
	}
}

func TestGetProductListWithFalseUrl(t *testing.T) {
	products := &Products{}
	_ = products.GetList("http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html2")

	if len(products.Products) != 0 {
		t.Error("Products list has to be empty")
	}
}

func TestGetProductListTotal(t *testing.T) {
	products := &Products{}
	_ = products.GetList("http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html")

	if products.Total == 0 {
		t.Error("Total shouldn't have be zero")
	}
}