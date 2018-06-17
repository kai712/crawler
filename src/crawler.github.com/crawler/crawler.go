package crawler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"crawler.github.com/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
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

// ParseJSON 解析json
// returnValue type 根据返回json自定义struct
func ParseJSON(url string, returnValue interface{}) error {
	resp, err := GetHTML(url)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp)
	bodystr := string(body)

	if err := json.Unmarshal([]byte(bodystr), &returnValue); err != nil {
		return err
	}

	return nil
}

// Save 存入mysql
func Save(result *models.APP, db *gorm.DB) error {
	// 有则改，无则增
	if err := db.Where("name = ?", result.Name).Assign(result).FirstOrCreate(&result).Error; err != nil {
		return err
	}
	return nil
}
