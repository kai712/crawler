package wandoujia

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"crawler.github.com/crawler"
	"crawler.github.com/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	baseURL = "http://www.wandoujia.com/category/app"
)

// Work 启动器
func Work(db *gorm.DB) {

	log.Printf("start crawler wandoujia APP")

	cateURLs, err := getCategoryURLs()
	if err != nil {
		log.Panicf("get category urls error: %s", err)
	}

	var wg sync.WaitGroup

	for _, url := range cateURLs {
		wg.Add(1)
		page := 1
		go parseList(page, url, &wg, db)
	}

	wg.Wait()

	log.Printf("done")
}

// getCategoryURLs 获取分类首页urls
func getCategoryURLs() ([]string, error) {
	doc, err := getDoc(baseURL)
	if err != nil {
		return nil, err
	}
	// 返回分类urls
	return crawler.GetHrefUrls(".tag-box .parent-cate .cate-link", doc), nil
}

// 解析每个分类下的列表页
func parseList(page int, url string, wg *sync.WaitGroup, db *gorm.DB) {
	for {
		newURL := fmt.Sprintf("%s/%d", url, page)
		doc, err := getDoc(newURL)
		if err != nil {
			log.Printf("Get list page error:%s", err)
			continue
		}
		content := doc.Find(".app-box clearfix").Text()
		if strings.Contains(content, "已经没有内容啦") {
			break
		}

		detailURLs := crawler.GetHrefUrls(".card .app-desc .name", doc)
		parsePage(detailURLs, db)
		page++
	}
	wg.Done()
}

// 解析每页应用详情
func parsePage(detailURLs []string, db *gorm.DB) {
	for _, detailURL := range detailURLs {
		doc, err := getDoc(detailURL)
		if err != nil {
			log.Printf("Get detail page error:%s", err)
			continue
		}
		data := getDetails(doc)
		log.Println(data)

		// if err := crawler.Save(data, db); err != nil {
		// 	log.Printf("save to mysql error: %s", err)
		// }
	}
}

// getDetails 获取应用信息
func getDetails(doc *goquery.Document) (result *models.APP) {

	name := doc.Find(".detail-top .app-info .title").Text()
	category := doc.Find(".infos .infos-list .tag-box").Eq(0).Find("a").Eq(1).Text()
	size := doc.Find(".infos .infos-list dd").Eq(0).Text()
	version := doc.Find(".infos .infos-list dd").Eq(4).Text()
	company := doc.Find(".infos .infos-list dd").Eq(6).Text()
	img, _ := doc.Find(".detail-top .app-icon img").Attr("src")

	result = &models.APP{
		Name:     name,
		Category: category,
		Size:     size,
		Version:  version,
		Company:  company,
		Img:      img,
	}

	return result
}

// getDoc 获取document
func getDoc(url string) (*goquery.Document, error) {
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
