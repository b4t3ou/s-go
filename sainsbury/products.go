package sainsbury

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"io"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
)

func (p *Products) GetList(url string) error {
	source, err := getRawHTML(url)
	defer source.Close()

	if err != nil {
		return err
	}

	p.parseProducts(source)

	return nil
}

func getRawHTML(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
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

type Product struct {
	Title       string  `json:"title"`
	Size        string  `json:"size"`
	UnitPrice   float64 `json:"unit_price"`
	Description string  `json:"description"`
	url         string
}

func (p *Product) setTitleAndProductUrl(token *html.Token, tokenizer *html.Tokenizer) {
	if token.Data == "a" {
		p.url = token.Attr[0].Val
		tokenizer.NextIsNotRawText()
		tokenizer.Next()
		p.Title = strings.TrimSpace(html.UnescapeString(tokenizer.Token().String()))
	}
}

func (p *Product) setUnitPrice(token *html.Token, tokenizer *html.Tokenizer) {
	tokenizer.NextIsNotRawText()
	tokenizer.Next()
	priceString := strings.TrimSpace(tokenizer.Token().String())
	floatPrice, err := strconv.ParseFloat(priceString[2:len(priceString)], 64)

	if err != nil {
		return
	}

	p.UnitPrice = floatPrice
}