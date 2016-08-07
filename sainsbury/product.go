package sainsbury

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

// Product is the base struct of the product response
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

func (p *Product) getExtendedData(wg *sync.WaitGroup) {
	p.getDescription()
	p.getSize()
	wg.Done()
}

func (p *Product) getDescription() {
	source, _ := getRawHTML(p.url)
	defer source.Close()

	node, err := html.Parse(source)

	if err != nil {
		fmt.Println(err)
		//wg.Done()
		return
	}

	p.findMetaDescription(node)
}

func (p *Product) findMetaDescription(n *html.Node) (ok bool) {
	if n.Data == "meta" && n.Attr[0].Val == "description" {
		p.Description = n.Attr[1].Val
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ok = p.findMetaDescription(c); ok {
			return
		}
	}
	return
}

func (p *Product) getSize() {
	source, _ := getRawHTML(p.url)
	defer source.Close()

	content, _ := ioutil.ReadAll(source)
	size := float64(len(content)) / 1024
	p.Size = strconv.FormatFloat(size, 'f', 2, 64) + "kb"
}
