package sainsbury

import (
	"code.google.com/p/go.net/html"
	"encoding/json"
	"io"
	"sync"
)

func (p *Products) GetList(url string) error {
	source, err := getRawHTML(url)
	defer source.Close()

	if err != nil {
		return err
	}

	p.parseProducts(source)
	p.getProductExtendedData()

	return nil
}

type Products struct {
	Products []Product `json:"results"`
	Total    float64   `json:"total"`
}

func (p *Products) ToJSON() ([]byte, error) {
	jsonString, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	return jsonString, nil
}

func (p *Products) parseProducts(source io.Reader) {
	tokenizer := html.NewTokenizer(source)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return
		}
		token := tokenizer.Token()
		switch tokenType {
		case html.StartTagToken:
			if token.Data == "div" && len(token.Attr) > 0 && token.Attr[0].Key == "class" && token.Attr[0].Val == "product " {
				p.setProduct(tokenizer)
			}
		case html.TextToken:
		case html.EndTagToken:
		case html.SelfClosingTagToken:

		}
	}
}

func (p *Products) setProduct(tokenizer *html.Tokenizer) {
	product := &Product{}

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		switch tokenType {
		case html.StartTagToken:
			if token.Data == "a" {
				product.setTitleAndProductUrl(&token, tokenizer)
			} else if token.Data == "p" && len(token.Attr) > 0 && token.Attr[0].Val == "pricePerUnit" {
				product.setUnitPrice(&token, tokenizer)
			}

		case html.TextToken:
		case html.EndTagToken:
			if token.Data == "li" {
				p.Products = append(p.Products, *product)
				p.Total += product.UnitPrice
				return
			}
		case html.SelfClosingTagToken:

		}
	}
}

func (p *Products) getProductExtendedData() {
	wg := sync.WaitGroup{}
	wg.Add(len(p.Products))

	for key := range p.Products {
		go p.Products[key].getExtendedData(&wg)
	}

	wg.Wait()
}
