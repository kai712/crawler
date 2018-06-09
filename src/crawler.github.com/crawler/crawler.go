package crawler

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

// GetHTML 获取html
func GetHTML(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// SetUTF8 将编码转为UTF-8
func SetUTF8(content io.Reader) (*transform.Reader, error) {
	bytes, err := bufio.NewReader(content).Peek(1024)
	if err != nil {
		return nil, err
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return transform.NewReader(content, e.NewDecoder()), nil
}

// GetHrefUrls 获取页面a标签href值
func GetHrefUrls(selector string, doc *goquery.Document) []string {
	var urls []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if len(url) > 0 {
			urls = append(urls, url)
		}
	})
	return urls
}
