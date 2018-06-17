package qimai

import (
	"context"
	"fmt"
	"log"

	"crawler.github.com/crawler"
	"crawler.github.com/models"
	"github.com/chromedp/chromedp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	baseURL = "https://www.qimai.cn/andapp/baseinfo/appid/"
	count   = 20
)

var page = 1

// Work start
func Work(db *gorm.DB) {

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := openChrome(ctxt)
	if err != nil {
		log.Panicf("create chrome instance error: %s", err)
	}

	for {
		// open chrome target
		// c, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)))
		// if err != nil {
		// 	log.Panicf("open chrome target error: %s", err)
		// }

		if isBreak := task(c, &ctxt, page, db); isBreak {
			break
		}

		page++
	}

	// remove chrome instance
	if err := closeChrome(c, &ctxt); err != nil {
		log.Panicf("remove chrome instance error: %s", err)
	}

}

// openChrome create chrome instance
func openChrome(ctxt context.Context) (*chromedp.CDP, error) {
	c, err := chromedp.New(ctxt)
	if err != nil {
		return c, err
	}
	return c, nil
}

// closeChrome remove chrome instance
func closeChrome(c *chromedp.CDP, ctxt *context.Context) error {
	// shutdown chrome
	if err := c.Shutdown(*ctxt); err != nil {
		return err
	}
	// wait for chrome to finish
	if err := c.Wait(); err != nil {
		return err
	}
	return nil
}

// task
func task(c *chromedp.CDP, ctxt *context.Context, page int, db *gorm.DB) (isBreak bool) {
	isBreak = false
	result := models.APP{}
	url := fmt.Sprintf(`%s%d`, baseURL, page)

	err := c.Run(*ctxt, parse(url, &result))

	if err != nil {
		log.Printf("get web page error: %s", err)
		return
	}

	if result.Name == "" {
		isBreak = true
		return
	}

	fmt.Println(&result)

	if err := crawler.Save(&result, db); err != nil {
		log.Printf("save to mysql error: %s", err)
	}
	return
}

// parse 解析页面
func parse(url string, result *models.APP) chromedp.Tasks {
	ok := true
	fmt.Println(url)
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.breadcrumb-wrap>ul li:nth-child(5)`, chromedp.ByQuery),
		chromedp.Text(`.breadcrumb-wrap>ul li:nth-child(5)`, &result.Name, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(`.app-header .auther .value`, chromedp.ByQuery),
		chromedp.Text(`.app-header .auther .value`, &result.Company, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(`.app-header .genre .value`, chromedp.ByQuery),
		chromedp.Text(`.app-header .genre .value`, &result.Category, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(`.app-baseinfo .baseinfo-list li:nth-child(4) .info`, chromedp.ByQuery),
		chromedp.Text(`.app-baseinfo .baseinfo-list li:nth-child(4) .info`, &result.Version, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(`.app-baseinfo .baseinfo-list li:nth-child(5) .info`, chromedp.ByQuery),
		chromedp.Text(`.app-baseinfo .baseinfo-list li:nth-child(5) .info`, &result.Size, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.WaitVisible(`.app-icon`, chromedp.ByQuery),
		chromedp.AttributeValue(`.app-icon`, `src`, &result.Img, &ok, chromedp.NodeVisible, chromedp.ByQuery),
	}
}
