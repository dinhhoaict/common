package webutil

import (
	"github.com/dinhhoaict/common/log"
	"io"
)
import "github.com/PuerkitoBio/goquery"

func GetDocument(b io.Reader) (*goquery.Document, error){
	doc, err := goquery.NewDocumentFromReader(b)
	return doc, err
}

func FindElements(doc *goquery.Document, query string) ([]*goquery.Selection, error){
	var ret []*goquery.Selection
	logger := log.Logger()
	doc.Find(query).Each(func(i int, selection *goquery.Selection) {
		logger.Infof("%#v", selection.Find("content").Text())
		content, _ := selection.Attr("content")
		logger.Infof("%#v", content)
		ret = append(ret, selection)
	})
	return ret, nil
}