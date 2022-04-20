package forum

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Zdfx struct {
	name,
	baseUrl,
	cookie string
}

func (zdfx *Zdfx) FormHash() (string, bool) {
	return FormHash(zdfx)
}

func (zdfx *Zdfx) Name() string {
	return zdfx.name
}

func (zdfx *Zdfx) BasicUrl() string {
	return zdfx.baseUrl
}

func (zdfx *Zdfx) Cookie() string {
	return zdfx.cookie
}

func NewZdfxClient() Sign {
	cookie := os.Getenv("ZDFX_COOKIE")
	name := "终点"
	baseUrl := "https://bbs.zdfx.net/"
	if len(cookie) == 0 {
		return NewNoCookieClient(name)
	}
	client := Zdfx{
		name,
		baseUrl,
		cookie,
	}
	return &client
}

func (zdfx *Zdfx) Sign() (<-chan string, bool) {
	if len(zdfx.cookie) == 0 {
		c := make(chan string, 1)
		c <- zdfx.name + "Cookie未设置！"
		close(c)
		return c, false
	}
	c := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		wg.Wait()
		close(c)
	}()

	go func() {
		ctx, _ := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		zdfx.signInternal(ctx, c)
		cancel()
		wg.Done()
	}()

	go func() {
		ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
		zdfx.lottery(ctx, c)
		cancel()
		wg.Done()
	}()
	return c, true
}

func (zdfx *Zdfx) cookieSlice() chromedp.Action {
	cookies := strings.Split(zdfx.cookie, ";")
	slice := make([]string, 0)
	for _, c := range cookies {
		if len(c) == 0 {
			continue
		}
		params := strings.Trim(c, " ")
		slice = append(slice, strings.Split(params, "=")...)
	}
	return chromedp.ActionFunc(
		func(ctx context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for i := 0; i < len(slice); i += 2 {
				err := network.SetCookie(slice[i], slice[i+1]).
					WithExpires(&expr).
					WithDomain("bbs.zdfx.net").
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func (zdfx *Zdfx) signInternal(ctx context.Context, c chan<- string) {
	err := chromedp.Run(ctx,
		zdfx.cookieSlice(),
		chromedp.Navigate(`https://bbs.zdfx.net/k_misign-sign.html`),
		chromedp.Click(`#JD_sign`),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			c <- "已签到"
		} else {
			c <- err.Error()
		}
	}
}

func (zdfx *Zdfx) lottery(ctx context.Context, c chan<- string) {
	var res string
	err := chromedp.Run(ctx,
		zdfx.cookieSlice(),
		chromedp.Navigate(zdfx.baseUrl+`plugin.php?id=yinxingfei_zzza:yaoyao`),
		chromedp.Click(`.num_box > .btn`, chromedp.NodeVisible),
		chromedp.Sleep(1500*time.Millisecond),
		chromedp.InnerHTML(`div #res`, &res),
	)
	if err != nil {
		c <- err.Error()
	} else {
		c <- res
	}
}
