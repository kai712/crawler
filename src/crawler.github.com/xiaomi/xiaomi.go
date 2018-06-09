package xiaomi

import (
	"fmt"
	"log"

	"crawler.github.com/crawler"
	"crawler.github.com/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	baseURL  = "http://app.mi.com"
	username = ""
	password = ""
	dbname   = ""
)

// Work 启动器
func Work() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbname))
	if err != nil {
		log.Panic("open database error:", err)
	}
	defer db.Close()

	log.Printf("start crawler xiaomi apps")

	cateURLs := GetCategoryURLs()
	done := make(chan bool, len(cateURLs))

	for _, url := range cateURLs {
		page := 0

		go func(page int, url string) {
			for {
				newURL := baseURL + url + "#page=" + fmt.Sprintf("%d", page)
				doc, err := GetDoc(newURL)
				if err != nil {
					log.Printf("Get list page error:%s", err)
					continue
				}
				content := doc.Find("#all-applist").Text()
				if content == "" {
					break
				}

				detailURLs := crawler.GetHrefUrls("#all-applist li h5 a", doc)

				for _, detailURL := range detailURLs {
					doc, err := GetDoc(baseURL + detailURL)
					if err != nil {
						log.Printf("Get detail page error:%s", err)
						continue
					}
					data := GetDetails(doc)
					db.Save(&data)
				}
				page++
			}
			done <- true
		}(page, url)
	}
	defer log.Printf("end crawler xiaomi apps")
	for i := 0; i < len(cateURLs); i++ {
		<-done
	}
}

// GetDoc 获取document
func GetDoc(url string) (*goquery.Document, error) {
	resp, err := crawler.GetHTML(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// GetCategoryURLs 获取分类urls
func GetCategoryURLs() []string {
	doc, err := GetDoc(baseURL)
	if err != nil {
		log.Panic(err)
	}
	// 返回分类urls
	return crawler.GetHrefUrls(".category-list li a", doc)
}

// GetDetails 获取应用信息
func GetDetails(doc *goquery.Document) (result *models.APP) {

	name := doc.Find(".app-info h3").Text()
	pkgName := doc.Find(".details ul.cf li").Eq(7).Text()
	category := doc.Find(".bread-crumb ul li").Eq(1).Find("a").Text()
	size := doc.Find(".details ul.cf li").Eq(1).Text()
	version := doc.Find(".details ul.cf li").Eq(3).Text()
	company := doc.Find(".app-info p").Eq(0).Text()

	result = &models.APP{
		Name:     name,
		PkgName:  pkgName,
		Category: category,
		Size:     size,
		Version:  version,
		Company:  company,
	}

	return result
}
