package util

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func GetText(res *http.Response, success string, fail string) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(res.Body)
	if err != nil {
		return err.Error()
	}
	html := buf.String()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	if element := doc.Find(success); element.Size() == 0 {
		if element = doc.Find(fail); element.Size() != 0 {
			return element.Text()
		} else {
			return "无法解析HTML"
		}
	} else {
		return element.Text()
	}
}
