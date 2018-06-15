package qimai

import (
	"context"
	"fmt"
	"log"

	"crawler.github.com/models"
	"github.com/chromedp/chromedp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	baseURL = "https://www.qimai.cn/andapp/baseinfo/appid/"
	count   = 5
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

	// var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		// wg.Add(1)

		task(c, &ctxt, page, db)

		page++
	}
	// wg.Wait()

	// remove chrome instance
	if err := closeChrome(c, &ctxt); err != nil {
		log.Panicf("remove chrome instance error: %s", err)
	}

}

// openChrome create chrome instance
func openChrome(ctxt context.Context) (*chromedp.CDP, error) {
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
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

// task parse web page
func task(c *chromedp.CDP, ctxt *context.Context, page int, db *gorm.DB) {
	result := models.APP{}

	// defer wg.Done()

	err := c.Run(*ctxt, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf(`%s%d`, baseURL, page)),
		chromedp.WaitVisible(`#appMain`, chromedp.ByID),
		chromedp.Text(`.breadcrumb-wrap>ul li:nth-child(5)`, &result.Name, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.app-header .auther .value`, &result.Company, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.app-header .genre .value`, &result.Category, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.app-header .appid .value`, &result.PkgName, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.app-baseinfo .baseinfo-list li:nth-child(4) .info`, &result.Version, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.app-baseinfo .baseinfo-list li:nth-child(5) .info`, &result.Size, chromedp.NodeVisible, chromedp.ByQuery),
	})

	fmt.Println(fmt.Sprintf(`====%s%d`, baseURL, page))
	fmt.Println(&result)

	if err != nil {
		log.Printf("get web page error: %s", err)
		return
	}

	done := make(chan bool)

	go save(result, db, done)

	<-done
}

// 写入mysql
func save(result models.APP, db *gorm.DB, done chan bool) error {
	if err := db.Save(&result).Error; err != nil {
		return err
	}
	done <- true
	return nil
}
