package forum

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type Discuz interface {
	Sign
	FormHash() (string, bool)
}

func FormHash(discuz Discuz) (string, bool) {
	req, _ := http.NewRequest("GET", discuz.BasicUrl(), nil)
	req.Header.Set("Cookie", discuz.Cookie())
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err.Error(), false
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err.Error(), false
	}
	return doc.Find("#scbar_form").Find("input:nth-child(2)").Attr("value")
}
